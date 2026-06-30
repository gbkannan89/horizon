package engine

import (
	"context"
	"fmt"
	"time"
)

// Service orchestrates risk assessment computation and persistence.
type Service struct {
	engine *Engine
	repo   Repository
}

// NewService creates a new risk assessment service.
func NewService(repo Repository) *Service {
	return &Service{engine: NewEngine(), repo: repo}
}

// RunAssessment executes a risk assessment and persists the result.
func (s *Service) RunAssessment(ctx context.Context, inputs Inputs) (*RiskOutput, error) {
	output := s.engine.Execute(inputs)
	output.AssessmentID = fmt.Sprintf("risk-%d", time.Now().UnixMilli())

	if err := s.repo.Save(ctx, output); err != nil {
		return nil, fmt.Errorf("save risk: %w", err)
	}
	return output, nil
}

// GetAssessment retrieves a risk assessment by ID.
func (s *Service) GetAssessment(ctx context.Context, id string) (*RiskOutput, error) {
	return s.repo.GetByID(ctx, id)
}

// GetLatestAssessment retrieves the most recent assessment.
func (s *Service) GetLatestAssessment(ctx context.Context) (*RiskOutput, error) {
	return s.repo.GetLatest(ctx)
}

// Repository persists risk assessment outputs.
type Repository interface {
	Save(ctx context.Context, output *RiskOutput) error
	GetByID(ctx context.Context, id string) (*RiskOutput, error)
	GetLatest(ctx context.Context) (*RiskOutput, error)
}
