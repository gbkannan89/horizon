package engine

// InsightCategory represents the category of insight.
type InsightCategory string
const (
	ICBehaviour   InsightCategory = "Behaviour"
	ICCashFlow    InsightCategory = "CashFlow"
	ICSavings     InsightCategory = "Savings"
	ICGoals       InsightCategory = "Goals"
	ICPortfolio   InsightCategory = "Portfolio"
	ICRisk        InsightCategory = "Risk"
	ICDebt        InsightCategory = "Debt"
	ICHealthScore InsightCategory = "HealthScore"
	ICNetWorth    InsightCategory = "NetWorth"
	ICInvestments InsightCategory = "Investments"
	ICLifestyle   InsightCategory = "Lifestyle"
	ICLifeEvents  InsightCategory = "LifeEvents"
	ICAchievements InsightCategory = "Achievements"
)

// InsightType represents the type of insight card.
type InsightType string
const (
	ITPattern     InsightType = "Pattern"
	ITTrend       InsightType = "Trend"
	ITOpportunity InsightType = "Opportunity"
	ITBehaviour   InsightType = "Behaviour"
	ITAchievement InsightType = "Achievement"
	ITWarning     InsightType = "Warning"
)

// InsightCard represents a single insight.
type InsightCard struct {
	InsightID     string          `json:"insight_id"`
	Category      InsightCategory `json:"category"`
	Type          InsightType     `json:"type"`
	Title         string          `json:"title"`
	Summary       string          `json:"summary"`
	Detail        string          `json:"detail"`
	Confidence    string          `json:"confidence"`
	IsAIGenerated bool            `json:"is_ai_generated"`
	Priority      int             `json:"priority"`
}

// InsightsFeed is the complete insights workspace view.
type InsightsFeed struct {
	Insights     []InsightCard `json:"insights"`
	Total        int           `json:"total"`
	View         string        `json:"view"` // today, weekly, monthly, quarterly, yearly
}

// Inputs for the insights experience.
type Inputs struct {
	UserID           string       `json:"user_id"`
	SavingsRate      float64      `json:"savings_rate"`
	SavingsTrend     string       `json:"savings_trend"` // improving, declining, stable
	CashFlowSurplus  int64        `json:"cash_flow_surplus"`
	GoalProgressAvg  float64      `json:"goal_progress_avg"`
	GoalsOnTrack     int          `json:"goals_on_track"`
	TotalGoals       int          `json:"total_goals"`
	RiskScore        int          `json:"risk_score"`
	DebtUtilization  float64      `json:"debt_utilization"`
	PortfolioValue   int64        `json:"portfolio_value"`
	NetWorth         int64        `json:"net_worth"`
	NetWorthChange   int64        `json:"net_worth_change"`
	HealthScore      int          `json:"health_score"`
	HealthChange     int          `json:"health_change"`
	AchievementCount int          `json:"achievement_count"`
	SpendingAnomaly  string       `json:"spending_anomaly,omitempty"`
	HasLifeEvent     bool         `json:"has_life_event"`
	View             string       `json:"view"`
}
