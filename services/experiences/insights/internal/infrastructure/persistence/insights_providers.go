package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InsightsPGProvider implements all insight data providers using PostgreSQL.
type InsightsPGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewInsightsPGProvider(pool *pgxpool.Pool) *InsightsPGProvider {
	return &InsightsPGProvider{pool: pool, now: time.Now}
}

// HealthProvider

func (p *InsightsPGProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var grade string
	err := p.pool.QueryRow(ctx,
		`SELECT overall_score, score_grade FROM health_scores ORDER BY created_at DESC LIMIT 1`).Scan(&score, &grade)
	if err != nil {
		return 0, "Pending", nil
	}
	return score, grade, nil
}

func (p *InsightsPGProvider) GetHealthChange(ctx context.Context, userID string) (int, error) {
	var scores []int
	rows, err := p.pool.Query(ctx,
		`SELECT overall_score FROM health_scores ORDER BY created_at DESC LIMIT 2`)
	if err != nil {
		return 0, nil
	}
	defer rows.Close()
	for rows.Next() {
		var s int
		rows.Scan(&s)
		scores = append(scores, s)
	}
	if len(scores) < 2 {
		return 0, nil
	}
	return scores[0] - scores[1], nil
}

// RiskProvider

func (p *InsightsPGProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var level string
	err := p.pool.QueryRow(ctx,
		`SELECT composite_score, risk_level FROM risk_assessments ORDER BY created_at DESC LIMIT 1`).Scan(&score, &level)
	if err != nil {
		return 0, "Minimal", nil
	}
	return score, level, nil
}

// ProjProvider

func (p *InsightsPGProvider) GetSavingsRate(ctx context.Context, userID string) (float64, string, error) {
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var income, expenses float64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id::text = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)

	if income <= 0 {
		return 0, "stable", nil
	}
	rate := (income - expenses) / income * 100
	if rate < 0 {
		rate = 0
	}

	// Compare with previous month for trend
	prevStart := som.AddDate(0, -1, 0)
	var prevIncome, prevExpenses float64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id::text = $1 AND state = 'POSTED' AND effective_date >= $2 AND effective_date < $3`,
		userID, prevStart, som).Scan(&prevIncome, &prevExpenses)

	trend := "stable"
	if prevIncome > 0 {
		prevRate := (prevIncome - prevExpenses) / prevIncome * 100
		if rate > prevRate+2 {
			trend = "improving"
		} else if rate < prevRate-2 {
			trend = "declining"
		}
	}

	return rate, trend, nil
}

func (p *InsightsPGProvider) GetNetWorthChange(ctx context.Context, userID string) (int64, int64, error) {
	var netWorth int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(unit_price * COALESCE(quantity, 1)), 0) FROM assets WHERE owner_id::text = $1 AND status = 'Active'`, userID).Scan(&netWorth)

	var liabilities int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(original_principal), 0) FROM liabilities WHERE owner_id::text = $1 AND status = 'Active'`, userID).Scan(&liabilities)
	netWorth -= liabilities

	// Net worth change: income - expenses from current month
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var income, expenses int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id::text = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)

	change := income - expenses
	return change, netWorth, nil
}

func (p *InsightsPGProvider) GetCashFlowSurplus(ctx context.Context, userID string) (int64, error) {
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var income, expenses int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id::text = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)

	return income - expenses, nil
}

// GoalProvider

func (p *InsightsPGProvider) GetGoalProgress(ctx context.Context, userID string) (int, int, float64, error) {
	var total int
	_ = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1 AND status NOT IN ('Archived')`, userID).Scan(&total)
	if total == 0 {
		return 0, 0, 0, nil
	}

	var onTrack int
	_ = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1 AND status = 'Active'`, userID).Scan(&onTrack)

	avgPct := float64(0)
	if total > 0 {
		avgPct = float64(onTrack) / float64(total) * 100
	}

	return onTrack, total, avgPct, nil
}

func (p *InsightsPGProvider) GetMilestoneCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

// PfProvider

func (p *InsightsPGProvider) GetPortfolioReturn(ctx context.Context, userID string) (int64, float64, error) {
	var value int64
	var costBasis float64

	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(unit_price * COALESCE(quantity, 1)), 0),
			COALESCE(SUM(cost_basis), 0)
		FROM assets a
		JOIN portfolio_members pm ON a.asset_id = pm.entity_id
		JOIN portfolios pf ON pm.portfolio_id = pf.portfolio_id
		WHERE pf.owner_id::text = $1 AND pm.entity_type = 'Asset'`, userID).Scan(&value, &costBasis)

	if costBasis > 0 {
		return value, (float64(value) - costBasis) / costBasis * 100, nil
	}
	return value, 0, nil
}

// RecProvider

func (p *InsightsPGProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM recommendations`).Scan(&count)
	if err != nil {
		return false, 0, nil
	}
	return count > 0, count, nil
}

// OptProvider

func (p *InsightsPGProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM optimizations`).Scan(&count)
	if err != nil {
		return false, 0, nil
	}
	return count > 0, count, nil
}

// SimProvider

func (p *InsightsPGProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM simulations`).Scan(&count)
	if err != nil {
		return false, 0, nil
	}
	return count > 0, count, nil
}

// EventProvider

func (p *InsightsPGProvider) GetAchievementCount(ctx context.Context, userID string) (int, error) {
	// Achievements are derived — count goals with Completed status as proxy
	var count int
	_ = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1 AND status = 'Completed'`, userID).Scan(&count)
	if count == 0 {
		// Also count milestones reached from simulations/projections
		count = 0
	}
	return count, nil
}

func (p *InsightsPGProvider) GetSpendingAnomaly(ctx context.Context, userID string) (string, error) {
	// Detect unusual spending: compare current month expenses to previous month
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	prevStart := som.AddDate(0, -1, 0)

	var currentExp, prevExp float64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(ABS(amount)), 0) FROM financial_events
		WHERE user_id::text = $1 AND amount < 0 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&currentExp)
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(ABS(amount)), 0) FROM financial_events
		WHERE user_id::text = $1 AND amount < 0 AND state = 'POSTED' AND effective_date >= $2 AND effective_date < $3`,
		userID, prevStart, som).Scan(&prevExp)

	if prevExp > 0 && currentExp > prevExp*1.3 {
		return "Spending this month is significantly higher than last month", nil
	}
	return "", nil
}

// AcctProvider

func (p *InsightsPGProvider) GetAccountCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM accounts WHERE owner_id::text = $1 AND status NOT IN ('Closed', 'Archived')`, userID).Scan(&count)
	if err != nil {
		return 0, nil
	}
	return count, nil
}
