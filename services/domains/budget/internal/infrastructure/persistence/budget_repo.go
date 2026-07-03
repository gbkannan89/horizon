package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/budget/internal/domain"
)

type BudgetRepository struct {
	pool *pgxpool.Pool
}

func NewBudgetRepository(pool *pgxpool.Pool) *BudgetRepository {
	return &BudgetRepository{pool: pool}
}

func (r *BudgetRepository) Save(ctx context.Context, b *domain.Budget) error {
	catsJSON, _ := json.Marshal(b.Categories())
	tagsJSON, _ := json.Marshal(b.Tags())

	query := `INSERT INTO budgets (id, user_id, household_id, name, period, start_date, end_date,
		status, total_budgeted, total_spent, total_remaining, currency, categories, tags, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		ON CONFLICT (id) DO UPDATE SET
			name=$4, period=$5, start_date=$6, end_date=$7, status=$8,
			total_budgeted=$9, total_spent=$10, total_remaining=$11, currency=$12,
			categories=$13, tags=$14, updated_at=$16`

	now := time.Now().UTC()
	if b.CreatedAt().IsZero() {
		id, err := r.generateID(ctx)
		if err != nil {
			return fmt.Errorf("generate id: %w", err)
		}
		*b = *domain.ReconstructFromDB(
			id, b.UserID(), b.HouseholdID(), b.Name(), b.Period(),
			b.StartDate(), b.EndDate(), b.Status(),
			b.TotalBudgeted(), b.TotalSpent(), b.TotalRemaining(),
			b.Currency(), b.Categories(), b.Tags(), now, now,
		)
		catsJSON, _ = json.Marshal(b.Categories())
	}

	_, err := r.pool.Exec(ctx, query,
		b.ID(), b.UserID(), b.HouseholdID(), b.Name(), string(b.Period()),
		b.StartDate(), b.EndDate(), string(b.Status()),
		b.TotalBudgeted(), b.TotalSpent(), b.TotalRemaining(),
		b.Currency(), catsJSON, tagsJSON,
		b.CreatedAt(), b.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("exec save: %w", err)
	}
	return nil
}

func (r *BudgetRepository) UpdateStatus(ctx context.Context, id string, from, to domain.BudgetStatus) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE budgets SET status=$2, updated_at=$3 WHERE id=$1 AND status=$4`,
		id, string(to), time.Now().UTC(), string(from))
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("budget not found or status mismatch")
	}
	return nil
}

func (r *BudgetRepository) GetByID(ctx context.Context, id string) (*domain.Budget, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, user_id, COALESCE(household_id,''), name, period, start_date, end_date,
		status, total_budgeted, total_spent, total_remaining, currency, categories, tags, created_at, updated_at
		FROM budgets WHERE id=$1`, id)
	return r.scanRow(row)
}

func (r *BudgetRepository) ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Budget, string, error) {
	if limit <= 0 { limit = 25 }
	if limit > 100 { limit = 100 }
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(household_id,''), name, period, start_date, end_date,
		status, total_budgeted, total_spent, total_remaining, currency, categories, tags, created_at, updated_at
		FROM budgets WHERE user_id=$1 AND ($2='' OR id>$2)
		ORDER BY id LIMIT $3`, userID, cursor, limit+1)
	if err != nil {
		return nil, "", fmt.Errorf("list by user: %w", err)
	}
	defer rows.Close()
	return r.scanRows(rows, limit)
}

func (r *BudgetRepository) ListByHousehold(ctx context.Context, householdID string, cursor string, limit int) ([]*domain.Budget, string, error) {
	if limit <= 0 { limit = 25 }
	if limit > 100 { limit = 100 }
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(household_id,''), name, period, start_date, end_date,
		status, total_budgeted, total_spent, total_remaining, currency, categories, tags, created_at, updated_at
		FROM budgets WHERE household_id=$1 AND ($2='' OR id>$2)
		ORDER BY id LIMIT $3`, householdID, cursor, limit+1)
	if err != nil {
		return nil, "", fmt.Errorf("list by household: %w", err)
	}
	defer rows.Close()
	return r.scanRows(rows, limit)
}

func (r *BudgetRepository) ListByPeriod(ctx context.Context, userID string, start, end time.Time, cursor string, limit int) ([]*domain.Budget, string, error) {
	if limit <= 0 { limit = 25 }
	if limit > 100 { limit = 100 }
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(household_id,''), name, period, start_date, end_date,
		status, total_budgeted, total_spent, total_remaining, currency, categories, tags, created_at, updated_at
		FROM budgets WHERE user_id=$1 AND start_date>=$2 AND end_date<=$3 AND ($4='' OR id>$4)
		ORDER BY id LIMIT $5`, userID, start, end, cursor, limit+1)
	if err != nil {
		return nil, "", fmt.Errorf("list by period: %w", err)
	}
	defer rows.Close()
	return r.scanRows(rows, limit)
}

func (r *BudgetRepository) GetByCategory(ctx context.Context, userID, category string, cursor string, limit int) ([]*domain.Budget, string, error) {
	if limit <= 0 { limit = 25 }
	if limit > 100 { limit = 100 }
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, COALESCE(household_id,''), name, period, start_date, end_date,
		status, total_budgeted, total_spent, total_remaining, currency, categories, tags, created_at, updated_at
		FROM budgets WHERE user_id=$1 AND categories @> $2::jsonb AND ($3='' OR id>$3)
		ORDER BY id LIMIT $4`,
		userID, fmt.Sprintf(`[{"category": "%s"}]`, category), cursor, limit+1)
	if err != nil {
		return nil, "", fmt.Errorf("get by category: %w", err)
	}
	defer rows.Close()
	return r.scanRows(rows, limit)
}

func (r *BudgetRepository) GetBudgetVsActual(ctx context.Context, budgetID string) ([]domain.BudgetCategory, error) {
	b, err := r.GetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}
	return b.Categories(), nil
}

func (r *BudgetRepository) generateID(ctx context.Context) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `SELECT gen_random_uuid()::text`).Scan(&id)
	return id, err
}

func (r *BudgetRepository) scanRow(row pgx.Row) (*domain.Budget, error) {
	var (
		id, userID, householdID, name, period, status, currency string
		startDate, endDate                                       time.Time
		totalBudgeted, totalSpent, totalRemaining                int64
		catsJSON, tagsJSON                                       []byte
		createdAt, updatedAt                                     time.Time
	)
	err := row.Scan(&id, &userID, &householdID, &name, &period, &startDate, &endDate,
		&status, &totalBudgeted, &totalSpent, &totalRemaining,
		&currency, &catsJSON, &tagsJSON, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, fmt.Errorf("budget not found") }
		return nil, fmt.Errorf("scan: %w", err)
	}
	var categories []domain.BudgetCategory
	if len(catsJSON) > 0 { json.Unmarshal(catsJSON, &categories) }
	if categories == nil { categories = []domain.BudgetCategory{} }
	var tags []string
	if len(tagsJSON) > 0 { json.Unmarshal(tagsJSON, &tags) }
	if tags == nil { tags = []string{} }
	return domain.ReconstructFromDB(id, userID, householdID, name,
		domain.BudgetPeriod(period), startDate, endDate,
		domain.BudgetStatus(status),
		totalBudgeted, totalSpent, totalRemaining, currency,
		categories, tags, createdAt, updatedAt), nil
}

func (r *BudgetRepository) scanRows(rows pgx.Rows, limit int) ([]*domain.Budget, string, error) {
	var budgets []*domain.Budget
	for rows.Next() {
		b, err := r.scanRow(rows)
		if err != nil { return nil, "", err }
		budgets = append(budgets, b)
	}
	var cursor string
	if len(budgets) > limit {
		cursor = budgets[len(budgets)-1].ID()
		budgets = budgets[:limit]
	}
	return budgets, cursor, nil
}

var _ domain.Repository = (*BudgetRepository)(nil)
