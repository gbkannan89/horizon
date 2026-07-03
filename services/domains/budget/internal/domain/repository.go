package domain

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, budget *Budget) error
	UpdateStatus(ctx context.Context, id string, from, to BudgetStatus) error
	GetByID(ctx context.Context, id string) (*Budget, error)
	ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*Budget, string, error)
	ListByPeriod(ctx context.Context, userID string, start, end time.Time, cursor string, limit int) ([]*Budget, string, error)
	GetByCategory(ctx context.Context, userID, category string, cursor string, limit int) ([]*Budget, string, error)
	GetBudgetVsActual(ctx context.Context, budgetID string) ([]BudgetCategory, error)
	ListByHousehold(ctx context.Context, householdID string, cursor string, limit int) ([]*Budget, string, error)
}
