package domain

import "context"

type Repository interface {
	Save(ctx context.Context, inst *Institution) error
	UpdateStatus(ctx context.Context, id string, from, to InstitutionStatus) error
	GetByID(ctx context.Context, id string) (*Institution, error)
	ListByType(ctx context.Context, instType InstitutionType, cursor string, limit int) ([]*Institution, string, error)
	ListByCountry(ctx context.Context, country string, cursor string, limit int) ([]*Institution, string, error)
	ListByStatus(ctx context.Context, status InstitutionStatus, cursor string, limit int) ([]*Institution, string, error)
	Search(ctx context.Context, query string, cursor string, limit int) ([]*Institution, string, error)
}
