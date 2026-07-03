package application

import (
	"context"
	"errors"

	"github.com/horizon/core/packages/events"
	"github.com/horizon/core/services/domains/recurring/internal/application/dto/command"
	"github.com/horizon/core/services/domains/recurring/internal/application/dto/query"
	"github.com/horizon/core/services/domains/recurring/internal/domain"
)

var ErrNotFound = errors.New("recurring transaction not found")

type Service struct {
	repo      domain.Repository
	publisher events.Publisher
}

func NewService(repo domain.Repository, pub events.Publisher) *Service {
	return &Service{
		repo:      repo,
		publisher: pub,
	}
}

func (s *Service) CreateRecurring(ctx context.Context, cmd command.CreateRecurringCommand) error {
	rt, err := domain.NewRecurringTransaction(
		"", cmd.UserID, cmd.HouseholdID, cmd.Name, cmd.Description,
		cmd.Amount, cmd.Currency, cmd.Frequency, cmd.Interval,
		cmd.StartDate, cmd.EndDate, cmd.EventTemplate,
		cmd.SkipHolidays, cmd.SkipWeekends, cmd.Tags, cmd.Metadata,
	)
	if err != nil {
		return err
	}

	if err := s.repo.Save(ctx, rt); err != nil {
		return err
	}

	data := []byte(`{"recurring_id":"` + rt.ID() + `","user_id":"` + rt.UserID() + `"}`)
	_ = s.publisher.Publish(ctx, events.NewEnvelope(string(domain.EventRecurringCreated), 1, data))

	return nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*query.RecurringTransactionDTO, error) {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return toDTO(rt), nil
}

func (s *Service) ListByUser(ctx context.Context, userID string) ([]query.RecurringTransactionDTO, error) {
	rts, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var dtos []query.RecurringTransactionDTO
	for _, rt := range rts {
		dtos = append(dtos, *toDTO(rt))
	}
	return dtos, nil
}

func (s *Service) Activate(ctx context.Context, id string) error {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	if err := rt.Activate(); err != nil {
		return err
	}
	if err := s.repo.Save(ctx, rt); err != nil {
		return err
	}
	_ = s.publisher.Publish(ctx, events.NewEnvelope(string(domain.EventRecurringActivated), 1, []byte(`{"recurring_id":"`+id+`"}`)))
	return nil
}

func (s *Service) Pause(ctx context.Context, id string) error {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	if err := rt.Pause(); err != nil {
		return err
	}
	if err := s.repo.Save(ctx, rt); err != nil {
		return err
	}
	_ = s.publisher.Publish(ctx, events.NewEnvelope(string(domain.EventRecurringPaused), 1, []byte(`{"recurring_id":"`+id+`"}`)))
	return nil
}

func (s *Service) Cancel(ctx context.Context, id string) error {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	if err := rt.Cancel(); err != nil {
		return err
	}
	if err := s.repo.Save(ctx, rt); err != nil {
		return err
	}
	_ = s.publisher.Publish(ctx, events.NewEnvelope(string(domain.EventRecurringCancelled), 1, []byte(`{"recurring_id":"`+id+`"}`)))
	return nil
}

func (s *Service) Archive(ctx context.Context, id string) error {
	rt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	if err := rt.Archive(); err != nil {
		return err
	}
	if err := s.repo.Save(ctx, rt); err != nil {
		return err
	}
	_ = s.publisher.Publish(ctx, events.NewEnvelope(string(domain.EventRecurringArchived), 1, []byte(`{"recurring_id":"`+id+`"}`)))
	return nil
}

func toDTO(rt *domain.RecurringTransaction) *query.RecurringTransactionDTO {
	return &query.RecurringTransactionDTO{
		ID:             rt.ID(),
		UserID:         rt.UserID(),
		HouseholdID:    rt.HouseholdID(),
		Name:           rt.Name(),
		Description:    rt.Description(),
		Amount:         rt.Amount(),
		Currency:       rt.Currency(),
		Frequency:      string(rt.Frequency()),
		Interval:       rt.Interval(),
		StartDate:      rt.StartDate(),
		EndDate:        rt.EndDate(),
		NextOccurrence: rt.NextOccurrence(),
		LastOccurrence: rt.LastOccurrence(),
		Status:         string(rt.Status()),
		SkipHolidays:   rt.SkipHolidays(),
		SkipWeekends:   rt.SkipWeekends(),
		Tags:           rt.Tags(),
		Metadata:       rt.Metadata(),
		CreatedAt:      rt.CreatedAt(),
		UpdatedAt:      rt.UpdatedAt(),
	}
}
