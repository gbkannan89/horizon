package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/allocation/internal/domain"
)

type AllocationRepository struct {
	pool *pgxpool.Pool
}

func NewAllocationRepository(pool *pgxpool.Pool) *AllocationRepository {
	return &AllocationRepository{pool: pool}
}

func (r *AllocationRepository) Save(ctx context.Context, a *domain.Allocation) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO allocations (allocation_id, goal_id, funding_source_id, funding_source_type, allocation_type,
			allocation_strategy, priority, currency, status, effective_date, reserved_amount, allocated_amount, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
		ON CONFLICT (allocation_id) DO UPDATE SET status=$9, reserved_amount=$11, allocated_amount=$12, updated_at=NOW()`,
		a.ID(), a.GoalID(), a.FundingSourceID(), string(a.FundingSourceType()), string(a.AllocationType()),
		string(a.Strategy()), a.Priority(), a.Currency(), string(a.Status()),
		a.EffectiveDate(), a.ReservedAmount(), a.AllocatedAmount())
	return err
}

func (r *AllocationRepository) UpdateStatus(ctx context.Context, id string, from, to domain.AllocationStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE allocations SET status=$1, updated_at=NOW() WHERE allocation_id=$2 AND status=$3`,
		string(to), id, string(from))
	return err
}

func (r *AllocationRepository) GetByID(ctx context.Context, id string) (*domain.Allocation, error) {
	var aid, gid, fsid, at, st, cur string
	var pri int
	var ra, aa int64
	var ca time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT allocation_id, goal_id, funding_source_id, allocation_type, priority, currency, status, reserved_amount, allocated_amount, created_at
		FROM allocations WHERE allocation_id = $1`, id).Scan(&aid, &gid, &fsid, &at, &pri, &cur, &st, &ra, &aa, &ca)
	if err != nil { return nil, fmt.Errorf("get allocation: %w", err) }
	return domain.ReconstructFromDB(aid, gid, fsid, domain.AllocationType(at), pri, cur, domain.AllocationStatus(st), ra, aa, nil, ca, ca), nil
}

func (r *AllocationRepository) ListByGoal(ctx context.Context, goalID string, cursor string, limit int) ([]*domain.Allocation, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT allocation_id, goal_id, funding_source_id, allocation_type, priority, currency, status, reserved_amount, allocated_amount, created_at
		FROM allocations WHERE goal_id = $1 ORDER BY priority LIMIT $2`, goalID, limit+1)
	defer rows.Close()
	return scanAllocs(rows, limit)
}

func (r *AllocationRepository) ListByFundingSource(ctx context.Context, sourceID string, cursor string, limit int) ([]*domain.Allocation, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT allocation_id, goal_id, funding_source_id, allocation_type, priority, currency, status, reserved_amount, allocated_amount, created_at
		FROM allocations WHERE funding_source_id = $1 ORDER BY priority LIMIT $2`, sourceID, limit+1)
	defer rows.Close()
	return scanAllocs(rows, limit)
}

func (r *AllocationRepository) ListByStatus(ctx context.Context, status domain.AllocationStatus, cursor string, limit int) ([]*domain.Allocation, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT allocation_id, goal_id, funding_source_id, allocation_type, priority, currency, status, reserved_amount, allocated_amount, created_at
		FROM allocations WHERE status = $1 ORDER BY priority LIMIT $2`, string(status), limit+1)
	defer rows.Close()
	return scanAllocs(rows, limit)
}

func (r *AllocationRepository) GetActiveAllocations(ctx context.Context, goalID string) ([]*domain.Allocation, error) {
	rows, _ := r.pool.Query(ctx,
		`SELECT allocation_id, goal_id, funding_source_id, allocation_type, priority, currency, status, reserved_amount, allocated_amount, created_at
		FROM allocations WHERE goal_id = $1 AND status = 'Active' ORDER BY priority`, goalID)
	defer rows.Close()
	allocs, _, _ := scanAllocs(rows, 100)
	return allocs, nil
}

func (r *AllocationRepository) GetTotalReserved(ctx context.Context, fundingSourceID string) (int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(reserved_amount), 0) FROM allocations WHERE funding_source_id = $1 AND status = 'Active'`,
		fundingSourceID).Scan(&total)
	return total, err
}

func scanAllocs(rows pgx.Rows, limit int) ([]*domain.Allocation, string, error) {
	var as []*domain.Allocation
	for rows.Next() {
		var aid, gid, fsid, at, st, cur string
		var pri int
		var ra, aa int64
		var ca time.Time
		rows.Scan(&aid, &gid, &fsid, &at, &pri, &cur, &st, &ra, &aa, &ca)
		as = append(as, domain.ReconstructFromDB(aid, gid, fsid, domain.AllocationType(at), pri, cur, domain.AllocationStatus(st), ra, aa, nil, ca, ca))
	}
	hasMore := len(as) > limit
	if hasMore { as = as[:limit] }
	return as, "", nil
}
