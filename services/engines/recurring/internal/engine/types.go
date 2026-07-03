package engine

import (
	"context"
	"time"
)

type RecurringTransaction struct {
	ID             string
	UserID         string
	NextOccurrence *time.Time
	Amount         float64
	Currency       string
	Description    string
	EventType      string
	Category       string
}

type RecurringProvider interface {
	GetDueByDate(ctx context.Context, date time.Time) ([]*RecurringTransaction, error)
	Save(ctx context.Context, r *RecurringTransaction) error
}

type FinancialEventCreator interface {
	CreateEvent(ctx context.Context, r *RecurringTransaction, occurrenceDate time.Time) error
}
