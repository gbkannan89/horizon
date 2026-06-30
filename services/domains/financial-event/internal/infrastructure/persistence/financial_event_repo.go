package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	fe "github.com/horizon/core/services/domains/financial-event/internal/domain"
)

type FinancialEventRepository struct {
	pool *pgxpool.Pool
}

func NewFinancialEventRepository(pool *pgxpool.Pool) *FinancialEventRepository {
	return &FinancialEventRepository{pool: pool}
}

func (r *FinancialEventRepository) Save(ctx context.Context, event *fe.FinancialEvent) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO financial_events (event_id, user_id, event_type, event_date, effective_date,
			currency, amount, source, destination, description, origin, confidence, created_by,
			state, sequence_number, order_index, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,NOW(),NOW())
		ON CONFLICT (event_id) DO UPDATE SET state=$14, updated_at=NOW()`,
		event.EventID(), event.UserID(), string(event.EventType()),
		event.EventDate(), event.EffectiveDate(), event.Currency(), event.Amount(),
		event.Source(), event.Destination(), event.Description(),
		string(event.Origin()), string(event.Confidence()), string(event.CreatedBy()),
		string(event.State()), event.SequenceNumber(), event.OrderIndex())
	return err
}

func (r *FinancialEventRepository) UpdateState(ctx context.Context, eventID string, fromState, toState fe.EventState) error {
	_, err := r.pool.Exec(ctx, `UPDATE financial_events SET state=$1, updated_at=NOW() WHERE event_id=$2 AND state=$3`,
		string(toState), eventID, string(fromState))
	return err
}

func (r *FinancialEventRepository) GetByID(ctx context.Context, eventID string) (*fe.FinancialEvent, error) {
	var id, uid, et, cur, desc, src, dst, orig, conf, cb, st, oi string
	var amt int64
	var seq int64
	var ed, efd, ca, ua time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description,
			source, destination, origin, confidence, created_by, state, sequence_number, order_index, created_at, updated_at
		FROM financial_events WHERE event_id = $1`, eventID).Scan(
		&id, &uid, &et, &amt, &cur, &ed, &efd, &desc, &src, &dst, &orig, &conf, &cb, &st, &seq, &oi, &ca, &ua)
	if err != nil { return nil, fmt.Errorf("get event: %w", err) }

	return fe.ReconstructFromDB(id, uid, "", fe.EventType(et), "", ed, efd, cur, amt,
		src, dst, desc, "", "", fe.EventOrigin(orig), fe.EventConfidence(conf), "", fe.CreatedBy(cb),
		"", "", "", nil, fe.EventState(st), seq, oi, ca, ua), nil
}

func (r *FinancialEventRepository) ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*fe.FinancialEvent, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description, state, created_at
		FROM financial_events WHERE user_id = $1 ORDER BY effective_date DESC LIMIT $2`, userID, limit+1)
	defer rows.Close()
	return scanFeList(rows, limit)
}

func (r *FinancialEventRepository) ListByAccount(ctx context.Context, accountID string, cursor string, limit int) ([]*fe.FinancialEvent, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description, state, created_at
		FROM financial_events WHERE source = $1 OR destination = $1 ORDER BY effective_date DESC LIMIT $2`, accountID, limit+1)
	defer rows.Close()
	return scanFeList(rows, limit)
}

func (r *FinancialEventRepository) ListByDateRange(ctx context.Context, userID string, start, end time.Time, cursor string, limit int) ([]*fe.FinancialEvent, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description, state, created_at
		FROM financial_events WHERE user_id = $1 AND effective_date >= $2 AND effective_date <= $3 ORDER BY effective_date LIMIT $4`,
		userID, start, end, limit+1)
	defer rows.Close()
	return scanFeList(rows, limit)
}

func (r *FinancialEventRepository) ListByType(ctx context.Context, userID string, eventType fe.EventType, cursor string, limit int) ([]*fe.FinancialEvent, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description, state, created_at
		FROM financial_events WHERE user_id = $1 AND event_type = $2 ORDER BY effective_date DESC LIMIT $3`,
		userID, string(eventType), limit+1)
	defer rows.Close()
	return scanFeList(rows, limit)
}

func (r *FinancialEventRepository) GetTimeline(ctx context.Context, userID string, cursor string, limit int) ([]*fe.FinancialEvent, string, error) {
	return r.ListByUser(ctx, userID, cursor, limit)
}

func (r *FinancialEventRepository) FindPotentialDuplicates(ctx context.Context, amount int64, currency string, source, destination string, effectiveDate time.Time, window time.Duration) ([]*fe.FinancialEvent, error) {
	start := effectiveDate.Add(-window)
	end := effectiveDate.Add(window)
	rows, _ := r.pool.Query(ctx,
		`SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description, state, created_at
		FROM financial_events WHERE amount = $1 AND currency = $2 AND effective_date >= $3 AND effective_date <= $4
		AND state IN ('CONFIRMED','POSTED') LIMIT 10`,
		amount, currency, start, end)
	defer rows.Close()
	events, _, _ := scanFeList(rows, 10)
	return events, nil
}

func (r *FinancialEventRepository) GetReversalChain(ctx context.Context, eventID string) ([]*fe.FinancialEvent, error) {
	rows, _ := r.pool.Query(ctx,
		`WITH RECURSIVE chain AS (
			SELECT * FROM financial_events WHERE reversal_of_event_id = $1
			UNION ALL
			SELECT fe.* FROM financial_events fe INNER JOIN chain c ON fe.reversal_of_event_id = c.event_id
		) SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description, state, created_at
		FROM chain LIMIT 100`, eventID)
	defer rows.Close()
	events, _, _ := scanFeList(rows, 100)
	return events, nil
}

func (r *FinancialEventRepository) GetOriginalEvent(ctx context.Context, reversalEventID string) (*fe.FinancialEvent, error) {
	var id, uid, et, cur, desc, st string
	var amt int64
	var ed, efd, ca time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT fe.event_id, fe.user_id, fe.event_type, fe.amount, fe.currency, fe.event_date, fe.effective_date, fe.description, fe.state, fe.created_at
		FROM financial_events fe INNER JOIN financial_events rev ON fe.event_id = rev.reversal_of_event_id
		WHERE rev.event_id = $1`, reversalEventID).Scan(&id, &uid, &et, &amt, &cur, &ed, &efd, &desc, &st, &ca)
	if err != nil { return nil, fmt.Errorf("get original: %w", err) }

	return fe.ReconstructFromDB(id, uid, "", fe.EventType(et), "", ed, efd, cur, amt,
		"", "", desc, "", "", "", "", "", "", "", "", "", nil, fe.EventState(st), 0, "", ca, ca), nil
}

func (r *FinancialEventRepository) FindByExternalReference(ctx context.Context, sourceOfTruth string, reference string) (*fe.FinancialEvent, error) {
	var id, uid, et, cur, desc, st string
	var amt int64
	var ed, efd, ca time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT event_id, user_id, event_type, amount, currency, event_date, effective_date, description, state, created_at
		FROM financial_events WHERE source_of_truth = $1 AND reference = $2 LIMIT 1`,
		sourceOfTruth, reference).Scan(&id, &uid, &et, &amt, &cur, &ed, &efd, &desc, &st, &ca)
	if err != nil { return nil, fmt.Errorf("find ref: %w", err) }

	return fe.ReconstructFromDB(id, uid, "", fe.EventType(et), "", ed, efd, cur, amt,
		"", "", desc, "", "", "", "", "", "", "", "", "", nil, fe.EventState(st), 0, "", ca, ca), nil
}

func (r *FinancialEventRepository) GetNextSequenceNumber(ctx context.Context, userID string) (int64, error) {
	var seq int64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(sequence_number), 0) + 1 FROM financial_events WHERE user_id = $1`, userID).Scan(&seq)
	return seq, err
}

func scanFeList(rows pgx.Rows, limit int) ([]*fe.FinancialEvent, string, error) {
	var events []*fe.FinancialEvent
	for rows.Next() {
		var id, uid, et, cur, desc, st string
		var amt int64
		var ed, efd, ca time.Time
		rows.Scan(&id, &uid, &et, &amt, &cur, &ed, &efd, &desc, &st, &ca)
		events = append(events, fe.ReconstructFromDB(id, uid, "", fe.EventType(et), "", ed, efd, cur, amt,
			"", "", desc, "", "", "", "", "", "", "", "", "", nil, fe.EventState(st), 0, "", ca, ca))
	}
	hasMore := len(events) > limit
	if hasMore { events = events[:limit] }
	return events, "", nil
}
