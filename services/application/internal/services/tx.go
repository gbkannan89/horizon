package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager manages database transactions across repositories and outbox.
type TxManager struct {
	pool *pgxpool.Pool
}

// NewTxManager creates a new transaction manager.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithinTx executes a function within a database transaction.
func (tm *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback: %v (from: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

type txKey struct{}

// ServiceRegistration holds all dependencies for the application layer.
type ServiceRegistration struct {
	UserService            interface{}
	InstitutionService     interface{}
	FinancialEventService  interface{}
	GoalService            interface{}
	AccountService         interface{}
	AllocationService      interface{}
	AssetService           interface{}
	LiabilityService       interface{}
	PortfolioService       interface{}
	TxManager              *TxManager
}

// NewServiceRegistration creates a new service registration.
func NewServiceRegistration(
	tm *TxManager,
) *ServiceRegistration {
	return &ServiceRegistration{
		TxManager: tm,
	}
}
