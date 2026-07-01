package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AdvisorPGProvider implements all advisor data providers using PostgreSQL.
type AdvisorPGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewAdvisorPGProvider(pool *pgxpool.Pool) *AdvisorPGProvider {
	return &AdvisorPGProvider{pool: pool, now: time.Now}
}

// DashProvider

func (p *AdvisorPGProvider) GetNetWorth(ctx context.Context, userID string) (int64, error) {
	var assets int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(unit_price * COALESCE(quantity, 1)), 0) FROM assets WHERE owner_id = $1 AND status = 'Active'`, userID).Scan(&assets)
	var liabilities int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(original_principal), 0) FROM liabilities WHERE owner_id = $1 AND status = 'Active'`, userID).Scan(&liabilities)
	return assets - liabilities, nil
}

func (p *AdvisorPGProvider) GetCashBalance(ctx context.Context, userID string) (int64, error) {
	var cash int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(unit_price * COALESCE(quantity, 1)), 0) FROM assets
		WHERE owner_id = $1 AND classification IN ('CashAndCashEquivalent', 'FixedDeposit') AND status = 'Active'`, userID).Scan(&cash)
	return cash, nil
}

func (p *AdvisorPGProvider) GetMonthlyFlow(ctx context.Context, userID string) (int64, int64, error) {
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var income, expenses int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)
	return income, expenses, nil
}

// GoalProvider

func (p *AdvisorPGProvider) GetGoalProgress(ctx context.Context, userID string) (int, int, int, int64, error) {
	var total int
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND status NOT IN ('Archived')`, userID).Scan(&total)
	if total == 0 {
		return 0, 0, 0, 0, nil
	}
	var onTrack int
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND status = 'Active'`, userID).Scan(&onTrack)
	var atRisk int
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND status = 'AtRisk'`, userID).Scan(&atRisk)
	var gap int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM((success_criteria->>'target_value')::numeric - COALESCE(a.alloc, 0)), 0)
		FROM goals g LEFT JOIN (SELECT goal_id, SUM(allocated_amount) as alloc FROM allocations WHERE status IN ('Active','Approved') GROUP BY goal_id) a ON g.goal_id = a.goal_id
		WHERE g.user_id = $1 AND g.status NOT IN ('Archived')`, userID).Scan(&gap)
	if gap < 0 { gap = 0 }
	return onTrack, atRisk, total, gap, nil
}

// AcctProvider

func (p *AdvisorPGProvider) GetAccountSummary(ctx context.Context, userID string) (int, int64, error) {
	var count int
	var balance int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(a.unit_price * COALESCE(a.quantity, 1)), 0)
		FROM accounts ac LEFT JOIN assets a ON ac.account_id = a.account_id
		WHERE ac.owner_id = $1 AND ac.status NOT IN ('Closed', 'Archived')`, userID).Scan(&count, &balance)
	return count, balance, nil
}

// PfProvider

func (p *AdvisorPGProvider) GetPortfolioSummary(ctx context.Context, userID string) (int64, float64, float64, error) {
	var value int64
	var costBasis float64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(a.unit_price * COALESCE(a.quantity, 1)), 0), COALESCE(SUM(a.cost_basis), 0)
		FROM portfolios pf JOIN portfolio_members pm ON pf.portfolio_id = pm.portfolio_id
		JOIN assets a ON pm.entity_id = a.asset_id
		WHERE pf.owner_id = $1 AND pm.entity_type = 'Asset' AND pf.status = 'Active'`, userID).Scan(&value, &costBasis)
	retPct := 0.0
	if costBasis > 0 {
		retPct = (float64(value) - costBasis) / costBasis * 100
	}
	return value, float64(value) - costBasis, retPct, nil
}

// HealthProvider

func (p *AdvisorPGProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var grade string
	err := p.pool.QueryRow(ctx, `SELECT overall_score, score_grade FROM health_scores ORDER BY created_at DESC LIMIT 1`).Scan(&score, &grade)
	if err != nil { return 0, "Pending", nil }
	return score, grade, nil
}

func (p *AdvisorPGProvider) GetHealthChange(ctx context.Context, userID string) (int, error) {
	var scores []int
	rows, err := p.pool.Query(ctx, `SELECT overall_score FROM health_scores ORDER BY created_at DESC LIMIT 2`)
	if err != nil { return 0, nil }
	defer rows.Close()
	for rows.Next() {
		var s int
		rows.Scan(&s)
		scores = append(scores, s)
	}
	if len(scores) < 2 { return 0, nil }
	return scores[0] - scores[1], nil
}

// RiskProvider

func (p *AdvisorPGProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var level string
	err := p.pool.QueryRow(ctx, `SELECT composite_score, risk_level FROM risk_assessments ORDER BY created_at DESC LIMIT 1`).Scan(&score, &level)
	if err != nil { return 0, "Minimal", nil }
	return score, level, nil
}

// ProjProvider

func (p *AdvisorPGProvider) GetProjectionSummary(ctx context.Context, userID string) (int64, bool, string, error) {
	var nw int64
	var assets int64
	var liabilities int64
	_ = p.pool.QueryRow(ctx, `SELECT COALESCE(SUM(unit_price * COALESCE(quantity, 1)), 0) FROM assets WHERE owner_id = $1 AND status = 'Active'`, userID).Scan(&assets)
	_ = p.pool.QueryRow(ctx, `SELECT COALESCE(SUM(original_principal), 0) FROM liabilities WHERE owner_id = $1 AND status = 'Active'`, userID).Scan(&liabilities)
	nw = assets - liabilities

	conf := "Medium"
	_ = p.pool.QueryRow(ctx, `SELECT COALESCE(output_data->>'confidence', 'Medium') FROM projection_outputs ORDER BY created_at DESC LIMIT 1`).Scan(&conf)
	return nw, nw > 0, conf, nil
}

// RecProvider

func (p *AdvisorPGProvider) GetRecommendationSummary(ctx context.Context, userID string) (bool, int, string, error) {
	var count int
	var title string
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*), COALESCE((SELECT rec_data->>'title' FROM recommendations ORDER BY priority LIMIT 1), '') FROM recommendations`).Scan(&count, &title)
	if err != nil { return false, 0, "", nil }
	return count > 0, count, title, nil
}

// OptProvider

func (p *AdvisorPGProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM optimizations`).Scan(&count)
	if err != nil { return false, 0, nil }
	return count > 0, count, nil
}

// SimProvider

func (p *AdvisorPGProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM simulations`).Scan(&count)
	if err != nil { return false, 0, nil }
	return count > 0, count, nil
}

// TimelineProvider

func (p *AdvisorPGProvider) GetEventCount(ctx context.Context, userID string) (int, error) {
	var count int
	_ = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM financial_events WHERE user_id = $1 AND effective_date >= NOW() - INTERVAL '30 days'`, userID).Scan(&count)
	return count, nil
}

// InsightProvider

func (p *AdvisorPGProvider) GetInsightSummary(ctx context.Context, userID string) (int, int, error) {
	return 0, 0, nil
}

// NotifProvider

func (p *AdvisorPGProvider) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
