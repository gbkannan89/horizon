package engine

import (
	"fmt"
	"time"
)

// DashboardMode represents the current dashboard mode.
type DashboardMode string
const (DMOverview DashboardMode = "Overview"; DMDecision DashboardMode = "Decision"; DMGoal DashboardMode = "Goal"; DMPortfolio DashboardMode = "Portfolio"; DMPlanning DashboardMode = "Planning"; DMHistorical DashboardMode = "Historical")

// DashboardState represents the overall dashboard state.
type DashboardState string
const (DSFirstTimeUser DashboardState = "FirstTimeUser"; DSNoData DashboardState = "NoData"; DSHealthy DashboardState = "Healthy"; DSAttentionNeeded DashboardState = "AttentionNeeded"; DSCritical DashboardState = "Critical"; DSGoalAchieved DashboardState = "GoalAchieved"; DSGoalAtRisk DashboardState = "GoalAtRisk"; DSOffline DashboardState = "Offline"; DSHistoricalView DashboardState = "HistoricalView")

// WidgetTier represents the priority tier.
type WidgetTier int
const (Tier1 WidgetTier = 1; Tier2 WidgetTier = 2; Tier3 WidgetTier = 3)

// Widget represents a single dashboard widget.
type Widget struct {
	WidgetID   string      `json:"widget_id"`
	WidgetType string      `json:"widget_type"`
	Title      string      `json:"title"`
	Tier       WidgetTier  `json:"tier"`
	Priority   int         `json:"priority"`
	Data       interface{} `json:"data"`
	Confidence string      `json:"confidence"`
	Visible    bool        `json:"visible"`
}

// DashboardOutput is the complete dashboard view model.
type DashboardOutput struct {
	DashboardID     string         `json:"dashboard_id"`
	Mode            DashboardMode  `json:"mode"`
	State           DashboardState `json:"state"`
	Tier1           []Widget       `json:"tier1"`
	Tier2           []Widget       `json:"tier2"`
	Tier3           []Widget       `json:"tier3"`
	CriticalAlert   *Widget        `json:"critical_alert,omitempty"`
	LastRefreshTime string         `json:"last_refresh_time"`
}

// SummaryOutput is the lightweight dashboard summary.
type SummaryOutput struct {
	NetWorth       int64  `json:"net_worth"`
	HealthScore    int    `json:"health_score"`
	HealthGrade    string `json:"health_grade"`
	RiskScore      int    `json:"risk_score"`
	RiskLevel      string `json:"risk_level"`
	GoalsOnTrack   int    `json:"goals_on_track"`
	TotalGoals     int    `json:"total_goals"`
	CashBalance    int64  `json:"cash_balance"`
	PortfolioValue int64  `json:"portfolio_value"`
	TotalDebt      int64  `json:"total_debt"`
}

// Inputs aggregates all data needed for the dashboard.
type Inputs struct {
	UserID            string `json:"user_id"`
	NetWorth          int64  `json:"net_worth"`
	TotalAssets       int64  `json:"total_assets"`
	TotalLiabilities  int64  `json:"total_liabilities"`
	CashBalance       int64  `json:"cash_balance"`
	PortfolioValue    int64  `json:"portfolio_value"`
	MonthlyIncome     int64  `json:"monthly_income"`
	MonthlyExpenses   int64  `json:"monthly_expenses"`
	SavingsRate       float64 `json:"savings_rate"`
	EmergencyMonths   float64 `json:"emergency_months"`
	BillsDue          int64  `json:"bills_due"`
	HealthScore       int    `json:"health_score"`
	HealthGrade       string `json:"health_grade"`
	RiskScore         int    `json:"risk_score"`
	RiskLevel         string `json:"risk_level"`
	GoalCount         int    `json:"goal_count"`
	GoalsOnTrack      int    `json:"goals_on_track"`
	GoalsAtRisk       int    `json:"goals_at_risk"`
	HasRecommendation bool   `json:"has_recommendation"`
	RecCount          int    `json:"rec_count"`
	HasSimulation     bool   `json:"has_simulation"`
	IsFirstTimeUser   bool   `json:"is_first_time_user"`
	HasAccounts       bool   `json:"has_accounts"`
	HasEvents         bool   `json:"has_events"`
	HasGoals          bool   `json:"has_goals"`
	RecentEvents      int    `json:"recent_events"`
	Achievements      int    `json:"achievements"`
}

// Composer builds the dashboard view model from domain data.
type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

func (c *Composer) Compose(inputs Inputs, mode DashboardMode) *DashboardOutput {
	if mode == "" { mode = DMOverview }
	state := c.detectState(inputs)
	t1, t2, t3, alert := c.buildWidgets(inputs, state, mode)
	return &DashboardOutput{
		DashboardID: "dash-" + inputs.UserID,
		Mode: mode, State: state,
		Tier1: t1, Tier2: t2, Tier3: t3,
		CriticalAlert: alert,
		LastRefreshTime: time.Now().UTC().Format(time.RFC3339),
	}
}

func (c *Composer) BuildSummary(inputs Inputs) *SummaryOutput {
	return &SummaryOutput{
		NetWorth: inputs.NetWorth, HealthScore: inputs.HealthScore,
		HealthGrade: inputs.HealthGrade, RiskScore: inputs.RiskScore,
		RiskLevel: inputs.RiskLevel, GoalsOnTrack: inputs.GoalsOnTrack,
		TotalGoals: inputs.GoalCount, CashBalance: inputs.CashBalance,
		PortfolioValue: inputs.PortfolioValue, TotalDebt: inputs.TotalLiabilities,
	}
}

func (c *Composer) detectState(inputs Inputs) DashboardState {
	if inputs.IsFirstTimeUser && !inputs.HasGoals { return DSFirstTimeUser }
	if inputs.HasAccounts && !inputs.HasEvents { return DSNoData }
	if inputs.HealthScore < 40 || inputs.RiskScore >= 80 { return DSCritical }
	if inputs.GoalsAtRisk > 0 { return DSGoalAtRisk }
	if inputs.HealthScore >= 60 { return DSHealthy }
	return DSAttentionNeeded
}

func (c *Composer) buildWidgets(inputs Inputs, state DashboardState, mode DashboardMode) ([]Widget, []Widget, []Widget, *Widget) {
	var t1, t2, t3 []Widget; var critical *Widget

	if state == DSCritical || inputs.RiskScore >= 80 {
		critical = &Widget{WidgetID: "alert-1", WidgetType: "critical_alert", Title: "Critical Alert — Attention Required", Tier: Tier1, Confidence: "High", Visible: true}
	}

	t1 = append(t1, Widget{WidgetID: "health", WidgetType: "health_score", Title: fmtHealth(inputs.HealthScore, inputs.HealthGrade), Tier: Tier1, Priority: 1, Data: inputs.HealthScore, Confidence: "High", Visible: true})

	if inputs.GoalCount > 0 {
		t1 = append(t1, Widget{WidgetID: "goals", WidgetType: "goal_progress", Title: fmt.Sprintf("%d/%d Goals On Track", inputs.GoalsOnTrack, inputs.GoalCount), Tier: Tier1, Priority: 2, Visible: true})
	}
	if inputs.HasRecommendation {
		t1 = append(t1, Widget{WidgetID: "rec", WidgetType: "top_recommendation", Title: "Top Recommendation Available", Tier: Tier1, Priority: 3, Confidence: "High", Visible: true})
	}

	t2 = append(t2, Widget{WidgetID: "nw", WidgetType: "net_worth", Title: fmtMoney(inputs.NetWorth), Tier: Tier2, Priority: 1, Data: inputs.NetWorth, Visible: true})
	t2 = append(t2, Widget{WidgetID: "cash", WidgetType: "cash_position", Title: fmtMoney(inputs.CashBalance), Tier: Tier2, Priority: 2, Visible: true})
	if inputs.PortfolioValue > 0 {
		t2 = append(t2, Widget{WidgetID: "pf", WidgetType: "portfolio_snapshot", Title: fmtMoney(inputs.PortfolioValue), Tier: Tier2, Priority: 3, Visible: true})
	}
	t2 = append(t2, Widget{WidgetID: "risk", WidgetType: "risk_summary", Title: fmt.Sprintf("Risk: %d — %s", inputs.RiskScore, inputs.RiskLevel), Tier: Tier2, Priority: 4, Visible: true})
	if inputs.BillsDue > 0 {
		t2 = append(t2, Widget{WidgetID: "bills", WidgetType: "upcoming_bills", Title: fmtMoney(inputs.BillsDue) + " due", Tier: Tier2, Priority: 5, Visible: true})
	}
	t2 = append(t2, Widget{WidgetID: "debt", WidgetType: "debt_summary", Title: fmtMoney(inputs.TotalLiabilities), Tier: Tier2, Priority: 6, Visible: inputs.TotalLiabilities > 0})

	t3 = append(t3, Widget{WidgetID: "cf", WidgetType: "cash_flow", Title: fmt.Sprintf("Income: %s | Expenses: %s", fmtMoney(inputs.MonthlyIncome), fmtMoney(inputs.MonthlyExpenses)), Tier: Tier3, Visible: true})
	if inputs.RecentEvents > 0 {
		t3 = append(t3, Widget{WidgetID: "events", WidgetType: "recent_events", Title: fmt.Sprintf("%d Recent Events", inputs.RecentEvents), Tier: Tier3, Visible: true})
	}
	if inputs.Achievements > 0 {
		t3 = append(t3, Widget{WidgetID: "achievements", WidgetType: "achievements", Title: fmt.Sprintf("%d Achievements", inputs.Achievements), Tier: Tier3, Visible: true})
	}

	return t1, t2, t3, critical
}

func fmtHealth(score int, grade string) string { return fmt.Sprintf("Health: %d (%s)", score, grade) }
func fmtMoney(v int64) string                  { return fmt.Sprintf("₹%d", v) }
