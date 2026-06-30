package engine

import "fmt"

type DashboardState string
const (
	DSFirstTimeUser   DashboardState = "FirstTimeUser"
	DSNoData          DashboardState = "NoData"
	DSHealthy         DashboardState = "Healthy"
	DSAttentionNeeded DashboardState = "AttentionNeeded"
	DSCritical        DashboardState = "Critical"
	DSGoalAchieved    DashboardState = "GoalAchieved"
	DSGoalAtRisk      DashboardState = "GoalAtRisk"
	DSOffline         DashboardState = "Offline"
	DSHistoricalView  DashboardState = "HistoricalView"
)

type DashboardMode string
const (DMOverview DashboardMode = "Overview"; DMDecision DashboardMode = "Decision"; DMGoal DashboardMode = "Goal"; DMPortfolio DashboardMode = "Portfolio"; DMPlanning DashboardMode = "Planning"; DMHistorical DashboardMode = "Historical")

type WidgetTier int
const (Tier1 WidgetTier = 1; Tier2 WidgetTier = 2; Tier3 WidgetTier = 3)

type Widget struct {
	WidgetID string `json:"widget_id"`
	WidgetType string `json:"widget_type"`
	Title string `json:"title"`
	Tier WidgetTier `json:"tier"`
	Priority int `json:"priority"`
	Visible bool `json:"visible"`
	Data interface{} `json:"data"`
	State string `json:"state"`
	Confidence string `json:"confidence"`
}

type DashboardOutput struct {
	DashboardID string `json:"dashboard_id"`
	Mode DashboardMode `json:"mode"`
	State DashboardState `json:"state"`
	Tier1 []Widget `json:"tier1"`
	Tier2 []Widget `json:"tier2"`
	Tier3 []Widget `json:"tier3"`
	CriticalAlert *Widget `json:"critical_alert,omitempty"`
	LastRefreshTime string `json:"last_refresh_time"`
}

type Inputs struct {
	UserID string `json:"user_id"`
	NetWorth int64 `json:"net_worth"`
	HealthScore int `json:"health_score"`
	HealthScoreGrade string `json:"health_score_grade"`
	RiskScore int `json:"risk_score"`
	RiskLevel string `json:"risk_level"`
	PortfolioValue int64 `json:"portfolio_value"`
	CashBalance int64 `json:"cash_balance"`
	TotalDebt int64 `json:"total_debt"`
	MonthlyIncome int64 `json:"monthly_income"`
	MonthlyExpenses int64 `json:"monthly_expenses"`
	SavingsRate float64 `json:"savings_rate"`
	EmergencyMonths float64 `json:"emergency_months"`
	BillsDue int64 `json:"bills_due"`
	GoalCount int `json:"goal_count"`
	GoalsOnTrack int `json:"goals_on_track"`
	GoalsAtRisk int `json:"goals_at_risk"`
	HasRecommendation bool `json:"has_recommendation"`
	IsFirstTimeUser bool `json:"is_first_time_user"`
	HasAccounts bool `json:"has_accounts"`
	HasEvents bool `json:"has_events"`
	HasGoals bool `json:"has_goals"`
	RecentEvents int `json:"recent_events"`
}

type Composer struct{}
func NewComposer() *Composer { return &Composer{} }

func (c *Composer) Compose(inputs Inputs, mode DashboardMode) *DashboardOutput {
	state := c.detectState(inputs)
	t1, t2, t3, alert := c.composeWidgets(inputs, state, mode)
	return &DashboardOutput{DashboardID: "dash-" + inputs.UserID, Mode: mode, State: state, Tier1: t1, Tier2: t2, Tier3: t3, CriticalAlert: alert}
}

func (c *Composer) detectState(inputs Inputs) DashboardState {
	if inputs.IsFirstTimeUser && !inputs.HasGoals { return DSFirstTimeUser }
	if inputs.HasAccounts && !inputs.HasEvents { return DSNoData }
	if inputs.GoalsAtRisk > 0 { return DSGoalAtRisk }
	if inputs.HealthScore < 40 || inputs.RiskScore >= 80 { return DSCritical }
	if inputs.HealthScore >= 60 && inputs.RiskScore < 60 { return DSHealthy }
	return DSAttentionNeeded
}

func (c *Composer) composeWidgets(inputs Inputs, state DashboardState, mode DashboardMode) ([]Widget, []Widget, []Widget, *Widget) {
	var t1, t2, t3 []Widget; var critical *Widget
	if state == DSCritical || inputs.RiskScore >= 80 { critical = &Widget{WidgetID: "alert-1", WidgetType: "critical_alert", Title: "Critical Alert", Tier: Tier1, Visible: true} }
	t1 = append(t1, Widget{WidgetID: "health", WidgetType: "health_score", Title: "Health Score: " + fmt.Sprint(inputs.HealthScore), Tier: Tier1, Priority: 1, Visible: true, Confidence: "High"})
	if inputs.GoalCount > 0 {
		t1 = append(t1, Widget{WidgetID: "goals", WidgetType: "goal_progress", Title: fmt.Sprintf("%d/%d goals on track", inputs.GoalsOnTrack, inputs.GoalCount), Tier: Tier1, Priority: 2, Visible: true})
	}
	if inputs.HasRecommendation {
		t1 = append(t1, Widget{WidgetID: "rec", WidgetType: "top_recommendation", Title: "Top Recommendation", Tier: Tier1, Priority: 3, Visible: true})
	}
	t2 = append(t2, Widget{WidgetID: "networth", WidgetType: "net_worth", Title: "Net Worth: " + fmtMoney(inputs.NetWorth), Tier: Tier2, Priority: 1, Visible: true})
	t2 = append(t2, Widget{WidgetID: "cash", WidgetType: "cash_position", Title: "Cash: " + fmtMoney(inputs.CashBalance), Tier: Tier2, Priority: 2, Visible: true})
	if inputs.PortfolioValue > 0 { t2 = append(t2, Widget{WidgetID: "portfolio", WidgetType: "portfolio_snapshot", Title: "Portfolio: " + fmtMoney(inputs.PortfolioValue), Tier: Tier2, Priority: 3, Visible: true}) }
	t2 = append(t2, Widget{WidgetID: "risk", WidgetType: "risk_summary", Title: "Risk: " + fmt.Sprint(inputs.RiskScore), Tier: Tier2, Priority: 4, Visible: true})
	if inputs.BillsDue > 0 { t2 = append(t2, Widget{WidgetID: "bills", WidgetType: "upcoming_bills", Title: "Bills due: " + fmtMoney(inputs.BillsDue), Tier: Tier2, Priority: 5, Visible: true}) }
	t3 = append(t3, Widget{WidgetID: "cashflow", WidgetType: "monthly_cash_flow", Title: "Cash Flow", Tier: Tier3, Visible: true})
	if inputs.RecentEvents > 0 { t3 = append(t3, Widget{WidgetID: "events", WidgetType: "recent_events", Title: fmt.Sprintf("%d recent events", inputs.RecentEvents), Tier: Tier3, Visible: true}) }
	return t1, t2, t3, critical
}

func fmtMoney(v int64) string { return fmt.Sprintf("₹%d", v) }
