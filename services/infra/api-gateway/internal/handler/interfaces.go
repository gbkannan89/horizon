package handler

import (
	"context"
	"time"
)

// --- User Service Interface ---
type UserService interface {
	Create(ctx context.Context, name, email string) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context) (interface{}, error)
	Activate(ctx context.Context, id string) error
	Suspend(ctx context.Context, id string) error
	Archive(ctx context.Context, id string) error
}

// --- Goal Service Interface ---
type GoalService interface {
	Create(ctx context.Context, userID, name string, targetAmount float64) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context, userID string) (interface{}, error)
	Activate(ctx context.Context, id string) error
	Complete(ctx context.Context, id string) error
	Archive(ctx context.Context, id string) error
}

// --- Account Service Interface ---
type AccountService interface {
	Create(ctx context.Context, userID, name, acctType, currency string) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context, userID string) (interface{}, error)
	Activate(ctx context.Context, id string) error
	Freeze(ctx context.Context, id string) error
	Close(ctx context.Context, id string) error
}

// --- Event Service Interface ---
type EventService interface {
	Create(ctx context.Context, userID, eventType string, amount float64, currency string) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context, userID string) (interface{}, error)
	Confirm(ctx context.Context, id string) error
	Post(ctx context.Context, id string) error
	Reverse(ctx context.Context, id string, reason string) error
	Archive(ctx context.Context, id string) error
}

// --- Institution Service Interface ---
type InstitutionService interface {
	Create(ctx context.Context, name, instType, country string) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context) (interface{}, error)
}

// --- Allocation Service Interface ---
type AllocationService interface {
	Create(ctx context.Context, goalID, sourceID, allocType, currency string, amount float64) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context) (interface{}, error)
}

// --- Asset Service Interface ---
type AssetService interface {
	Create(ctx context.Context, ownerID, name, cls, currency string, value float64) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context, ownerID string) (interface{}, error)
}

// --- Liability Service Interface ---
type LiabilityService interface {
	Create(ctx context.Context, ownerID, name, cls, currency string, principal float64, rate float64, maturityDate time.Time) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context, ownerID string) (interface{}, error)
	Settle(ctx context.Context, id string) error
}

// --- Portfolio Service Interface ---
type PortfolioService interface {
	Create(ctx context.Context, ownerID, name, pfType, currency string) (string, error)
	Get(ctx context.Context, id string) (interface{}, error)
	List(ctx context.Context, ownerID string) (interface{}, error)
}
