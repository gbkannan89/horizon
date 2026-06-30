package engine

import "fmt"

type AdvisorMode string
const (AMOverview AdvisorMode = "Overview"; AMDecision AdvisorMode = "Decision"; AMPlanning AdvisorMode = "Planning"; AMLearning AdvisorMode = "Learning"; AMGoalCoaching AdvisorMode = "GoalCoaching"; AMRiskReview AdvisorMode = "RiskReview"; AMRetirementPlan AdvisorMode = "RetirementPlanning"; AMPortfolioCoach AdvisorMode = "PortfolioCoaching"; AMHistoricalReview AdvisorMode = "HistoricalReview")

type CardType string
const (CTTodayBestAction CardType = "TodayBestAction"; CTGoalCoach CardType = "GoalCoach"; CTHealthSummary CardType = "HealthSummary"; CTRiskAlert CardType = "RiskAlert"; CTProjectionSummary CardType = "ProjectionSummary"; CTSimulationResult CardType = "SimulationResult"; CTOptimizationStr CardType = "OptimizationStrategy"; CTFinancialInsight CardType = "FinancialInsight"; CTEducationCard CardType = "EducationCard"; CTAchievementCard CardType = "AchievementCard")

type AdvisorCard struct {
	CardID string `json:"card_id"`; CardType CardType `json:"card_type"`
	Title string `json:"title"`; Summary string `json:"summary"`
	Detail string `json:"detail"`; Priority int `json:"priority"`
	Confidence string `json:"confidence"`; IsAIGenerated bool `json:"is_ai_generated"`
	Actions []string `json:"actions"`
}

type AdvisorOutput struct {
	SessionID string `json:"session_id"`; Mode AdvisorMode `json:"mode"`
	Cards []AdvisorCard `json:"cards"`
}

type Inputs struct {
	UserID string `json:"user_id"`; HealthScore int `json:"health_score"`
	HealthGrade string `json:"health_grade"`; RiskScore int `json:"risk_score"`
	RiskLevel string `json:"risk_level"`; HasRecommendation bool `json:"has_recommendation"`
	RecTitle string `json:"rec_title"`; RecExpected string `json:"rec_expected"`
	HasSimulation bool `json:"has_simulation"`; GoalCount int `json:"goal_count"`
	GoalsOnTrack int `json:"goals_on_track"`; GoalsAtRisk int `json:"goals_at_risk"`
	ProjectionSummary string `json:"projection_summary"`
	HasOptimization bool `json:"has_optimization"`; HasAchievement bool `json:"has_achievement"`
}

type Engine struct{}
func NewEngine() *Engine { return &Engine{} }

func (e *Engine) Execute(inputs Inputs, mode AdvisorMode) *AdvisorOutput {
	if mode == "" { mode = AMOverview }
	return &AdvisorOutput{SessionID: fmt.Sprintf("adv-%s", inputs.UserID), Mode: mode, Cards: e.buildCards(inputs, mode)}
}

func (e *Engine) buildCards(inputs Inputs, mode AdvisorMode) []AdvisorCard {
	var cards []AdvisorCard
	if inputs.HasRecommendation {
		cards = append(cards, AdvisorCard{
			CardID: "c-1", CardType: CTTodayBestAction, Title: inputs.RecTitle,
			Summary: inputs.RecExpected, Priority: 1, Confidence: "High", Actions: []string{"Tell me more", "Accept"},
		})
	}
	cards = append(cards, AdvisorCard{
		CardID: "c-2", CardType: CTHealthSummary, Title: fmt.Sprintf("Health: %d (%s)", inputs.HealthScore, inputs.HealthGrade),
		Summary: "Your financial health overview", Priority: 2, Confidence: "High",
	})
	if inputs.RiskScore >= 60 {
		cards = append(cards, AdvisorCard{
			CardID: "c-3", CardType: CTRiskAlert, Title: fmt.Sprintf("Risk: %d (%s)", inputs.RiskScore, inputs.RiskLevel),
			Summary: "Review your risk exposure", Priority: 3, Confidence: "High",
		})
	}
	return cards
}
