package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/portfolio/internal/application/dto/command"
	"github.com/horizon/core/services/domains/portfolio/internal/application/dto/query"
	"github.com/horizon/core/services/domains/portfolio/internal/domain"
)

type PortfolioService struct {
	repo domain.Repository
	pub  domain.Publisher
	fac  *domain.PortfolioFactory
	now  func() time.Time
}

func NewPortfolioService(r domain.Repository, p domain.Publisher, n func() time.Time) *PortfolioService {
	return &PortfolioService{repo: r, pub: p, fac: domain.NewPortfolioFactory(), now: n}
}

func (s *PortfolioService) Create(ctx context.Context, cmd command.CreatePortfolioCommand) (*command.PfResult, error) {
	pt := cmd.PortfolioType
	if pt == "" { pt = "Investment" }
	mm := cmd.MembershipModel
	if mm == "" { mm = "Manual" }

	p, err := s.fac.Create("", cmd.Name, domain.PortfolioType(pt), cmd.OwnerID,
		cmd.BaseCurrency, domain.MembershipModel(mm), domain.RiskProfile(cmd.RiskProfile),
		domain.BenchmarkProfile(cmd.BenchmarkProfile), cmd.Benchmark, cmd.SourceOfTruth,
		cmd.Tags, cmd.Notes)
	if err != nil { return nil, err }
	if err := s.repo.Save(ctx, p); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.pub.Publish(domain.NewPortfolioCreated(p)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.PfResult{PortfolioID: p.ID(), Status: string(p.Status()), Success: true}, nil
}

func (s *PortfolioService) exec(ctx context.Context, id string, fn func(*domain.Portfolio) error, pubFn func(*domain.Portfolio) domain.DomainEvent) (*command.PfResult, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	if err := fn(p); err != nil { return nil, err }
	if err := s.repo.Save(ctx, p); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if pubFn != nil {
		if err := s.pub.Publish(pubFn(p)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	}
	return &command.PfResult{PortfolioID: p.ID(), Status: string(p.Status()), Success: true}, nil
}

func (s *PortfolioService) Activate(ctx context.Context, cmd command.IDCmd) (*command.PfResult, error) {
	return s.exec(ctx, cmd.PortfolioID, domain.ActivatePortfolio, func(p *domain.Portfolio) domain.DomainEvent { return domain.NewPortfolioActivated(p) })
}
func (s *PortfolioService) Rebalance(ctx context.Context, cmd command.IDCmd) (*command.PfResult, error) {
	return s.exec(ctx, cmd.PortfolioID, domain.RebalancePortfolio, func(p *domain.Portfolio) domain.DomainEvent { return domain.NewPortfolioRebalanced(p) })
}
func (s *PortfolioService) Archive(ctx context.Context, cmd command.IDCmd) (*command.PfResult, error) {
	return s.exec(ctx, cmd.PortfolioID, domain.ArchivePortfolio, func(p *domain.Portfolio) domain.DomainEvent { return domain.NewPortfolioArchived(p) })
}

func toView(p *domain.Portfolio) *query.PfView {
	return &query.PfView{
		PortfolioID: p.ID(), Name: p.Name(), PortfolioType: string(p.PortfolioType()),
		BaseCurrency: p.BaseCurrency(), Status: string(p.Status()),
		CurrentValue: p.CurrentValue(), CostBasis: p.CostBasis(),
		Health: string(p.Health()), RiskProfile: string(p.RiskProfile()),
		MemberCount: len(p.Members()), CreatedAt: p.CreatedAt().Format(time.RFC3339),
	}
}
func toPag(ps []*domain.Portfolio, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{Portfolios: make([]query.PfView, 0, len(ps)), NextCursor: cursor, HasMore: cursor != ""}
	for _, p := range ps { r.Portfolios = append(r.Portfolios, *toView(p)) }; return r
}

func (s *PortfolioService) GetByID(ctx context.Context, q query.GetPfQuery) (*query.PfView, error) {
	p, err := s.repo.GetByID(ctx, q.PortfolioID)
	if err != nil { return nil, err }; return toView(p), nil
}
func (s *PortfolioService) ListByUser(ctx context.Context, q query.ListByUserQuery) (*query.PaginatedResult, error) {
	ps, c, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil { return nil, err }; return toPag(ps, c), nil
}
func (s *PortfolioService) ListByType(ctx context.Context, q query.ListByTypeQuery) (*query.PaginatedResult, error) {
	ps, c, err := s.repo.ListByType(ctx, q.UserID, domain.PortfolioType(q.PfType), q.Cursor, q.Limit)
	if err != nil { return nil, err }; return toPag(ps, c), nil
}
