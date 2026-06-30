package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRow struct {
	ID              string
	EventType       string
	EventPayload    []byte
	IdempotencyKey  string
	CorrelationID   string
	PartitionKey    string
	Status          string
	RetryCount      int
	CreatedAt       time.Time
	LastAttemptAt   *time.Time
}

type OutboxRepository struct {
	pool *pgxpool.Pool
}

func NewOutboxRepository(pool *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{pool: pool}
}

func (r *OutboxRepository) Insert(ctx context.Context, row *OutboxRow) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO outbox (id, event_type, event_payload, idempotency_key, correlation_id, partition_key, status, retry_count, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, 'PENDING', 0, NOW())`,
		row.ID, row.EventType, row.EventPayload, row.IdempotencyKey, row.CorrelationID, row.PartitionKey)
	return err
}

func (r *OutboxRepository) FetchPending(ctx context.Context, limit int) ([]OutboxRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_type, event_payload, idempotency_key, correlation_id, partition_key, status, retry_count, created_at, last_attempt_at
		 FROM outbox WHERE status = 'PENDING' ORDER BY created_at LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []OutboxRow
	for rows.Next() {
		var o OutboxRow
		err := rows.Scan(&o.ID, &o.EventType, &o.EventPayload, &o.IdempotencyKey,
			&o.CorrelationID, &o.PartitionKey, &o.Status, &o.RetryCount, &o.CreatedAt, &o.LastAttemptAt)
		if err != nil {
			return nil, err
		}
		results = append(results, o)
	}
	return results, nil
}

func (r *OutboxRepository) MarkCompleted(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM outbox WHERE id = $1`, id)
	return err
}

func (r *OutboxRepository) MarkFailed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE outbox SET retry_count = retry_count + 1, last_attempt_at = NOW(), status = CASE WHEN retry_count >= 4 THEN 'FAILED' ELSE 'PENDING' END WHERE id = $1`, id)
	return err
}

type EventStoreRepository struct {
	pool *pgxpool.Pool
}

func NewEventStoreRepository(pool *pgxpool.Pool) *EventStoreRepository {
	return &EventStoreRepository{pool: pool}
}

func (r *EventStoreRepository) Append(ctx context.Context, eventID, eventType, aggregateType, aggregateID string, data []byte, version int, correlationID, userID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO event_store (event_id, event_type, aggregate_type, aggregate_id, event_data, version, correlation_id, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`,
		eventID, eventType, aggregateType, aggregateID, data, version, correlationID, userID)
	return err
}

func (r *EventStoreRepository) GetEvents(ctx context.Context, aggregateType, aggregateID string) ([]pgx.Row, error) {
	rows, _ := r.pool.Query(ctx,
		`SELECT event_id, event_type, event_data, version, correlation_id, created_at
		 FROM event_store WHERE aggregate_type = $1 AND aggregate_id = $2 ORDER BY version`,
		aggregateType, aggregateID)
	defer rows.Close()
	return nil, rows.Err()
}
