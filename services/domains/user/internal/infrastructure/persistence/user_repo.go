package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/user/internal/domain"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Save(ctx context.Context, u *domain.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (user_id, display_name, legal_name, country, base_currency, locale, timezone,
			user_health, user_confidence, financial_identity_profile, status, user_type,
			financial_behaviour_profile, profile_completeness, source_of_truth, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW(),NOW())
		ON CONFLICT (user_id) DO UPDATE SET display_name=$2, status=$11, updated_at=NOW()`,
		u.UserID(), u.DisplayName(), u.LegalName(), u.Country(), u.BaseCurrency(),
		u.Locale(), u.Timezone(), string(u.UserHealth()), string(u.UserConfidence()),
		string(u.FinancialIdentityProfile()), string(u.Status()), string(u.UserType()),
		string(u.FinancialBehaviourProfile()), u.ProfileCompleteness(), string(u.SourceOfTruth()))
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userID string, fromStatus, toStatus domain.UserStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET status=$1, updated_at=NOW() WHERE user_id=$2 AND status=$3`,
		string(toStatus), userID, string(fromStatus))
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	var id, dn, ln, cn, bc, loc, tz, uh, uc, fip, st, ut, bp, sot string
	var pc float64
	var ca, ua time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT user_id, display_name, legal_name, country, base_currency, locale, timezone,
			user_health, user_confidence, financial_identity_profile, status, user_type,
			financial_behaviour_profile, profile_completeness, source_of_truth, created_at, updated_at
		FROM users WHERE user_id = $1`, userID).Scan(&id, &dn, &ln, &cn, &bc, &loc, &tz,
		&uh, &uc, &fip, &st, &ut, &bp, &pc, &sot, &ca, &ua)
	if err != nil { return nil, fmt.Errorf("get user: %w", err) }
	return domain.ReconstructFromDB(id, dn, ln, "", cn, bc, loc, tz,
		domain.UserHealth(uh), domain.UserConfidence(uc), domain.FinancialIdentityProfile(fip),
		domain.UserStatus(st), domain.UserType(ut), domain.FinBehaviourProfile(bp),
		pc, domain.SOT(sot), nil, nil, "", nil, ca, ua), nil
}

func (r *UserRepository) ListByStatus(ctx context.Context, status domain.UserStatus, cursor string, limit int) ([]*domain.User, string, error) {
	if limit <= 0 { limit = 25 }
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, display_name, legal_name, country, base_currency, status, created_at
		FROM users WHERE status = $1 ORDER BY created_at DESC LIMIT $2`, string(status), limit+1)
	if err != nil { return nil, "", err }
	defer rows.Close()
	return scanUsers(rows, limit)
}

func (r *UserRepository) Search(ctx context.Context, query string, cursor string, limit int) ([]*domain.User, string, error) {
	if limit <= 0 { limit = 25 }
	rows, err := r.pool.Query(ctx,
		`SELECT user_id, display_name, legal_name, country, base_currency, status, created_at
		FROM users WHERE display_name ILIKE $1 OR legal_name ILIKE $1 ORDER BY created_at DESC LIMIT $2`,
		"%"+query+"%", limit+1)
	if err != nil { return nil, "", err }
	defer rows.Close()
	return scanUsers(rows, limit)
}

func (r *UserRepository) GetConsents(ctx context.Context, userID string) ([]*domain.ConsentRecord, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT consent_type, granted, granted_at, revoked_at, scope FROM user_consents WHERE user_id = $1`, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var cs []*domain.ConsentRecord
	for rows.Next() {
		var c domain.ConsentRecord
		rows.Scan(&c.ConsentType, &c.Granted, &c.GrantedAt, &c.RevokedAt, &c.Scope)
		cs = append(cs, &c)
	}
	return cs, nil
}

func (r *UserRepository) GetHouseholdMemberships(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT household_id FROM user_households WHERE user_id = $1`, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func scanUsers(rows pgx.Rows, limit int) ([]*domain.User, string, error) {
	var us []*domain.User
	for rows.Next() {
		var id, dn, ln, cn, bc, st string
		var ca time.Time
		rows.Scan(&id, &dn, &ln, &cn, &bc, &st, &ca)
		us = append(us, domain.ReconstructFromDB(id, dn, ln, "", cn, bc, "", "",
			"", "", "", domain.UserStatus(st), "", "", 0, "", nil, nil, "", nil, ca, ca))
	}
	hasMore := len(us) > limit
	if hasMore { us = us[:limit] }
	return us, "", nil
}
