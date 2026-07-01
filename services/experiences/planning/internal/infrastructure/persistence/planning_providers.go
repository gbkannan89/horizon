package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/planning/internal/engine"
)

// PlanningPGProvider implements all planning data providers using PostgreSQL.
type PlanningPGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewPlanningPGProvider(pool *pgxpool.Pool) *PlanningPGProvider {
	return &PlanningPGProvider{pool: pool, now: time.Now}
}

// GoalProvider

func (p *PlanningPGProvider) GetGoalSummary(ctx context.Context, userID string) (onTrack, total int, fundingGap int64, err error) {
	err = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1 AND status NOT IN ('Archived')`, userID).Scan(&total)
	if err != nil || total == 0 {
		return 0, 0, 0, nil
	}

	_ = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1 AND status = 'Active'`, userID).Scan(&onTrack)

	var gap int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM((success_criteria->>'target_value')::numeric - COALESCE(a.alloc, 0)), 0)
		FROM goals g
		LEFT JOIN (
			SELECT goal_id, SUM(allocated_amount) as alloc FROM allocations
			WHERE status IN ('Active','Approved') GROUP BY goal_id
		) a ON g.goal_id = a.goal_id
		WHERE g.user_id::text = $1::text AND g.status NOT IN ('Archived')`, userID).Scan(&gap)
	if gap < 0 {
		gap = 0
	}

	return onTrack, total, gap, nil
}

// AccountProvider

func (p *PlanningPGProvider) GetCashBalance(ctx context.Context, userID string) (int64, int64, int64, error) {
	var cash int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(unit_price * COALESCE(quantity, 1)), 0) FROM assets
		WHERE owner_id::text = $1 AND classification IN ('CashAndCashEquivalent', 'FixedDeposit') AND status = 'Active'`, userID).Scan(&cash)

	var income, expenses int64
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id::text = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)

	return cash, income, expenses, nil
}

// ProjProvider

func (p *PlanningPGProvider) GetPlanningProjection(ctx context.Context, userID string) (netWorth, income, expenses, portfolio int64, conf string, err error) {
	// Net worth: assets - liabilities
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(unit_price * COALESCE(quantity, 1)), 0) FROM assets WHERE owner_id::text = $1 AND status = 'Active'`, userID).Scan(&netWorth)

	var liabilities int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(original_principal), 0) FROM liabilities WHERE owner_id::text = $1 AND status = 'Active'`, userID).Scan(&liabilities)
	netWorth -= liabilities

	// Portfolio value
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(a.unit_price * COALESCE(a.quantity, 1)), 0)
		FROM portfolio_members pm JOIN assets a ON pm.entity_id = a.asset_id
		JOIN portfolios pf ON pm.portfolio_id = pf.portfolio_id
		WHERE pf.owner_id::text = $1 AND pm.entity_type = 'Asset'`, userID).Scan(&portfolio)

	// Monthly income/expenses
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id::text = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)

	conf = "Medium"
	// Try to find confidence from projection
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(output_data->>'confidence', 'Medium') FROM projection_outputs ORDER BY created_at DESC LIMIT 1`).Scan(&conf)

	return
}

// RiskProvider

func (p *PlanningPGProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var level string
	err := p.pool.QueryRow(ctx,
		`SELECT composite_score, risk_level FROM risk_assessments ORDER BY created_at DESC LIMIT 1`).Scan(&score, &level)
	if err != nil {
		return 0, "Minimal", nil
	}
	return score, level, nil
}

// HealthProvider

func (p *PlanningPGProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var grade string
	err := p.pool.QueryRow(ctx,
		`SELECT overall_score, score_grade FROM health_scores ORDER BY created_at DESC LIMIT 1`).Scan(&score, &grade)
	if err != nil {
		return 0, "Pending", nil
	}
	return score, grade, nil
}

// RecProvider

func (p *PlanningPGProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM recommendations`).Scan(&count)
	if err != nil {
		return false, 0, nil
	}
	return count > 0, count, nil
}

// OptProvider

func (p *PlanningPGProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM optimizations`).Scan(&count)
	if err != nil {
		return false, 0, nil
	}
	return count > 0, count, nil
}

// SimProvider

func (p *PlanningPGProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM simulations`).Scan(&count)
	if err != nil {
		return false, 0, nil
	}
	return count > 0, count, nil
}

// EventProvider

func (p *PlanningPGProvider) GetEventCount(ctx context.Context, userID string) (int, int, error) {
	var evtCnt int
	_ = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM financial_events WHERE user_id::text = $1 AND effective_date >= NOW() - INTERVAL '30 days'`, userID).Scan(&evtCnt)
	return evtCnt, 0, nil
}

// ScenarioProvider

func (p *PlanningPGProvider) GetScenarios(ctx context.Context, userID string) ([]engine.Scenario, error) {
	return nil, nil
}
