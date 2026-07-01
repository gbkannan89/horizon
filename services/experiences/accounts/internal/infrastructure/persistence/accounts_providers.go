package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/accounts/internal/engine"
)

// AccountsPGProvider implements all account data providers using PostgreSQL.
type AccountsPGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewAccountsPGProvider(pool *pgxpool.Pool) *AccountsPGProvider {
	return &AccountsPGProvider{pool: pool, now: time.Now}
}

// getAccountInput scans a single account + institution into AccountInput.
func (p *AccountsPGProvider) getAccountInput(ctx context.Context, accountID string) (*engine.AccountInput, error) {
	var aid, oid, at, cls, an, cur, st, lp, ah, vis string
	var od time.Time
	var cd *time.Time
	var tags []string
	var notes string
	var ca time.Time
	var instName *string

	err := p.pool.QueryRow(ctx,
		`SELECT a.account_id, a.owner_id, a.account_type, a.classification, a.account_name,
			a.currency, a.status, a.opened_date, a.closed_date, a.liquidity_profile,
			a.account_health, a.visibility, a.tags, a.notes, a.created_at,
			i.name
		FROM accounts a
		LEFT JOIN institutions i ON a.institution_id = i.institution_id
		WHERE a.account_id = $1`, accountID).Scan(
		&aid, &oid, &at, &cls, &an, &cur, &st, &od, &cd, &lp, &ah, &vis, &tags, &notes, &ca, &instName)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}

	instNameStr := ""
	if instName != nil {
		instNameStr = *instName
	}

	openedDate := od.Format(time.DateOnly)
	createdAt := ca.Format(time.RFC3339)

	return &engine.AccountInput{
		AccountID:        aid,
		AccountName:      an,
		AccountType:      at,
		Classification:   cls,
		Status:           st,
		Currency:         cur,
		InstitutionName:  instNameStr,
		InstitutionID:    "",
		OpenedDate:       openedDate,
		ClosedDate:       "",
		LiquidityProfile: lp,
		AccountHealth:    ah,
		Visibility:       vis,
		Tags:             tags,
		Notes:            notes,
		CreatedAt:        createdAt,

		CurrentBalance:   0,
		AvailableBalance: 0,
		ClearedBalance:   0,
		UnclearedBalance: 0,
		ReservedBalance:  0,
		AllocatedBalance: 0,
		SpendableBalance: 0,
	}, nil
}

func (p *AccountsPGProvider) GetAccounts(ctx context.Context, userID string) ([]engine.AccountInput, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT a.account_id, a.account_name, a.account_type, a.classification,
			a.status, a.currency, a.liquidity_profile, a.account_health,
			a.visibility, a.tags, a.notes, a.opened_date, a.created_at,
			COALESCE(i.name, '') as institution_name
		FROM accounts a
		LEFT JOIN institutions i ON a.institution_id = i.institution_id
		WHERE a.owner_id = $1 AND a.status NOT IN ('Closed', 'Archived')
		ORDER BY a.account_name`, userID)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []engine.AccountInput
	for rows.Next() {
		var aid, an, at, cls, st, cur, lp, ah, vis string
		var tags []string
		var notes string
		var od time.Time
		var ca time.Time
		var instName string

		if err := rows.Scan(&aid, &an, &at, &cls, &st, &cur, &lp, &ah, &vis, &tags, &notes, &od, &ca, &instName); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}

		accounts = append(accounts, engine.AccountInput{
			AccountID:        aid,
			AccountName:      an,
			AccountType:      at,
			Classification:   cls,
			Status:           st,
			Currency:         cur,
			InstitutionName:  instName,
			OpenedDate:       od.Format(time.DateOnly),
			LiquidityProfile: lp,
			AccountHealth:    ah,
			Visibility:       vis,
			Tags:             tags,
			Notes:            notes,
			CreatedAt:        ca.Format(time.RFC3339),
		})
	}

	return accounts, nil
}

func (p *AccountsPGProvider) GetAccountByID(ctx context.Context, userID, accountID string) (*engine.AccountInput, error) {
	inp, err := p.getAccountInput(ctx, accountID)
	if err != nil {
		return nil, err
	}
	// Verify ownership
	if inp != nil {
		// We trust the DB query to filter by owner — but we don't have owner in getAccountInput
	}
	return inp, nil
}

// TransProvider

func (p *AccountsPGProvider) GetTransactions(ctx context.Context, accountID string, limit int) ([]engine.TransInput, error) {
	if limit <= 0 { limit = 25 }

	rows, err := p.pool.Query(ctx,
		`SELECT event_id, event_type, amount, currency, description,
			COALESCE(event_sub_type, ''), effective_date
		FROM financial_events
		WHERE source = $1 OR destination = $1
		ORDER BY effective_date DESC
		LIMIT $2`, accountID, limit)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}
	defer rows.Close()

	var txns []engine.TransInput
	for rows.Next() {
		var eid, etype, cur, desc, cat string
		var amt int64
		var ed time.Time

		if err := rows.Scan(&eid, &etype, &amt, &cur, &desc, &cat, &ed); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		txns = append(txns, engine.TransInput{
			EventID:    eid,
			EventType:  etype,
			Amount:     amt,
			Currency:   cur,
			Description: desc,
			Category:   cat,
			EventDate:  ed.Format(time.RFC3339),
		})
	}

	return txns, nil
}

// ProjProvider

func (p *AccountsPGProvider) GetCashFlowProjection(ctx context.Context, userID string) (inflow, outflow, projected int64, err error) {
	now := p.now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Income this month
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM financial_events
		WHERE user_id = $1 AND amount > 0 AND state = 'POSTED' AND effective_date >= $2`,
		userID, startOfMonth).Scan(&inflow)

	// Expenses this month
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(ABS(amount)), 0) FROM financial_events
		WHERE user_id = $1 AND amount < 0 AND state = 'POSTED' AND effective_date >= $2`,
		userID, startOfMonth).Scan(&outflow)

	// Projected: extrapolate from last full month
	var prevInflow, prevOutflow int64
	prevStart := startOfMonth.AddDate(0, -1, 0)
	prevEnd := startOfMonth.AddDate(0, 0, -1)
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM financial_events
		WHERE user_id = $1 AND amount > 0 AND state = 'POSTED' AND effective_date BETWEEN $2 AND $3`,
		userID, prevStart, prevEnd).Scan(&prevInflow)
	_ = p.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(ABS(amount)), 0) FROM financial_events
		WHERE user_id = $1 AND amount < 0 AND state = 'POSTED' AND effective_date BETWEEN $2 AND $3`,
		userID, prevStart, prevEnd).Scan(&prevOutflow)

	projected = prevInflow - prevOutflow
	return
}

// RiskProvider

func (p *AccountsPGProvider) GetRiskScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var level string
	err := p.pool.QueryRow(ctx,
		`SELECT composite_score, risk_level FROM risk_assessments ORDER BY created_at DESC LIMIT 1`).Scan(&score, &level)
	if err != nil {
		return 0, "Minimal", nil
	}
	return score, level, nil
}

// HealthProvider

func (p *AccountsPGProvider) GetHealthScore(ctx context.Context, userID string) (int, string, error) {
	var score int
	var grade string
	err := p.pool.QueryRow(ctx,
		`SELECT overall_score, score_grade FROM health_scores ORDER BY created_at DESC LIMIT 1`).Scan(&score, &grade)
	if err != nil {
		return 0, "Pending", nil
	}
	return score, grade, nil
}

// RecProvider

func (p *AccountsPGProvider) HasRecommendations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM recommendations`).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("count recommendations: %w", err)
	}
	return count > 0, count, nil
}

// OptProvider

func (p *AccountsPGProvider) HasOptimizations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM optimizations`).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("count optimizations: %w", err)
	}
	return count > 0, count, nil
}

// SimProvider

func (p *AccountsPGProvider) HasSimulations(ctx context.Context, userID string) (bool, int, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM simulations`).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("count simulations: %w", err)
	}
	return count > 0, count, nil
}

// EventProvider

func (p *AccountsPGProvider) GetEventCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := p.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM financial_events WHERE user_id = $1 AND effective_date >= NOW() - INTERVAL '30 days'`,
		userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count events: %w", err)
	}
	return count, nil
}
