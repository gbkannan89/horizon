package provider

import (
	"context"
	"fmt"
	"strings"
)

// MockProvider returns deterministic placeholder responses.
// No network calls, no model inference, no external dependencies.
type MockProvider struct {
	Name string
}

func NewMock() *MockProvider {
	return &MockProvider{Name: "mock"}
}

func (m *MockProvider) Explain(ctx context.Context, req ExplainRequest) (*ExplainResponse, error) {
	explanation := m.generateExplanation(req.PromptType, req.Data)
	return &ExplainResponse{
		Explanation: explanation,
		Confidence:  "Medium",
		Provider:    m.Name,
		Version:     req.Version,
	}, nil
}

func (m *MockProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	reply := m.generateReply(req.Message)
	return &ChatResponse{
		Reply:      reply,
		Confidence: "Medium",
		Provider:   m.Name,
		SessionID:  req.SessionID,
	}, nil
}

func (m *MockProvider) Summarize(ctx context.Context, req SummarizeRequest) (*SummarizeResponse, error) {
	return &SummarizeResponse{
		Summary:    fmt.Sprintf("[Mock] Summary of '%s' based on provided financial data. When connected to a real AI provider, this would generate a comprehensive summary.", req.Topic),
		Confidence: "Medium",
		Provider:   m.Name,
	}, nil
}

func (m *MockProvider) Capabilities() Capabilities {
	return Capabilities{
		Provider: m.Name, Chat: true, Explanation: true, Summarize: true,
		Streaming: false, ToolCalling: false, Embeddings: false, Vision: false,
		MaxContext: 4096, Models: []string{"mock-1"},
	}
}

func (m *MockProvider) Health(ctx context.Context) (*HealthResponse, error) {
	return &HealthResponse{
		Status:   "operational",
		Provider: m.Name,
		Message:  "Mock AI provider — no external dependencies. Ready for real provider replacement.",
	}, nil
}

func (m *MockProvider) generateExplanation(promptType string, data map[string]interface{}) string {
	switch promptType {
	case "explain_health":
		return m.explainHealth(data)
	case "explain_risk":
		return m.explainRisk(data)
	case "explain_goal":
		return m.explainGoal(data)
	case "explain_portfolio":
		return m.explainPortfolio(data)
	case "explain_recommendation":
		return m.explainRecommendation(data)
	case "explain_optimization":
		return m.explainOptimization(data)
	case "explain_projection":
		return m.explainProjection(data)
	case "explain_simulation":
		return m.explainSimulation(data)
	case "financial_summary":
		return m.financialSummary(data)
	case "advisor_summary":
		return m.advisorSummary(data)
	default:
		return fmt.Sprintf("[Mock] Explanation for %s is not yet implemented. When connected to a real AI provider, this will provide a detailed analysis based on your financial data.", promptType)
	}
}

func (m *MockProvider) explainHealth(data map[string]interface{}) string {
	score, _ := data["score"].(int)
	grade, _ := data["grade"].(string)
	change, _ := data["change"].(int)
	trend := "stable"
	if change > 0 { trend = "improving" } else if change < 0 { trend = "declining" }
	return fmt.Sprintf("Your financial health score is %d (%s), which is %s. This score is calculated from 13 dimensions including savings rate, debt levels, emergency fund coverage, and goal progress. %s.",
		score, grade, trend, m.trendAdvice(change, "health"))
}

func (m *MockProvider) explainRisk(data map[string]interface{}) string {
	score, _ := data["score"].(int)
	level, _ := data["level"].(string)
	return fmt.Sprintf("Your composite risk score is %d, which is classified as %s. This assessment evaluates 21 indicators across market exposure, concentration risk, liquidity risk, and debt risk. %s",
		score, level, m.riskAdvice(score))
}

func (m *MockProvider) explainGoal(data map[string]interface{}) string {
	name, _ := data["name"].(string)
	progress, _ := data["progress"].(float64)
	onTrack, _ := data["on_track"].(bool)
	status := "on track"
	if !onTrack { status = "needs attention" }
	return fmt.Sprintf("Your goal '%s' is at %.0f%% progress and is %s. This projection is based on your current contribution rate of ₹%.0f per month and expected return of %.1f%%.",
		name, progress, status, data["monthly"], data["return"])
}

func (m *MockProvider) explainPortfolio(data map[string]interface{}) string {
	val, _ := data["value"].(float64)
	ret, _ := data["return"].(float64)
	return fmt.Sprintf("Your portfolio is valued at ₹%.0f with a period return of %.1f%%. The allocation is diversified across multiple asset classes. Regular rebalancing helps maintain your target allocation.",
		val, ret)
}

func (m *MockProvider) explainRecommendation(data map[string]interface{}) string {
	title, _ := data["title"].(string)
	count, _ := data["count"].(int)
	if count > 0 {
		return fmt.Sprintf("Top recommendation: '%s'. This recommendation is selected based on your financial goals, current progress, and projected outcomes. Reviewing and acting on recommendations can improve your financial health score.", title)
	}
	return "No active recommendations at this time. Your financial plan is on track."
}

func (m *MockProvider) explainOptimization(data map[string]interface{}) string {
	count, _ := data["count"].(int)
	if count > 0 {
		return fmt.Sprintf("There are %d optimization strategies available for your financial plan. These strategies evaluate trade-offs between contribution rates, asset allocation, and goal prioritization to suggest the most efficient path forward.", count)
	}
	return "No optimization strategies are currently available."
}

func (m *MockProvider) explainProjection(data map[string]interface{}) string {
	nw, _ := data["net_worth_projected"].(float64)
	onTrack, _ := data["on_track"].(bool)
	status := "on track"
	if !onTrack { status = "off track" }
	return fmt.Sprintf("Your projected net worth is ₹%.0f, and your financial plan is %s. This projection assumes an annual return of %.1f%% and inflation rate of %.1f%%. Actual results may vary.",
		nw, status, data["annual_return"], data["inflation"])
}

func (m *MockProvider) explainSimulation(data map[string]interface{}) string {
	count, _ := data["count"].(int)
	if count > 0 {
		return fmt.Sprintf("You have %d simulation scenarios available. Simulations let you explore 'what if' scenarios by adjusting contribution rates, retirement age, or investment strategy to see how different choices affect your financial outcomes.", count)
	}
	return "No simulations are currently available."
}

func (m *MockProvider) financialSummary(data map[string]interface{}) string {
	nw, _ := data["net_worth"].(float64)
	income, _ := data["income"].(float64)
	expenses, _ := data["expenses"].(float64)
	savings := income - expenses
	return fmt.Sprintf("Your current net worth is ₹%.0f. Your monthly income is ₹%.0f against expenses of ₹%.0f, giving you a monthly surplus of ₹%.0f. Your savings rate is approximately %.0f%%.",
		nw, income, expenses, savings, m.savingsRate(income, expenses))
}

func (m *MockProvider) advisorSummary(data map[string]interface{}) string {
	health, _ := data["health_score"].(int)
	risk, _ := data["risk_score"].(int)
	goals, _ := data["goals_on_track"].(int)
	totalGoals, _ := data["total_goals"].(int)
	return fmt.Sprintf("Here is your financial overview: Health Score %d/100, Risk Score %d/100, %d/%d goals on track. Your financial plan is being monitored and we will alert you if any metrics require attention.",
		health, risk, goals, totalGoals)
}

func (m *MockProvider) trendAdvice(change int, metric string) string {
	if change > 0 { return "Keep up the good financial habits that are driving this improvement." }
	if change < 0 { return "Review the declining areas in your financial profile and consider adjustments." }
	return "Maintain your current financial habits to keep your score stable."
}

func (m *MockProvider) riskAdvice(score int) string {
	if score >= 80 { return "Immediate attention is recommended to address high-risk indicators." }
	if score >= 60 { return "Review the identified risk factors and consider mitigation strategies." }
	return "Your risk profile is within a healthy range. Continue monitoring."
}

func (m *MockProvider) savingsRate(income, expenses float64) float64 {
	if income <= 0 { return 0 }
	return (income - expenses) / income * 100
}

func (m *MockProvider) generateReply(msg string) string {
	msg = strings.ToLower(msg)
	switch {
	case strings.Contains(msg, "health"):
		return "[Mock] Your financial health is calculated from 13 dimensions. Your Health Score Engine provides the exact score. I can explain what each dimension means if you'd like."
	case strings.Contains(msg, "risk"):
		return "[Mock] Your risk assessment evaluates 21 indicators across multiple categories. The Risk Engine produces a composite score. Would you like a detailed breakdown?"
	case strings.Contains(msg, "goal") || strings.Contains(msg, "goal"):
		return "[Mock] Your goals are managed by the Goal domain and evaluated by the Projection Engine. I can explain how your progress is calculated and what factors affect your goal timeline."
	case strings.Contains(msg, "hello") || strings.Contains(msg, "hi"):
		return "[Mock] Hello! I'm your Horizon financial advisor. I can help explain your health score, risk assessment, goal progress, portfolio performance, and recommendations. What would you like to know?"
	default:
		return fmt.Sprintf("[Mock] I understand you're asking about '%s'. When connected to a real AI provider, I would provide a detailed, personalized response based on your financial data. For now, please explore the deterministic dashboards for accurate information.", msg)
	}
}
