package domain

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, event *FinancialEvent) error
	UpdateState(ctx context.Context, eventID string, fromState, toState EventState) error

	GetByID(ctx context.Context, eventID string) (*FinancialEvent, error)
	ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*FinancialEvent, string, error)
	ListByAccount(ctx context.Context, accountID string, cursor string, limit int) ([]*FinancialEvent, string, error)
	ListByDateRange(ctx context.Context, userID string, start, end time.Time, cursor string, limit int) ([]*FinancialEvent, string, error)
	ListByType(ctx context.Context, userID string, eventType EventType, cursor string, limit int) ([]*FinancialEvent, string, error)
	GetTimeline(ctx context.Context, userID string, cursor string, limit int) ([]*FinancialEvent, string, error)

	FindPotentialDuplicates(ctx context.Context, amount int64, currency string, source, destination string, effectiveDate time.Time, window time.Duration) ([]*FinancialEvent, error)

	GetReversalChain(ctx context.Context, eventID string) ([]*FinancialEvent, error)
	GetOriginalEvent(ctx context.Context, reversalEventID string) (*FinancialEvent, error)

	FindByExternalReference(ctx context.Context, sourceOfTruth string, reference string) (*FinancialEvent, error)
	GetNextSequenceNumber(ctx context.Context, userID string) (int64, error)
}
