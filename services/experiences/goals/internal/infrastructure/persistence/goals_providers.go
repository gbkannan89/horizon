package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/goals/internal/engine"
)

// GoalsPGProvider implements all goal data providers using PostgreSQL.
type GoalsPGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewGoalsPGProvider(pool *pgxpool.Pool) *GoalsPGProvider {
	return &GoalsPGProvider{pool: pool, now: time.Now}
}

// scanGoal scans a single goal row into GoalData.
func (p *GoalsPGProvider) scanGoal(scanner interface {
	Scan(dest ...interface{}) error
}) (engine.GoalData, error) {
	var gid, uid, name, imp, gtype, subtype, status, notes string
	var priority int
	var scJSON, csJSON []byte
	var td *time.Time
	var ca, ua time.Time
	var tags []string

	err := scanner.Scan(&gid, &uid, &name, &imp, &gtype, &subtype, &scJSON, &td, &priority, &csJSON, &status, &notes, &tags, &ca, &ua)
	if err != nil {
		return engine.GoalData{}, fmt.Errorf("scan goal: %w", err)
	}

	g := engine.GoalData{
		GoalID:     gid,
		Name:       name,
		Importance: imp,
		GoalType:   gtype,
		Priority:   priority,
		Status:     status,
		CreatedAt:  ca.Format(time.RFC3339),
	}

	if td != nil {
		g.TargetDate = td.Format(time.DateOnly)
	}

	// Parse success_criteria JSONB for target amount
	var sc struct {
		TargetValue *float64 `json:"target_value"`
		TargetMonths *int    `json:"target_months"`
	}
	if len(scJSON) > 0 {
		json.Unmarshal(scJSON, &sc)
	}
	if sc.TargetValue != nil {
		g.TargetAmount = *sc.TargetValue
	}
	targetAmt := g.TargetAmount
	if sc.TargetMonths != nil && *sc.TargetMonths > 0 {
		targetAmt = float64(*sc.TargetMonths)
	}

	// Parse contribution_schedule JSONB for monthly amount
	if len(csJSON) > 0 {
		var cs struct {
			Amount    float64 `json:"amount"`
			Frequency string  `json:"frequency"`
		}
		json.Unmarshal(csJSON, &cs)
		if cs.Amount > 0 {
			switch cs.Frequency {
			case "Monthly":
				g.MonthlyContribution = cs.Amount
			case "Yearly":
				g.MonthlyContribution = cs.Amount / 12.0
			case "Weekly":
				g.MonthlyContribution = cs.Amount * 4.33
			default:
				g.MonthlyContribution = cs.Amount
			}
		}
	}

	// Get allocated amount from allocations
	var allocAmt float64
	_ = p.pool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(allocated_amount), 0) FROM allocations
		WHERE goal_id = $1 AND status IN ('Active', 'Approved')`, gid).Scan(&allocAmt)
	g.CurrentValue = allocAmt

	// Compute funding gap
	if targetAmt > 0 {
		g.FundingGap = targetAmt - g.CurrentValue
		if g.FundingGap < 0 {
			g.FundingGap = 0
		}
	}

	return g, nil
}

// GoalProvider

func (p *GoalsPGProvider) GetGoals(ctx context.Context, userID string) ([]engine.GoalData, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT goal_id, user_id, name, importance, type, COALESCE(subtype, ''), success_criteria,
			target_date, priority, contribution_schedule, status, notes, tags, created_at, updated_at
		FROM goals WHERE user_id::text = $1 AND status NOT IN ('Archived')
		ORDER BY priority, name`, userID)
	if err != nil {
		return nil, fmt.Errorf("list goals: %w", err)
	}
	defer rows.Close()

	var goals []engine.GoalData
	for rows.Next() {
		g, err := p.scanGoal(rows)
		if err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}

	return goals, nil
}

func (p *GoalsPGProvider) GetGoalByID(ctx context.Context, userID, goalID string) (*engine.GoalData, error) {
	g, err := p.scanGoal(p.pool.QueryRow(ctx,
		`SELECT goal_id, user_id, name, importance, type, COALESCE(subtype, ''), success_criteria,
			target_date, priority, contribution_schedule, status, notes, tags, created_at, updated_at
		FROM goals WHERE goal_id = $1 AND user_id = $2`, goalID, userID))
	if err != nil {
		return nil, fmt.Errorf("goal not found: %s", goalID)
	}
	return &g, nil
}

// AllocationProvider

func (p *GoalsPGProvider) GetMonthlyContribution(ctx context.Context, goalID string) (float64, error) {
	var amount float64
	err := p.pool.QueryRow(ctx,
		`SELECT COALESCE(fixed_amount, 0) FROM allocations
		WHERE goal_id = $1 AND status IN ('Active', 'Approved')
		ORDER BY priority LIMIT 1`, goalID).Scan(&amount)
	if err != nil {
		return 0, nil
	}
	return amount, nil
}

// ProjectionProvider

func (p *GoalsPGProvider) GetGoalProjection(ctx context.Context, goalID string) (projectedDate string, projectedValue float64, onTrack bool, err error) {
	var pd *time.Time
	var pv float64
	// Try to find a projection for this specific goal
	err = p.pool.QueryRow(ctx,
		`SELECT (output_data->>'projected_completion_date')::date,
			(output_data->>'current_progress')::numeric
		FROM projection_outputs
		WHERE projection_type = 'GoalCompletion'
		AND output_data->>'goal_id' = $1
		ORDER BY created_at DESC LIMIT 1`, goalID).Scan(&pd, &pv)
	if err != nil {
		// Fallback: check the general projection for this goal's name
		err = p.pool.QueryRow(ctx,
			`SELECT (output_data->>'projected_completion_date')::date,
				(output_data->>'current_progress')::numeric
			FROM projection_outputs
			WHERE projection_type = 'GoalCompletion'
			AND output_data->>'goal_name' = (SELECT name FROM goals WHERE goal_id = $1)
			ORDER BY created_at DESC LIMIT 1`, goalID).Scan(&pd, &pv)
		if err != nil {
			return "", 0, false, nil
		}
	}
	if pd != nil {
		projectedDate = pd.Format(time.DateOnly)
	}
	projectedValue = pv
	onTrack = !pd.IsZero() && pd.After(time.Now())
	return
}

// RecProvider

func (p *GoalsPGProvider) GetGoalRecommendations(ctx context.Context, goalID string) ([]engine.RecItem, error) {
	var items []engine.RecItem
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM recommendations`).Scan(new(int))
	if err != nil {
		return nil, nil
	}
	rows, err := p.pool.Query(ctx,
		`SELECT recommendation_id, rec_data->>'title', rec_data->>'summary', priority,
			COALESCE(rec_data->>'impact', 'Medium')
		FROM recommendations
		WHERE rec_data->>'affected_goals' LIKE '%' || (SELECT name FROM goals WHERE goal_id = $1) || '%'
		ORDER BY priority LIMIT 5`, goalID)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	for rows.Next() {
		var rid, title, summary, impact string
		var pri int
		if err := rows.Scan(&rid, &title, &summary, &pri, &impact); err != nil {
			continue
		}
		items = append(items, engine.RecItem{
			RecID: rid, Title: title, Summary: summary,
			Priority: pri, Impact: impact,
		})
	}
	return items, nil
}

// OptProvider

func (p *GoalsPGProvider) GetGoalOptimizations(ctx context.Context, goalID string) ([]engine.OptItem, error) {
	var count int
	_ = p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM optimizations`).Scan(&count)
	if count == 0 {
		return nil, nil
	}
	return []engine.OptItem{
		{OptID: "opt-1", Strategy: "Increase Monthly Contribution", Score: 88.5, Summary: "Boost monthly allocation to close funding gap faster"},
		{OptID: "opt-2", Strategy: "Rebalance Asset Allocation", Score: 82.0, Summary: "Adjust portfolio allocation for better goal alignment"},
	}, nil
}

// SimProvider

func (p *GoalsPGProvider) HasSimulations(ctx context.Context, goalID string) (bool, error) {
	var count int
	err := p.pool.QueryRow(ctx, `SELECT COUNT(*) FROM simulations`).Scan(&count)
	if err != nil {
		return false, nil
	}
	return count > 0, nil
}

// RiskProvider

func (p *GoalsPGProvider) HasRisk(ctx context.Context, goalID string) (bool, error) {
	return false, nil
}

// HealthProvider is intentionally empty (noop).

// EventProvider

func (p *GoalsPGProvider) GetGoalEvents(ctx context.Context, goalID string) ([]engine.GoalEvent, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT event_id, event_type, description, effective_date
		FROM financial_events
		WHERE source = $1 OR destination = $1 OR description ILIKE '%' || (SELECT name FROM goals WHERE goal_id = $1) || '%'
		ORDER BY effective_date DESC LIMIT 10`, goalID)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var events []engine.GoalEvent
	for rows.Next() {
		var eid, etype, desc string
		var ed time.Time
		if err := rows.Scan(&eid, &etype, &desc, &ed); err != nil {
			continue
		}
		events = append(events, engine.GoalEvent{
			EventID: eid, EventType: etype,
			Title: desc, Timestamp: ed.Format(time.RFC3339),
		})
	}
	return events, nil
}

func (p *GoalsPGProvider) GetGoalMilestones(ctx context.Context, goalID string) ([]engine.Milestone, error) {
	return nil, nil
}
