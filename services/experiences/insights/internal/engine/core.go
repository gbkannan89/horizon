package engine

import "fmt"

// Engine generates deterministic insights from engine outputs.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// BuildFeed composes the insights feed.
func (e *Engine) BuildFeed(inputs Inputs, maxInsights int) *InsightsFeed {
	var cards []InsightCard

	// Savings insight
	if inputs.SavingsRate > 0 {
		dir := "stable"
		if inputs.SavingsTrend == "improving" { dir = "improving" } else if inputs.SavingsTrend == "declining" { dir = "declining" }
		cards = append(cards, InsightCard{
			InsightID: "ins-savings", Category: ICSavings, Type: ITTrend,
			Title: "Savings Rate Trend",
			Summary: fmt.Sprintf("Your savings rate is %.0f%% and %s", inputs.SavingsRate, dir),
			Confidence: "High", Priority: 1,
		})
	}

	// Cash flow insight
	if inputs.CashFlowSurplus > 0 {
		cards = append(cards, InsightCard{
			InsightID: "ins-cf", Category: ICCashFlow, Type: ITTrend,
			Title: "Positive Cash Flow",
			Summary: fmt.Sprintf("Monthly surplus of %s available", fmtMoney(inputs.CashFlowSurplus)),
			Confidence: "High", Priority: 2,
		})
	}

	// Goal progress insight
	if inputs.TotalGoals > 0 {
		status := "on track"
		it := ITAchievement
		if float64(inputs.GoalsOnTrack)/float64(inputs.TotalGoals) < 0.5 {
			status = "needs attention"
			it = ITWarning
		}
		cards = append(cards, InsightCard{
			InsightID: "ins-goals", Category: ICGoals, Type: it,
			Title: "Goal Progress",
			Summary: fmt.Sprintf("%d/%d goals %s", inputs.GoalsOnTrack, inputs.TotalGoals, status),
			Confidence: "High", Priority: 3,
		})
	}

	// Health score insight
	if inputs.HealthChange > 0 {
		cards = append(cards, InsightCard{
			InsightID: "ins-health", Category: ICHealthScore, Type: ITTrend,
			Title: "Health Score Improved",
			Summary: fmt.Sprintf("Up %d points this period", inputs.HealthChange),
			Confidence: "High", Priority: 4,
		})
	} else if inputs.HealthChange < 0 {
		cards = append(cards, InsightCard{
			InsightID: "ins-health-d", Category: ICHealthScore, Type: ITWarning,
			Title: "Health Score Declined",
			Summary: fmt.Sprintf("Down %d points this period", -inputs.HealthChange),
			Confidence: "High", Priority: 4,
		})
	}

	// Net worth insight
	if inputs.NetWorthChange > 0 {
		cards = append(cards, InsightCard{
			InsightID: "ins-nw", Category: ICNetWorth, Type: ITTrend,
			Title: "Net Worth Growing",
			Summary: fmt.Sprintf("Increased by %s this period", fmtMoney(inputs.NetWorthChange)),
			Confidence: "High", Priority: 5,
		})
	}

	// Achievements insight
	if inputs.AchievementCount > 0 {
		cards = append(cards, InsightCard{
			InsightID: "ins-ach", Category: ICAchievements, Type: ITAchievement,
			Title: fmt.Sprintf("%d Achievements Unlocked", inputs.AchievementCount),
			Summary: "Keep up the good financial habits!",
			Confidence: "High", Priority: 6,
		})
	}

	// Spending anomaly
	if inputs.SpendingAnomaly != "" {
		cards = append(cards, InsightCard{
			InsightID: "ins-anomaly", Category: ICBehaviour, Type: ITWarning,
			Title: "Spending Pattern Change",
			Summary: inputs.SpendingAnomaly,
			Confidence: "Medium", Priority: 7,
		})
	}

	// Risk insight
	if inputs.RiskScore >= 60 {
		cards = append(cards, InsightCard{
			InsightID: "ins-risk", Category: ICRisk, Type: ITWarning,
			Title: fmt.Sprintf("Risk Score: %d", inputs.RiskScore),
			Summary: "Your risk indicators need attention",
			Confidence: "High", Priority: 8,
		})
	}

	if maxInsights <= 0 { maxInsights = 10 }
	if len(cards) > maxInsights { cards = cards[:maxInsights] }

	return &InsightsFeed{
		Insights: cards,
		Total:    len(cards),
		View:     inputs.View,
	}
}

func fmtMoney(v int64) string { return fmt.Sprintf("₹%d", v) }
