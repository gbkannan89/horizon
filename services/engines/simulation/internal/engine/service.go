package engine

import (
	"context"
	"fmt"
	"time"
)

// Service orchestrates simulation execution and persistence.
type Service struct {
	engine *Engine
	repo   Repository
}

// NewService creates a new simulation service.
func NewService(repo Repository) *Service {
	return &Service{engine: NewEngine(), repo: repo}
}

// RunScenario executes a simulation and persists the result.
func (s *Service) RunScenario(ctx context.Context, baseline Inputs, simType SimType, params map[string]interface{}) (*SimulationOutput, error) {
	output, err := s.engine.Execute(baseline, simType, params)
	if err != nil {
		return nil, fmt.Errorf("simulation: %w", err)
	}
	output.ScenarioID = fmt.Sprintf("sim-%d", time.Now().UnixMilli())

	if err := s.repo.Save(ctx, output); err != nil {
		return nil, fmt.Errorf("save simulation: %w", err)
	}
	return output, nil
}

// CompareScenarios compares two simulations side by side.
func (s *Service) CompareScenarios(ctx context.Context, idA, idB string) (*ComparisonResult, error) {
	a, err := s.repo.GetByID(ctx, idA)
	if err != nil { return nil, fmt.Errorf("get %s: %w", idA, err) }
	b, err := s.repo.GetByID(ctx, idB)
	if err != nil { return nil, fmt.Errorf("get %s: %w", idB, err) }

	return &ComparisonResult{
		ScenarioA: a,
		ScenarioB: b,
		Deltas:    computeCrossDeltas(a, b),
	}, nil
}

// GetScenario retrieves a simulation by ID.
func (s *Service) GetScenario(ctx context.Context, id string) (*SimulationOutput, error) {
	return s.repo.GetByID(ctx, id)
}

// Repository persists simulation outputs.
type Repository interface {
	Save(ctx context.Context, output *SimulationOutput) error
	GetByID(ctx context.Context, id string) (*SimulationOutput, error)
	GetLatest(ctx context.Context) ([]*SimulationOutput, error)
}

// ComparisonResult holds a side-by-side comparison.
type ComparisonResult struct {
	ScenarioA *SimulationOutput `json:"scenario_a"`
	ScenarioB *SimulationOutput `json:"scenario_b"`
	Deltas    []DeltaEntry      `json:"deltas"`
}

func computeCrossDeltas(a, b *SimulationOutput) []DeltaEntry {
	return []DeltaEntry{
		{Metric: "Net Worth", BaselineValue: a.DeltaSummary[0].ScenarioValue, ScenarioValue: b.DeltaSummary[0].ScenarioValue},
		{Metric: "Monthly Income", BaselineValue: a.DeltaSummary[1].ScenarioValue, ScenarioValue: b.DeltaSummary[1].ScenarioValue},
	}
}
