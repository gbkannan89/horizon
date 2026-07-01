package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/notifications/internal/engine"
)

// NotifsPGProvider implements all notification data providers using PostgreSQL.
type NotifsPGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewNotifsPGProvider(pool *pgxpool.Pool) *NotifsPGProvider {
	return &NotifsPGProvider{pool: pool, now: time.Now}
}

// HealthProvider

func (p *NotifsPGProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var grade string
	err := p.pool.QueryRow(ctx, `SELECT overall_score, score_grade FROM health_scores ORDER BY created_at DESC LIMIT 1`).Scan(&score, &grade)
	if err != nil { return 0, "Pending", nil }
	return score, grade, nil
}

func (p *NotifsPGProvider) GetHealthChange(ctx context.Context, userID string) (int, error) {
	var scores []int
	rows, err := p.pool.Query(ctx, `SELECT overall_score FROM health_scores ORDER BY created_at DESC LIMIT 2`)
	if err != nil { return 0, nil }
	defer rows.Close()
	for rows.Next() {
		var s int
		rows.Scan(&s)
		scores = append(scores, s)
	}
	if len(scores) < 2 { return 0, nil }
	return scores[0] - scores[1], nil
}

// RiskProvider

func (p *NotifsPGProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var level string
	err := p.pool.QueryRow(ctx, `SELECT composite_score, risk_level FROM risk_assessments ORDER BY created_at DESC LIMIT 1`).Scan(&score, &level)
	if err != nil { return 0, "Minimal", nil }
	return score, level, nil
}

// ProjProvider

func (p *NotifsPGProvider) GetCashFlowSurplus(ctx context.Context, userID string) (int64, error) {
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var income, expenses int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)
	return income - expenses, nil
}

func (p *NotifsPGProvider) GetNetWorthChange(ctx context.Context, userID string) (int64, error) {
	now := p.now()
	som := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var income, expenses int64
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0)
		FROM financial_events WHERE user_id = $1 AND state = 'POSTED' AND effective_date >= $2`, userID, som).Scan(&income, &expenses)
	return income - expenses, nil
}

// GoalProvider

func (p *NotifsPGProvider) GetGoalProgress(ctx context.Context, userID string) (onTrack, total, atRisk int, err error) {
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND status NOT IN ('Archived')`, userID).Scan(&total)
	if total == 0 { return 0, 0, 0, nil }
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND status = 'Active'`, userID).Scan(&onTrack)
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND status = 'AtRisk'`, userID).Scan(&atRisk)
	return onTrack, total, atRisk, nil
}

// RecProvider

func (p *NotifsPGProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM recommendations`).Scan(&count)
	if err != nil { return false, 0, nil }
	return count > 0, count, nil
}

// OptProvider

func (p *NotifsPGProvider) HasOptimizations(ctx context.Context, userID string) (bool, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM optimizations`).Scan(&count)
	if err != nil { return false, nil }
	return count > 0, nil
}

// EventProvider

func (p *NotifsPGProvider) GetAchievementCount(ctx context.Context, userID string) (int, error) {
	var count int
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM goals WHERE user_id = $1 AND status = 'Completed'`, userID).Scan(&count)
	return count, nil
}

// PrefProvider

func (p *NotifsPGProvider) GetPreferences(ctx context.Context, userID string) ([]engine.Preference, error) {
	return []engine.Preference{}, nil
}
