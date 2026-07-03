package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/household/internal/domain"
)

type HouseholdRepository struct {
	pool *pgxpool.Pool
}

func NewHouseholdRepository(pool *pgxpool.Pool) *HouseholdRepository {
	return &HouseholdRepository{pool: pool}
}

func (r *HouseholdRepository) Save(ctx context.Context, h *domain.Household) error {
	membersJSON, err := json.Marshal(h.Members())
	if err != nil {
		return fmt.Errorf("marshal members: %w", err)
	}

	linkedAccountsJSON, _ := json.Marshal(h.LinkedAccounts())
	linkedGoalsJSON, _ := json.Marshal(h.LinkedGoals())
	linkedBudgetsJSON, _ := json.Marshal(h.LinkedBudgets())
	goalContributionsJSON, _ := json.Marshal(h.GoalContributions())

	query := `INSERT INTO households (id, name, household_type, head_of_household_id, members, status, currency, country,
		total_assets, total_liabilities, total_net_worth, health, tags, notes, 
		linked_accounts, linked_goals, linked_budgets, goal_contributions, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)
		ON CONFLICT (id) DO UPDATE SET
			name=$2, household_type=$3, head_of_household_id=$4, members=$5, status=$6,
			currency=$7, country=$8, total_assets=$9, total_liabilities=$10, total_net_worth=$11,
			health=$12, tags=$13, notes=$14, linked_accounts=$15, linked_goals=$16, 
			linked_budgets=$17, goal_contributions=$18, updated_at=$20`

	now := time.Now().UTC()
	if h.CreatedAt().IsZero() {
		id, err := r.generateID(ctx)
		if err != nil {
			return fmt.Errorf("generate id: %w", err)
		}
		*h = *domain.ReconstructFromDB(
			id, h.Name(), h.HouseholdType(), h.HeadOfHouseholdID(),
			h.Members(), h.Status(), h.Currency(), h.Country(),
			h.TotalAssets(), h.TotalLiabilities(), h.TotalNetWorth(),
			h.Health(), h.Tags(), h.Notes(),
			h.LinkedAccounts(), h.LinkedGoals(), h.LinkedBudgets(), h.GoalContributions(),
			now, now,
		)
		membersJSON, _ = json.Marshal(h.Members())
	}

	tagsJSON, _ := json.Marshal(h.Tags())
	if tagsJSON == nil { tagsJSON = []byte("[]") }

	_, err = r.pool.Exec(ctx, query,
		h.ID(), h.Name(), string(h.HouseholdType()), h.HeadOfHouseholdID(),
		membersJSON, string(h.Status()), h.Currency(), h.Country(),
		h.TotalAssets(), h.TotalLiabilities(), h.TotalNetWorth(),
		string(h.Health()), tagsJSON, h.Notes(),
		linkedAccountsJSON, linkedGoalsJSON, linkedBudgetsJSON, goalContributionsJSON,
		h.CreatedAt(), h.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("exec save: %w", err)
	}
	return nil
}

func (r *HouseholdRepository) UpdateStatus(ctx context.Context, id string, from, to domain.HouseholdStatus) error {
	query := `UPDATE households SET status=$2, updated_at=$3 WHERE id=$1 AND status=$4`
	tag, err := r.pool.Exec(ctx, query, id, string(to), time.Now().UTC(), string(from))
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("household not found or status mismatch")
	}
	return nil
}

func (r *HouseholdRepository) GetByID(ctx context.Context, id string) (*domain.Household, error) {
	query := `SELECT id, name, household_type, head_of_household_id, members, status, currency, country,
		total_assets, total_liabilities, total_net_worth, health, tags, notes, 
		linked_accounts, linked_goals, linked_budgets, goal_contributions, created_at, updated_at
		FROM households WHERE id=$1`

	row := r.pool.QueryRow(ctx, query, id)
	return r.scanRow(row)
}

func (r *HouseholdRepository) ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Household, string, error) {
	if limit <= 0 { limit = 25 }
	if limit > 100 { limit = 100 }

	query := `SELECT id, name, household_type, head_of_household_id, members, status, currency, country,
		total_assets, total_liabilities, total_net_worth, health, tags, notes,
		linked_accounts, linked_goals, linked_budgets, goal_contributions, created_at, updated_at
		FROM households
		WHERE members @> $1::jsonb
		AND ($2 = '' OR id > $2)
		ORDER BY id
		LIMIT $3`

	memberFilter, _ := json.Marshal([]map[string]string{{"user_id": userID}})

	rows, err := r.pool.Query(ctx, query, memberFilter, cursor, limit+1)
	if err != nil {
		return nil, "", fmt.Errorf("list by user: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows, limit)
}

func (r *HouseholdRepository) ListByStatus(ctx context.Context, status domain.HouseholdStatus, cursor string, limit int) ([]*domain.Household, string, error) {
	if limit <= 0 { limit = 25 }
	if limit > 100 { limit = 100 }

	query := `SELECT id, name, household_type, head_of_household_id, members, status, currency, country,
		total_assets, total_liabilities, total_net_worth, health, tags, notes,
		linked_accounts, linked_goals, linked_budgets, goal_contributions, created_at, updated_at
		FROM households
		WHERE status=$1 AND ($2 = '' OR id > $2)
		ORDER BY id
		LIMIT $3`

	rows, err := r.pool.Query(ctx, query, string(status), cursor, limit+1)
	if err != nil {
		return nil, "", fmt.Errorf("list by status: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows, limit)
}



func (r *HouseholdRepository) generateID(ctx context.Context) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `SELECT gen_random_uuid()::text`).Scan(&id)
	return id, err
}

func (r *HouseholdRepository) scanRow(row pgx.Row) (*domain.Household, error) {
	var (
		id, name, hhType, headID, status, currency, country, health string
		membersJSON, tagsJSON                                       []byte
		linkedAccountsJSON, linkedGoalsJSON, linkedBudgetsJSON      []byte
		goalContributionsJSON                                       []byte
		notes                                                        string
		totalAssets, totalLiabilities, totalNetWorth                 int64
		createdAt, updatedAt                                         time.Time
	)

	err := row.Scan(&id, &name, &hhType, &headID, &membersJSON, &status,
		&currency, &country, &totalAssets, &totalLiabilities, &totalNetWorth,
		&health, &tagsJSON, &notes, 
		&linkedAccountsJSON, &linkedGoalsJSON, &linkedBudgetsJSON, &goalContributionsJSON,
		&createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("household not found")
		}
		return nil, fmt.Errorf("scan: %w", err)
	}

	var members []domain.HouseholdMember
	if len(membersJSON) > 0 {
		json.Unmarshal(membersJSON, &members)
	}
	if members == nil { members = []domain.HouseholdMember{} }

	var tags []string
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &tags)
	}
	if tags == nil { tags = []string{} }

	var linkedAccounts []domain.LinkedAccount
	if len(linkedAccountsJSON) > 0 { json.Unmarshal(linkedAccountsJSON, &linkedAccounts) }
	
	var linkedGoals []domain.LinkedGoal
	if len(linkedGoalsJSON) > 0 { json.Unmarshal(linkedGoalsJSON, &linkedGoals) }
	
	var linkedBudgets []domain.LinkedBudget
	if len(linkedBudgetsJSON) > 0 { json.Unmarshal(linkedBudgetsJSON, &linkedBudgets) }
	
	var goalContribs []domain.GoalContribution
	if len(goalContributionsJSON) > 0 { json.Unmarshal(goalContributionsJSON, &goalContribs) }

	return domain.ReconstructFromDB(
		id, name, domain.HouseholdType(hhType), headID, members,
		domain.HouseholdStatus(status), currency, country,
		totalAssets, totalLiabilities, totalNetWorth,
		domain.HouseholdHealth(health), tags, notes,
		linkedAccounts, linkedGoals, linkedBudgets, goalContribs,
		createdAt, updatedAt,
	), nil
}

func (r *HouseholdRepository) scanRows(rows pgx.Rows, limit int) ([]*domain.Household, string, error) {
	var households []*domain.Household
	for rows.Next() {
		h, err := r.scanRow(rows)
		if err != nil {
			return nil, "", err
		}
		households = append(households, h)
	}

	var cursor string
	if len(households) > limit {
		cursor = households[len(households)-1].ID()
		households = households[:limit]
	}

	return households, cursor, nil
}

// Ensure compliance
var _ domain.Repository = (*HouseholdRepository)(nil)
