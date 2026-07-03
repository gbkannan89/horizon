package persistence

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DataProvider struct {
	pool *pgxpool.Pool
}

func NewDataProvider(pool *pgxpool.Pool) *DataProvider {
	return &DataProvider{pool: pool}
}

func (p *DataProvider) queryJSON(ctx context.Context, query string, args ...any) []map[string]interface{} {
	row := p.pool.QueryRow(ctx, query, args...)
	var result []byte
	if err := row.Scan(&result); err != nil {
		return []map[string]interface{}{}
	}
	if result == nil {
		return []map[string]interface{}{}
	}
	var data []map[string]interface{}
	if err := json.Unmarshal(result, &data); err != nil {
		return []map[string]interface{}{}
	}
	if data == nil {
		return []map[string]interface{}{}
	}
	return data
}

func (p *DataProvider) GetAccounts(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(accounts.*) ORDER BY account_name), '[]'::json) FROM accounts WHERE owner_id = $1`, userID)
}

func (p *DataProvider) GetTransactions(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(financial_events.*) ORDER BY effective_date DESC), '[]'::json) FROM financial_events WHERE user_id = $1`, userID)
}

func (p *DataProvider) GetGoals(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(goals.*) ORDER BY priority), '[]'::json) FROM goals WHERE user_id = $1`, userID)
}

func (p *DataProvider) GetAllocations(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(a.*) ORDER BY a.created_at), '[]'::json) FROM allocations a JOIN goals g ON a.goal_id = g.id WHERE g.user_id = $1`, userID)
}

func (p *DataProvider) GetAssets(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(assets.*) ORDER BY asset_name), '[]'::json) FROM assets WHERE owner_id = $1`, userID)
}

func (p *DataProvider) GetLiabilities(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(liabilities.*) ORDER BY created_at), '[]'::json) FROM liabilities WHERE owner_id = $1`, userID)
}

func (p *DataProvider) GetPortfolios(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(portfolios.*) ORDER BY name), '[]'::json) FROM portfolios WHERE owner_id = $1`, userID)
}

func (p *DataProvider) GetInstitutions(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(i.*) ORDER BY i.name), '[]'::json) FROM institutions i WHERE EXISTS (SELECT 1 FROM accounts a WHERE a.institution_id = i.id AND a.owner_id = $1)`, userID)
}

func (p *DataProvider) GetHouseholds(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(households.*) ORDER BY name), '[]'::json) FROM households WHERE members @> $1::jsonb`,
		[]map[string]string{{"user_id": userID}})
}

func (p *DataProvider) GetHealthScores(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(h.*) ORDER BY h.calculated_at DESC), '[]'::json) FROM health_scores h WHERE h.user_id = $1`, userID)
}

func (p *DataProvider) GetRiskAssessments(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(r.*) ORDER BY r.assessed_at DESC), '[]'::json) FROM risk_assessments r WHERE r.user_id = $1`, userID)
}

func (p *DataProvider) GetRecommendations(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(r.*) ORDER BY r.created_at DESC), '[]'::json) FROM recommendations r WHERE r.user_id = $1`, userID)
}

func (p *DataProvider) GetOptimizations(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(o.*) ORDER BY o.created_at DESC), '[]'::json) FROM optimizations o WHERE o.user_id = $1`, userID)
}

func (p *DataProvider) GetProjections(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(p.*) ORDER BY p.created_at DESC), '[]'::json) FROM projection_outputs p WHERE p.user_id = $1`, userID)
}

func (p *DataProvider) GetSimulations(ctx context.Context, userID string) []map[string]interface{} {
	return p.queryJSON(ctx,
		`SELECT COALESCE(json_agg(to_json(s.*) ORDER BY s.created_at DESC), '[]'::json) FROM simulations s WHERE s.user_id = $1`, userID)
}
