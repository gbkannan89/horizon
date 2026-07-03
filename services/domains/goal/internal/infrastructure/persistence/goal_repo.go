package persistence

import (
	"context"
	"fmt"
	"time"

	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/goal/internal/domain"
)

type GoalRepository struct {
	pool *pgxpool.Pool
}

func NewGoalRepository(pool *pgxpool.Pool) *GoalRepository {
	return &GoalRepository{pool: pool}
}

func (r *GoalRepository) Save(ctx context.Context, g *domain.Goal) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO goals (goal_id, user_id, household_id, name, importance, type, subtype, priority, status, notes, tags, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW())
		ON CONFLICT (goal_id) DO UPDATE SET priority=$8, status=$9, updated_at=NOW()`,
		g.ID(), g.UserID(), g.HouseholdID(), g.Name(), string(g.Importance()), string(g.GoalType()),
		string(g.Subtype()), g.Priority(), string(g.Status()), g.Notes(), g.Tags())
	return err
}

func (r *GoalRepository) UpdateStatus(ctx context.Context, goalID string, fromStatus, toStatus domain.GoalStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE goals SET status=$1, updated_at=NOW() WHERE goal_id=$2 AND status=$3`,
		string(toStatus), goalID, string(fromStatus))
	return err
}

func (r *GoalRepository) UpdatePriority(ctx context.Context, goalID string, newPriority int) error {
	_, err := r.pool.Exec(ctx, `UPDATE goals SET priority=$1, updated_at=NOW() WHERE goal_id=$2`, newPriority, goalID)
	return err
}

func (r *GoalRepository) GetByID(ctx context.Context, goalID string) (*domain.Goal, error) {
	var id, uid, name, imp, gt, st string
	var hhid *string
	var pri int
	var ca, ua time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at, updated_at
		FROM goals WHERE goal_id = $1`, goalID).Scan(&id, &uid, &hhid, &name, &imp, &gt, &pri, &st, &ca, &ua)
	if err != nil { return nil, fmt.Errorf("get goal: %w", err) }

	return domain.ReconstructFromDB(id, uid, hhid, name, domain.GoalImportance(imp), domain.GoalType(gt), "",
		domain.SuccessCriteria{}, pri, domain.GoalStatus(st), "", nil, ca, ua, nil), nil
}

func (r *GoalRepository) ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Goal, string, error) {
	if limit <= 0 { limit = 25 }
	return scanGoalList(r.pool, ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at
		FROM goals WHERE user_id = $1 ORDER BY priority LIMIT $2`, userID, limit)
}

func (r *GoalRepository) ListByStatus(ctx context.Context, userID string, status domain.GoalStatus, cursor string, limit int) ([]*domain.Goal, string, error) {
	if limit <= 0 { limit = 25 }
	return scanGoalList(r.pool, ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at
		FROM goals WHERE user_id = $1 AND status = $2 ORDER BY priority LIMIT $3`,
		userID, string(status), limit)
}

func (r *GoalRepository) ListByImportance(ctx context.Context, userID string, importance domain.GoalImportance, cursor string, limit int) ([]*domain.Goal, string, error) {
	if limit <= 0 { limit = 25 }
	return scanGoalList(r.pool, ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at
		FROM goals WHERE user_id = $1 AND importance = $2 ORDER BY priority LIMIT $3`,
		userID, string(importance), limit)
}

func (r *GoalRepository) ListByType(ctx context.Context, userID string, goalType domain.GoalType, cursor string, limit int) ([]*domain.Goal, string, error) {
	if limit <= 0 { limit = 25 }
	return scanGoalList(r.pool, ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at
		FROM goals WHERE user_id = $1 AND type = $2 ORDER BY priority LIMIT $3`,
		userID, string(goalType), limit)
}

func (r *GoalRepository) ListByHousehold(ctx context.Context, householdID string, cursor string, limit int) ([]*domain.Goal, string, error) {
	if limit <= 0 { limit = 25 }
	return scanGoalList(r.pool, ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at
		FROM goals WHERE household_id = $1 ORDER BY priority LIMIT $2`,
		householdID, limit)
}

func (r *GoalRepository) GetActiveByPriority(ctx context.Context, userID string) ([]*domain.Goal, error) {
	goals, _, _ := scanGoalList(r.pool, ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at
		FROM goals WHERE user_id = $1 AND status = 'Active' ORDER BY priority`, userID, 100)
	return goals, nil
}

func (r *GoalRepository) GetMaxPriority(ctx context.Context, userID string) (int, error) {
	var max int
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(priority), 0) FROM goals WHERE user_id = $1 AND status = 'Active'`, userID).Scan(&max)
	return max, err
}

func (r *GoalRepository) FindConflict(ctx context.Context, userID string, priority int, excludeGoalID string) (*domain.Goal, error) {
	var id, uid, name, imp, gt, st string
	var hhid *string
	var pri int
	var ca time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT goal_id, user_id, household_id, name, importance, type, priority, status, created_at
		FROM goals WHERE user_id = $1 AND priority = $2 AND status = 'Active' AND goal_id != $3 LIMIT 1`,
		userID, priority, excludeGoalID).Scan(&id, &uid, &hhid, &name, &imp, &gt, &pri, &st, &ca)
	if err != nil { return nil, nil }
	return domain.ReconstructFromDB(id, uid, hhid, name, domain.GoalImportance(imp), domain.GoalType(gt), "",
		domain.SuccessCriteria{}, pri, domain.GoalStatus(st), "", nil, ca, ca, nil), nil
}

func (r *GoalRepository) RebalancePriorities(ctx context.Context, userID string, afterPriority int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE goals SET priority = priority - 1, updated_at = NOW()
		WHERE user_id = $1 AND priority > $2 AND status = 'Active'`, userID, afterPriority)
	return err
}

func scanGoalList(pool *pgxpool.Pool, ctx context.Context, query string, args ...interface{}) ([]*domain.Goal, string, error) {
	args[len(args)-1] = args[len(args)-1].(int) + 1
	rows, _ := pool.Query(ctx, query, args...)
	defer rows.Close()

	var goals []*domain.Goal
	for rows.Next() {
		var id, uid, name, imp, gt, st string
		var hhid *string
		var pri int
		var ca time.Time
		rows.Scan(&id, &uid, &hhid, &name, &imp, &gt, &pri, &st, &ca)
		goals = append(goals, domain.ReconstructFromDB(id, uid, hhid, name, domain.GoalImportance(imp),
			domain.GoalType(gt), "", domain.SuccessCriteria{}, pri, domain.GoalStatus(st), "", nil, ca, ca, nil))
	}
	hasMore := len(goals) > args[len(args)-1].(int)-1
	if hasMore { goals = goals[:len(goals)-1] }
	return goals, "", nil
}

