package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/risk/internal/engine"
)

type RiskRepository struct {
	pool *pgxpool.Pool
}

func NewRiskRepository(pool *pgxpool.Pool) *RiskRepository {
	return &RiskRepository{pool: pool}
}

func (r *RiskRepository) Save(ctx context.Context, output *engine.RiskOutput) error {
	data, err := json.Marshal(output)
	if err != nil { return fmt.Errorf("marshal: %w", err) }

	_, err = r.pool.Exec(ctx,
		`INSERT INTO risk_assessments (assessment_id, composite_score, risk_level, output_data, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (assessment_id) DO UPDATE SET composite_score=$2, risk_level=$3, output_data=$4`,
		output.AssessmentID, output.CompositeScore, string(output.CompositeLevel), data)
	return err
}

func (r *RiskRepository) GetByID(ctx context.Context, id string) (*engine.RiskOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT output_data FROM risk_assessments WHERE assessment_id = $1`, id).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get risk %s: %w", id, err) }

	var o engine.RiskOutput
	if err := json.Unmarshal(data, &o); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	return &o, nil
}

func (r *RiskRepository) GetLatest(ctx context.Context) (*engine.RiskOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT output_data FROM risk_assessments ORDER BY created_at DESC LIMIT 1`).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get latest risk: %w", err) }

	var o engine.RiskOutput
	if err := json.Unmarshal(data, &o); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	return &o, nil
}
