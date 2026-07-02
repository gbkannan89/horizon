package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/rules/internal/domain"
)

type PostgresRuleRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRuleRepository(pool *pgxpool.Pool) *PostgresRuleRepository {
	return &PostgresRuleRepository{pool: pool}
}

func (r *PostgresRuleRepository) Save(ctx context.Context, rule *domain.Rule) error {
	id := rule.ID()
	conds, _ := json.Marshal(rule.Conditions())
	acts, _ := json.Marshal(rule.Actions())
	if id == "" {
		err := r.pool.QueryRow(ctx, `SELECT gen_random_uuid()::text`).Scan(&id)
		if err != nil {
			return fmt.Errorf("generate id: %w", err)
		}
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO rules (id, user_id, household_id, name, description, category, priority, enabled, conditions, actions, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (id) DO UPDATE SET
			name=$4, description=$5, category=$6, priority=$7, enabled=$8,
			conditions=$9, actions=$10, updated_at=$12`,
		id, rule.UserID(), nullOrEmpty(rule.HouseholdID()), rule.Name(), rule.Description(),
		rule.Category(), rule.Priority(), rule.Enabled(), conds, acts,
		rule.CreatedAt(), rule.UpdatedAt(),
	)
	return err
}

func (r *PostgresRuleRepository) GetByID(ctx context.Context, id string) (*domain.Rule, error) {
	var (
		ruleID, userID, name, desc, category, createdAt, updatedAt string
		householdID, condsJSON, actsJSON                           *string
		priority                                                   int
		enabled                                                    bool
	)
	err := r.pool.QueryRow(ctx, `SELECT id, user_id, household_id, name, description, category, priority, enabled, conditions, actions, created_at, updated_at FROM rules WHERE id=$1`, id).
		Scan(&ruleID, &userID, &householdID, &name, &desc, &category, &priority, &enabled, &condsJSON, &actsJSON, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return scanRule(ruleID, userID, householdID, name, desc, category, priority, enabled, condsJSON, actsJSON, createdAt, updatedAt)
}

func (r *PostgresRuleRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Rule, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, household_id, name, description, category, priority, enabled, conditions, actions, created_at, updated_at FROM rules WHERE user_id=$1 ORDER BY priority`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []*domain.Rule
	for rows.Next() {
		var (
			ruleID, name, desc, cat, createdAt, updatedAt string
			hhID, condsJSON, actsJSON                    *string
			priority                                      int
			enabled                                       bool
		)
		if err := rows.Scan(&ruleID, &userID, &hhID, &name, &desc, &cat, &priority, &enabled, &condsJSON, &actsJSON, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		r, err := scanRule(ruleID, userID, hhID, name, desc, cat, priority, enabled, condsJSON, actsJSON, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func (r *PostgresRuleRepository) ListByCategory(ctx context.Context, category string) ([]*domain.Rule, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, household_id, name, description, category, priority, enabled, conditions, actions, created_at, updated_at FROM rules WHERE category=$1 ORDER BY priority`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []*domain.Rule
	for rows.Next() {
		var (
			ruleID, userID, name, desc, cat, createdAt, updatedAt string
			hhID, condsJSON, actsJSON                            *string
			priority                                              int
			enabled                                               bool
		)
		if err := rows.Scan(&ruleID, &userID, &hhID, &name, &desc, &cat, &priority, &enabled, &condsJSON, &actsJSON, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		r, err := scanRule(ruleID, userID, hhID, name, desc, cat, priority, enabled, condsJSON, actsJSON, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func (r *PostgresRuleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM rules WHERE id=$1`, id)
	return err
}

func scanRule(id, userID string, householdID *string, name, desc, category string, priority int, enabled bool, condsJSON, actsJSON *string, createdAt, updatedAt string) (*domain.Rule, error) {
	var conds []domain.Condition
	if condsJSON != nil && *condsJSON != "" {
		json.Unmarshal([]byte(*condsJSON), &conds)
	}
	var acts []domain.Action
	if actsJSON != nil && *actsJSON != "" {
		json.Unmarshal([]byte(*actsJSON), &acts)
	}
	hhID := ""
	if householdID != nil {
		hhID = *householdID
	}
	return domain.NewRule(id, userID, hhID, name, desc, category, priority, enabled, conds, acts, createdAt, updatedAt)
}

func nullOrEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
