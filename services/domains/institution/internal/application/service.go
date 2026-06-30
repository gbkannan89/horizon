package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/institution/internal/application/dto/command"
	"github.com/horizon/core/services/domains/institution/internal/application/dto/query"
	"github.com/horizon/core/services/domains/institution/internal/domain"
)

type InstService struct {
	repo domain.Repository
	pub  domain.Publisher
	fac  *domain.InstitutionFactory
	now  func() time.Time
}

func NewInstService(r domain.Repository, p domain.Publisher, n func() time.Time) *InstService {
	return &InstService{repo: r, pub: p, fac: domain.NewInstitutionFactory(), now: n}
}

func (s *InstService) Register(ctx context.Context, cmd command.RegisterInstCommand) (*command.InstResult, error) {
	i, err := s.fac.Create("", cmd.Name, cmd.Country, domain.InstitutionType(cmd.InstType),
		nil, nil, cmd.Website, cmd.Phone, cmd.Email, cmd.Address, cmd.HQ,
		cmd.Regulator, cmd.RegLicense, cmd.SourceOfTruth, nil, cmd.Tags, cmd.Notes)
	if err != nil {
		return nil, err
	}
	if err := domain.RegisterInstitution(i); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, i); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(domain.NewInstRegistered(i)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.InstResult{InstitutionID: i.ID(), Status: string(i.Status()), Success: true}, nil
}

func (s *InstService) execState(ctx context.Context, id string, fn func(*domain.Institution) error, pubFn func(*domain.Institution) domain.DomainEvent) (*command.InstResult, error) {
	i, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := fn(i); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, i); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.pub.Publish(pubFn(i)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}
	return &command.InstResult{InstitutionID: i.ID(), Status: string(i.Status()), Success: true}, nil
}

func (s *InstService) Verify(ctx context.Context, cmd command.InstIDCmd) (*command.InstResult, error) {
	return s.execState(ctx, cmd.InstitutionID, domain.VerifyInstitution, func(i *domain.Institution) domain.DomainEvent { return domain.NewInstVerified(i) })
}
func (s *InstService) Activate(ctx context.Context, cmd command.InstIDCmd) (*command.InstResult, error) {
	return s.execState(ctx, cmd.InstitutionID, domain.ActivateInstitution, func(i *domain.Institution) domain.DomainEvent { return domain.NewInstActivated(i) })
}
func (s *InstService) Suspend(ctx context.Context, cmd command.InstIDCmd) (*command.InstResult, error) {
	return s.execState(ctx, cmd.InstitutionID, domain.SuspendInstitution, func(i *domain.Institution) domain.DomainEvent { return domain.NewInstSuspended(i) })
}
func (s *InstService) Close(ctx context.Context, cmd command.InstIDCmd) (*command.InstResult, error) {
	return s.execState(ctx, cmd.InstitutionID, domain.CloseInstitution, func(i *domain.Institution) domain.DomainEvent { return domain.NewInstClosed(i) })
}
func (s *InstService) Archive(ctx context.Context, cmd command.InstIDCmd) (*command.InstResult, error) {
	return s.execState(ctx, cmd.InstitutionID, domain.ArchiveInstitution, func(i *domain.Institution) domain.DomainEvent { return domain.NewInstArchived(i) })
}

func toView(i *domain.Institution) *query.InstView {
	return &query.InstView{
		InstitutionID: i.ID(), Name: i.Name(), InstType: string(i.Type()),
		Category: i.Category(), Status: string(i.Status()), Country: i.Country(),
		TrustLevel: string(i.TrustLevel()), Health: string(i.Health()),
		Confidence: string(i.Confidence()), CreatedAt: i.CreatedAt().Format(time.RFC3339),
	}
}

func toPag(insts []*domain.Institution, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{
		Institutions: make([]query.InstView, 0, len(insts)),
		NextCursor:   cursor, HasMore: cursor != "",
	}
	for _, i := range insts {
		r.Institutions = append(r.Institutions, *toView(i))
	}
	return r
}

func (s *InstService) GetByID(ctx context.Context, q query.GetByIDQuery) (*query.InstView, error) {
	i, err := s.repo.GetByID(ctx, q.InstitutionID)
	if err != nil {
		return nil, err
	}
	return toView(i), nil
}

func (s *InstService) ListByType(ctx context.Context, q query.ListByTypeQuery) (*query.PaginatedResult, error) {
	is, c, err := s.repo.ListByType(ctx, domain.InstitutionType(q.InstType), q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPag(is, c), nil
}

func (s *InstService) ListByCountry(ctx context.Context, q query.ListByCountryQuery) (*query.PaginatedResult, error) {
	is, c, err := s.repo.ListByCountry(ctx, q.Country, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPag(is, c), nil
}

func (s *InstService) ListByStatus(ctx context.Context, q query.ListByStatusQuery) (*query.PaginatedResult, error) {
	is, c, err := s.repo.ListByStatus(ctx, domain.InstitutionStatus(q.Status), q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPag(is, c), nil
}

func (s *InstService) Search(ctx context.Context, q query.SearchQuery) (*query.PaginatedResult, error) {
	is, c, err := s.repo.Search(ctx, q.Query, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPag(is, c), nil
}
