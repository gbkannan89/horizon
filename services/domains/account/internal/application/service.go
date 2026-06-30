package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/account/internal/application/dto/command"
	"github.com/horizon/core/services/domains/account/internal/application/dto/query"
	"github.com/horizon/core/services/domains/account/internal/domain"
)

type AccountService struct {
	repo      domain.Repository
	publisher domain.Publisher
	factory   *domain.AccountFactory
	now       func() time.Time
}

func NewAccountService(repo domain.Repository, publisher domain.Publisher, now func() time.Time) *AccountService {
	return &AccountService{repo: repo, publisher: publisher, factory: domain.NewAccountFactory(), now: now}
}

func (s *AccountService) Create(ctx context.Context, cmd command.CreateAccountCommand) (*command.AccountResult, error) {
	opened := s.now()
	if cmd.OpenedDate != "" {
		t, err := time.Parse("2006-01-02", cmd.OpenedDate)
		if err == nil { opened = t }
	}
	vis := domain.Visibility(cmd.Visibility)
	if vis == "" { vis = domain.VisPrivate }
	own := domain.OwnershipModel(cmd.Ownership)
	if own == "" { own = domain.OMPersonal }

	a, err := s.factory.Create("", cmd.OwnerID, cmd.AccountName, cmd.Currency,
		domain.AccountType(cmd.AccountType), domain.AccountClassification(cmd.Classification),
		opened, vis, own, cmd.CreditLimit, cmd.InterestRate,
		cmd.InstitutionID, cmd.HouseholdID, cmd.SubType, cmd.Country,
		cmd.SourceOfTruth, cmd.ExtRef, cmd.Tags, nil, cmd.Notes)
	if err != nil { return nil, err }

	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAccountCreated(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AccountResult{AccountID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AccountService) Activate(ctx context.Context, cmd command.ActivateAccountCommand) (*command.AccountResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AccountID); if err != nil { return nil, err }
	if err := domain.ActivateAccount(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAccountActivated(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AccountResult{AccountID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AccountService) Freeze(ctx context.Context, cmd command.FreezeAccountCommand) (*command.AccountResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AccountID); if err != nil { return nil, err }
	if err := domain.FreezeAccount(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAccountFrozen(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AccountResult{AccountID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AccountService) Unfreeze(ctx context.Context, cmd command.UnfreezeAccountCommand) (*command.AccountResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AccountID); if err != nil { return nil, err }
	if err := domain.UnfreezeAccount(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAccountUnfrozen(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AccountResult{AccountID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AccountService) Close(ctx context.Context, cmd command.CloseAccountCommand) (*command.AccountResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AccountID); if err != nil { return nil, err }
	cd, _ := time.Parse("2006-01-02", cmd.ClosedDate)
	if err := domain.CloseAccount(a, cd); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAccountClosed(a, cmd.Reason)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AccountResult{AccountID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AccountService) Archive(ctx context.Context, cmd command.ArchiveAccountCommand) (*command.AccountResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AccountID); if err != nil { return nil, err }
	if err := domain.ArchiveAccount(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAccountArchived(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AccountResult{AccountID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func toAccountView(a *domain.Account) *query.AccountView {
	return &query.AccountView{
		AccountID: a.ID(), AccountName: a.AccountName(), AccountType: string(a.AccountType()),
		Classification: string(a.Classification()), Currency: a.Currency(),
		Status: string(a.Status()), OwnerID: a.OwnerID(),
		Liquidity: string(a.LiquidityProfile()), Health: string(a.AccountHealth()),
		CreatedAt: a.CreatedAt().Format(time.RFC3339),
	}
}

func toPaginated(accts []*domain.Account, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{Accounts: make([]query.AccountView, 0, len(accts)), NextCursor: cursor, HasMore: cursor != ""}
	for _, a := range accts { r.Accounts = append(r.Accounts, *toAccountView(a)) }
	return r
}

func (s *AccountService) GetByID(ctx context.Context, q query.GetAccountQuery) (*query.AccountView, error) {
	a, err := s.repo.GetByID(ctx, q.AccountID); if err != nil { return nil, err }
	return toAccountView(a), nil
}
func (s *AccountService) ListByUser(ctx context.Context, q query.ListByUserQuery) (*query.PaginatedResult, error) {
	accts, c, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil { return nil, err }; return toPaginated(accts, c), nil
}
func (s *AccountService) ListByType(ctx context.Context, q query.ListByTypeQuery) (*query.PaginatedResult, error) {
	accts, c, err := s.repo.ListByType(ctx, q.UserID, domain.AccountType(q.AcctType), q.Cursor, q.Limit)
	if err != nil { return nil, err }; return toPaginated(accts, c), nil
}
func (s *AccountService) ListByStatus(ctx context.Context, q query.ListByStatusQuery) (*query.PaginatedResult, error) {
	accts, c, err := s.repo.ListByStatus(ctx, q.UserID, domain.AccountStatus(q.Status), q.Cursor, q.Limit)
	if err != nil { return nil, err }; return toPaginated(accts, c), nil
}
