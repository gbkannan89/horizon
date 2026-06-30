package engine

import "fmt"

type AdvisorMode string

const (
	AMOverview          AdvisorMode = "Overview"
	AMDecision          AdvisorMode = "Decision"
	AMPlanning          AdvisorMode = "Planning"
	AMLearning          AdvisorMode = "Learning"
	AMGoalCoaching      AdvisorMode = "GoalCoaching"
	AMRiskReview        AdvisorMode = "RiskReview"
	AMRetirementPlan    AdvisorMode = "RetirementPlanning"
	AMPortfolioCoach    AdvisorMode = "PortfolioCoaching"
	AMHistoricalReview  AdvisorMode = "HistoricalReview"
)

type CardType string

const (
	CTBestAction    CardType = "TodayBestAction"
	CTGoalCoach     CardType = "GoalCoach"
	CTHealth        CardType = "HealthSummary"
	CTRisk          CardType = "RiskAlert"
	CTProjection    CardType = "ProjectionSummary"
	CTSimulation    CardType = "SimulationResult"
	CTOptimization  CardType = "OptimizationStrategy"
	CTInsight       CardType = "FinancialInsight"
	CTEducation     CardType = "EducationCard"
	CTAchievement   CardType = "AchievementCard"
	CTTimeline      CardType = "TimelineSummary"
	CTNotifications CardType = "NotificationsSummary"
)

type Card struct {
	CardID        string   `json:"card_id"`
	CardType      CardType `json:"card_type"`
	Title         string   `json:"title"`
	Summary       string   `json:"summary"`
	Priority      int      `json:"priority"`
	Confidence    string   `json:"confidence"`
	IsAIGenerated bool     `json:"is_ai_generated"`
	Actions       []string `json:"actions"`
	Data          interface{} `json:"data,omitempty"`
}

type AdvisorWorkspace struct {
	SessionID string       `json:"session_id"`
	Mode      AdvisorMode  `json:"mode"`
	Cards     []Card       `json:"cards"`
}

// AdvisorContext is the structured, deterministic context for Phase 7 AI consumption.
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

type UserSummaryCtx struct {
	UserID string `json:"user_id"`
}
type FinancialSummaryCtx struct {
	NetWorth      int64  `json:"net_worth"`
	CashBalance   int64  `json:"cash_balance"`
	MonthlyIncome int64  `json:"monthly_income"`
	MonthlyExp    int64  `json:"monthly_expenses"`
}
type GoalSummaryCtx struct {
	TotalGoals   int `json:"total_goals"`
	OnTrack      int `json:"on_track"`
	AtRisk       int `json:"at_risk"`
	FundingGap   int64 `json:"funding_gap"`
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
	Score   int    `json:"score"`
	Grade   string `json:"grade"`
	Change  int    `json:"change"`
}
type RiskCtx struct {
	Score int    `json:"score"`
	Level string `json:"level"`
}
type ProjectionCtx struct {
	NetWorthProjected int64  `json:"net_worth_projected"`
	OnTrack           bool   `json:"on_track"`
	Confidence        string `json:"confidence"`
}
type RecsCtx struct {
	HasRecs bool `json:"has_recommendations"`
	Count   int  `json:"count"`
}
type OptsCtx struct {
	HasOpts bool `json:"has_optimizations"`
	Count   int  `json:"count"`
}
type SimsCtx struct {
	HasSims bool `json:"has_simulations"`
	Count   int  `json:"count"`
}
type EventsCtx struct {
	RecentCount int `json:"recent_count"`
}
type NotifsCtx struct {
	UnreadCount int `json:"unread_count"`
}
type InsightsCtx struct {
	TotalCount    int `json:"total_count"`
	CriticalCount int `json:"critical_count"`
}

// Inputs aggregates across all 8 experiences + 6 engines.
type Inputs struct {
	UserID string `json:"user_id"`

	// Financial
	NetWorth      int64 `json:"net_worth"`
	CashBalance   int64 `json:"cash_balance"`
	MonthlyIncome int64 `json:"monthly_income"`
	MonthlyExp    int64 `json:"monthly_expenses"`

	// Goals
	TotalGoals int `json:"total_goals"`
	OnTrack    int `json:"on_track"`
	AtRisk     int `json:"at_risk"`
	FundingGap int64 `json:"funding_gap"`

	// Accounts
	TotalAccounts int   `json:"total_accounts"`
	TotalBalance  int64 `json:"total_balance"`

	// Portfolio
	PortfolioValue int64   `json:"portfolio_value"`
	TotalReturn    float64 `json:"total_return"`
	ReturnPct      float64 `json:"return_pct"`

	// Health
	HealthScore int    `json:"health_score"`
	HealthGrade string `json:"health_grade"`
	HealthChg   int    `json:"health_change"`

	// Risk
	RiskScore int    `json:"risk_score"`
	RiskLevel string `json:"risk_level"`

	// Projection
	NWProjected int64  `json:"net_worth_projected"`
	ProjOnTrack bool   `json:"projection_on_track"`
	ProjConf    string `json:"projection_confidence"`

	// Recommendations
	HasRecs bool `json:"has_recommendations"`
	RecCount int `json:"rec_count"`
	RecTitle string `json:"rec_title"`

	// Optimization
	HasOpts bool `json:"has_optimizations"`
	OptCount int `json:"opt_count"`

	// Simulation
	HasSims bool `json:"has_simulations"`
	SimCount int `json:"sim_count"`

	// Events
	EventCount int `json:"event_count"`

	// Notifications
	NotifUnread int `json:"notif_unread_count"`

	// Insights
	InsightTotal int `json:"insight_total"`
	InsightCrit  int `json:"insight_critical"`
}

type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

func (c *Composer) BuildWorkspace(inputs Inputs, mode AdvisorMode) *AdvisorWorkspace {
	if mode == "" { mode = AMOverview }
	cards := c.buildCards(inputs, mode)
	return &AdvisorWorkspace{
		SessionID: fmt.Sprintf("adv-%s", inputs.UserID),
		Mode:      mode,
		Cards:     cards,
	}
}

func (c *Composer) BuildContext(inputs Inputs) *AdvisorContext {
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

func (c *Composer) BuildSummary(inputs Inputs) *FinancialSummaryCtx {
	return &FinancialSummaryCtx{
		NetWorth: inputs.NetWorth, CashBalance: inputs.CashBalance,
		MonthlyIncome: inputs.MonthlyIncome, MonthlyExp: inputs.MonthlyExp,
	}
}

func CtxSummary(score int, grade string, change int) string {
	s := fmt.Sprintf("Health: %d (%s)", score, grade)
	if change > 0 { s += fmt.Sprintf(" ↑%d", change) } else if change < 0 { s += fmt.Sprintf(" ↓%d", -change) }
	return s
}
func CtxRiskSummary(score int, level string) string { return fmt.Sprintf("Score: %d — %s", score, level) }
func CtxRecSummary(count int) string { return fmt.Sprintf("%d recommendation(s) available", count) }
func CtxInsightSummary(total, critical int) string {
	if critical > 0 { return fmt.Sprintf("%d critical out of %d", critical, total) }
	return fmt.Sprintf("%d insights", total)
}
func CtxEventSummary(count int) string { return fmt.Sprintf("%d recent events", count) }
func CtxNotifSummary(unread int) string { return fmt.Sprintf("%d unread", unread) }

func (c *Composer) buildCards(inputs Inputs, mode AdvisorMode) []Card {
	var cards []Card

	if inputs.HasRecs {
		title := inputs.RecTitle
		if title == "" { title = "Recommendation Available" }
		cards = append(cards, Card{
			CardID: "c-best", CardType: CTBestAction, Title: title,
			Summary:  fmt.Sprintf("%d recommendation(s)", inputs.RecCount),
			Priority: 1, Confidence: "High", Actions: []string{"Tell me more", "Accept"},
		})
	}

	cards = append(cards, Card{
		CardID: "c-health", CardType: CTHealth, Title: fmt.Sprintf("Health: %d (%s)", inputs.HealthScore, inputs.HealthGrade),
		Summary: "Your financial health overview", Priority: 2, Confidence: "High",
	})

	if inputs.RiskScore >= 60 {
		cards = append(cards, Card{
			CardID: "c-risk", CardType: CTRisk, Title: fmt.Sprintf("Risk: %d (%s)", inputs.RiskScore, inputs.RiskLevel),
			Summary: "Review your risk exposure", Priority: 3, Confidence: "High",
		})
	}

	if inputs.TotalGoals > 0 {
		cards = append(cards, Card{
			CardID: "c-goals", CardType: CTGoalCoach, Title: "Goal Progress",
			Summary:  fmt.Sprintf("%d/%d on track", inputs.OnTrack, inputs.TotalGoals),
			Priority: 4, Confidence: "High",
		})
	}

	if inputs.NWProjected > 0 {
		cards = append(cards, Card{
			CardID: "c-proj", CardType: CTProjection, Title: "Projection",
			Summary:  fmt.Sprintf("Projected net worth: ₹%d", inputs.NWProjected),
			Priority: 5, Confidence: inputs.ProjConf,
		})
	}

	if inputs.HasOpts {
		cards = append(cards, Card{
			CardID: "c-opt", CardType: CTOptimization, Title: "Optimization",
			Summary: fmt.Sprintf("%d strategies", inputs.OptCount), Priority: 6, Confidence: "High",
		})
	}

	if inputs.HasSims {
		cards = append(cards, Card{
			CardID: "c-sim", CardType: CTSimulation, Title: "Simulation",
			Summary: fmt.Sprintf("%d scenarios", inputs.SimCount), Priority: 7, Confidence: "High",
		})
	}

	if inputs.InsightTotal > 0 {
		summary := fmt.Sprintf("%d total", inputs.InsightTotal)
		if inputs.InsightCrit > 0 { summary = fmt.Sprintf("%d critical, %d total", inputs.InsightCrit, inputs.InsightTotal) }
		cards = append(cards, Card{
			CardID: "c-insight", CardType: CTInsight, Title: "Insights",
			Summary: summary, Priority: 8, Confidence: "High",
		})
	}

	if inputs.EventCount > 0 {
		cards = append(cards, Card{
			CardID: "c-timeline", CardType: CTTimeline, Title: "Recent Activity",
			Summary: fmt.Sprintf("%d events", inputs.EventCount), Priority: 9, Confidence: "High",
		})
	}

	if inputs.NotifUnread > 0 {
		cards = append(cards, Card{
			CardID: "c-notif", CardType: CTNotifications, Title: "Notifications",
			Summary: fmt.Sprintf("%d unread", inputs.NotifUnread), Priority: 10, Confidence: "High",
		})
	}

	return cards
}
