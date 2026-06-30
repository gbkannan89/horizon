package domain

import "context"

type Repository interface {
	Save(ctx context.Context, goal *Goal) error
	UpdateStatus(ctx context.Context, goalID string, fromStatus, toStatus GoalStatus) error
	UpdatePriority(ctx context.Context, goalID string, newPriority int) error

	GetByID(ctx context.Context, goalID string) (*Goal, error)
	ListByUser(ctx context.Context, userID string, cursor string, limit int) ([]*Goal, string, error)
	ListByStatus(ctx context.Context, userID string, status GoalStatus, cursor string, limit int) ([]*Goal, string, error)
	ListByImportance(ctx context.Context, userID string, importance GoalImportance, cursor string, limit int) ([]*Goal, string, error)
	ListByType(ctx context.Context, userID string, goalType GoalType, cursor string, limit int) ([]*Goal, string, error)
	GetActiveByPriority(ctx context.Context, userID string) ([]*Goal, error)

	GetMaxPriority(ctx context.Context, userID string) (int, error)
	FindConflict(ctx context.Context, userID string, priority int, excludeGoalID string) (*Goal, error)
	RebalancePriorities(ctx context.Context, userID string, afterPriority int) error
}
