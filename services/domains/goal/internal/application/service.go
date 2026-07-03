package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/goal/internal/application/dto/command"
	"github.com/horizon/core/services/domains/goal/internal/application/dto/query"
	"github.com/horizon/core/services/domains/goal/internal/domain"
)

type GoalService struct {
	repo      domain.Repository
	publisher domain.Publisher
	factory   *domain.GoalFactory
	now       func() time.Time
}

func NewGoalService(repo domain.Repository, publisher domain.Publisher, now func() time.Time) *GoalService {
	return &GoalService{repo: repo, publisher: publisher, factory: domain.NewGoalFactory(), now: now}
}

func (s *GoalService) Create(ctx context.Context, cmd command.CreateGoalCommand) (*command.GoalResult, error) {
	imp := domain.GoalImportance(cmd.Importance)
	gtype := domain.GoalType(cmd.GoalType)
	subtype := domain.GoalSubtype(cmd.Subtype)

	var td *time.Time
	if cmd.TargetDate != nil {
		t, err := time.Parse("2006-01-02", *cmd.TargetDate)
		if err != nil {
			return nil, fmt.Errorf("invalid target date: %w", err)
		}
		td = &t
	}
	sc := domain.SuccessCriteria{Model: domain.SuccessModel(cmd.SuccessModel), TargetValue: cmd.TargetValue, TargetMonths: cmd.TargetMonths, CustomDesc: cmd.CustomDesc, CustomTargetVal: cmd.CustomTargetVal}

	var hhid *string // M-027 TODO: Set from command if provided
	now := s.now()
	goal, err := s.factory.Create("", cmd.UserID, hhid, cmd.Name, imp, gtype, subtype, sc, cmd.Priority, td, true, cmd.RiskTolerance, cmd.Notes, nil, cmd.Tags, nil, now)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, goal); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalCreated(goal)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: goal.ID(), Status: string(goal.Status()), Success: true}, nil
}

func (s *GoalService) Activate(ctx context.Context, cmd command.ActivateGoalCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	if err := domain.ActivateGoal(g); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalActivated(g)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) Pause(ctx context.Context, cmd command.PauseGoalCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	if err := domain.PauseGoal(g); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalPaused(g, cmd.Reason)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) Resume(ctx context.Context, cmd command.ResumeGoalCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	if err := domain.ResumeGoal(g); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalResumed(g)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) Restructure(ctx context.Context, cmd command.RestructureGoalCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	sc := domain.SuccessCriteria{Model: domain.SuccessModel(cmd.SuccessModel), TargetValue: cmd.TargetValue, TargetMonths: cmd.TargetMonths, CustomDesc: cmd.CustomDesc, CustomTargetVal: cmd.CustomTargetVal}
	var td *time.Time
	if cmd.TargetDate != nil {
		t, err := time.Parse("2006-01-02", *cmd.TargetDate)
		if err != nil {
			return nil, fmt.Errorf("invalid target date: %w", err)
		}
		td = &t
	}
	if err := domain.RestructureGoal(g, sc, td, cmd.Priority, s.now()); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalRestructured(g)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) Complete(ctx context.Context, cmd command.CompleteGoalCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	if err := domain.CompleteGoal(g); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalCompleted(g)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) Miss(ctx context.Context, cmd command.MissGoalCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	if err := domain.MissGoal(g); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalMissed(g)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) Archive(ctx context.Context, cmd command.ArchiveGoalCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	if err := domain.ArchiveGoal(g); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalArchived(g)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) UpdatePriority(ctx context.Context, cmd command.UpdatePriorityCommand) (*command.GoalResult, error) {
	g, err := s.repo.GetByID(ctx, cmd.GoalID)
	if err != nil {
		return nil, err
	}
	if err := domain.ValidatePriority(cmd.NewPriority); err != nil {
		return nil, err
	}
	conflict, _ := s.repo.FindConflict(ctx, g.UserID(), cmd.NewPriority, g.ID())
	if conflict != nil {
		return nil, errors.New("priority conflict: another active goal has this priority")
	}
	old := g.Priority()
	domain.UpdatePriority(g, cmd.NewPriority)
	if err := s.repo.Save(ctx, g); err != nil {
		return nil, fmt.Errorf("save goal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewGoalPriorityChanged(g, old)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.GoalResult{GoalID: g.ID(), Status: string(g.Status()), Success: true}, nil
}

func (s *GoalService) GetGoal(ctx context.Context, q query.GetGoalQuery) (*query.GoalView, error) {
	g, err := s.repo.GetByID(ctx, q.GoalID)
	if err != nil {
		return nil, err
	}
	return toGoalView(g), nil
}

func (s *GoalService) ListByUser(ctx context.Context, q query.ListByUserQuery) (*query.PaginatedResult, error) {
	goals, cursor, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginated(goals, cursor), nil
}

func (s *GoalService) ListByStatus(ctx context.Context, q query.ListByStatusQuery) (*query.PaginatedResult, error) {
	goals, cursor, err := s.repo.ListByStatus(ctx, q.UserID, domain.GoalStatus(q.Status), q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginated(goals, cursor), nil
}

func (s *GoalService) ListByImportance(ctx context.Context, q query.ListByImportanceQuery) (*query.PaginatedResult, error) {
	goals, cursor, err := s.repo.ListByImportance(ctx, q.UserID, domain.GoalImportance(q.Importance), q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginated(goals, cursor), nil
}

func (s *GoalService) ListByType(ctx context.Context, q query.ListByTypeQuery) (*query.PaginatedResult, error) {
	goals, cursor, err := s.repo.ListByType(ctx, q.UserID, domain.GoalType(q.GoalType), q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginated(goals, cursor), nil
}

func (s *GoalService) ListByHousehold(ctx context.Context, q query.ListByHouseholdQuery) (*query.PaginatedResult, error) {
	goals, cursor, err := s.repo.ListByHousehold(ctx, q.HouseholdID, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginated(goals, cursor), nil
}

func toGoalView(g *domain.Goal) *query.GoalView {
	v := &query.GoalView{
		GoalID: g.ID(), UserID: g.UserID(), Name: g.Name(),
		Importance: string(g.Importance()), GoalType: string(g.GoalType()),
		Subtype: string(g.Subtype()), Priority: g.Priority(),
		Progress: domain.CalculateProgress(g.CurrentValue(), g.SuccessCriteria()),
		Remaining: domain.CalculateRemaining(g.CurrentValue(), g.SuccessCriteria()),
		Status: string(g.Status()), HasTargetDate: g.TargetDate() != nil,
		CreatedAt: g.CreatedAt().Format(time.RFC3339),
	}
	if g.HouseholdID() != nil {
		v.HouseholdID = *g.HouseholdID()
	}
	return v
}

func toPaginated(goals []*domain.Goal, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{Goals: make([]query.GoalView, 0, len(goals)), NextCursor: cursor, HasMore: cursor != ""}
	for _, g := range goals {
		r.Goals = append(r.Goals, *toGoalView(g))
	}
	return r
}
