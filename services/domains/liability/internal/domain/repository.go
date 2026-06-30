package domain

import "context"

type Repository interface {
	Save(ctx context.Context, liability *Liability) error
	UpdateStatus(ctx context.Context, id string, from, to LiabilityStatus) error
	GetByID(ctx context.Context, id string) (*Liability, error)
	ListByUser(ctx context.Context, ownerID string, cursor string, limit int) ([]*Liability, string, error)
	ListByClassification(ctx context.Context, ownerID string, cls string, cursor string, limit int) ([]*Liability, string, error)
	ListByStatus(ctx context.Context, ownerID string, status LiabilityStatus, cursor string, limit int) ([]*Liability, string, error)
	ListByInstitution(ctx context.Context, institutionID string, cursor string, limit int) ([]*Liability, string, error)
}
