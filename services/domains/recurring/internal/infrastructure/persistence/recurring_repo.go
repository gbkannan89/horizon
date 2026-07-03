package persistence

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/recurring/internal/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Save(ctx context.Context, rt *domain.RecurringTransaction) error {
	tags, _ := json.Marshal(rt.Tags())
	meta, _ := json.Marshal(rt.Metadata())
	tmpl, _ := json.Marshal(rt.EventTemplate())
	_, err := r.pool.Exec(ctx, `
		INSERT INTO recurring_transactions (id, user_id, household_id, name, description, amount, currency,
			frequency, interval, start_date, end_date, next_occurrence, last_occurrence, status,
			event_template, skip_holidays, skip_weekends, tags, metadata, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)
		ON CONFLICT (id) DO UPDATE SET
			name=$4, description=$5, amount=$6, frequency=$8, interval=$9,
			end_date=$11, next_occurrence=$12, last_occurrence=$13, status=$14,
			event_template=$15, skip_holidays=$16, skip_weekends=$17, tags=$18, metadata=$19, updated_at=$21`,
		rt.ID(), rt.UserID(), nullOrEmpty(rt.HouseholdID()), rt.Name(), rt.Description(),
		rt.Amount(), rt.Currency(), string(rt.Frequency()), rt.Interval(),
		rt.StartDate(), rt.EndDate(), rt.NextOccurrence(), rt.LastOccurrence(),
		string(rt.Status()), tmpl, rt.SkipHolidays(), rt.SkipWeekends(), tags, meta,
		rt.CreatedAt(), rt.UpdatedAt(),
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.RecurringTransaction, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, COALESCE(household_id,''), name, COALESCE(description,''),
			amount, currency, frequency, interval, start_date, end_date,
			next_occurrence, last_occurrence, status,
			event_template->>'type', COALESCE(event_template->>'category',''),
			COALESCE(event_template->>'source',''), COALESCE(event_template->>'destination',''),
			skip_holidays, skip_weekends, COALESCE(tags::text,'[]'), COALESCE(metadata::text,'{}'),
			created_at, updated_at
		FROM recurring_transactions WHERE id=$1`, id)
	return scanRow(row)
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]*domain.RecurringTransaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, COALESCE(household_id,''), name, COALESCE(description,''),
			amount, currency, frequency, interval, start_date, end_date,
			next_occurrence, last_occurrence, status,
			event_template->>'type', COALESCE(event_template->>'category',''),
			COALESCE(event_template->>'source',''), COALESCE(event_template->>'destination',''),
			skip_holidays, skip_weekends, COALESCE(tags::text,'[]'), COALESCE(metadata::text,'{}'),
			created_at, updated_at
		FROM recurring_transactions WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *PostgresRepository) ListByStatus(ctx context.Context, status domain.RecurringStatus) ([]*domain.RecurringTransaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, COALESCE(household_id,''), name, COALESCE(description,''),
			amount, currency, frequency, interval, start_date, end_date,
			next_occurrence, last_occurrence, status,
			event_template->>'type', COALESCE(event_template->>'category',''),
			COALESCE(event_template->>'source',''), COALESCE(event_template->>'destination',''),
			skip_holidays, skip_weekends, COALESCE(tags::text,'[]'), COALESCE(metadata::text,'{}'),
			created_at, updated_at
		FROM recurring_transactions WHERE status=$1 ORDER BY created_at DESC`, string(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *PostgresRepository) GetDueByDate(ctx context.Context, date time.Time) ([]*domain.RecurringTransaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, COALESCE(household_id,''), name, COALESCE(description,''),
			amount, currency, frequency, interval, start_date, end_date,
			next_occurrence, last_occurrence, status,
			event_template->>'type', COALESCE(event_template->>'category',''),
			COALESCE(event_template->>'source',''), COALESCE(event_template->>'destination',''),
			skip_holidays, skip_weekends, COALESCE(tags::text,'[]'), COALESCE(metadata::text,'{}'),
			created_at, updated_at
		FROM recurring_transactions WHERE status='Active' AND next_occurrence <= $1 ORDER BY next_occurrence`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *PostgresRepository) GetActiveRecurring(ctx context.Context, userID string) ([]*domain.RecurringTransaction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, COALESCE(household_id,''), name, COALESCE(description,''),
			amount, currency, frequency, interval, start_date, end_date,
			next_occurrence, last_occurrence, status,
			event_template->>'type', COALESCE(event_template->>'category',''),
			COALESCE(event_template->>'source',''), COALESCE(event_template->>'destination',''),
			skip_holidays, skip_weekends, COALESCE(tags::text,'[]'), COALESCE(metadata::text,'{}'),
			created_at, updated_at
		FROM recurring_transactions WHERE user_id=$1 AND status='Active' ORDER BY next_occurrence`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRow(row scanner) (*domain.RecurringTransaction, error) {
	var (
		id, userID, householdID, name, desc string
		amount                              int64
		currency, freq                      string
		interval                            int
		startDate                           time.Time
		endDate, nextOcc, lastOcc           *time.Time
		status, evType, evCat, evSrc, evDst string
		skipHols, skipWknds                bool
		tagsJSON, metaJSON                  string
		createdAt, updatedAt                time.Time
	)
	err := row.Scan(&id, &userID, &householdID, &name, &desc,
		&amount, &currency, &freq, &interval, &startDate,
		&endDate, &nextOcc, &lastOcc, &status,
		&evType, &evCat, &evSrc, &evDst,
		&skipHols, &skipWknds, &tagsJSON, &metaJSON,
		&createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	var tags []string
	json.Unmarshal([]byte(tagsJSON), &tags)
	var meta map[string]string
	json.Unmarshal([]byte(metaJSON), &meta)
	return domain.ReconstructFromDB(id, userID, householdID, name, desc,
		amount, currency, freq, interval, startDate, endDate, nextOcc, lastOcc,
		status, evType, evCat, evSrc, evDst,
		skipHols, skipWknds, tags, meta, createdAt, updatedAt), nil
}

func scanRows(rows pgxRows) ([]*domain.RecurringTransaction, error) {
	var result []*domain.RecurringTransaction
	for rows.Next() {
		rt, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, rt)
	}
	return result, nil
}

type pgxRows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
}

func nullOrEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func init() { _ = time.UTC }
