package domain

import (
	"context"
	"time"
)

// Repository defines persistence operations for assets.
type Repository interface {
	Save(ctx context.Context, asset *Asset) error
	UpdateStatus(ctx context.Context, assetID string, fromStatus, toStatus string) error
	UpdateValuation(ctx context.Context, assetID string, value int64, method string, date time.Time) error
	UpdateQuantity(ctx context.Context, assetID string, quantity float64) error
	UpdateCostBasis(ctx context.Context, assetID string, costBasis int64) error

	GetByID(ctx context.Context, assetID string) (*Asset, error)
	ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*Asset, string, error)
	ListByClassification(ctx context.Context, userID string, classification string, cursor string, limit int) ([]*Asset, string, error)
	ListByStatus(ctx context.Context, userID string, status string, cursor string, limit int) ([]*Asset, string, error)
	ListByInstitution(ctx context.Context, institutionID string, cursor string, limit int) ([]*Asset, string, error)
	ListByAccount(ctx context.Context, accountID string, cursor string, limit int) ([]*Asset, string, error)
	GetValuationHistory(ctx context.Context, assetID string, cursor string, limit int) ([]*ValuationRecord, string, error)
}
