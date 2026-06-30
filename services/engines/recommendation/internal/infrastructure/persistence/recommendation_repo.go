package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/recommendation/internal/engine"
)

type RecommendationRepository struct {
	pool *pgxpool.Pool
}

func NewRecommendationRepository(pool *pgxpool.Pool) *RecommendationRepository {
	return &RecommendationRepository{pool: pool}
}

func (r *RecommendationRepository) Save(ctx context.Context, rec *engine.Recommendation) error {
	data, err := json.Marshal(rec)
	if err != nil { return fmt.Errorf("marshal: %w", err) }

	_, err = r.pool.Exec(ctx,
		`INSERT INTO recommendations (recommendation_id, category, priority, score, status, rec_data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (recommendation_id) DO UPDATE SET status=$5, rec_data=$6`,
		rec.ID, string(rec.Category), rec.Priority, rec.Score, string(rec.Status), data)
	return err
}

func (r *RecommendationRepository) UpdateStatus(ctx context.Context, id string, status engine.RecStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE recommendations SET status=$1 WHERE recommendation_id=$2`, string(status), id)
	return err
}

func (r *RecommendationRepository) GetActive(ctx context.Context) ([]engine.Recommendation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT rec_data FROM recommendations WHERE status IN ('Generated','Viewed') ORDER BY priority LIMIT 50`)
	if err != nil { return nil, err }
	defer rows.Close()

	var recs []engine.Recommendation
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil { return nil, err }
		var rec engine.Recommendation
		if err := json.Unmarshal(data, &rec); err != nil { return nil, err }
		recs = append(recs, rec)
	}
	return recs, nil
}
