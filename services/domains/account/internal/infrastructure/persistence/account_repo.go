package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/account/internal/domain"
)

type AccountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

func (r *AccountRepository) Save(ctx context.Context, a *domain.Account) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO accounts (account_id, owner_id, household_id, account_type, classification, account_name, currency, status,
			opened_date, liquidity_profile, account_health, visibility, tags, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,NOW(),NOW())
		ON CONFLICT (account_id) DO UPDATE SET account_name=$6, status=$8, visibility=$12, updated_at=NOW()`,
		a.ID(), a.OwnerID(), a.HouseholdID(), string(a.AccountType()), string(a.Classification()), a.AccountName(),
		a.Currency(), string(a.Status()), a.OpenedDate(), string(a.LiquidityProfile()),
		string(a.AccountHealth()), string(a.Visibility()), a.Tags())
	return err
}

func (r *AccountRepository) UpdateStatus(ctx context.Context, id string, from, to domain.AccountStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE accounts SET status=$1, updated_at=NOW() WHERE account_id=$2 AND status=$3`,
		string(to), id, string(from))
	return err
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	var aid, oid, at, cls, an, cur, st, lp, ah, vis string
	var hhid *string
	var tags []string
	var od, ca, ua time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT account_id, owner_id, household_id, account_type, classification, account_name, currency, status,
			opened_date, liquidity_profile, account_health, visibility, tags, created_at, updated_at
		FROM accounts WHERE account_id = $1`, id).Scan(&aid, &oid, &hhid, &at, &cls, &an, &cur, &st, &od, &lp, &ah, &vis, &tags, &ca, &ua)
	if err != nil { return nil, fmt.Errorf("get account: %w", err) }
	return domain.ReconstructFromDB(aid, oid, hhid, an, cur, domain.AccountType(at), domain.AccountClassification(cls),
		domain.AccountStatus(st), od, domain.LiquidityProfile(lp), domain.AccountHealth(ah),
		domain.Visibility(vis), tags, ca, ua), nil
}

func (r *AccountRepository) ListByUser(ctx context.Context, ownerID string, cursor string, limit int) ([]*domain.Account, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT account_id, owner_id, household_id, account_type, classification, account_name, currency, status, liquidity_profile, account_health, visibility, created_at
		FROM accounts WHERE owner_id = $1 ORDER BY account_name LIMIT $2`, ownerID, limit+1)
	defer rows.Close()
	return scanAccountList(rows, limit)
}

func (r *AccountRepository) ListByType(ctx context.Context, ownerID string, at domain.AccountType, cursor string, limit int) ([]*domain.Account, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT account_id, owner_id, household_id, account_type, classification, account_name, currency, status, liquidity_profile, account_health, visibility, created_at
		FROM accounts WHERE owner_id = $1 AND account_type = $2 ORDER BY account_name LIMIT $3`, ownerID, string(at), limit+1)
	defer rows.Close()
	return scanAccountList(rows, limit)
}

func (r *AccountRepository) ListByStatus(ctx context.Context, ownerID string, status domain.AccountStatus, cursor string, limit int) ([]*domain.Account, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT account_id, owner_id, household_id, account_type, classification, account_name, currency, status, liquidity_profile, account_health, visibility, created_at
		FROM accounts WHERE owner_id = $1 AND status = $2 ORDER BY account_name LIMIT $3`, ownerID, string(status), limit+1)
	defer rows.Close()
	return scanAccountList(rows, limit)
}

func (r *AccountRepository) ListByInstitution(ctx context.Context, institutionID string, cursor string, limit int) ([]*domain.Account, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT account_id, owner_id, household_id, account_type, classification, account_name, currency, status, liquidity_profile, account_health, visibility, created_at
		FROM accounts WHERE institution_id = $1 ORDER BY account_name LIMIT $2`, institutionID, limit+1)
	defer rows.Close()
	return scanAccountList(rows, limit)
}

func (r *AccountRepository) ListByHousehold(ctx context.Context, householdID string, cursor string, limit int) ([]*domain.Account, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT account_id, owner_id, household_id, account_type, classification, account_name, currency, status, liquidity_profile, account_health, visibility, created_at
		FROM accounts WHERE household_id = $1 AND visibility IN ('Household', 'Shared') ORDER BY account_name LIMIT $2`, householdID, limit+1)
	defer rows.Close()
	return scanAccountList(rows, limit)
}

func scanAccountList(rows pgx.Rows, limit int) ([]*domain.Account, string, error) {
	var accts []*domain.Account
	for rows.Next() {
		var aid, oid, at, cls, an, cur, st, lp, ah, vis string
		var hhid *string
		var ca time.Time
		rows.Scan(&aid, &oid, &hhid, &at, &cls, &an, &cur, &st, &lp, &ah, &vis, &ca)
		accts = append(accts, domain.ReconstructFromDB(aid, oid, hhid, an, cur,
			domain.AccountType(at), domain.AccountClassification(cls),
			domain.AccountStatus(st), time.Time{}, domain.LiquidityProfile(lp), domain.AccountHealth(ah),
			domain.Visibility(vis), nil, ca, ca))
	}
	hasMore := len(accts) > limit
	if hasMore { accts = accts[:limit] }
	return accts, "", nil
}

