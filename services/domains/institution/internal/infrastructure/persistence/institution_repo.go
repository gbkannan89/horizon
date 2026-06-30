package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/institution/internal/domain"
)

type InstitutionRepository struct {
	pool *pgxpool.Pool
}

func NewInstitutionRepository(pool *pgxpool.Pool) *InstitutionRepository {
	return &InstitutionRepository{pool: pool}
}

func (r *InstitutionRepository) Save(ctx context.Context, inst *domain.Institution) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO institutions (institution_id, name, institution_type, category, status, country,
			trust_level, health, confidence, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW(),NOW())
		ON CONFLICT (institution_id) DO UPDATE SET name=$2, status=$5, trust_level=$7, health=$8, confidence=$9`,
		inst.ID(), inst.Name(), string(inst.Type()), inst.Category(), string(inst.Status()),
		inst.Country(), string(inst.TrustLevel()), string(inst.Health()), string(inst.Confidence()))
	return err
}

func (r *InstitutionRepository) UpdateStatus(ctx context.Context, id string, from, to domain.InstitutionStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE institutions SET status=$1 WHERE institution_id=$2 AND status=$3`,
		string(to), id, string(from))
	return err
}

func (r *InstitutionRepository) GetByID(ctx context.Context, id string) (*domain.Institution, error) {
	var iid, name, it, cat, st, cntry, tl, h, conf string
	var ca time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT institution_id, name, institution_type, category, status, country, trust_level, health, confidence, created_at
		FROM institutions WHERE institution_id = $1`, id).Scan(&iid, &name, &it, &cat, &st, &cntry, &tl, &h, &conf, &ca)
	if err != nil { return nil, fmt.Errorf("get institution: %w", err) }
	return domain.ReconstructFromDB(iid, name, domain.InstitutionType(it), cat, st, cntry,
		domain.TrustLevel(tl), domain.InstitutionHealth(h), domain.InstitutionConfidence(conf), nil, ca), nil
}

func (r *InstitutionRepository) ListByType(ctx context.Context, instType domain.InstitutionType, cursor string, limit int) ([]*domain.Institution, string, error) {
	return scanInstList(r.pool, ctx, `institution_type = $1 ORDER BY name LIMIT $2`, string(instType), limit)
}

func (r *InstitutionRepository) ListByCountry(ctx context.Context, country string, cursor string, limit int) ([]*domain.Institution, string, error) {
	return scanInstList(r.pool, ctx, `country = $1 ORDER BY name LIMIT $2`, country, limit)
}

func (r *InstitutionRepository) ListByStatus(ctx context.Context, status domain.InstitutionStatus, cursor string, limit int) ([]*domain.Institution, string, error) {
	return scanInstList(r.pool, ctx, `status = $1 ORDER BY name LIMIT $2`, string(status), limit)
}

func (r *InstitutionRepository) Search(ctx context.Context, query string, cursor string, limit int) ([]*domain.Institution, string, error) {
	return scanInstList(r.pool, ctx, `name ILIKE $1 ORDER BY name LIMIT $2`, "%"+query+"%", limit)
}

func scanInstList(pool *pgxpool.Pool, ctx context.Context, where string, arg interface{}, limit int) ([]*domain.Institution, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := pool.Query(ctx,
		`SELECT institution_id, name, institution_type, category, status, country, trust_level, health, confidence, created_at
		FROM institutions WHERE `+where, arg, limit+1)
	defer rows.Close()
	var insts []*domain.Institution
	for rows.Next() {
		var iid, name, it, cat, st, cntry, tl, h, conf string
		var ca time.Time
		rows.Scan(&iid, &name, &it, &cat, &st, &cntry, &tl, &h, &conf, &ca)
		insts = append(insts, domain.ReconstructFromDB(iid, name, domain.InstitutionType(it), cat, st, cntry,
			domain.TrustLevel(tl), domain.InstitutionHealth(h), domain.InstitutionConfidence(conf), nil, ca))
	}
	hasMore := len(insts) > limit
	if hasMore { insts = insts[:limit] }
	return insts, "", nil
}
