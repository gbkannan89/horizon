package prompts

import "fmt"

// Template represents a versioned prompt template.
type Template struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	Text    string `json:"text"`
}

// Manager holds all prompt templates.
type Manager struct {
	templates map[string]map[int]Template // type -> version -> template
}

func NewManager() *Manager {
	m := &Manager{templates: make(map[string]map[int]Template)}
	m.registerDefaults()
	return m
}

func (m *Manager) Get(promptType string, version int) (Template, bool) {
	versions, ok := m.templates[promptType]
	if !ok { return Template{}, false }
	t, ok := versions[version]
	return t, ok
}

func (m *Manager) LatestVersion(promptType string) int {
	versions, ok := m.templates[promptType]
	if !ok { return 1 }
	maxV := 1
	for v := range versions { if v > maxV { maxV = v } }
	return maxV
}

func (m *Manager) List() []TemplateSummary {
	var summaries []TemplateSummary
	for t, versions := range m.templates {
		maxV := 0
		for v := range versions { if v > maxV { maxV = v } }
		summaries = append(summaries, TemplateSummary{Type: t, Version: maxV, AvailableVersions: len(versions)})
	}
	return summaries
}

type TemplateSummary struct {
	Type             string `json:"type"`
	Version          int    `json:"version"`
	AvailableVersions int   `json:"available_versions"`
}

// Render substitutes placeholders in a template with data values.
func Render(t Template, data map[string]interface{}) string {
	rendered := t.Text
	for k, v := range data {
		placeholder := fmt.Sprintf("{{%s}}", k)
		val := fmt.Sprintf("%v", v)
		rendered = replaceAll(rendered, placeholder, val)
	}
	return rendered
}

func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := indexOf(s, old)
		if idx < 0 { break }
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result + s
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr { return i }
	}
	return -1
}

func (m *Manager) registerDefaults() {
	versions := map[string]map[int]string{
		"explain_health": {
			1: "The user's financial health score is {{score}} out of 100 (grade: {{grade}}). This represents a change of {{change}} points. The score is calculated from 13 dimensions including savings rate, debt management, emergency fund coverage, insurance adequacy, investment diversification, and goal progress tracking. Explain what this score means in simple terms and what the user can do to improve it.",
			2: "Financial Health Score: {{score}}/100 ({{grade}}). Trend: {{change}} points. This comprehensive score evaluates 13 financial health dimensions. Provide a personalized interpretation of this score, highlighting the top 2 strengths and top 2 weaknesses based on the available data. Suggest specific, actionable steps for improvement.",
		},
		"explain_risk": {
			1: "The user's composite risk score is {{score}} out of 100, classified as {{level}}. Explain what this risk level means for their financial plan, which risk indicators are most relevant, and what actions could help reduce their risk exposure.",
		},
		"explain_goal": {
			1: "The goal '{{name}}' has {{progress}}% progress and is currently {{status}}. The target amount is {{target}} with current savings of {{current}}. Monthly contribution is {{monthly}}. Explain the goal status and what factors will affect timely completion.",
		},
		"explain_portfolio": {
			1: "Portfolio value: {{value}}. Period return: {{return}}%. Allocation: {{allocation}}. Provide a summary of portfolio performance and rebalancing suggestions if applicable.",
		},
		"explain_recommendation": {
			1: "The top recommendation is: '{{title}}'. There are {{count}} active recommendation(s) in total. Explain why this recommendation is prioritised, how it connects to the user's goals, and what the expected impact would be.",
		},
		"explain_optimization": {
			1: "There are {{count}} optimization strategies available. Explain what optimization means in a financial planning context and how these strategies can help improve outcomes.",
		},
		"explain_projection": {
			1: "Projected net worth: {{net_worth_projected}}. On track: {{on_track}}. Annual return assumption: {{annual_return}}%. Inflation: {{inflation}}%. Explain the projection methodology and key assumptions.",
		},
		"explain_simulation": {
			1: "There are {{count}} simulation scenario(s) available. Explain how simulations help with financial planning and what types of 'what if' scenarios can be explored.",
		},
		"financial_summary": {
			1: "Net worth: {{net_worth}}. Income: {{income}}. Expenses: {{expenses}}. Provide a brief financial summary highlighting the most important metrics and any areas that need attention.",
		},
		"advisor_summary": {
			1: "Health Score: {{health_score}}/100. Risk Score: {{risk_score}}/100. Goals on track: {{goals_on_track}}/{{total_goals}}. Provide a comprehensive financial advisor summary covering all aspects of the user's financial life.",
		},
		"insight_spending": {
			1: "Analyze the spending pattern anomaly: {{anomaly}}. Provide a personalized insight with actionable advice to help the user manage their discretionary spending better.",
		},
		"insight_saving": {
			1: "The user has a current savings rate of {{savings_rate}}% and uninvested surplus cash. Explain this savings opportunity and suggest specific actions to optimize their cash flow for better long-term returns.",
		},
		"insight_risk": {
			1: "The user's recent activity indicates a change in risk exposure. Identify key risk indicators from the provided data and explain how to mitigate potential negative impacts on their financial plan.",
		},
		"insight_goal": {
			1: "Based on the current trajectory, provide a deep dive into the goal progress for the user. Highlight any risks of falling short and offer concrete steps to get back on track.",
		},
	}

	for ptype, vers := range versions {
		m.templates[ptype] = make(map[int]Template)
		for v, text := range vers {
			m.templates[ptype][v] = Template{Type: ptype, Version: v, Text: text}
		}
	}
}
