package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/portfolio/internal/engine"
)

// PortfolioPGProvider implements all portfolio data providers using PostgreSQL.
type PortfolioPGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewPortfolioPGProvider(pool *pgxpool.Pool) *PortfolioPGProvider {
	return &PortfolioPGProvider{pool: pool, now: time.Now}
}

// PortfolioProvider

func (p *PortfolioPGProvider) GetPortfolioSummary(ctx context.Context, userID string) (value int64, totalReturn float64, returnPct float64, err error) {
	err = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(a.unit_price * COALESCE(a.quantity, 1)), 0),
			COALESCE(SUM(a.unit_price * COALESCE(a.quantity, 1) - a.cost_basis), 0)
		FROM portfolios pf
		JOIN portfolio_members pm ON pf.portfolio_id = pm.portfolio_id
		JOIN assets a ON pm.entity_id = a.asset_id
		WHERE pf.owner_id = $1 AND pm.entity_type = 'Asset' AND pf.status = 'Active'`, userID).Scan(&value, &totalReturn)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("get portfolio summary: %w", err)
	}
	if value > 0 && totalReturn > 0 {
		costBasis := float64(value) - totalReturn
		if costBasis > 0 {
			returnPct = (totalReturn / costBasis) * 100.0
		}
	}
	return
}

// AssetProvider

func (p *PortfolioPGProvider) GetAllocations(ctx context.Context, userID string) ([]engine.AllocEntry, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT a.classification, SUM(a.unit_price * COALESCE(a.quantity, 1)) as total_value
		FROM portfolios pf
		JOIN portfolio_members pm ON pf.portfolio_id = pm.portfolio_id
		JOIN assets a ON pm.entity_id = a.asset_id
		WHERE pf.owner_id = $1 AND pm.entity_type = 'Asset' AND pf.status = 'Active'
		GROUP BY a.classification
		ORDER BY total_value DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("get allocations: %w", err)
	}
	defer rows.Close()

	var total int64
	var entries []engine.AllocEntry
	for rows.Next() {
		var label string
		var value int64
		if err := rows.Scan(&label, &value); err != nil {
			return nil, fmt.Errorf("scan allocation: %w", err)
		}
		total += value
		entries = append(entries, engine.AllocEntry{Label: label, Value: value})
	}

	for i := range entries {
		if total > 0 {
			entries[i].Percent = float64(entries[i].Value) / float64(total) * 100.0
		}
	}

	return entries, nil
}

// PerfProvider

func (p *PortfolioPGProvider) GetPerformance(ctx context.Context, userID string) (periodReturn, periodReturnPct, benchmarkRet float64, unrealizedGL, realizedGL int64, err error) {
	pfValue, _, _, err := p.GetPortfolioSummary(ctx, userID)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}
	unrealizedGL = pfValue

	now := p.now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var income, expenses int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events
		WHERE user_id = $1 AND effective_date >= $2 AND state = 'POSTED'`,
		userID, startOfMonth).Scan(&income, &expenses)

	periodReturn = float64(income - expenses)
	if float64(pfValue) > 0 {
		periodReturnPct = (periodReturn / float64(pfValue)) * 100.0
	}

	return periodReturn, periodReturnPct, 0, unrealizedGL, 0, nil
}

// RiskProvider

func (p *PortfolioPGProvider) GetPortfolioRisk(ctx context.Context, userID string) (score int, level string, sharpe, vol, mdd, valueAtRisk float64, err error) {
	err = p.pool.QueryRow(ctx,
		`SELECT composite_score, risk_level FROM risk_assessments ORDER BY created_at DESC LIMIT 1`).Scan(&score, &level)
	if err != nil {
		return 0, "Minimal", 0, 0, 0, 0, nil
	}
	return score, level, 0, 0, 0, 0, nil
}

// ProjProvider

func (p *PortfolioPGProvider) GetPortfolioProjection(ctx context.Context, userID string) (projectedVal float64, confidence string, horizonYears int, annualReturn float64, err error) {
	err = p.pool.QueryRow(ctx,
		`SELECT (output_data->>'projected_value')::numeric,
			COALESCE(output_data->>'confidence', 'Medium'),
			COALESCE((output_data->>'horizon_years')::int, 0),
			COALESCE((output_data->>'annual_return')::numeric, 0)
		FROM projection_outputs
		WHERE projection_type = 'GoalCompletion'
		ORDER BY created_at DESC LIMIT 1`).Scan(&projectedVal, &confidence, &horizonYears, &annualReturn)
	if err != nil {
		return 0, "Medium", 0, 0, nil
	}
	return
}

// RecProvider

func (p *PortfolioPGProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM recommendations`).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("count recommendations: %w", err)
	}
	return count > 0, count, nil
}

// OptProvider

func (p *PortfolioPGProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM optimizations`).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("count optimizations: %w", err)
	}
	return count > 0, count, nil
}

// SimProvider

func (p *PortfolioPGProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM simulations`).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("count simulations: %w", err)
	}
	return count > 0, count, nil
}

// EventProvider

func (p *PortfolioPGProvider) GetEventCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM financial_events
		WHERE user_id = $1 AND effective_date >= NOW() - INTERVAL '30 days'`, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count events: %w", err)
	}
	return count, nil
}

func (p *PortfolioPGProvider) GetMilestoneCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
