package persistence

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/engines/simulation/internal/engine"
)

type SimulationRepository struct {
	pool *pgxpool.Pool
}

func NewSimulationRepository(pool *pgxpool.Pool) *SimulationRepository {
	return &SimulationRepository{pool: pool}
}

func (r *SimulationRepository) Save(ctx context.Context, output *engine.SimulationOutput) error {
	data, err := json.Marshal(output)
	if err != nil { return fmt.Errorf("marshal: %w", err) }

	_, err = r.pool.Exec(ctx,
		`INSERT INTO simulations (scenario_id, sim_type, status, sim_data, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (scenario_id) DO UPDATE SET status=$3, sim_data=$4`,
		output.ScenarioID, string(output.SimType), string(output.Status), data)
	return err
}

func (r *SimulationRepository) GetByID(ctx context.Context, id string) (*engine.SimulationOutput, error) {
	var data []byte
	err := r.pool.QueryRow(ctx,
		`SELECT sim_data FROM simulations WHERE scenario_id = $1`, id).Scan(&data)
	if err != nil { return nil, fmt.Errorf("get simulation %s: %w", id, err) }

	var o engine.SimulationOutput
	if err := json.Unmarshal(data, &o); err != nil { return nil, fmt.Errorf("unmarshal: %w", err) }
	return &o, nil
}

func (r *SimulationRepository) GetLatest(ctx context.Context) ([]*engine.SimulationOutput, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT sim_data FROM simulations ORDER BY created_at DESC LIMIT 10`)
	if err != nil { return nil, err }
	defer rows.Close()

	var outputs []*engine.SimulationOutput
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil { return nil, err }
		var o engine.SimulationOutput
		if err := json.Unmarshal(data, &o); err != nil { return nil, err }
		outputs = append(outputs, &o)
	}
	return outputs, nil
}
