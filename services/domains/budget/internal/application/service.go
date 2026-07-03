package application

import (
	"context"
	"fmt"
	"time"

	"github.com/horizon/core/services/domains/budget/internal/application/dto/command"
	"github.com/horizon/core/services/domains/budget/internal/application/dto/query"
	"github.com/horizon/core/services/domains/budget/internal/domain"
)

type BudgetService struct {
	repo      domain.Repository
	publisher domain.Publisher
	factory   *domain.BudgetFactory
	now       func() time.Time
}

func NewBudgetService(repo domain.Repository, publisher domain.Publisher, now func() time.Time) *BudgetService {
	return &BudgetService{
		repo:      repo,
		publisher: publisher,
		factory:   domain.NewBudgetFactory(),
		now:       now,
	}
}

func (s *BudgetService) Create(ctx context.Context, cmd command.CreateBudgetCommand) (*command.BudgetResult, error) {
	startDate, _ := time.Parse("2006-01-02", cmd.StartDate)
	endDate, _ := time.Parse("2006-01-02", cmd.EndDate)
	if endDate.IsZero() {
		endDate = startDate.AddDate(0, 1, 0)
	}
	if startDate.IsZero() {
		startDate = s.now().UTC()
	}

	currency := cmd.Currency
	if currency == "" {
		currency = "INR"
	}

	var categories []domain.BudgetCategory
	for _, c := range cmd.Categories {
		categories = append(categories, domain.BudgetCategory{
			Category:       c.Category,
			Subcategory:    c.Subcategory,
			BudgetedAmount: c.BudgetedAmount,
			Rollover:       c.Rollover,
		})
	}

	b, err := s.factory.Create(
		"", cmd.UserID, cmd.HouseholdID, cmd.Name,
		domain.BudgetPeriod(cmd.Period),
		startDate, endDate,
		currency, categories, cmd.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("factory: %w", err)
	}

	if err := s.repo.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	if err := s.publisher.Publish(domain.NewBudgetCreated(b)); err != nil {
		return nil, fmt.Errorf("publish: %w", err)
	}

	return &command.BudgetResult{
		BudgetID: b.ID(),
		Status:   string(b.Status()),
		Success:  true,
	}, nil
}

func (s *BudgetService) Activate(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
	b, err := s.repo.GetByID(ctx, cmd.BudgetID)
	if err != nil { return nil, err }
	if err := domain.ActivateBudget(b); err != nil { return nil, err }
	if err := s.repo.Save(ctx, b); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewBudgetActivated(b)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.BudgetResult{BudgetID: b.ID(), Status: string(b.Status()), Success: true}, nil
}

func (s *BudgetService) Pause(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
	b, err := s.repo.GetByID(ctx, cmd.BudgetID)
	if err != nil { return nil, err }
	if err := domain.PauseBudget(b); err != nil { return nil, err }
	if err := s.repo.Save(ctx, b); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewBudgetPaused(b)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.BudgetResult{BudgetID: b.ID(), Status: string(b.Status()), Success: true}, nil
}

func (s *BudgetService) Resume(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
	b, err := s.repo.GetByID(ctx, cmd.BudgetID)
	if err != nil { return nil, err }
	if err := domain.ResumeBudget(b); err != nil { return nil, err }
	if err := s.repo.Save(ctx, b); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewBudgetResumed(b)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.BudgetResult{BudgetID: b.ID(), Status: string(b.Status()), Success: true}, nil
}

func (s *BudgetService) Complete(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
	b, err := s.repo.GetByID(ctx, cmd.BudgetID)
	if err != nil { return nil, err }
	if err := domain.CompleteBudget(b); err != nil { return nil, err }
	if err := s.repo.Save(ctx, b); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewBudgetCompleted(b)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.BudgetResult{BudgetID: b.ID(), Status: string(b.Status()), Success: true}, nil
}

func (s *BudgetService) Archive(ctx context.Context, cmd command.BudgetIDCommand) (*command.BudgetResult, error) {
	b, err := s.repo.GetByID(ctx, cmd.BudgetID)
	if err != nil { return nil, err }
	if err := domain.ArchiveBudget(b); err != nil { return nil, err }
	if err := s.repo.Save(ctx, b); err != nil { return nil, fmt.Errorf("save: %w", err) }
	if err := s.publisher.Publish(domain.NewBudgetArchived(b)); err != nil { return nil, fmt.Errorf("publish: %w", err) }
	return &command.BudgetResult{BudgetID: b.ID(), Status: string(b.Status()), Success: true}, nil
}

func (s *BudgetService) UpdateCategory(ctx context.Context, cmd command.UpdateCategoryCommand) (*command.BudgetResult, error) {
	b, err := s.repo.GetByID(ctx, cmd.BudgetID)
	if err != nil { return nil, err }
	updated := false
	for i, c := range b.Categories() {
		if c.ID == cmd.CategoryID {
			if cmd.BudgetedAmount != nil {
				b.Categories()[i].BudgetedAmount = *cmd.BudgetedAmount
			}
			if cmd.SpentAmount != nil {
				b.Categories()[i].SpentAmount = *cmd.SpentAmount
			}
			cat := b.Categories()[i]
			cat.RemainingAmount = domain.CalculateCategoryRemaining(cat.BudgetedAmount, cat.SpentAmount, cat.Rollover, 0)
			b.Categories()[i].RemainingAmount = cat.RemainingAmount
			updated = true
			break
		}
	}
	if !updated { return nil, fmt.Errorf("category not found") }

	// Recalculate totals
	var totalBudgeted, totalSpent int64
	for _, c := range b.Categories() {
		totalBudgeted += c.BudgetedAmount
		totalSpent += c.SpentAmount
	}
	b.Categories() //nolint

	if err := s.repo.Save(ctx, b); err != nil { return nil, fmt.Errorf("save: %w", err) }
	return &command.BudgetResult{BudgetID: b.ID(), Status: string(b.Status()), Success: true}, nil
}

func (s *BudgetService) GetByID(ctx context.Context, q query.GetBudgetQuery) (*query.BudgetView, error) {
	b, err := s.repo.GetByID(ctx, q.BudgetID)
	if err != nil { return nil, err }
	return toBudgetView(b), nil
}

func (s *BudgetService) ListByUser(ctx context.Context, q query.ListByUserQuery) (*query.PaginatedResult, error) {
	budgets, cursor, err := s.repo.ListByUser(ctx, q.UserID, q.Cursor, q.Limit)
	if err != nil { return nil, err }
	return toPaginated(budgets, cursor), nil
}

func (s *BudgetService) ListByPeriod(ctx context.Context, q query.ListByPeriodQuery) (*query.PaginatedResult, error) {
	start, _ := time.Parse("2006-01-02", q.StartDate)
	end, _ := time.Parse("2006-01-02", q.EndDate)
	budgets, cursor, err := s.repo.ListByPeriod(ctx, q.UserID, start, end, q.Cursor, q.Limit)
	if err != nil { return nil, err }
	return toPaginated(budgets, cursor), nil
}

func (s *BudgetService) ListByHousehold(ctx context.Context, q query.ListByHouseholdQuery) (*query.PaginatedResult, error) {
	budgets, cursor, err := s.repo.ListByHousehold(ctx, q.HouseholdID, q.Cursor, q.Limit)
	if err != nil { return nil, err }
	return toPaginated(budgets, cursor), nil
}

func toBudgetView(b *domain.Budget) *query.BudgetView {
	categories := make([]query.CategoryView, len(b.Categories()))
	for i, c := range b.Categories() {
		spentPct := domain.CalculateSpentPercentage(c.SpentAmount, c.BudgetedAmount)
		categories[i] = query.CategoryView{
			ID: c.ID, Category: c.Category, Subcategory: c.Subcategory,
			BudgetedAmount: c.BudgetedAmount, SpentAmount: c.SpentAmount,
			RemainingAmount: c.RemainingAmount, Rollover: c.Rollover,
			SpentPct: spentPct,
		}
	}
	return &query.BudgetView{
		BudgetID: b.ID(), HouseholdID: b.HouseholdID(), Name: b.Name(), Period: string(b.Period()),
		StartDate: b.StartDate().Format("2006-01-02"),
		EndDate:   b.EndDate().Format("2006-01-02"),
		Status: string(b.Status()),
		TotalBudgeted: b.TotalBudgeted(), TotalSpent: b.TotalSpent(),
		TotalRemaining: b.TotalRemaining(), Currency: b.Currency(),
		Categories: categories, Tags: b.Tags(),
		CreatedAt: b.CreatedAt().Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt().Format(time.RFC3339),
	}
}

func toPaginated(budgets []*domain.Budget, cursor string) *query.PaginatedResult {
	r := &query.PaginatedResult{
		Budgets: make([]query.BudgetView, 0, len(budgets)),
		NextCursor: cursor, HasMore: cursor != "",
	}
	for _, b := range budgets {
		r.Budgets = append(r.Budgets, *toBudgetView(b))
	}
	return r
}
