package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/domains/asset/internal/domain"
)

type AssetRepository struct {
	pool *pgxpool.Pool
}

func NewAssetRepository(pool *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{pool: pool}
}

func (r *AssetRepository) Save(ctx context.Context, a *domain.Asset) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO assets (asset_id, asset_name, classification, ownership_model, owner_id, currency,
			cost_basis, valuation_profile, valuation_method, valuation_date, quantity, unit_price,
			liquidity_profile, ownership_percentage, status, source_of_truth, tags, version, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,NOW(),$10,$11,$12,$13,$14,$15,$16,1,NOW(),NOW())
		ON CONFLICT (asset_id) DO UPDATE SET cost_basis=$7, status=$14, updated_at=NOW()`,
		a.ID(), a.AssetName(), string(a.Classification()), string(a.OwnershipModel()),
		a.OwnerID(), a.Currency(), a.CostBasis(), string(a.ValuationProfile()),
		string(a.ValuationMethod()), a.Quantity(), a.UnitPrice(), string(a.LiquidityProfile()),
		a.OwnershipPercentage(), string(a.Status()), 		a.SourceOfTruth())
	return err
}

func (r *AssetRepository) UpdateStatus(ctx context.Context, assetID string, fromStatus, toStatus string) error {
	_, err := r.pool.Exec(ctx, `UPDATE assets SET status=$1, updated_at=NOW() WHERE asset_id=$2 AND status=$3`,
		toStatus, assetID, fromStatus)
	return err
}

func (r *AssetRepository) UpdateValuation(ctx context.Context, assetID string, value int64, method string, date time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE assets SET unit_price=$1, valuation_method=$2, valuation_date=$3, updated_at=NOW() WHERE asset_id=$4`,
		value, method, date, assetID)
	_, err2 := r.pool.Exec(ctx,
		`INSERT INTO valuation_history (asset_id, value, method, valuation_date) VALUES ($1,$2,$3,$4)`,
		assetID, value, method, date)
	if err != nil { return err }
	return err2
}

func (r *AssetRepository) UpdateQuantity(ctx context.Context, assetID string, quantity float64) error {
	_, err := r.pool.Exec(ctx, `UPDATE assets SET quantity=$1, updated_at=NOW() WHERE asset_id=$2`, quantity, assetID)
	return err
}

func (r *AssetRepository) UpdateCostBasis(ctx context.Context, assetID string, costBasis int64) error {
	_, err := r.pool.Exec(ctx, `UPDATE assets SET cost_basis=$1, updated_at=NOW() WHERE asset_id=$2`, costBasis, assetID)
	return err
}

func (r *AssetRepository) GetByID(ctx context.Context, assetID string) (*domain.Asset, error) {
	var id, an, cls, oid, cur, vp, vm, lp, st, sot string
	var tags []string
	var cb, op float64
	var qty *float64
	var up float64
	var vd, ca time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT asset_id, asset_name, classification, owner_id, currency, cost_basis,
			valuation_profile, valuation_method, valuation_date, quantity, unit_price,
			liquidity_profile, ownership_percentage, status, source_of_truth, tags, created_at
		FROM assets WHERE asset_id = $1`, assetID).Scan(&id, &an, &cls, &oid, &cur, &cb,
		&vp, &vm, &vd, &qty, &up, &lp, &op, &st, &sot, &tags, &ca)
	if err != nil { return nil, fmt.Errorf("get asset: %w", err) }

	return domain.ReconstructFromDB(id, an, domain.AssetClassification(cls), oid, cur,
		int64(cb), domain.ValuationProfile(vp), domain.ValuationMethod(vm), vd,
		qty, up, domain.LiquidityProfile(lp), op, domain.AssetStatus(st), sot, tags, ca), nil
}

func (r *AssetRepository) ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*domain.Asset, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT asset_id, asset_name, classification, owner_id, currency, cost_basis, status, created_at
		FROM assets WHERE owner_id = $1 ORDER BY asset_name LIMIT $2`, userID, limit+1)
	defer rows.Close()
	return scanAssetList(rows, limit)
}

func (r *AssetRepository) ListByClassification(ctx context.Context, userID string, classification string, cursor string, limit int) ([]*domain.Asset, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT asset_id, asset_name, classification, owner_id, currency, cost_basis, status, created_at
		FROM assets WHERE owner_id = $1 AND classification = $2 ORDER BY asset_name LIMIT $3`,
		userID, classification, limit+1)
	defer rows.Close()
	return scanAssetList(rows, limit)
}

func (r *AssetRepository) ListByStatus(ctx context.Context, userID string, status string, cursor string, limit int) ([]*domain.Asset, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT asset_id, asset_name, classification, owner_id, currency, cost_basis, status, created_at
		FROM assets WHERE owner_id = $1 AND status = $2 ORDER BY asset_name LIMIT $3`,
		userID, status, limit+1)
	defer rows.Close()
	return scanAssetList(rows, limit)
}

func (r *AssetRepository) ListByInstitution(ctx context.Context, institutionID string, cursor string, limit int) ([]*domain.Asset, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT asset_id, asset_name, classification, owner_id, currency, cost_basis, status, created_at
		FROM assets WHERE institution_id = $1 ORDER BY asset_name LIMIT $2`, institutionID, limit+1)
	defer rows.Close()
	return scanAssetList(rows, limit)
}

func (r *AssetRepository) ListByAccount(ctx context.Context, accountID string, cursor string, limit int) ([]*domain.Asset, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT asset_id, asset_name, classification, owner_id, currency, cost_basis, status, created_at
		FROM assets WHERE account_id = $1 ORDER BY asset_name LIMIT $2`, accountID, limit+1)
	defer rows.Close()
	return scanAssetList(rows, limit)
}

func (r *AssetRepository) GetValuationHistory(ctx context.Context, assetID string, cursor string, limit int) ([]*domain.ValuationRecord, string, error) {
	if limit <= 0 { limit = 25 }
	rows, _ := r.pool.Query(ctx,
		`SELECT asset_id, value, method, valuation_date, source_of_truth, recorded_at::text
		FROM valuation_history WHERE asset_id = $1 ORDER BY valuation_date DESC LIMIT $2`, assetID, limit+1)
	defer rows.Close()

	var records []*domain.ValuationRecord
	for rows.Next() {
		var r domain.ValuationRecord
		if err := rows.Scan(&r.AssetID, &r.Value, &r.Method, &r.ValuationDate, &r.SourceOfTruth, &r.RecordedAt); err != nil {
			return nil, "", err
		}
		records = append(records, &r)
	}
	hasMore := len(records) > limit
	if hasMore { records = records[:limit] }
	return records, "", nil
}

func scanAssetList(rows pgx.Rows, limit int) ([]*domain.Asset, string, error) {
	var assets []*domain.Asset
	for rows.Next() {
		var id, an, cls, oid, cur, st string
		var cb int64
		var ca time.Time
		if err := rows.Scan(&id, &an, &cls, &oid, &cur, &cb, &st, &ca); err != nil {
			return nil, "", err
		}
		assets = append(assets, domain.ReconstructFromDB(id, an, domain.AssetClassification(cls), oid, cur,
			cb, "", "", time.Time{}, nil, 0, "", 100, domain.AssetStatus(st), "", nil, ca))
	}
	hasMore := len(assets) > limit
	if hasMore { assets = assets[:limit] }
	return assets, "", nil
}
