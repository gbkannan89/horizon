package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/notifications/internal/engine"
)

type PGStateRepository struct {
	pool *pgxpool.Pool
}

func NewPGStateRepository(pool *pgxpool.Pool) *PGStateRepository {
	return &PGStateRepository{pool: pool}
}

func (r *PGStateRepository) Get(ctx context.Context, userID, notifID string) (engine.State, bool) {
	var state engine.State
	var snoozeUntil *time.Time
	
	err := r.pool.QueryRow(ctx, "SELECT state, snooze_until FROM notification_state WHERE user_id = $1 AND notif_id = $2", userID, notifID).Scan(&state, &snoozeUntil)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false
		}
		// Log error in a real app, returning false for now
		return "", false
	}
	
	if state == engine.StateSnoozed && snoozeUntil != nil {
		if time.Now().After(*snoozeUntil) {
			// Snooze expired, return unread
			return engine.StateUnread, true
		}
	}
	
	return state, true
}

func (r *PGStateRepository) SetState(ctx context.Context, userID, notifID string, s engine.State) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO notification_state (user_id, notif_id, state, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, notif_id)
		DO UPDATE SET state = EXCLUDED.state, snooze_until = NULL, updated_at = NOW()
	`, userID, notifID, s)
	return err
}

func (r *PGStateRepository) Snooze(ctx context.Context, userID, notifID string, until string) error {
	snoozeTime, err := time.Parse(time.RFC3339, until)
	if err != nil {
		snoozeTime = time.Now().Add(24 * time.Hour)
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO notification_state (user_id, notif_id, state, snooze_until, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, notif_id)
		DO UPDATE SET state = EXCLUDED.state, snooze_until = EXCLUDED.snooze_until, updated_at = NOW()
	`, userID, notifID, engine.StateSnoozed, snoozeTime)
	return err
}

type PGPreferenceRepository struct {
	pool *pgxpool.Pool
}

func NewPGPreferenceRepository(pool *pgxpool.Pool) *PGPreferenceRepository {
	return &PGPreferenceRepository{pool: pool}
}

func (r *PGPreferenceRepository) GetPreferences(ctx context.Context, userID string) ([]engine.Preference, error) {
	rows, err := r.pool.Query(ctx, "SELECT category, in_app, push, email, min_priority, enabled FROM notification_preferences WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefs []engine.Preference
	for rows.Next() {
		var p engine.Preference
		if err := rows.Scan(&p.Category, &p.InApp, &p.Push, &p.Email, &p.MinPriority, &p.Enabled); err != nil {
			return nil, err
		}
		prefs = append(prefs, p)
	}
	return prefs, nil
}

func (r *PGPreferenceRepository) SavePreferences(ctx context.Context, userID string, prefs []engine.Preference) error {
	batch := &pgx.Batch{}
	
	for _, p := range prefs {
		batch.Queue(`
			INSERT INTO notification_preferences (user_id, category, in_app, push, email, min_priority, enabled, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			ON CONFLICT (user_id, category)
			DO UPDATE SET
				in_app = EXCLUDED.in_app,
				push = EXCLUDED.push,
				email = EXCLUDED.email,
				min_priority = EXCLUDED.min_priority,
				enabled = EXCLUDED.enabled,
				updated_at = NOW()
		`, userID, p.Category, p.InApp, p.Push, p.Email, p.MinPriority, p.Enabled)
	}
	
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	
	for i := 0; i < len(prefs); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}
