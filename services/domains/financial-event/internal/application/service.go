package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/financial-event/internal/application/dto/command"
	"github.com/horizon/core/services/domains/financial-event/internal/application/dto/query"
	"github.com/horizon/core/services/domains/financial-event/internal/domain"
)

const duplicateThreshold = 0.7

type EventService struct {
	repo      domain.Repository
	publisher domain.Publisher
	factory   *domain.FinancialEventFactory
	now       func() time.Time
}

func NewEventService(repo domain.Repository, publisher domain.Publisher, now func() time.Time) *EventService {
	return &EventService{
		repo:      repo,
		publisher: publisher,
		factory:   domain.NewFinancialEventFactory(),
		now:       now,
	}
}

func (s *EventService) CreateDraft(ctx context.Context, cmd command.CreateDraftCommand) (*command.CreateDraftResult, error) {
	if cmd.UserID == "" {
		return nil, errors.New("user_id is required")
	}

	eventID := domain.NewEventID()
	createdBy := domain.CreatedBy(cmd.CreatedBy)
	if !domain.AllCreatedBy[createdBy] {
		return nil, fmt.Errorf("invalid created_by: %s", cmd.CreatedBy)
	}

	origin := domain.EventOrigin(cmd.Origin)
	if !domain.AllOrigins[origin] {
		return nil, fmt.Errorf("invalid origin: %s", cmd.Origin)
	}

	confidence := domain.EventConfidence(cmd.Confidence)
	event, err := s.factory.CreateDraft(
		eventID, cmd.UserID,
		domain.EventType(cmd.EventType),
		cmd.Amount, cmd.Currency,
		cmd.EventDate, cmd.EffectiveDate,
		cmd.Description,
		origin, confidence, createdBy,
		cmd.Source, cmd.Destination, cmd.Reference,
		cmd.Notes, cmd.HouseholdID, cmd.EventSubType, nil,
	)
	if err != nil {
		return nil, err
	}

	if cmd.ImportedFrom != "" {
		existing, _ := s.repo.FindByExternalReference(ctx, cmd.ImportedFrom, cmd.Reference)
		if existing != nil {
			return nil, fmt.Errorf("duplicate import: event %s already exists with same reference", existing.EventID())
		}
	}

	if err := s.repo.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}

	if err := s.publisher.Publish(domain.NewFinancialEventCreated(event)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}

	return &command.CreateDraftResult{
		EventID:   eventID,
		State:     string(event.State()),
	}, nil
}

func (s *EventService) Submit(ctx context.Context, cmd command.SubmitCommand) (*command.CommandResult, error) {
	event, err := s.repo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", cmd.EventID)
	}
	if err := domain.SubmitEvent(event); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}
	return &command.CommandResult{Success: true, EventID: event.EventID(), State: string(event.State())}, nil
}

func (s *EventService) Confirm(ctx context.Context, cmd command.ConfirmCommand) (*command.CommandResult, error) {
	event, err := s.repo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", cmd.EventID)
	}
	if err := domain.ConfirmEvent(event); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}
	if err := s.publisher.Publish(domain.NewFinancialEventConfirmed(event)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	seq, err := s.repo.GetNextSequenceNumber(ctx, event.UserID())
	if err == nil {
		event.SetSequenceNumber(seq)
		event.SetOrderIndex(domain.ComputeOrderIndex(event.EffectiveDate(), event.EventDate(), seq))
	}
	return &command.CommandResult{Success: true, EventID: event.EventID(), State: string(event.State())}, nil
}

func (s *EventService) Post(ctx context.Context, cmd command.PostCommand) (*command.CommandResult, error) {
	event, err := s.repo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", cmd.EventID)
	}
	if err := domain.PostEvent(event); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}
	if err := s.publisher.Publish(domain.NewFinancialEventPosted(event)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, EventID: event.EventID(), State: string(event.State())}, nil
}

func (s *EventService) Reverse(ctx context.Context, cmd command.ReverseCommand) (*command.CommandResult, error) {
	original, err := s.repo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", cmd.EventID)
	}
	if original.State() != domain.StateConfirmed && original.State() != domain.StatePosted {
		return nil, errors.New("can only reverse a Confirmed or Posted event")
	}

	reversalID := domain.NewEventID()
	eventDate := s.now()
	effectiveDate := s.now()

	reversal := domain.NewBaseEvent(
		reversalID, original.UserID(),
		original.EventType(), -original.Amount(), original.Currency(),
		eventDate, effectiveDate,
		fmt.Sprintf("Reversal of %s: %s", original.EventID(), cmd.Description),
		original.Origin(), original.Confidence(), domain.CreatedByUser, domain.StateDraft,
	)

	if err := domain.ReverseEvent(original, reversal); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, original); err != nil {
		return nil, fmt.Errorf("save original: %w", err)
	}
	if err := s.repo.Save(ctx, reversal); err != nil {
		return nil, fmt.Errorf("save reversal: %w", err)
	}
	if err := s.publisher.Publish(domain.NewFinancialEventReversed(original, reversal)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, EventID: reversalID, State: string(original.State())}, nil
}

func (s *EventService) Cancel(ctx context.Context, cmd command.CancelCommand) (*command.CommandResult, error) {
	event, err := s.repo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", cmd.EventID)
	}
	if err := domain.CancelEvent(event); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}
	if err := s.publisher.Publish(domain.NewFinancialEventCancelled(event, cmd.Reason)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, EventID: event.EventID(), State: string(event.State())}, nil
}

func (s *EventService) Archive(ctx context.Context, cmd command.ArchiveCommand) (*command.CommandResult, error) {
	event, err := s.repo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", cmd.EventID)
	}
	if err := domain.ArchiveEvent(event); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, event); err != nil {
		return nil, fmt.Errorf("save event: %w", err)
	}
	if err := s.publisher.Publish(domain.NewFinancialEventArchived(event)); err != nil {
		return nil, fmt.Errorf("publish event: %w", err)
	}
	return &command.CommandResult{Success: true, EventID: event.EventID(), State: string(event.State())}, nil
}

func (s *EventService) Import(ctx context.Context, cmd command.ImportCommand) (*command.CommandResult, error) {
	if cmd.ImportedFrom == "" {
		return nil, errors.New("imported_from is required for imported events")
	}

	draftCmd := command.CreateDraftCommand{
		UserID:        cmd.UserID,
		EventType:     cmd.EventType,
		Amount:        cmd.Amount,
		Currency:      cmd.Currency,
		EventDate:     cmd.EventDate,
		EffectiveDate: cmd.EffectiveDate,
		Source:        cmd.Source,
		Destination:   cmd.Destination,
		Description:   cmd.Description,
		Reference:     cmd.Reference,
		Origin:        string(domain.OriginImport),
		Confidence:    string(domain.ConfidenceImported),
		CreatedBy:     string(domain.CreatedByImport),
		ImportedFrom:  cmd.ImportedFrom,
	}

	result, err := s.CreateDraft(ctx, draftCmd)
	if err != nil {
		return nil, fmt.Errorf("create import: %w", err)
	}

	if _, err := s.Submit(ctx, command.SubmitCommand{EventID: result.EventID}); err != nil {
		return nil, fmt.Errorf("submit import: %w", err)
	}
	if _, err := s.Confirm(ctx, command.ConfirmCommand{EventID: result.EventID}); err != nil {
		return nil, fmt.Errorf("confirm import: %w", err)
	}

	if err := s.publisher.Publish(domain.NewFinancialEventImported(
		domain.NewBaseEvent(result.EventID, cmd.UserID, domain.EventType(cmd.EventType),
			cmd.Amount, cmd.Currency, cmd.EventDate, cmd.EffectiveDate, cmd.Description,
			domain.OriginImport, domain.ConfidenceImported, domain.CreatedByImport, domain.StateConfirmed,
		),
	)); err != nil {
		return nil, fmt.Errorf("publish imported: %w", err)
	}

	return &command.CommandResult{Success: true, EventID: result.EventID, State: string(domain.StateConfirmed)}, nil
}

func (s *EventService) DetectDuplicate(ctx context.Context, cmd command.DetectDuplicateCommand) (*command.DetectDuplicateResult, error) {
	candidate, err := s.repo.GetByID(ctx, cmd.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", cmd.EventID)
	}
	if candidate.State() != domain.StatePending {
		return nil, errors.New("duplicate detection only applies to pending events")
	}

	existing, err := s.repo.FindPotentialDuplicates(ctx,
		candidate.Amount(), candidate.Currency(),
		candidate.Source(), candidate.Destination(),
		candidate.EffectiveDate(), 3*24*time.Hour,
	)
	if err != nil {
		return nil, fmt.Errorf("find duplicates: %w", err)
	}

	score, matchedID := domain.ComputeDuplicateScore(candidate, existing)
	result := &command.DetectDuplicateResult{
		IsDuplicate: score >= duplicateThreshold,
		MatchedID:   matchedID,
		Confidence:  score,
	}
	if result.IsDuplicate {
		s.publisher.Publish(domain.NewFinancialEventDuplicateDetected(cmd.EventID, matchedID, score))
	}
	return result, nil
}

func (s *EventService) GetEvent(ctx context.Context, q query.GetEventQuery) (*query.EventResult, error) {
	event, err := s.repo.GetByID(ctx, q.EventID)
	if err != nil {
		return nil, fmt.Errorf("event not found: %s", q.EventID)
	}
	return toEventResult(event), nil
}

func (s *EventService) ListByUser(ctx context.Context, q query.ListEventsByUserQuery) (*query.PaginatedResult, error) {
	events, cursor, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginatedResult(events, cursor), nil
}

func (s *EventService) ListByAccount(ctx context.Context, q query.ListEventsByAccountQuery) (*query.PaginatedResult, error) {
	events, cursor, err := s.repo.ListByAccount(ctx, q.AccountID, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginatedResult(events, cursor), nil
}

func (s *EventService) ListByDateRange(ctx context.Context, q query.ListEventsByDateRangeQuery) (*query.PaginatedResult, error) {
	events, cursor, err := s.repo.ListByDateRange(ctx, q.UserID, q.Start, q.End, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginatedResult(events, cursor), nil
}

func (s *EventService) ListByType(ctx context.Context, q query.ListEventsByTypeQuery) (*query.PaginatedResult, error) {
	events, cursor, err := s.repo.ListByType(ctx, q.UserID, domain.EventType(q.EventType), q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginatedResult(events, cursor), nil
}

func (s *EventService) GetTimeline(ctx context.Context, q query.GetTimelineQuery) (*query.PaginatedResult, error) {
	events, cursor, err := s.repo.GetTimeline(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil {
		return nil, err
	}
	return toPaginatedResult(events, cursor), nil
}

func toEventResult(e *domain.FinancialEvent) *query.EventResult {
	return &query.EventResult{
		EventID:       e.EventID(),
		UserID:        e.UserID(),
		EventType:     string(e.EventType()),
		Amount:        e.Amount(),
		Currency:      e.Currency(),
		EventDate:     e.EventDate(),
		EffectiveDate: e.EffectiveDate(),
		Description:   e.Description(),
		State:         string(e.State()),
		Origin:        string(e.Origin()),
		Confidence:    string(e.Confidence()),
		CreatedBy:     string(e.CreatedBy()),
		Source:        e.Source(),
		Destination:   e.Destination(),
		Reference:     e.Reference(),
		ImportedFrom:  e.ImportedFrom(),
		ReversalOfID:  e.ReversalOfEventID(),
		CorrelationID: e.CorrelationID(),
		OrderIndex:    e.OrderIndex(),
		CreatedAt:     e.CreatedAt(),
		UpdatedAt:     e.UpdatedAt(),
	}
}

func toPaginatedResult(events []*domain.FinancialEvent, cursor string) *query.PaginatedResult {
	result := &query.PaginatedResult{
		Events:     make([]query.EventResult, 0, len(events)),
		NextCursor: cursor,
		HasMore:    cursor != "",
	}
	for _, e := range events {
		result.Events = append(result.Events, *toEventResult(e))
	}
	return result
}
