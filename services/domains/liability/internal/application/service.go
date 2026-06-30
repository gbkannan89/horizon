package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/liability/internal/application/dto/command"
	"github.com/horizon/core/services/domains/liability/internal/application/dto/query"
	"github.com/horizon/core/services/domains/liability/internal/domain"
)

type LiabilityService struct {
	repo domain.Repository
	pub  domain.Publisher
	fac  *domain.LiabilityFactory
	now  func() time.Time
}

func NewLiabilityService(r domain.Repository, p domain.Publisher, n func() time.Time) *LiabilityService {
	return &LiabilityService{repo: r, pub: p, fac: domain.NewLiabilityFactory(), now: n}
}

func (s *LiabilityService) Create(ctx context.Context, cmd command.CreateLiabilityCommand) (*command.LiabResult, error) {
	md, err := time.Parse("2006-01-02", cmd.MaturityDate)
	if err != nil { return nil, fmt.Errorf("invalid date: %w", err) }
	im := cmd.InterestModel; if im == "" { im = "Simple" }
	rp := cmd.RepaymentProfile; if rp == "" { rp = "FixedInstallment" }
	rm := cmd.RepaymentModel; if rm == "" { rm = "FixedEMI" }

	l, err := s.fac.Create("", cmd.Name, domain.LiabilityClassification(cmd.Classification),
		cmd.OwnerID, cmd.Currency, cmd.OriginalPrincipal, cmd.InterestRate,
		domain.InterestModel(im), domain.RepaymentProfile(rp), domain.RepaymentModel(rm),
		cmd.InstallmentAmount, cmd.InstallmentFreq, cmd.RemainingInstallments, md,
		cmd.SourceOfTruth, cmd.LiabilityType, cmd.HouseholdID, cmd.InstitutionID,
		cmd.ServicingAccountID, cmd.Collateral, cmd.ExtRef, cmd.Tags, cmd.Notes)
	if err != nil { return nil, err }
	if err := s.repo.Save(ctx, l); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.pub.Publish(domain.NewLiabilityCreated(l)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.LiabResult{LiabilityID: l.ID(), Status: string(l.Status()), Success: true}, nil
}

func (s *LiabilityService) exec(ctx context.Context, id string, fn func(*domain.Liability) error, pubFn func(*domain.Liability) domain.DomainEvent) (*command.LiabResult, error) {
	l, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	if err := fn(l); err != nil { return nil, err }
	if err := s.repo.Save(ctx, l); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if pf := pubFn; pf != nil {
		if err := s.pub.Publish(pf(l)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	}
	return &command.LiabResult{LiabilityID: l.ID(), Status: string(l.Status()), Success: true}, nil
}

func (s *LiabilityService) Activate(ctx context.Context, cmd command.IDCmd) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, domain.ActivateLiability, func(l *domain.Liability) domain.DomainEvent { return domain.NewLiabilityActivated(l) })
}
func (s *LiabilityService) MakePayment(ctx context.Context, cmd command.PaymentCommand) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, func(l *domain.Liability) error { return domain.ApplyPayment(l, cmd.Amount) }, func(l *domain.Liability) domain.DomainEvent { return domain.NewPaymentApplied(l, cmd.Amount) })
}
func (s *LiabilityService) ApplyPrepayment(ctx context.Context, cmd command.PrepaymentCommand) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, func(l *domain.Liability) error { return domain.ApplyPrepayment(l, cmd.Amount) }, func(l *domain.Liability) domain.DomainEvent { return domain.NewPrepaymentApplied(l, cmd.Amount) })
}
func (s *LiabilityService) GracePeriod(ctx context.Context, cmd command.IDCmd) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, domain.EnterGracePeriod, func(l *domain.Liability) domain.DomainEvent { return domain.NewGracePeriodStarted(l) })
}
func (s *LiabilityService) MarkDelinquent(ctx context.Context, cmd command.IDCmd) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, domain.MarkDelinquent, func(l *domain.Liability) domain.DomainEvent { return domain.NewLiabilityDelinquent(l) })
}
func (s *LiabilityService) Settle(ctx context.Context, cmd command.IDCmd) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, domain.SettleLiability, func(l *domain.Liability) domain.DomainEvent { return domain.NewLiabilitySettled(l) })
}
func (s *LiabilityService) WriteOff(ctx context.Context, cmd command.IDCmd) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, domain.WriteOffLiability, func(l *domain.Liability) domain.DomainEvent { return domain.NewLiabilityWrittenOff(l) })
}
func (s *LiabilityService) Archive(ctx context.Context, cmd command.IDCmd) (*command.LiabResult, error) {
	return s.exec(ctx, cmd.LiabilityID, domain.ArchiveLiability, func(l *domain.Liability) domain.DomainEvent { return domain.NewLiabilityArchived(l) })
}

func toView(l *domain.Liability) *query.LiabView {
	return &query.LiabView{
		LiabilityID: l.ID(), Name: l.Name(), Classification: string(l.Classification()),
		Currency: l.Currency(), OriginalPrincipal: l.OriginalPrincipal(),
		OutstandingBalance: l.OutstandingBalance(), InterestRate: l.InterestRate(),
		InterestModel: string(l.InterestModel()), Status: string(l.Status()),
		RemainingInstallments: l.RemainingInstallments(),
		MaturityDate: l.MaturityDate().Format("2006-01-02"),
		Health: string(l.LiabilityHealth()), CreatedAt: l.CreatedAt().Format(time.RFC3339),
	}
}

func toPag(ls []*domain.Liability, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{Liabilities: make([]query.LiabView, 0, len(ls)), NextCursor: cursor, HasMore: cursor != ""}
	for _, l := range ls { r.Liabilities = append(r.Liabilities, *toView(l)) }; return r
}

func (s *LiabilityService) GetByID(ctx context.Context, q query.GetLiabilityQuery) (*query.LiabView, error) {
	l, err := s.repo.GetByID(ctx, q.LiabilityID)
	if err != nil { return nil, err }; return toView(l), nil
}
func (s *LiabilityService) ListByUser(ctx context.Context, q query.ListByUserQuery) (*query.PaginatedResult, error) {
	ls, c, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil { return nil, err }; return toPag(ls, c), nil
}
