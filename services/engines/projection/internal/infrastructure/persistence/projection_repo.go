package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/projection/internal/engine"
)

type ProjectionRepository struct {
	pool *pgxpool.Pool
}

func NewProjectionRepository(pool *pgxpool.Pool) *ProjectionRepository {
	return &ProjectionRepository{pool: pool}
}

func (r *ProjectionRepository) Save(ctx context.Context, output *engine.ProjectionOutput) error {
	data, err := json.Marshal(output)
	if err != nil { return fmt.Errorf("marshal: %w", err) }

	_, err = r.pool.Exec(ctx,
		`INSERT INTO projection_outputs (output_id, projection_type, version, status, output_data, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (output_id) DO UPDATE SET status=$4, output_data=$5`,
		output.OutputID, string(output.ProjectionType), output.Version, string(output.Status), data)
	return err
}

func (r *ProjectionRepository) GetByID(ctx context.Context, id string) (*engine.ProjectionOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT output_data FROM projection_outputs WHERE output_id = $1`, id).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get projection %s: %w", id, err) }

	var output engine.ProjectionOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &output, nil
}

func (r *ProjectionRepository) GetLatest(ctx context.Context, projType engine.ProjectionType) (*engine.ProjectionOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT output_data FROM projection_outputs
		WHERE projection_type = $1 AND status = 'Completed'
		ORDER BY created_at DESC LIMIT 1`, string(projType)).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get latest %s: %w", string(projType), err) }

	var output engine.ProjectionOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &output, nil
}

func (r *ProjectionRepository) ListByUser(ctx context.Context, userID string, limit int) ([]*engine.ProjectionOutput, error) {
	if limit <= 0 { limit = 10 }
	rows, err := r.pool.Query(ctx,
		`SELECT output_data FROM projection_outputs ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil { return nil, err }
	defer rows.Close()

	var outputs []*engine.ProjectionOutput
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil { return nil, err }
		var o engine.ProjectionOutput
		if err := json.Unmarshal(data, &o); err != nil { return nil, err }
		outputs = append(outputs, &o)
	}
	return outputs, nil
}
