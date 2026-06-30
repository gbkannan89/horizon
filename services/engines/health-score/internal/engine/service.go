package engine

import (
	"context"
	"fmt"
	"time"
)

// Service orchestrates health score computation and persistence.
type Service struct {
	engine *Engine
	repo   Repository
}

// NewService creates a new health score service.
func NewService(repo Repository) *Service {
	return &Service{engine: NewEngine(), repo: repo}
}

// RunScore computes a health score and persists the result.
func (s *Service) RunScore(ctx context.Context, inputs Inputs, previousScore int) (*HealthScoreOutput, error) {
	inputs.PreviousOverallScore = previousScore
	output := s.engine.Execute(inputs)
	output.ScoreID = fmt.Sprintf("hs-%d", time.Now().UnixMilli())

	if err := s.repo.Save(ctx, output); err != nil {
		return nil, fmt.Errorf("save health score: %w", err)
	}
	return output, nil
}

// GetScore retrieves a health score by ID.
func (s *Service) GetScore(ctx context.Context, id string) (*HealthScoreOutput, error) {
	return s.repo.GetByID(ctx, id)
}

// GetLatestScore retrieves the most recent health score.
func (s *Service) GetLatestScore(ctx context.Context) (*HealthScoreOutput, error) {
	return s.repo.GetLatest(ctx)
}

// Repository persists health score outputs.
type Repository interface {
	Save(ctx context.Context, output *HealthScoreOutput) error
	GetByID(ctx context.Context, id string) (*HealthScoreOutput, error)
	GetLatest(ctx context.Context) (*HealthScoreOutput, error)
}
