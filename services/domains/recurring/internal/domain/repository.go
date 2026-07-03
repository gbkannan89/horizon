package domain

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, r *RecurringTransaction) error
	GetByID(ctx context.Context, id string) (*RecurringTransaction, error)
	ListByUser(ctx context.Context, userID string) ([]*RecurringTransaction, error)
	ListByStatus(ctx context.Context, status RecurringStatus) ([]*RecurringTransaction, error)
	GetDueByDate(ctx context.Context, date time.Time) ([]*RecurringTransaction, error)
	GetActiveRecurring(ctx context.Context, userID string) ([]*RecurringTransaction, error)
}
