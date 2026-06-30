package engine

import (
	"context"
	"fmt"
	"time"
)

// Service orchestrates recommendation generation and lifecycle management.
type Service struct {
	engine *Engine
	repo   Repository
}

// NewService creates a new recommendation service.
func NewService(cfg Config, repo Repository) *Service {
	return &Service{engine: NewEngine(cfg), repo: repo}
}

// GenerateRecommendations runs the engine and persists results.
func (s *Service) GenerateRecommendations(ctx context.Context, inputs Inputs) ([]Recommendation, error) {
	recs, err := s.engine.Execute(inputs)
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	for i := range recs {
		if recs[i].ID == "" {
			recs[i].ID = fmt.Sprintf("rec-%d-%d", time.Now().UnixMilli(), i)
		}
		recs[i].Status = RSGenerated
		if err := s.repo.Save(ctx, &recs[i]); err != nil {
			return nil, fmt.Errorf("save rec %s: %w", recs[i].ID, err)
		}
	}
	return recs, nil
}

// AcceptRecommendation marks a recommendation as accepted.
func (s *Service) AcceptRecommendation(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, RSAccepted)
}

// DismissRecommendation marks a recommendation as dismissed.
func (s *Service) DismissRecommendation(ctx context.Context, id string) error {
	return s.repo.UpdateStatus(ctx, id, RSRejected)
}

// GetRecommendations retrieves active recommendations.
func (s *Service) GetRecommendations(ctx context.Context) ([]Recommendation, error) {
	return s.repo.GetActive(ctx)
}

// Repository persists recommendations.
type Repository interface {
	Save(ctx context.Context, rec *Recommendation) error
	UpdateStatus(ctx context.Context, id string, status RecStatus) error
	GetActive(ctx context.Context) ([]Recommendation, error)
}
