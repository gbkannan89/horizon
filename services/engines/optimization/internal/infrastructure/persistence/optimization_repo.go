package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/optimization/internal/engine"
)

type OptimizationRepository struct {
	pool *pgxpool.Pool
}

func NewOptimizationRepository(pool *pgxpool.Pool) *OptimizationRepository {
	return &OptimizationRepository{pool: pool}
}

func (r *OptimizationRepository) Save(ctx context.Context, output *engine.OptimizationOutput) error {
	data, err := json.Marshal(output)
	if err != nil { return fmt.Errorf("marshal: %w", err) }

	_, err = r.pool.Exec(ctx,
		`INSERT INTO optimizations (optimization_id, status, candidates_evaluated, candidates_valid, opt_data, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (optimization_id) DO UPDATE SET status=$2, opt_data=$5`,
		output.OptimizationID, string(output.Status), output.CandidatesEvaluated, output.CandidatesValid, data)
	return err
}

func (r *OptimizationRepository) GetByID(ctx context.Context, id string) (*engine.OptimizationOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT opt_data FROM optimizations WHERE optimization_id = $1`, id).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get optimization %s: %w", id, err) }

	var o engine.OptimizationOutput
	if err := json.Unmarshal(data, &o); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	return &o, nil
}

func (r *OptimizationRepository) GetLatest(ctx context.Context) (*engine.OptimizationOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT opt_data FROM optimizations ORDER BY created_at DESC LIMIT 1`).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get latest optimization: %w", err) }

	var o engine.OptimizationOutput
	if err := json.Unmarshal(data, &o); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	return &o, nil
}
