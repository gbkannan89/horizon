package engine

import (
	"context"
	"fmt"
	"time"
)

// Repository persists projection outputs.
type Repository interface {
	Save(ctx context.Context, output *ProjectionOutput) error
	GetByID(ctx context.Context, id string) (*ProjectionOutput, error)
	GetLatest(ctx context.Context, projType ProjectionType) (*ProjectionOutput, error)
	ListByUser(ctx context.Context, userID string, limit int) ([]*ProjectionOutput, error)
}

// Service orchestrates projection computation and persistence.
type Service struct {
	engine *Engine
	repo   Repository
}

// NewService creates a new projection service.
func NewService(assumptions Assumptions, repo Repository) *Service {
	return &Service{
		engine: NewEngine(assumptions),
		repo:   repo,
	}
}

// RunProjection executes a projection and persists the result.
func (s *Service) RunProjection(ctx context.Context, inputs Inputs, projType ProjectionType) (*ProjectionOutput, error) {
	output := s.engine.Execute(inputs, projType)
	output.OutputID = fmt.Sprintf("proj-%s-%d", projType, time.Now().UnixMilli())

	if err := s.repo.Save(ctx, output); err != nil {
		return nil, fmt.Errorf("save projection: %w", err)
	}
	return output, nil
}

// GetProjection retrieves a projection by ID.
func (s *Service) GetProjection(ctx context.Context, id string) (*ProjectionOutput, error) {
	return s.repo.GetByID(ctx, id)
}

// GetLatestProjection retrieves the latest projection of a given type.
func (s *Service) GetLatestProjection(ctx context.Context, projType ProjectionType) (*ProjectionOutput, error) {
	return s.repo.GetLatest(ctx, projType)
}
