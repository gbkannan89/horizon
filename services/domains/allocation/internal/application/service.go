package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/allocation/internal/application/dto/command"
	"github.com/horizon/core/services/domains/allocation/internal/application/dto/query"
	"github.com/horizon/core/services/domains/allocation/internal/domain"
)

type AllocationService struct {
	repo      domain.Repository
	publisher domain.Publisher
	factory   *domain.AllocationFactory
	now       func() time.Time
}

func NewAllocationService(repo domain.Repository, pub domain.Publisher, now func() time.Time) *AllocationService {
	return &AllocationService{repo: repo, publisher: pub, factory: domain.NewAllocationFactory(), now: now}
}

func (s *AllocationService) Create(ctx context.Context, cmd command.CreateAllocationCommand) (*command.AllocResult, error) {
	eff, _ := time.Parse("2006-01-02", cmd.EffectiveDate)
	var exp *time.Time
	if cmd.ExpirationDate != nil { t, _ := time.Parse("2006-01-02", *cmd.ExpirationDate); exp = &t }
	at := domain.AllocationType(cmd.AllocationType)
	if at == "" { at = domain.ATFixedAmount }
	strat := domain.AllocationStrategy(cmd.Strategy)
	if strat == "" { strat = domain.ASGoalPriority }
	sot := cmd.SourceOfTruth; if sot == "" { sot = "User" }
	cb := cmd.CreatedBy; if cb == "" { cb = "User" }

	a, err := s.factory.Create("", cmd.GoalID, cmd.FundingSourceID, domain.FSTAccount,
		at, strat, cmd.Priority, cmd.Currency, eff, cmd.Weight, cmd.FixedAmount, exp, sot, cb, cmd.Tags, cmd.Notes)
	if err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAllocationCreated(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AllocationService) Plan(ctx context.Context, cmd command.PlanAllocationCommand) (*command.AllocResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AllocationID); if err != nil { return nil, err }
	if a.Status() != domain.StDraft { return nil, fmt.Errorf("only draft allocations can be planned") }
	a.SetStatus(domain.StPlanned)
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func (s *AllocationService) Activate(ctx context.Context, cmd command.ActivateAllocationCommand) (*command.AllocResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AllocationID); if err != nil { return nil, err }
	if err := domain.ActivateAlloc(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAllocationActivated(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}
func (s *AllocationService) Pause(ctx context.Context, cmd command.PauseAllocationCommand) (*command.AllocResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AllocationID); if err != nil { return nil, err }
	if err := domain.PauseAlloc(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAllocationPaused(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}
func (s *AllocationService) Resume(ctx context.Context, cmd command.ResumeAllocationCommand) (*command.AllocResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AllocationID); if err != nil { return nil, err }
	if err := domain.ResumeAlloc(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAllocationResumed(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}
func (s *AllocationService) Complete(ctx context.Context, cmd command.CompleteAllocationCommand) (*command.AllocResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AllocationID); if err != nil { return nil, err }
	if err := domain.CompleteAlloc(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAllocationCompleted(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}
func (s *AllocationService) Cancel(ctx context.Context, cmd command.CancelAllocationCommand) (*command.AllocResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AllocationID); if err != nil { return nil, err }
	if err := domain.CancelAlloc(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAllocationCancelled(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}
func (s *AllocationService) Archive(ctx context.Context, cmd command.ArchiveAllocationCommand) (*command.AllocResult, error) {
	a, err := s.repo.GetByID(ctx, cmd.AllocationID); if err != nil { return nil, err }
	if err := domain.ArchiveAlloc(a); err != nil { return nil, err }
	if err := s.repo.Save(ctx, a); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewAllocationArchived(a)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.AllocResult{AllocationID: a.ID(), Status: string(a.Status()), Success: true}, nil
}

func toView(a *domain.Allocation) *query.AllocView {
	return &query.AllocView{
		AllocationID: a.ID(), GoalID: a.GoalID(), FundingSourceID: a.FundingSourceID(),
		AllocationType: string(a.AllocationType()), Priority: a.Priority(), Currency: a.Currency(),
		Status: string(a.Status()), Health: string(a.Health()),
		ReservedAmount: a.ReservedAmount(), AllocatedAmount: a.AllocatedAmount(),
		CreatedAt: a.CreatedAt().Format(time.RFC3339),
	}
}
func toPag(allocs []*domain.Allocation, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{Allocations: make([]query.AllocView, 0, len(allocs)), NextCursor: cursor, HasMore: cursor != ""}
	for _, a := range allocs { r.Allocations = append(r.Allocations, *toView(a)) }; return r
}

func (s *AllocationService) GetByID(ctx context.Context, q query.GetByIDQuery) (*query.AllocView, error) {
	a, err := s.repo.GetByID(ctx, q.AllocationID); if err != nil { return nil, err }; return toView(a), nil
}
func (s *AllocationService) ListByGoal(ctx context.Context, q query.ListByGoalQuery) (*query.PaginatedResult, error) {
	as, c, err := s.repo.ListByGoal(ctx, q.GoalID, q.Cursor, q.Limit); if err != nil { return nil, err }; return toPag(as, c), nil
}
func (s *AllocationService) ListBySource(ctx context.Context, q query.ListBySourceQuery) (*query.PaginatedResult, error) {
	as, c, err := s.repo.ListByFundingSource(ctx, q.SourceID, q.Cursor, q.Limit); if err != nil { return nil, err }; return toPag(as, c), nil
}
