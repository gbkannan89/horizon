package context

// AdvisorContext is the structured, deterministic context for AI consumption.
// Every value comes from an existing deterministic engine or domain service.
type AdvisorContext struct {
	UserSummary      UserSummaryCtx      `json:"user_summary"`
	FinancialSummary FinancialSummaryCtx `json:"financial_summary"`
	GoalSummary      GoalSummaryCtx      `json:"goal_summary"`
	AccountSummary   AccountSummaryCtx   `json:"account_summary"`
	PortfolioSummary PortfolioSummaryCtx `json:"portfolio_summary"`
	Health           HealthCtx           `json:"health"`
	Risk             RiskCtx             `json:"risk"`
	Projection       ProjectionCtx       `json:"projection"`
	Recommendations  RecsCtx             `json:"recommendations"`
	Optimizations    OptsCtx             `json:"optimizations"`
	Simulations      SimsCtx             `json:"simulations"`
	TimelineEvents   EventsCtx           `json:"timeline_events"`
	Notifications    NotifsCtx           `json:"notifications"`
	Insights         InsightsCtx         `json:"insights"`
}

type UserSummaryCtx struct { UserID string `json:"user_id"` }
type FinancialSummaryCtx struct {
	NetWorth      int64  `json:"net_worth"`
	CashBalance   int64  `json:"cash_balance"`
	MonthlyIncome int64  `json:"monthly_income"`
	MonthlyExp    int64  `json:"monthly_expenses"`
}
type GoalSummaryCtx struct {
	TotalGoals int   `json:"total_goals"`
	OnTrack    int   `json:"on_track"`
	AtRisk     int   `json:"at_risk"`
	FundingGap int64 `json:"funding_gap"`
}
type AccountSummaryCtx struct {
	TotalAccounts int   `json:"total_accounts"`
	TotalBalance  int64 `json:"total_balance"`
}
type PortfolioSummaryCtx struct {
	PortfolioValue int64   `json:"portfolio_value"`
	TotalReturn    float64 `json:"total_return"`
	ReturnPct      float64 `json:"return_pct"`
	RiskScore      int     `json:"risk_score"`
}
type HealthCtx struct {
	Score  int    `json:"score"`
	Grade  string `json:"grade"`
	Change int    `json:"change"`
}
type RiskCtx struct { Score int `json:"score"`; Level string `json:"level"` }
type ProjectionCtx struct {
	NetWorthProjected int64  `json:"net_worth_projected"`
	OnTrack           bool   `json:"on_track"`
	Confidence        string `json:"confidence"`
}
type RecsCtx struct { HasRecs bool `json:"has_recommendations"`; Count int `json:"count"` }
type OptsCtx struct { HasOpts bool `json:"has_optimizations"`; Count int `json:"count"` }
type SimsCtx struct { HasSims bool `json:"has_simulations"`; Count int `json:"count"` }
type EventsCtx struct { RecentCount int `json:"recent_count"` }
type NotifsCtx struct { UnreadCount int `json:"unread_count"` }
type InsightsCtx struct { TotalCount int `json:"total_count"`; CriticalCount int `json:"critical_count"` }

// Builder constructs the AI context from raw inputs.
type Builder struct{}

func NewBuilder() *Builder { return &Builder{} }

func (b *Builder) Build(inputs Inputs) *AdvisorContext {
	return &AdvisorContext{
		UserSummary:      UserSummaryCtx{UserID: inputs.UserID},
		FinancialSummary: FinancialSummaryCtx{NetWorth: inputs.NetWorth, CashBalance: inputs.CashBalance, MonthlyIncome: inputs.MonthlyIncome, MonthlyExp: inputs.MonthlyExp},
		GoalSummary:      GoalSummaryCtx{TotalGoals: inputs.TotalGoals, OnTrack: inputs.OnTrack, AtRisk: inputs.AtRisk, FundingGap: inputs.FundingGap},
		AccountSummary:   AccountSummaryCtx{TotalAccounts: inputs.TotalAccounts, TotalBalance: inputs.TotalBalance},
		PortfolioSummary: PortfolioSummaryCtx{PortfolioValue: inputs.PortfolioValue, TotalReturn: inputs.TotalReturn, ReturnPct: inputs.ReturnPct, RiskScore: inputs.RiskScore},
		Health:           HealthCtx{Score: inputs.HealthScore, Grade: inputs.HealthGrade, Change: inputs.HealthChg},
		Risk:             RiskCtx{Score: inputs.RiskScore, Level: inputs.RiskLevel},
		Projection:       ProjectionCtx{NetWorthProjected: inputs.NWProjected, OnTrack: inputs.ProjOnTrack, Confidence: inputs.ProjConf},
		Recommendations:  RecsCtx{HasRecs: inputs.HasRecs, Count: inputs.RecCount},
		Optimizations:    OptsCtx{HasOpts: inputs.HasOpts, Count: inputs.OptCount},
		Simulations:      SimsCtx{HasSims: inputs.HasSims, Count: inputs.SimCount},
		TimelineEvents:   EventsCtx{RecentCount: inputs.EventCount},
		Notifications:    NotifsCtx{UnreadCount: inputs.NotifUnread},
		Insights:         InsightsCtx{TotalCount: inputs.InsightTotal, CriticalCount: inputs.InsightCrit},
	}
}

// ToFlatMap serializes the context into a flat map for prompt injection.
func (b *Builder) ToFlatMap(ctx *AdvisorContext) map[string]interface{} {
	return map[string]interface{}{
		"net_worth":            ctx.FinancialSummary.NetWorth,
		"cash_balance":         ctx.FinancialSummary.CashBalance,
		"income":               ctx.FinancialSummary.MonthlyIncome,
		"expenses":             ctx.FinancialSummary.MonthlyExp,
		"total_goals":          ctx.GoalSummary.TotalGoals,
		"goals_on_track":       ctx.GoalSummary.OnTrack,
		"goals_at_risk":        ctx.GoalSummary.AtRisk,
		"funding_gap":          ctx.GoalSummary.FundingGap,
		"total_accounts":       ctx.AccountSummary.TotalAccounts,
		"total_balance":        ctx.AccountSummary.TotalBalance,
		"portfolio_value":      ctx.PortfolioSummary.PortfolioValue,
		"portfolio_return":     ctx.PortfolioSummary.TotalReturn,
		"portfolio_return_pct": ctx.PortfolioSummary.ReturnPct,
		"portfolio_risk":       ctx.PortfolioSummary.RiskScore,
		"health_score":         ctx.Health.Score,
		"health_grade":         ctx.Health.Grade,
		"health_change":        ctx.Health.Change,
		"risk_score":           ctx.Risk.Score,
		"risk_level":           ctx.Risk.Level,
		"net_worth_projected":  ctx.Projection.NetWorthProjected,
		"projection_on_track":  ctx.Projection.OnTrack,
		"projection_confidence": ctx.Projection.Confidence,
		"has_recommendations":  ctx.Recommendations.HasRecs,
		"rec_count":            ctx.Recommendations.Count,
		"has_optimizations":    ctx.Optimizations.HasOpts,
		"opt_count":            ctx.Optimizations.Count,
		"has_simulations":      ctx.Simulations.HasSims,
		"sim_count":            ctx.Simulations.Count,
		"event_count":          ctx.TimelineEvents.RecentCount,
		"notif_unread":         ctx.Notifications.UnreadCount,
		"insight_total":        ctx.Insights.TotalCount,
		"insight_critical":     ctx.Insights.CriticalCount,
	}
}

// Inputs for context building.
type Inputs struct {
	UserID       string `json:"user_id"`
	NetWorth     int64  `json:"net_worth"`
	CashBalance  int64  `json:"cash_balance"`
	MonthlyIncome int64 `json:"monthly_income"`
	MonthlyExp   int64  `json:"monthly_expenses"`
	TotalGoals   int    `json:"total_goals"`
	OnTrack      int    `json:"on_track"`
	AtRisk       int    `json:"at_risk"`
	FundingGap   int64  `json:"funding_gap"`
	TotalAccounts int   `json:"total_accounts"`
	TotalBalance int64  `json:"total_balance"`
	PortfolioValue int64 `json:"portfolio_value"`
	TotalReturn  float64 `json:"total_return"`
	ReturnPct    float64 `json:"return_pct"`
	RiskScore    int    `json:"risk_score"`
	HealthScore  int    `json:"health_score"`
	HealthGrade  string `json:"health_grade"`
	HealthChg    int    `json:"health_change"`
	RiskLevel    string `json:"risk_level"`
	NWProjected  int64  `json:"net_worth_projected"`
	ProjOnTrack  bool   `json:"projection_on_track"`
	ProjConf     string `json:"projection_confidence"`
	HasRecs      bool   `json:"has_recommendations"`
	RecCount     int    `json:"rec_count"`
	HasOpts      bool   `json:"has_optimizations"`
	OptCount     int    `json:"opt_count"`
	HasSims      bool   `json:"has_simulations"`
	SimCount     int    `json:"sim_count"`
	EventCount   int    `json:"event_count"`
	NotifUnread  int    `json:"notif_unread_count"`
	InsightTotal int    `json:"insight_total"`
	InsightCrit  int    `json:"insight_critical"`
}
