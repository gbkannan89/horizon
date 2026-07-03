package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/horizon/core/services/ai/provider"
)

type InsightCategory string
const (
	ICOpportunity   InsightCategory = "Opportunity"
	ICWarning       InsightCategory = "Warning"
	ICAchievement   InsightCategory = "Achievement"
	ICMilestone     InsightCategory = "Milestone"
	ICForecast      InsightCategory = "Forecast"
	ICOptimization  InsightCategory = "Optimization"
	ICRecommend     InsightCategory = "Recommendation"
	ICRisk          InsightCategory = "Risk"
	ICHealth        InsightCategory = "Health"
	ICPortfolio     InsightCategory = "Portfolio"
	ICGoal          InsightCategory = "Goal"
	ICCashFlow      InsightCategory = "CashFlow"
	ICNetWorth      InsightCategory = "NetWorth"
	ICSavings       InsightCategory = "Savings"
	ICBehaviour     InsightCategory = "Behaviour"
	ICDebt          InsightCategory = "Debt"
)

type Priority string
const (PriorityCritical Priority = "Critical"; PriorityHigh Priority = "High"; PriorityMedium Priority = "Medium"; PriorityLow Priority = "Low"; PriorityInfo Priority = "Informational")

type Insight struct {
	InsightID     string          `json:"insight_id"`
	Category      InsightCategory `json:"category"`
	Priority      Priority        `json:"priority"`
	Title         string          `json:"title"`
	Summary       string          `json:"summary"`
	Description   string          `json:"description"`
	SourceEngine  string          `json:"source_engine"`
	RelatedEntity string          `json:"related_entity,omitempty"`
	RelatedGoal   string          `json:"related_goal,omitempty"`
	RelatedAcct   string          `json:"related_account,omitempty"`
	RelatedPf     string          `json:"related_portfolio,omitempty"`
	Timestamp     string          `json:"timestamp"`
	Severity      string          `json:"severity"`
	Confidence    string          `json:"confidence"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type InsightsDashboard struct {
	Total     int                `json:"total"`
	ByPriority map[Priority]int  `json:"by_priority"`
	ByCategory map[InsightCategory]int `json:"by_category"`
	Critical   []Insight         `json:"critical,omitempty"`
	High       []Insight         `json:"high,omitempty"`
	Others     []Insight         `json:"others,omitempty"`
}

type InsightsFeed struct {
	Insights []Insight `json:"insights"`
	Total    int       `json:"total"`
	Cursor   string    `json:"cursor,omitempty"`
	HasMore  bool      `json:"has_more"`
}

type Inputs struct {
	UserID           string  `json:"user_id"`
	HealthScore      int     `json:"health_score"`
	HealthGrade      string  `json:"health_grade"`
	HealthChange     int     `json:"health_change"`
	RiskScore        int     `json:"risk_score"`
	RiskLevel        string  `json:"risk_level"`
	NetWorth         int64   `json:"net_worth"`
	NetWorthChange   int64   `json:"net_worth_change"`
	SavingsRate      float64 `json:"savings_rate"`
	SavingsTrend     string  `json:"savings_trend"`
	CashFlowSurplus  int64   `json:"cash_flow_surplus"`
	GoalsOnTrack     int     `json:"goals_on_track"`
	TotalGoals       int     `json:"total_goals"`
	GoalProgressAvg  float64 `json:"goal_progress_avg"`
	PortfolioValue   int64   `json:"portfolio_value"`
	PortfolioReturn  float64 `json:"portfolio_return"`
	DebtUtilization  float64 `json:"debt_utilization"`
	AchievementCount int     `json:"achievement_count"`
	SpendingAnomaly  string  `json:"spending_anomaly"`
	MilestoneCount   int     `json:"milestone_count"`
	HasRecs          bool    `json:"has_recommendations"`
	RecCount         int     `json:"rec_count"`
	HasOpts          bool    `json:"has_optimizations"`
	HasSims          bool    `json:"has_simulations"`
	EventCount       int     `json:"event_count"`
}

type Composer struct{
	ai provider.AIProvider
}

func NewComposer(ai provider.AIProvider) *Composer { return &Composer{ai: ai} }

func (c *Composer) BuildDashboard(ctx context.Context, inputs Inputs) *InsightsDashboard {
	insights := c.generate(ctx, inputs)
	byPriority := map[Priority]int{}
	byCategory := map[InsightCategory]int{}
	var critical, high, others []Insight
	for _, ins := range insights {
		byPriority[ins.Priority]++
		byCategory[ins.Category]++
		switch ins.Priority {
		case PriorityCritical: critical = append(critical, ins)
		case PriorityHigh: high = append(high, ins)
		default: others = append(others, ins)
		}
	}
	return &InsightsDashboard{
		Total: len(insights), ByPriority: byPriority, ByCategory: byCategory,
		Critical: critical, High: high, Others: others,
	}
}

func (c *Composer) BuildFeed(ctx context.Context, inputs Inputs, max int, cursor string) *InsightsFeed {
	insights := c.generate(ctx, inputs)
	if max <= 0 { max = 20 }
	hasMore := len(insights) > max
	if len(insights) > max { insights = insights[:max] }
	nextCursor := ""
	if hasMore && len(insights) > 0 { nextCursor = insights[len(insights)-1].InsightID }
	return &InsightsFeed{Insights: insights, Total: len(insights), Cursor: nextCursor, HasMore: hasMore}
}

func (c *Composer) FilterByCategory(ctx context.Context, inputs Inputs, cat InsightCategory) []Insight {
	all := c.generate(ctx, inputs)
	var filtered []Insight
	for _, ins := range all { if ins.Category == cat { filtered = append(filtered, ins) } }
	return filtered
}

func (c *Composer) FilterByPriority(ctx context.Context, inputs Inputs, pri Priority) []Insight {
	all := c.generate(ctx, inputs)
	var filtered []Insight
	for _, ins := range all { if ins.Priority == pri { filtered = append(filtered, ins) } }
	return filtered
}

func (c *Composer) Search(ctx context.Context, inputs Inputs, q string) []Insight {
	all := c.generate(ctx, inputs)
	var matched []Insight
	q = strings.ToLower(q)
	for _, ins := range all {
		if strings.Contains(strings.ToLower(ins.Title), q) || strings.Contains(strings.ToLower(ins.Summary), q) || strings.Contains(strings.ToLower(ins.Description), q) {
			matched = append(matched, ins)
		}
	}
	return matched
}

func (c *Composer) GetByID(ctx context.Context, inputs Inputs, id string) *Insight {
	for _, ins := range c.generate(ctx, inputs) {
		if ins.InsightID == id { return &ins }
	}
	return nil
}

func (c *Composer) addMeta(ins Insight, inputs Inputs, sev string) Insight {
	if ins.Timestamp == "" { ins.Timestamp = time.Now().UTC().Format(time.RFC3339) }
	if ins.Severity == "" { ins.Severity = sev }
	if ins.Confidence == "" { ins.Confidence = "Medium" }
	return ins
}

func (c *Composer) generate(ctx context.Context, inputs Inputs) []Insight {
	now := time.Now().UTC().Format(time.RFC3339)
	var insights []Insight

	add := func(ins Insight, sev string) {
		if ins.Timestamp == "" { ins.Timestamp = now }
		if ins.Severity == "" { ins.Severity = sev }
		if ins.Confidence == "" { ins.Confidence = "Medium" }
		if ins.Metadata == nil { ins.Metadata = map[string]interface{}{} }
		insights = append(insights, ins)
	}

	if inputs.HealthScore > 0 {
		if inputs.HealthChange > 0 {
			add(Insight{InsightID: "ins-health-up", Category: ICHealth, Title: "Health Score Improving", Priority: PriorityLow, SourceEngine: "HealthScore", Confidence: "High",
				Summary: fmt.Sprintf("Up %d points to %d (%s)", inputs.HealthChange, inputs.HealthScore, inputs.HealthGrade),
				Metadata: map[string]interface{}{"score": inputs.HealthScore, "change": inputs.HealthChange, "grade": inputs.HealthGrade}}, "positive")
		} else if inputs.HealthChange < 0 {
			add(Insight{InsightID: "ins-health-down", Category: ICHealth, Title: "Health Score Declined", Priority: PriorityHigh, SourceEngine: "HealthScore", Confidence: "High",
				Summary: fmt.Sprintf("Down %d points to %d (%s)", -inputs.HealthChange, inputs.HealthScore, inputs.HealthGrade),
				Metadata: map[string]interface{}{"score": inputs.HealthScore, "change": inputs.HealthChange, "grade": inputs.HealthGrade}}, "warning")
		}
		if inputs.HealthScore < 40 {
			add(Insight{InsightID: "ins-health-critical", Category: ICHealth, Title: "Critical Health Score", Priority: PriorityCritical, SourceEngine: "HealthScore", Confidence: "High",
				Summary: fmt.Sprintf("Score %d — immediate attention needed", inputs.HealthScore),
				Metadata: map[string]interface{}{"score": inputs.HealthScore}}, "critical")
		}
	}

	if inputs.RiskScore >= 80 {
		add(Insight{InsightID: "ins-risk-critical", Category: ICRisk, Title: "Critical Risk Level", Priority: PriorityCritical, SourceEngine: "Risk", Confidence: "High",
			Summary: fmt.Sprintf("Risk score %d — %s", inputs.RiskScore, inputs.RiskLevel),
			Metadata: map[string]interface{}{"score": inputs.RiskScore, "level": inputs.RiskLevel}}, "critical")
	} else if inputs.RiskScore >= 60 {
		add(Insight{InsightID: "ins-risk-high", Category: ICRisk, Title: "Elevated Risk", Priority: PriorityHigh, SourceEngine: "Risk", Confidence: "High",
			Summary: fmt.Sprintf("Risk score %d — review recommended", inputs.RiskScore),
			Metadata: map[string]interface{}{"score": inputs.RiskScore}}, "warning")
	} else if inputs.RiskScore < 30 {
		add(Insight{InsightID: "ins-risk-low", Category: ICRisk, Title: "Low Risk Profile", Priority: PriorityLow, SourceEngine: "Risk", Confidence: "High",
			Summary: fmt.Sprintf("Risk score %d — %s", inputs.RiskScore, inputs.RiskLevel),
			Metadata: map[string]interface{}{"score": inputs.RiskScore, "level": inputs.RiskLevel}}, "positive")
	}

	if inputs.SavingsRate > 0 {
		dir := "stable"
		if inputs.SavingsTrend == "improving" { dir = "improving" } else if inputs.SavingsTrend == "declining" { dir = "declining" }
		add(Insight{InsightID: "ins-savings", Category: ICSavings, Title: "Savings Rate Trend", Priority: PriorityMedium, SourceEngine: "Projection", Confidence: "High",
			Summary: fmt.Sprintf("%.0f%% savings rate — %s", inputs.SavingsRate, dir),
			Metadata: map[string]interface{}{"rate": inputs.SavingsRate, "trend": inputs.SavingsTrend}}, dir)
	}

	if inputs.CashFlowSurplus > 0 {
		add(Insight{InsightID: "ins-cf-surplus", Category: ICCashFlow, Title: "Positive Cash Flow", Priority: PriorityLow, SourceEngine: "Projection", Confidence: "High",
			Summary: fmt.Sprintf("Monthly surplus of %s", fmtMoney(inputs.CashFlowSurplus)),
			Metadata: map[string]interface{}{"surplus": inputs.CashFlowSurplus}}, "positive")
	} else if inputs.CashFlowSurplus < 0 {
		add(Insight{InsightID: "ins-cf-deficit", Category: ICCashFlow, Title: "Cash Flow Deficit", Priority: PriorityHigh, SourceEngine: "Projection", Confidence: "High",
			Summary: fmt.Sprintf("Monthly shortfall of %s", fmtMoney(-inputs.CashFlowSurplus)),
			Metadata: map[string]interface{}{"shortfall": -inputs.CashFlowSurplus}}, "warning")
	}

	if inputs.TotalGoals > 0 {
		ratio := float64(inputs.GoalsOnTrack) / float64(inputs.TotalGoals)
		if ratio >= 1.0 {
			add(Insight{InsightID: "ins-goals-all", Category: ICGoal, Title: "All Goals On Track", Priority: PriorityLow, SourceEngine: "Projection", Confidence: "High",
				Summary: fmt.Sprintf("%d/%d goals progressing well", inputs.GoalsOnTrack, inputs.TotalGoals),
				RelatedGoal: "all", Metadata: map[string]interface{}{"on_track": inputs.GoalsOnTrack, "total": inputs.TotalGoals}}, "achievement")
		} else if ratio < 0.5 {
			add(Insight{InsightID: "ins-goals-risk", Category: ICGoal, Title: "Goals Need Attention", Priority: PriorityHigh, SourceEngine: "Projection", Confidence: "High",
				Summary: fmt.Sprintf("Only %d/%d goals on track", inputs.GoalsOnTrack, inputs.TotalGoals),
				RelatedGoal: "all", Metadata: map[string]interface{}{"on_track": inputs.GoalsOnTrack, "total": inputs.TotalGoals}}, "warning")
		}
	}

	if inputs.NetWorthChange > 0 {
		add(Insight{InsightID: "ins-nw-up", Category: ICNetWorth, Title: "Net Worth Growing", Priority: PriorityLow, SourceEngine: "Projection", Confidence: "High",
			Summary: fmt.Sprintf("Increased by %s", fmtMoney(inputs.NetWorthChange)),
			Metadata: map[string]interface{}{"change": inputs.NetWorthChange, "net_worth": inputs.NetWorth}}, "positive")
	} else if inputs.NetWorthChange < 0 {
		add(Insight{InsightID: "ins-nw-down", Category: ICNetWorth, Title: "Net Worth Declined", Priority: PriorityMedium, SourceEngine: "Projection", Confidence: "High",
			Summary: fmt.Sprintf("Decreased by %s", fmtMoney(-inputs.NetWorthChange)),
			Metadata: map[string]interface{}{"change": inputs.NetWorthChange, "net_worth": inputs.NetWorth}}, "warning")
	}

	if inputs.PortfolioReturn > 0 {
		add(Insight{InsightID: "ins-pf-up", Category: ICPortfolio, Title: "Portfolio Performing Well", Priority: PriorityLow, SourceEngine: "Projection", Confidence: "Medium",
			Summary: fmt.Sprintf("Return: %.1f%%", inputs.PortfolioReturn),
			Metadata: map[string]interface{}{"return": inputs.PortfolioReturn, "value": inputs.PortfolioValue}}, "positive")
	}

	if inputs.AchievementCount > 0 {
		add(Insight{InsightID: "ins-ach", Category: ICAchievement, Title: "Achievements Unlocked", Priority: PriorityLow, SourceEngine: "HealthScore", Confidence: "High",
			Summary: fmt.Sprintf("%d achievement(s)", inputs.AchievementCount),
			Metadata: map[string]interface{}{"count": inputs.AchievementCount}}, "achievement")
	}

	if inputs.SpendingAnomaly != "" {
		add(Insight{InsightID: "ins-anomaly", Category: ICBehaviour, Title: "Spending Pattern Change", Priority: PriorityMedium, SourceEngine: "FinancialEvent", Confidence: "Medium",
			Summary: inputs.SpendingAnomaly}, "info")
	}

	if inputs.HasRecs {
		add(Insight{InsightID: "ins-recs", Category: ICRecommend, Title: "Recommendations Available", Priority: PriorityMedium, SourceEngine: "Recommendation", Confidence: "High",
			Summary: fmt.Sprintf("%d recommendation(s) ready for review", inputs.RecCount),
			Metadata: map[string]interface{}{"count": inputs.RecCount}}, "info")
	}

	if inputs.HasOpts {
		add(Insight{InsightID: "ins-opts", Category: ICOptimization, Title: "Optimization Strategies Available", Priority: PriorityMedium, SourceEngine: "Optimization", Confidence: "High",
			Summary: "Review optimization strategies to improve outcomes"}, "info")
	}

	if inputs.MilestoneCount > 0 {
		add(Insight{InsightID: "ins-ms", Category: ICMilestone, Title: "Milestones Reached", Priority: PriorityLow, SourceEngine: "Goal", Confidence: "High",
			Summary: fmt.Sprintf("%d milestone(s)", inputs.MilestoneCount),
			Metadata: map[string]interface{}{"count": inputs.MilestoneCount}}, "achievement")
	}

	if c.ai != nil {
		data := map[string]interface{}{
			"health_score": inputs.HealthScore,
			"savings_rate": inputs.SavingsRate,
			"net_worth":    inputs.NetWorth,
			"spending_anomaly": inputs.SpendingAnomaly,
		}
		resp, err := c.ai.GenerateInsights(ctx, provider.GenerateInsightsRequest{Data: data})
		if err == nil && resp != nil {
			for i, aiInsight := range resp.Insights {
				cat := ICBehaviour
				if aiInsight.Type == "savings_opportunity" { cat = ICSavings }
				add(Insight{
					InsightID:    fmt.Sprintf("ins-ai-%d", i),
					Category:     cat,
					Title:        aiInsight.Title,
					Priority:     PriorityMedium,
					SourceEngine: "AI",
					Confidence:   aiInsight.Confidence,
					Summary:      aiInsight.Summary,
					Description:  aiInsight.Explanation,
					Metadata:     map[string]interface{}{"ai_generated": true},
				}, "info")
			}
		}
	}

	return insights
}

func fmtMoney(v int64) string {
	if v >= 10000000 { return fmt.Sprintf("₹%.2fCr", float64(v)/10000000) }
	if v >= 100000 { return fmt.Sprintf("₹%.2fL", float64(v)/100000) }
	return fmt.Sprintf("₹%d", v)
}
