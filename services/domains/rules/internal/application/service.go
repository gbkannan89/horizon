package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/rules/internal/application/dto/command"
	"github.com/horizon/core/services/domains/rules/internal/application/dto/query"
	"github.com/horizon/core/services/domains/rules/internal/domain"
)

type RuleService struct {
	repo domain.Repository
	now  func() time.Time
}

func NewRuleService(repo domain.Repository, now func() time.Time) *RuleService {
	return &RuleService{repo: repo, now: now}
}

func (s *RuleService) Create(ctx context.Context, cmd command.CreateRuleCommand) (*command.RuleResult, error) {
	created := s.now().UTC().Format(time.RFC3339)
	rule, err := domain.NewRule("", cmd.UserID, cmd.HouseholdID, cmd.Name, cmd.Description,
		cmd.Category, cmd.Priority, cmd.Enabled, cmd.Conditions, cmd.Actions, created, created)
	if err != nil {
		return nil, fmt.Errorf("new rule: %w", err)
	}
	if err := domain.ValidateRule(rule); err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, rule); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	return &command.RuleResult{RuleID: rule.ID(), Success: true}, nil
}

func (s *RuleService) GetByID(ctx context.Context, id string) (*query.RuleView, error) {
	rule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toView(rule), nil
}

func (s *RuleService) ListByUser(ctx context.Context, userID string) ([]query.RuleView, error) {
	rules, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	views := make([]query.RuleView, 0, len(rules))
	for _, r := range rules {
		views = append(views, *toView(r))
	}
	return views, nil
}

func (s *RuleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func toView(r *domain.Rule) *query.RuleView {
	return &query.RuleView{
		RuleID: r.ID(), UserID: r.UserID(), HouseholdID: r.HouseholdID(),
		Name: r.Name(), Description: r.Description(), Category: r.Category(),
		Priority: r.Priority(), Enabled: r.Enabled(),
		Conditions: r.Conditions(), Actions: r.Actions(),
		CreatedAt: r.CreatedAt(), UpdatedAt: r.UpdatedAt(),
	}
}
