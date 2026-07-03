package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/recurring/internal/engine"
)

type PostgresRecurringProvider struct {
	pool *pgxpool.Pool
}

func NewPostgresRecurringProvider(pool *pgxpool.Pool) *PostgresRecurringProvider {
	return &PostgresRecurringProvider{pool: pool}
}

func (p *PostgresRecurringProvider) GetDueByDate(ctx context.Context, date time.Time) ([]*engine.RecurringTransaction, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT recurring_id, user_id, next_occurrence, amount, currency, description, event_type, category 
		 FROM recurring_transactions 
		 WHERE status = 'active' AND next_occurrence <= $1`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*engine.RecurringTransaction
	for rows.Next() {
		var r engine.RecurringTransaction
		var next *time.Time
		if err := rows.Scan(&r.ID, &r.UserID, &next, &r.Amount, &r.Currency, &r.Description, &r.EventType, &r.Category); err != nil {
			return nil, err
		}
		r.NextOccurrence = next
		result = append(result, &r)
	}
	return result, nil
}

func (p *PostgresRecurringProvider) Save(ctx context.Context, r *engine.RecurringTransaction) error {
	_, err := p.pool.Exec(ctx,
		`UPDATE recurring_transactions SET next_occurrence = $1 WHERE recurring_id = $2`, r.NextOccurrence, r.ID)
	return err
}

func (p *PostgresRecurringProvider) CreateEvent(ctx context.Context, r *engine.RecurringTransaction, occurrenceDate time.Time) error {
	eventID := "tx-" + r.ID + "-" + occurrenceDate.Format("20060102")
	_, err := p.pool.Exec(ctx,
		`INSERT INTO financial_events (event_id, user_id, event_type, amount, currency, event_date, description, state, category, tags)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, 'POSTED', $8, $9)
		 ON CONFLICT (event_id) DO NOTHING`,
		eventID, r.UserID, r.EventType, r.Amount, r.Currency, occurrenceDate, r.Description, r.Category, []string{"recurring"})
	return err
}
