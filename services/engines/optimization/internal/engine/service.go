package engine

import (
	"context"
	"fmt"
	"time"
)

// Service orchestrates optimization computation and persistence.
type Service struct {
	engine *Engine
	repo   Repository
}

// NewService creates a new optimization service.
func NewService(repo Repository) *Service {
	return &Service{engine: NewEngine(), repo: repo}
}

// RunOptimization executes the engine and persists the result.
func (s *Service) RunOptimization(ctx context.Context, inputs Inputs, objectives []Objective, constraints []Constraint) (*OptimizationOutput, error) {
	output := s.engine.Execute(inputs, objectives, constraints)
	output.OptimizationID = fmt.Sprintf("opt-%d", time.Now().UnixMilli())

	if err := s.repo.Save(ctx, output); err != nil {
		return nil, fmt.Errorf("save optimization: %w", err)
	}
	return output, nil
}

// GetOptimization retrieves an optimization by ID.
func (s *Service) GetOptimization(ctx context.Context, id string) (*OptimizationOutput, error) {
	return s.repo.GetByID(ctx, id)
}

// GetLatestOptimization retrieves the most recent optimization.
func (s *Service) GetLatestOptimization(ctx context.Context) (*OptimizationOutput, error) {
	return s.repo.GetLatest(ctx)
}

// Repository persists optimization outputs.
type Repository interface {
	Save(ctx context.Context, output *OptimizationOutput) error
	GetByID(ctx context.Context, id string) (*OptimizationOutput, error)
	GetLatest(ctx context.Context) (*OptimizationOutput, error)
}
