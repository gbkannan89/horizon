package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/health-score/internal/engine"
)

type HealthScoreRepository struct {
	pool *pgxpool.Pool
}

func NewHealthScoreRepository(pool *pgxpool.Pool) *HealthScoreRepository {
	return &HealthScoreRepository{pool: pool}
}

func (r *HealthScoreRepository) Save(ctx context.Context, output *engine.HealthScoreOutput) error {
	data, err := json.Marshal(output)
	if err != nil { return fmt.Errorf("marshal: %w", err) }

	_, err = r.pool.Exec(ctx,
		`INSERT INTO health_scores (score_id, overall_score, score_grade, score_data, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (score_id) DO UPDATE SET overall_score=$2, score_grade=$3, score_data=$4`,
		output.ScoreID, output.OverallScore, string(output.ScoreGrade), data)
	return err
}

func (r *HealthScoreRepository) GetByID(ctx context.Context, id string) (*engine.HealthScoreOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT score_data FROM health_scores WHERE score_id = $1`, id).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get health score %s: %w", id, err) }

	var o engine.HealthScoreOutput
	if err := json.Unmarshal(data, &o); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	return &o, nil
}

func (r *HealthScoreRepository) GetLatest(ctx context.Context) (*engine.HealthScoreOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT score_data FROM health_scores ORDER BY created_at DESC LIMIT 1`).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get latest health score: %w", err) }

	var o engine.HealthScoreOutput
	if err := json.Unmarshal(data, &o); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	return &o, nil
}
