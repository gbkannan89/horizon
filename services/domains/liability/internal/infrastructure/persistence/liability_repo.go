package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/liability/internal/domain"
)

type LiabilityRepository struct {
	pool *pgxpool.Pool
}

func NewLiabilityRepository(pool *pgxpool.Pool) *LiabilityRepository {
	return &LiabilityRepository{pool: pool}
}

func (r *LiabilityRepository) Save(ctx context.Context, l *domain.Liability) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO liabilities (liability_id, name, classification, owner_id, currency, original_principal,
			interest_rate, interest_model, repayment_model, remaining_installments, maturity_date, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
		ON CONFLICT (liability_id) DO UPDATE SET status=$12`,
		l.ID(), l.Name(), string(l.Classification()), l.OwnerID(), l.Currency(),
		l.OriginalPrincipal(), l.InterestRate(), string(l.InterestModel()),
		string(l.RepaymentModel()), l.RemainingInstallments(), l.MaturityDate(),
		string(l.Status()))
	return err
}

func (r *LiabilityRepository) UpdateStatus(ctx context.Context, id string, from, to domain.LiabilityStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE liabilities SET status=$1 WHERE liability_id=$2 AND status=$3`,
		string(to), id, string(from))
	return err
}

func (r *LiabilityRepository) GetByID(ctx context.Context, id string) (*domain.Liability, error) {
	var lid, name, cls, oid, cur, im, rm, st string
	var op, ir float64
	var rem int
	var md, ca time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT liability_id, name, classification, owner_id, currency, original_principal, interest_rate,
			interest_model, repayment_model, remaining_installments, maturity_date, status, created_at
		FROM liabilities WHERE liability_id = $1`, id).Scan(&lid, &name, &cls, &oid, &cur, &op, &ir, &im, &rm, &rem, &md, &st, &ca)
	if err != nil { return nil, fmt.Errorf("get liability: %w", err) }
	return domain.ReconstructFromDB(lid, name, domain.LiabilityClassification(cls), oid, cur,
		op, ir, domain.InterestModel(im), domain.RepaymentModel(rm),
		rem, md, "", "", domain.LiabilityStatus(st), "", nil, ca), nil
}

func (r *LiabilityRepository) ListByUser(ctx context.Context, ownerID string, cursor string, limit int) ([]*domain.Liability, string, error) {
	return scanLiabList(r.pool, ctx, `owner_id = $1 ORDER BY name LIMIT $2`, ownerID, limit)
}

func (r *LiabilityRepository) ListByClassification(ctx context.Context, ownerID string, cls string, cursor string, limit int) ([]*domain.Liability, string, error) {
	return scanLiabList(r.pool, ctx, `owner_id = $1 AND classification = $2 ORDER BY name LIMIT $3`, ownerID, cls, limit)
}

func (r *LiabilityRepository) ListByStatus(ctx context.Context, ownerID string, status domain.LiabilityStatus, cursor string, limit int) ([]*domain.Liability, string, error) {
	return scanLiabList(r.pool, ctx, `owner_id = $1 AND status = $2 ORDER BY name LIMIT $3`, ownerID, string(status), limit)
}

func (r *LiabilityRepository) ListByInstitution(ctx context.Context, institutionID string, cursor string, limit int) ([]*domain.Liability, string, error) {
	return scanLiabList(r.pool, ctx, `institution_id = $1 ORDER BY name LIMIT $2`, institutionID, limit)
}

func scanLiabList(pool *pgxpool.Pool, ctx context.Context, where string, args ...interface{}) ([]*domain.Liability, string, error) {
	limit := args[len(args)-1].(int)
	if limit <= 0 { limit = 25 }
	args[len(args)-1] = limit + 1
	rows, _ := pool.Query(ctx,
		`SELECT liability_id, name, classification, owner_id, currency, original_principal, interest_rate, status, created_at
		FROM liabilities WHERE `+where, args...)
	defer rows.Close()
	var ls []*domain.Liability
	for rows.Next() {
		var lid, name, cls, oid, cur, st string
		var op, ir float64
		var ca time.Time
		rows.Scan(&lid, &name, &cls, &oid, &cur, &op, &ir, &st, &ca)
		ls = append(ls, domain.ReconstructFromDB(lid, name, domain.LiabilityClassification(cls), oid, cur,
			op, ir, "", "", 0, time.Time{}, "", "", domain.LiabilityStatus(st), "", nil, ca))
	}
	hasMore := len(ls) > limit
	if hasMore { ls = ls[:limit] }
	return ls, "", nil
}
