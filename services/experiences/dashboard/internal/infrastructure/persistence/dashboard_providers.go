package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/dashboard/internal/aggregator"
)

// DashboardPGProvider implements all dashboard data providers using PostgreSQL.
type DashboardPGProvider struct {
	pool *pgxpool.Pool
}

// Ensure interface compliance.
var _ aggregator.AccountProvider = (*DashboardPGProvider)(nil)
var _ aggregator.AssetProvider = (*DashboardPGProvider)(nil)
var _ aggregator.LiabilityProvider = (*DashboardPGProvider)(nil)
var _ aggregator.GoalProvider = (*DashboardPGProvider)(nil)
var _ aggregator.EventProvider = (*DashboardPGProvider)(nil)
var _ aggregator.PortfolioProvider = (*DashboardPGProvider)(nil)
var _ aggregator.HealthProvider = (*DashboardPGProvider)(nil)
var _ aggregator.RiskProvider = (*DashboardPGProvider)(nil)
var _ aggregator.RecommendationProvider = (*DashboardPGProvider)(nil)
var _ aggregator.ProjectionProvider = (*DashboardPGProvider)(nil)
var _ aggregator.SimulationProvider = (*DashboardPGProvider)(nil)

func NewDashboardPGProvider(pool *pgxpool.Pool) *DashboardPGProvider {
	return &DashboardPGProvider{pool: pool}
}

// AccountProvider implementation

func (p *DashboardPGProvider) GetCashBalance(ctx context.Context, userID string) (int64, error) {
	var balance int64
	err := p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM financial_events
		WHERE user_id::text = $1 AND state = 'POSTED' AND event_type IN ('Salary', 'Income', 'Credit')
		AND effective_date >= DATE_TRUNC('month', CURRENT_DATE)`, userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("get cash balance: %w", err)
	}
	return balance, nil
}

func (p *DashboardPGProvider) HasAccounts(ctx context.Context, userID string) (bool, error) {
	var count int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM accounts WHERE owner_id::text = $1 AND status = 'Active'`, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("has accounts: %w", err)
	}
	return count > 0, nil
}

// AssetProvider implementation

func (p *DashboardPGProvider) GetTotalAssets(ctx context.Context, userID string) (int64, error) {
	var total int64
	err := p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(COALESCE(quantity, 1) * unit_price), 0) FROM assets
		WHERE owner_id::text = $1 AND status = 'Active'`, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get total assets: %w", err)
	}
	return total, nil
}

// LiabilityProvider implementation

func (p *DashboardPGProvider) GetTotalLiabilities(ctx context.Context, userID string) (int64, error) {
	var total int64
	err := p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(original_principal), 0) FROM liabilities
		WHERE owner_id::text = $1 AND status = 'Active'`, userID).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("get total liabilities: %w", err)
	}
	return total, nil
}

// GoalProvider implementation

func (p *DashboardPGProvider) GetGoalCounts(ctx context.Context, userID string) (onTrack, atRisk, total int, hasGoals bool, err error) {
	var goalCount int
	err = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1`, userID).Scan(&goalCount)
	if err != nil {
		return 0, 0, 0, false, fmt.Errorf("get goal counts: %w", err)
	}
	if goalCount == 0 {
		return 0, 0, 0, false, nil
	}

	var activeCount int
	err = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1 AND status = 'Active'`, userID).Scan(&activeCount)
	if err != nil {
		return 0, 0, 0, true, fmt.Errorf("get active goals: %w", err)
	}

	var atRiskCount int
	err = p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM goals WHERE user_id::text = $1 AND status = 'AtRisk'`, userID).Scan(&atRiskCount)
	if err != nil {
		return 0, 0, 0, true, fmt.Errorf("get at-risk goals: %w", err)
	}

	onTrack = activeCount - atRiskCount
	if onTrack < 0 { onTrack = 0 }
	return onTrack, atRiskCount, goalCount, true, nil
}

// EventProvider implementation

func (p *DashboardPGProvider) GetRecentCount(ctx context.Context, userID string) (int, bool, error) {
	var count int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM financial_events
		WHERE user_id::text = $1 AND effective_date >= NOW() - INTERVAL '30 days'`, userID).Scan(&count)
	if err != nil {
		return 0, false, fmt.Errorf("get recent events: %w", err)
	}
	return count, count > 0, nil
}

// PortfolioProvider implementation

func (p *DashboardPGProvider) GetPortfolioValue(ctx context.Context, userID string) (int64, error) {
	var value int64
	err := p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(a.unit_price * COALESCE(a.quantity, 1)), 0)
		FROM portfolio_members pm
		JOIN portfolios pf ON pm.portfolio_id = pf.portfolio_id
		JOIN assets a ON pm.entity_id::uuid = a.asset_id
		WHERE pf.owner_id::text = $1 AND pm.entity_type = 'Asset' AND pf.status = 'Active'`, userID).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("get portfolio value: %w", err)
	}
	return value, nil
}

// HealthProvider implementation

func (p *DashboardPGProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var grade string
	err := p.pool.QueryRow(ctx,
		`SELECT overall_score, score_grade FROM health_scores
		WHERE score_id = (SELECT MAX(score_id) FROM health_scores WHERE score_id LIKE $1 || '%')
		ORDER BY created_at DESC LIMIT 1`,
		safePrefix(userID, 8), // prefix match on score_id
	).Scan(&score, &grade)
	if err != nil {
		// Return default values if no health score exists
		return 0, "Pending", nil
	}
	return score, grade, nil
}

// RiskProvider implementation

func (p *DashboardPGProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var level string
	err := p.pool.QueryRow(ctx,
		`SELECT composite_score, risk_level FROM risk_assessments
		ORDER BY created_at DESC LIMIT 1`).Scan(&score, &level)
	if err != nil {
		return 0, "Minimal", nil
	}
	return score, level, nil
}

// RecommendationProvider implementation

func (p *DashboardPGProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM recommendations`).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("get recommendations: %w", err)
	}
	return count > 0, count, nil
}

// ProjectionProvider implementation

func (p *DashboardPGProvider) GetMonthlyIncomeExpenses(ctx context.Context, userID string) (income, expenses int64, err error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var inc, exp int64
	incErr := p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM financial_events
		WHERE user_id::text = $1 AND state = 'POSTED' AND amount > 0
		AND effective_date >= $2`, userID, startOfMonth).Scan(&inc)
	if incErr != nil {
		return 0, 0, fmt.Errorf("get monthly income: %w", incErr)
	}

	expErr := p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(ABS(amount)), 0) FROM financial_events
		WHERE user_id::text = $1 AND state = 'POSTED' AND amount < 0
		AND effective_date >= $2`, userID, startOfMonth).Scan(&exp)
	if expErr != nil {
		return 0, 0, fmt.Errorf("get monthly expenses: %w", expErr)
	}

	return inc, exp, nil
}

// SimulationProvider implementation

func (p *DashboardPGProvider) HasSimulations(ctx context.Context, userID string) (bool, error) {
	var count int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM simulations`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("get simulations: %w", err)
	}
	return count > 0, nil
}

func safePrefix(s string, n int) string {
	if len(s) <= n { return s }
	return s[:n]
}
