package infrastructure

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGDataProvider struct {
	pool *pgxpool.Pool
}

func NewPGDataProvider(pool *pgxpool.Pool) *PGDataProvider {
	return &PGDataProvider{pool: pool}
}

func (p *PGDataProvider) FetchSectionData(ctx context.Context, userID string, sectionType string) (map[string]interface{}, error) {
	switch sectionType {
	case "net_worth":
		return p.fetchNetWorth(ctx, userID)
	case "income_expense":
		return p.fetchIncomeExpense(ctx, userID)
	case "budget":
		return p.fetchBudgetOverview(ctx, userID)
	case "goals":
		return p.fetchGoalProgress(ctx, userID)
	default:
		return map[string]interface{}{"section": sectionType, "data": nil}, nil
	}
}

func (p *PGDataProvider) fetchNetWorth(ctx context.Context, userID string) (map[string]interface{}, error) {
	var assets, liabilities float64
	p.pool.QueryRow(ctx, `SELECT COALESCE(SUM(CASE WHEN type='asset' THEN current_balance ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN type='liability' THEN current_balance ELSE 0 END),0)
		FROM accounts WHERE user_id::text=$1`, userID).Scan(&assets, &liabilities)
	return map[string]interface{}{
		"total_assets":   assets,
		"total_liabilities": liabilities,
		"net_worth":      assets - liabilities,
	}, nil
}

func (p *PGDataProvider) fetchIncomeExpense(ctx context.Context, userID string) (map[string]interface{}, error) {
	var income, expenses float64
	p.pool.QueryRow(ctx, `SELECT COALESCE(SUM(CASE WHEN amount>0 THEN amount ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN amount<0 THEN ABS(amount) ELSE 0 END),0)
		FROM financial_events WHERE user_id::text=$1 AND event_date >= now() - interval '30 days'`, userID).Scan(&income, &expenses)
	return map[string]interface{}{
		"total_income":   income,
		"total_expenses": expenses,
		"net_savings":    income - expenses,
	}, nil
}

func (p *PGDataProvider) fetchBudgetOverview(ctx context.Context, userID string) (map[string]interface{}, error) {
	var budgeted, spent float64
	p.pool.QueryRow(ctx, `SELECT COALESCE(SUM(total_budgeted),0), COALESCE(SUM(total_spent),0)
		FROM budgets WHERE user_id::text=$1 AND status='Active'`, userID).Scan(&budgeted, &spent)
	return map[string]interface{}{
		"total_budgeted": budgeted,
		"total_spent":    spent,
		"remaining":      budgeted - spent,
	}, nil
}

func (p *PGDataProvider) fetchGoalProgress(ctx context.Context, userID string) (map[string]interface{}, error) {
	rows, err := p.pool.Query(ctx, `SELECT name, progress, status FROM goals WHERE user_id::text=$1 ORDER BY priority LIMIT 10`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var goals []map[string]interface{}
	for rows.Next() {
		var name, status string
		var progress float64
		if err := rows.Scan(&name, &progress, &status); err != nil {
			continue
		}
		goals = append(goals, map[string]interface{}{
			"name": name, "progress": progress, "status": status,
		})
	}
	return map[string]interface{}{"goals": goals, "total": len(goals)}, nil
}
