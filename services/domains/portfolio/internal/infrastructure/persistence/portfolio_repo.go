package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/portfolio/internal/domain"
)

type PortfolioRepository struct {
	pool *pgxpool.Pool
}

func NewPortfolioRepository(pool *pgxpool.Pool) *PortfolioRepository {
	return &PortfolioRepository{pool: pool}
}

func (r *PortfolioRepository) Save(ctx context.Context, p *domain.Portfolio) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO portfolios (portfolio_id, name, portfolio_type, owner_id, base_currency, membership_model,
			status, portfolio_health, portfolio_confidence, risk_profile, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		ON CONFLICT (portfolio_id) DO UPDATE SET name=$2, status=$7`,
		p.ID(), p.Name(), string(p.PortfolioType()), p.OwnerID(), p.BaseCurrency(),
		string(p.MembershipModel()), string(p.Status()), string(p.Health()),
		string(p.Confidence()), string(p.RiskProfile()))
	return err
}

func (r *PortfolioRepository) UpdateStatus(ctx context.Context, id string, from, to domain.PortfolioStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE portfolios SET status=$1 WHERE portfolio_id=$2 AND status=$3`,
		string(to), id, string(from))
	return err
}

func (r *PortfolioRepository) GetByID(ctx context.Context, id string) (*domain.Portfolio, error) {
	var pid, name, pt, oid, bc, mm, st, ph, pc, rp string
	var ca time.Time
	err := r.pool.QueryRow(ctx,
		`SELECT portfolio_id, name, portfolio_type, owner_id, base_currency, membership_model, status,
			portfolio_health, portfolio_confidence, risk_profile, created_at
		FROM portfolios WHERE portfolio_id = $1`, id).Scan(&pid, &name, &pt, &oid, &bc, &mm, &st, &ph, &pc, &rp, &ca)
	if err != nil { return nil, fmt.Errorf("get portfolio: %w", err) }
	return domain.ReconstructFromDB(pid, name, domain.PortfolioType(pt), oid, bc,
		domain.MembershipModel(mm), domain.PortfolioStatus(st),
		domain.PortfolioHealth(ph), domain.PortfolioConfidence(pc),
		domain.RiskProfile(rp), nil, ca), nil
}

func (r *PortfolioRepository) ListByUser(ctx context.Context, ownerID string, cursor string, limit int) ([]*domain.Portfolio, string, error) {
	return scanPfList(r.pool, ctx, `owner_id = $1 ORDER BY name LIMIT $2`, ownerID, limit)
}

func (r *PortfolioRepository) ListByType(ctx context.Context, ownerID string, pt domain.PortfolioType, cursor string, limit int) ([]*domain.Portfolio, string, error) {
	return scanPfList(r.pool, ctx, `owner_id = $1 AND portfolio_type = $2 ORDER BY name LIMIT $3`, ownerID, string(pt), limit)
}

func (r *PortfolioRepository) GetPortfolioHistory(ctx context.Context, id string, cursor string, limit int) ([]*domain.Portfolio, string, error) {
	return scanPfList(r.pool, ctx, `portfolio_id = $1 ORDER BY created_at DESC LIMIT $2`, id, limit)
}

func scanPfList(pool *pgxpool.Pool, ctx context.Context, where string, args ...interface{}) ([]*domain.Portfolio, string, error) {
	limit := args[len(args)-1].(int)
	if limit <= 0 { limit = 25 }
	args[len(args)-1] = limit + 1
	rows, _ := pool.Query(ctx,
		`SELECT portfolio_id, name, portfolio_type, owner_id, base_currency, status, portfolio_health, portfolio_confidence, risk_profile, created_at
		FROM portfolios WHERE `+where, args...)
	defer rows.Close()
	var ps []*domain.Portfolio
	for rows.Next() {
		var pid, name, pt, oid, bc, st, ph, pc, rp string
		var ca time.Time
		rows.Scan(&pid, &name, &pt, &oid, &bc, &st, &ph, &pc, &rp, &ca)
		ps = append(ps, domain.ReconstructFromDB(pid, name, domain.PortfolioType(pt), oid, bc,
			"", domain.PortfolioStatus(st), domain.PortfolioHealth(ph),
			domain.PortfolioConfidence(pc), domain.RiskProfile(rp), nil, ca))
	}
	hasMore := len(ps) > limit
	if hasMore { ps = ps[:limit] }
	return ps, "", nil
}
