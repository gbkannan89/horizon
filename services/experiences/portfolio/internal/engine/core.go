package engine

import "fmt"

// CardType represents portfolio workspace card types.
type CardType string
const (CTNetWorth CardType = "NetWorth"; CTAssetAllocation CardType = "AssetAllocation"; CTLiabilitySummary CardType = "LiabilitySummary"; CTLiquidity CardType = "Liquidity"; CTPerformance CardType = "Performance"; CTDiversification CardType = "Diversification"; CTGoalMapping CardType = "GoalMapping"; CTProjectionSummary CardType = "ProjectionSummary"; CTRiskSummary CardType = "RiskSummary"; CTOptimizationSummary CardType = "OptimizationSummary")

type PortfolioCard struct {
	CardID   string   `json:"card_id"`
	CardType CardType `json:"card_type"`
	Title    string   `json:"title"`
	Value    string   `json:"value"`
	Detail   string   `json:"detail"`
	Priority int      `json:"priority"`
}

type PortfolioWorkspace struct {
	Cards []PortfolioCard `json:"cards"`
}

type AllocationEntry struct {
	Classification string  `json:"classification"`
	Value          int64   `json:"value"`
	Percent        float64 `json:"percent"`
}

type PerformanceData struct {
	PeriodReturn     float64 `json:"period_return"`
	PeriodReturnPct  float64 `json:"period_return_pct"`
	BenchmarkReturn  float64 `json:"benchmark_return,omitempty"`
}

type Inputs struct {
	NetWorth         int64             `json:"net_worth"`
	TotalAssets      int64             `json:"total_assets"`
	TotalLiabilities int64             `json:"total_liabilities"`
	TotalCash        int64             `json:"total_cash"`
	PortfolioValue   int64             `json:"portfolio_value"`
	Allocations      []AllocationEntry `json:"allocations"`
	Performance      PerformanceData   `json:"performance"`
	RiskScore        int               `json:"risk_score"`
	LiquidityRatio   float64           `json:"liquidity_ratio"`
	EmergencyMonths  float64           `json:"emergency_months"`
	GoalMappedPct    float64           `json:"goal_mapped_pct"`
}

type Engine struct{}
func NewEngine() *Engine { return &Engine{} }

func (e *Engine) BuildWorkspace(inputs Inputs) *PortfolioWorkspace {
	cards := []PortfolioCard{
		{CardID: "nw", CardType: CTNetWorth, Title: "Net Worth", Value: fmtMoney(inputs.NetWorth), Detail: fmt.Sprintf("Assets: %s | Liabilities: %s", fmtMoney(inputs.TotalAssets), fmtMoney(inputs.TotalLiabilities)), Priority: 1},
		{CardID: "alloc", CardType: CTAssetAllocation, Title: "Asset Allocation", Value: fmt.Sprintf("%d assets", len(inputs.Allocations)), Detail: e.allocationSummary(inputs.Allocations), Priority: 2},
		{CardID: "liab", CardType: CTLiabilitySummary, Title: "Liabilities", Value: fmtMoney(inputs.TotalLiabilities), Priority: 3},
		{CardID: "liq", CardType: CTLiquidity, Title: "Liquidity", Value: fmt.Sprintf("%.1f ratio", inputs.LiquidityRatio), Detail: fmt.Sprintf("Emergency: %.0f months", inputs.EmergencyMonths), Priority: 4},
	}
	if inputs.Performance.PeriodReturn != 0 {
		cards = append(cards, PortfolioCard{CardID: "perf", CardType: CTPerformance, Title: "Performance", Value: fmt.Sprintf("%.1f%%", inputs.Performance.PeriodReturnPct), Priority: 5})
	}
	cards = append(cards, PortfolioCard{CardID: "risk", CardType: CTRiskSummary, Title: "Risk", Value: fmt.Sprintf("Score: %d", inputs.RiskScore), Priority: 6})
	return &PortfolioWorkspace{Cards: cards}
}

func (e *Engine) allocationSummary(allocations []AllocationEntry) string {
	if len(allocations) == 0 { return "No allocation data" }
	return fmt.Sprintf("%d classifications", len(allocations))
}

func fmtMoney(v int64) string { return fmt.Sprintf("₹%d", v) }
