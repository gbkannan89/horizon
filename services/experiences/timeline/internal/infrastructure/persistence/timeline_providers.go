package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/horizon/core/services/experiences/timeline/internal/engine"
)

// TimelinePGProvider implements all timeline event providers using PostgreSQL.
type TimelinePGProvider struct {
	pool *pgxpool.Pool
	now  func() time.Time
}

func NewTimelinePGProvider(pool *pgxpool.Pool) *TimelinePGProvider {
	return &TimelinePGProvider{pool: pool, now: time.Now}
}

func (p *TimelinePGProvider) GetEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	return nil, nil
}

// getLimit returns a safe limit.
func getLimit(limit int) int {
	if limit <= 0 { return 50 }
	return limit
}

// Financial events: all posted financial transactions
func (p *TimelinePGProvider) getFinancialEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT event_id, event_type, amount, currency, description, source, destination,
			effective_date, state, confidence
		FROM financial_events WHERE user_id::text = $1 AND state = 'POSTED'
		ORDER BY effective_date DESC LIMIT $2`, userID, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var eid, etype, cur, desc, src, dst, st, conf string
		var amt int64
		var ed time.Time
		if err := rows.Scan(&eid, &etype, &amt, &cur, &desc, &src, &dst, &ed, &st, &conf); err != nil {
			continue
		}
		summary := fmt.Sprintf("₹%d", amt)
		if amt > 0 { summary = fmt.Sprintf("+₹%d", amt) }
		items = append(items, engine.RawItem{
			TimelineID:  "fin-" + eid,
			Timestamp:   ed.Format(time.RFC3339),
			EventType:   etype,
			Category:    engine.ECFinancial,
			Title:       fmt.Sprintf("%s: %s", etype, desc),
			Summary:     summary,
			Description: desc,
			Severity:    severityForAmount(amt),
			RelatedAgg:  userID,
			Amount:      amt,
			Icon:        "currency_rupee", Color: "green",
		})
	}
	return items, nil
}

func severityForAmount(amt int64) engine.Severity {
	if amt < 0 { return engine.SevWarning }
	return engine.SevInfo
}

// Goal events: goals created, status changes
func (p *TimelinePGProvider) getGoalEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT goal_id, name, status, created_at, updated_at
		FROM goals WHERE user_id::text = $1 AND status NOT IN ('Archived')
		ORDER BY updated_at DESC LIMIT $2`, userID, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var gid, name, st string
		var ca, ua time.Time
		if err := rows.Scan(&gid, &name, &st, &ca, &ua); err != nil {
			continue
		}

		sev := engine.SevInfo
		icon, color := "track_changes", "blue"
		if st == "Completed" { sev, icon, color = engine.SevSuccess, "celebration", "green" }
		if st == "AtRisk" { sev, icon, color = engine.SevWarning, "warning", "amber" }

		ts := ua.Format(time.RFC3339)
		items = append(items, engine.RawItem{
			TimelineID:  "goal-" + gid + "-" + st,
			Timestamp:   ts,
			EventType:   "Goal" + st,
			Category:    engine.ECGoal,
			Title:       fmt.Sprintf("%s: %s", st, name),
			Summary:     fmt.Sprintf("Goal %s", st),
			Description: name,
			Severity:    sev,
			RelatedGoal: gid,
			RelatedAgg:  userID,
			Icon:        icon, Color: color,
		})
	}
	return items, nil
}

// Account events: accounts created or status changed
func (p *TimelinePGProvider) getAccountEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT account_id, account_name, account_type, status, created_at, updated_at
		FROM accounts WHERE owner_id::text = $1 AND status NOT IN ('Draft')
		ORDER BY updated_at DESC LIMIT $2`, userID, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var aid, an, at, st string
		var ca, ua time.Time
		if err := rows.Scan(&aid, &an, &at, &st, &ca, &ua); err != nil {
			continue
		}
		ts := ua.Format(time.RFC3339)
		if st == "Active" && ca.Equal(ua) {
			ts = ca.Format(time.RFC3339)
		}
		items = append(items, engine.RawItem{
			TimelineID:   "acct-" + aid + "-" + st,
			Timestamp:    ts,
			EventType:    "Account" + st,
			Category:     engine.ECAccount,
			Title:        fmt.Sprintf("Account %s: %s", st, an),
			Summary:      fmt.Sprintf("%s - %s", at, st),
			Description:  an,
			Severity:     engine.SevInfo,
			RelatedAcct:  aid,
			RelatedAgg:   userID,
			Icon:         "account_balance", Color: "indigo",
		})
	}
	return items, nil
}

// Asset events: assets created, revalued
func (p *TimelinePGProvider) getAssetEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT asset_id, asset_name, classification, unit_price, created_at, updated_at
		FROM assets WHERE owner_id::text = $1 AND status = 'Active'
		ORDER BY updated_at DESC LIMIT $2`, userID, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var aid, an, cls string
		var up float64
		var ca, ua time.Time
		if err := rows.Scan(&aid, &an, &cls, &up, &ca, &ua); err != nil {
			continue
		}
		items = append(items, engine.RawItem{
			TimelineID:   "ast-" + aid,
			Timestamp:    ua.Format(time.RFC3339),
			EventType:    "AssetUpdated",
			Category:     engine.ECAsset,
			Title:        fmt.Sprintf("Asset: %s", an),
			Summary:      fmt.Sprintf("%s - ₹%.0f", cls, up),
			Description:  an,
			Severity:     engine.SevInfo,
			RelatedAsset: aid,
			RelatedAgg:   userID,
			Amount:       int64(up),
			Icon:         "trending_up", Color: "teal",
		})
	}
	return items, nil
}

// Liability events
func (p *TimelinePGProvider) getLiabilityEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT liability_id, name, classification, original_principal, status, created_at
		FROM liabilities WHERE owner_id::text = $1 AND status = 'Active'
		ORDER BY created_at DESC LIMIT $2`, userID, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var lid, n, cls, st string
		var op int64
		var ca time.Time
		if err := rows.Scan(&lid, &n, &cls, &op, &st, &ca); err != nil { continue }
		items = append(items, engine.RawItem{
			TimelineID:    "liab-" + lid,
			Timestamp:     ca.Format(time.RFC3339),
			EventType:     "LiabilityRecorded",
			Category:      engine.ECLiability,
			Title:         fmt.Sprintf("Liability: %s", n),
			Summary:       fmt.Sprintf("₹%d", op),
			Severity:      engine.SevWarning,
			RelatedEntity: lid,
			Amount:        op,
			Icon:          "credit_card", Color: "purple",
		})
	}
	return items, nil
}

// Portfolio events
func (p *TimelinePGProvider) getPortfolioEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT p.portfolio_id, p.name, p.status, p.created_at,
			COALESCE(SUM(a.unit_price * COALESCE(a.quantity, 1)), 0) as value
		FROM portfolios p
		LEFT JOIN portfolio_members pm ON p.portfolio_id = pm.portfolio_id
		LEFT JOIN assets a ON pm.entity_id = a.asset_id
		WHERE p.owner_id::text = $1 AND p.status NOT IN ('Archived')
		GROUP BY p.portfolio_id, p.name, p.status, p.created_at
		ORDER BY p.created_at DESC LIMIT $2`, userID, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var pid, n, st string
		var ca time.Time
		var val int64
		if err := rows.Scan(&pid, &n, &st, &ca, &val); err != nil { continue }
		items = append(items, engine.RawItem{
			TimelineID:   "pf-" + pid,
			Timestamp:    ca.Format(time.RFC3339),
			EventType:    "Portfolio" + st,
			Category:     engine.ECPortfolio,
			Title:        fmt.Sprintf("Portfolio: %s (%s)", n, st),
			Summary:      fmt.Sprintf("Value: ₹%d", val),
			RelatedAgg:   pid,
			Amount:       val,
			Severity:     engine.SevInfo,
			Icon:         "pie_chart", Color: "blue",
		})
	}
	return items, nil
}

// Health events
func (p *TimelinePGProvider) getHealthEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT score_id, overall_score, score_grade, created_at
		FROM health_scores ORDER BY created_at DESC LIMIT $2`, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var sid, grade string
		var score int
		var ca time.Time
		if err := rows.Scan(&sid, &score, &grade, &ca); err != nil { continue }

		sev, icon, color := engine.SevInfo, "favorite", "pink"
		if score < 40 { sev, icon, color = engine.SevCritical, "heart_broken", "red" }
		if score < 60 { sev, icon, color = engine.SevWarning, "heart_broken", "amber" }

		items = append(items, engine.RawItem{
			TimelineID:  "hlt-" + sid,
			Timestamp:   ca.Format(time.RFC3339),
			EventType:   "HealthScoreUpdated",
			Category:    engine.ECHealth,
			Title:       fmt.Sprintf("Health Score: %d (%s)", score, grade),
			Summary:     fmt.Sprintf("Score: %d/100", score),
			Severity:    sev,
			RelatedAgg:  userID,
			Icon:        icon, Color: color,
		})
	}
	return items, nil
}

// Risk events
func (p *TimelinePGProvider) getRiskEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT assessment_id, composite_score, risk_level, created_at
		FROM risk_assessments ORDER BY created_at DESC LIMIT $2`, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var aid, level string
		var score int
		var ca time.Time
		if err := rows.Scan(&aid, &score, &level, &ca); err != nil { continue }

		sev, icon, color := engine.SevInfo, "shield", "green"
		if score >= 80 { sev, icon, color = engine.SevCritical, "gpp_bad", "red" }
		if score >= 60 { sev, icon, color = engine.SevWarning, "shield", "amber" }

		items = append(items, engine.RawItem{
			TimelineID:  "risk-" + aid,
			Timestamp:   ca.Format(time.RFC3339),
			EventType:   "RiskAssessmentGenerated",
			Category:    engine.ECRisk,
			Title:       fmt.Sprintf("Risk Score: %d (%s)", score, level),
			Summary:     fmt.Sprintf("Risk: %s", level),
			Severity:    sev,
			RelatedAgg:  userID,
			Icon:        icon, Color: color,
		})
	}
	return items, nil
}

// Recommendation events
func (p *TimelinePGProvider) getRecommendationEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT recommendation_id, category, priority, rec_data->>'title', created_at
		FROM recommendations ORDER BY created_at DESC LIMIT $2`, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var rid, cat, title string
		var pri int
		var ca time.Time
		if err := rows.Scan(&rid, &cat, &pri, &title, &ca); err != nil { continue }
		items = append(items, engine.RawItem{
			TimelineID:  "rec-" + rid,
			Timestamp:   ca.Format(time.RFC3339),
			EventType:   "RecommendationsGenerated",
			Category:    engine.ECRecommend,
			Title:       fmt.Sprintf("Recommendation: %s", title),
			Summary:     fmt.Sprintf("Category: %s", cat),
			Severity:    engine.SevInfo,
			RelatedAgg:  userID,
			Icon:        "lightbulb", Color: "amber",
		})
	}
	return items, nil
}

// Simulation events
func (p *TimelinePGProvider) getSimulationEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT scenario_id, sim_type, created_at FROM simulations
		ORDER BY created_at DESC LIMIT $2`, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var sid, st string
		var ca time.Time
		if err := rows.Scan(&sid, &st, &ca); err != nil { continue }
		items = append(items, engine.RawItem{
			TimelineID:  "sim-" + sid,
			Timestamp:   ca.Format(time.RFC3339),
			EventType:   "SimulationCompleted",
			Category:    engine.ECSimulation,
			Title:       fmt.Sprintf("Simulation: %s", st),
			Summary:     "Scenario evaluated",
			Severity:    engine.SevInfo,
			RelatedAgg:  userID,
			Icon:        "science", Color: "cyan",
		})
	}
	return items, nil
}

// Optimization events
func (p *TimelinePGProvider) getOptimizationEvents(ctx context.Context, userID string, limit int) ([]engine.RawItem, error) {
	rows, err := p.pool.Query(ctx,
		`SELECT optimization_id, status, created_at FROM optimizations
		ORDER BY created_at DESC LIMIT $2`, getLimit(limit))
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var items []engine.RawItem
	for rows.Next() {
		var oid, st string
		var ca time.Time
		if err := rows.Scan(&oid, &st, &ca); err != nil { continue }
		items = append(items, engine.RawItem{
			TimelineID:  "opt-" + oid,
			Timestamp:   ca.Format(time.RFC3339),
			EventType:   "OptimizationGenerated",
			Category:    engine.ECOptimization,
			Title:       fmt.Sprintf("Optimization: %s", st),
			Summary:     "Strategy evaluated",
			Severity:    engine.SevInfo,
			RelatedAgg:  userID,
			Icon:        "auto_graph", Color: "cyan",
		})
	}
	return items, nil
}

// Named provider types — each implements aggregator.EventProvider by calling a specific source method.
// These are exported so main.go can wire them into the aggregator.DataProviders struct.

type FinEventsProvider struct{ P *TimelinePGProvider }
func (f *FinEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return f.P.getFinancialEvents(ctx, uid, l) }

type GoalEventsProvider struct{ P *TimelinePGProvider }
func (g *GoalEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return g.P.getGoalEvents(ctx, uid, l) }

type AcctEventsProvider struct{ P *TimelinePGProvider }
func (a *AcctEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return a.P.getAccountEvents(ctx, uid, l) }

type AssetEventsProvider struct{ P *TimelinePGProvider }
func (a *AssetEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return a.P.getAssetEvents(ctx, uid, l) }

type LiabEventsProvider struct{ P *TimelinePGProvider }
func (l *LiabEventsProvider) GetEvents(ctx context.Context, uid string, lim int) ([]engine.RawItem, error) { return l.P.getLiabilityEvents(ctx, uid, lim) }

type PfEventsProvider struct{ P *TimelinePGProvider }
func (p *PfEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return p.P.getPortfolioEvents(ctx, uid, l) }

type HealthEventsProvider struct{ P *TimelinePGProvider }
func (h *HealthEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return h.P.getHealthEvents(ctx, uid, l) }

type RiskEventsProvider struct{ P *TimelinePGProvider }
func (r *RiskEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return r.P.getRiskEvents(ctx, uid, l) }

type RecEventsProvider struct{ P *TimelinePGProvider }
func (r *RecEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return r.P.getRecommendationEvents(ctx, uid, l) }

type SimEventsProvider struct{ P *TimelinePGProvider }
func (s *SimEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return s.P.getSimulationEvents(ctx, uid, l) }

type OptEventsProvider struct{ P *TimelinePGProvider }
func (o *OptEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return o.P.getOptimizationEvents(ctx, uid, l) }

type UserEventsProvider struct{}
func (u *UserEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return nil, nil }

type AchieveEventsProvider struct{}
func (a *AchieveEventsProvider) GetEvents(ctx context.Context, uid string, l int) ([]engine.RawItem, error) { return nil, nil }
