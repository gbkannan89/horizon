package engine

import "fmt"

// CardType represents portfolio dashboard card types.
type CardType string

const (
	CTOverview      CardType = "Overview"
	CTAllocation    CardType = "Allocation"
	CTPerformance   CardType = "Performance"
	CTHoldings      CardType = "Holdings"
	CTProjection    CardType = "Projection"
	CTRisk          CardType = "Risk"
	CTRecommend     CardType = "Recommendation"
	CTOptimization  CardType = "Optimization"
	CTSimulation    CardType = "Simulation"
	CTTimeline      CardType = "Timeline"
)

// Card represents a portfolio dashboard card.
type Card struct {
	CardID   string      `json:"card_id"`
	CardType CardType    `json:"card_type"`
	Title    string      `json:"title"`
	Summary  string      `json:"summary"`
	Data     interface{} `json:"data,omitempty"`
	Priority int         `json:"priority"`
}

// PortfolioDashboard is the full portfolio experience view.
type PortfolioDashboard struct {
	PortfolioValue int64         `json:"portfolio_value"`
	TotalReturn    float64       `json:"total_return"`
	TotalReturnPct float64       `json:"total_return_pct"`
	RiskScore      int           `json:"risk_score"`
	RiskLevel      string        `json:"risk_level"`
	Cards          []Card        `json:"cards"`
}

// AllocationView represents portfolio allocation breakdown.
type AllocationView struct {
	Allocations     []AllocEntry `json:"allocations"`
	TotalValue      int64        `json:"total_value"`
	Diversification float64      `json:"diversification_score"`
}

type AllocEntry struct {
	Label   string  `json:"label"`
	Value   int64   `json:"value"`
	Percent float64 `json:"percent"`
	Target  float64 `json:"target,omitempty"`
	Drift   float64 `json:"drift,omitempty"`
}

// PerformanceView represents portfolio performance.
type PerformanceView struct {
	PeriodReturn     float64 `json:"period_return"`
	PeriodReturnPct  float64 `json:"period_return_pct"`
	BenchmarkReturn  float64 `json:"benchmark_return,omitempty"`
	UnrealizedGL     int64   `json:"unrealized_gain_loss"`
	RealizedGL       int64   `json:"realized_gain_loss"`
	Period           string  `json:"period"`
}

// RiskView represents portfolio risk assessment.
type RiskView struct {
	RiskScore   int        `json:"risk_score"`
	RiskLevel   string     `json:"risk_level"`
	Var         float64    `json:"var,omitempty"`
	SharpeRatio float64    `json:"sharpe_ratio,omitempty"`
	Volatility  float64    `json:"volatility,omitempty"`
	MaxDrawdown float64    `json:"max_drawdown,omitempty"`
}

// ProjectionView represents portfolio projection.
type ProjectionView struct {
	ProjectedValue   float64 `json:"projected_value"`
	Confidence       string  `json:"confidence"`
	HorizonYears     int     `json:"horizon_years"`
	AnnualReturn     float64 `json:"annual_return,omitempty"`
}

// SimView represents portfolio simulation data.
type SimView struct {
	Simulations []SimEntry `json:"simulations"`
	TotalCount  int        `json:"total_count"`
}

type SimEntry struct {
	SimID    string `json:"sim_id"`
	Scenario string `json:"scenario"`
	Outcome  string `json:"outcome"`
}

// CardView represents portfolio timeline card data.
type CardView struct {
	Items []CardEntry `json:"items"`
	Count int         `json:"count"`
}

type CardEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Value string `json:"value"`
}

// Inputs aggregates all data for the portfolio experience.
type Inputs struct {
	UserID         string       `json:"user_id"`
	PortfolioValue int64        `json:"portfolio_value"`
	TotalReturn    float64      `json:"total_return"`
	TotalReturnPct float64      `json:"total_return_pct"`
	Allocations    []AllocEntry `json:"allocations"`
	RiskScore      int          `json:"risk_score"`
	RiskLevel      string       `json:"risk_level"`
	PeriodReturn   float64      `json:"period_return"`
	PeriodReturnPct float64     `json:"period_return_pct"`
	UnrealizedGL   int64        `json:"unrealized_gain_loss"`
	RealizedGL     int64        `json:"realized_gain_loss"`
	BenchmarkRet   float64      `json:"benchmark_return"`
	SharpeRatio    float64      `json:"sharpe_ratio"`
	Volatility     float64      `json:"volatility"`
	MaxDrawdown    float64      `json:"max_drawdown"`
	Var            float64      `json:"var"`
	ProjectedVal   float64      `json:"projected_value"`
	ProjectionConf string       `json:"projection_confidence"`
	HorizonYears   int          `json:"horizon_years"`
	AnnualReturn   float64      `json:"annual_return"`
	HasRecs        bool         `json:"has_recommendations"`
	HasOpts        bool         `json:"has_optimizations"`
	HasSims        bool         `json:"has_simulations"`
	RecCount       int          `json:"rec_count"`
	OptCount       int          `json:"opt_count"`
	SimCount       int          `json:"sim_count"`
	EventCount     int          `json:"event_count"`
	MilestoneCount int          `json:"milestone_count"`
}

// Composer builds portfolio experience view models.
type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

// BuildDashboard returns the full portfolio dashboard.
func (c *Composer) BuildDashboard(inputs Inputs) *PortfolioDashboard {
	cards := c.buildCards(inputs)
	return &PortfolioDashboard{
		PortfolioValue: inputs.PortfolioValue,
		TotalReturn:    inputs.TotalReturn,
		TotalReturnPct: inputs.TotalReturnPct,
		RiskScore:      inputs.RiskScore,
		RiskLevel:      inputs.RiskLevel,
		Cards:          cards,
	}
}

// BuildAllocation returns allocation breakdown.
func (c *Composer) BuildAllocation(inputs Inputs) *AllocationView {
	total := int64(0)
	for _, a := range inputs.Allocations { total += a.Value }
	div := c.diversificationScore(inputs.Allocations)
	return &AllocationView{Allocations: inputs.Allocations, TotalValue: total, Diversification: div}
}

// BuildPerformance returns performance data.
func (c *Composer) BuildPerformance(inputs Inputs) *PerformanceView {
	return &PerformanceView{
		PeriodReturn: inputs.PeriodReturn, PeriodReturnPct: inputs.PeriodReturnPct,
		BenchmarkReturn: inputs.BenchmarkRet,
		UnrealizedGL: inputs.UnrealizedGL, RealizedGL: inputs.RealizedGL,
		Period: "1M",
	}
}

// BuildRisk returns risk assessment.
func (c *Composer) BuildRisk(inputs Inputs) *RiskView {
	return &RiskView{
		RiskScore: inputs.RiskScore, RiskLevel: inputs.RiskLevel,
		Var: inputs.Var, SharpeRatio: inputs.SharpeRatio,
		Volatility: inputs.Volatility, MaxDrawdown: inputs.MaxDrawdown,
	}
}

// BuildProjection returns projection view.
func (c *Composer) BuildProjection(inputs Inputs) *ProjectionView {
	return &ProjectionView{
		ProjectedValue: inputs.ProjectedVal, Confidence: inputs.ProjectionConf,
		HorizonYears: inputs.HorizonYears, AnnualReturn: inputs.AnnualReturn,
	}
}

// BuildRecommendations returns recommendation card data.
func (c *Composer) BuildRecommendations(inputs Inputs) *CardView {
	count := 0
	if inputs.HasRecs { count = inputs.RecCount }
	return &CardView{
		Items: []CardEntry{{ID: "recs", Title: "Active Recommendations", Value: fmt.Sprintf("%d", count)}},
		Count: count,
	}
}

// BuildOptimization returns optimization card data.
func (c *Composer) BuildOptimization(inputs Inputs) *CardView {
	count := 0
	if inputs.HasOpts { count = inputs.OptCount }
	return &CardView{
		Items: []CardEntry{{ID: "opts", Title: "Strategies Available", Value: fmt.Sprintf("%d", count)}},
		Count: count,
	}
}

// BuildSimulations returns simulation card data.
func (c *Composer) BuildSimulations(inputs Inputs) *SimView {
	return &SimView{
		Simulations: []SimEntry{}, TotalCount: inputs.SimCount,
	}
}

// BuildCardTimeline returns timeline card data.
func (c *Composer) BuildCardTimeline(inputs Inputs) *CardView {
	return &CardView{
		Items: []CardEntry{{ID: "events", Title: "Recent Events", Value: fmt.Sprintf("%d", inputs.EventCount)}},
		Count: inputs.EventCount,
	}
}

func (c *Composer) buildCards(inputs Inputs) []Card {
	cards := []Card{
		{CardID: "ov", CardType: CTOverview, Title: "Portfolio Overview",
			Summary:  fmt.Sprintf("₹%d total", inputs.PortfolioValue),
			Data:     map[string]interface{}{"value": inputs.PortfolioValue, "return": inputs.TotalReturnPct},
			Priority: 1},
		{CardID: "alloc", CardType: CTAllocation, Title: "Asset Allocation",
			Summary:  fmt.Sprintf("%d classes", len(inputs.Allocations)),
			Data:     inputs.Allocations,
			Priority: 2},
		{CardID: "perf", CardType: CTPerformance, Title: "Performance",
			Summary:  fmt.Sprintf("%.1f%% return", inputs.PeriodReturnPct),
			Data:     map[string]interface{}{"return": inputs.PeriodReturnPct, "unrealized": inputs.UnrealizedGL},
			Priority: 3},
		{CardID: "hold", CardType: CTHoldings, Title: "Holdings",
			Summary:  fmt.Sprintf("₹%d total", inputs.PortfolioValue),
			Priority: 4},
		{CardID: "proj", CardType: CTProjection, Title: "Projection",
			Summary:  fmt.Sprintf("₹%.0f projected", inputs.ProjectedVal),
			Data:     map[string]interface{}{"projected": inputs.ProjectedVal, "horizon": inputs.HorizonYears},
			Priority: 5},
	}
	if inputs.RiskScore > 0 {
		cards = append(cards, Card{CardID: "risk", CardType: CTRisk, Title: "Risk Assessment",
			Summary: fmt.Sprintf("Score: %d — %s", inputs.RiskScore, inputs.RiskLevel),
			Data:    map[string]interface{}{"score": inputs.RiskScore, "level": inputs.RiskLevel},
			Priority: 6})
	}
	if inputs.HasRecs {
		cards = append(cards, Card{CardID: "rec", CardType: CTRecommend, Title: "Recommendations",
			Summary: fmt.Sprintf("%d available", inputs.RecCount), Priority: 7})
	}
	if inputs.HasOpts {
		cards = append(cards, Card{CardID: "opt", CardType: CTOptimization, Title: "Optimization",
			Summary: fmt.Sprintf("%d strategies", inputs.OptCount), Priority: 8})
	}
	if inputs.HasSims {
		cards = append(cards, Card{CardID: "sim", CardType: CTSimulation, Title: "Simulation",
			Summary: fmt.Sprintf("%d scenarios", inputs.SimCount), Priority: 9})
	}
	if inputs.EventCount > 0 {
		cards = append(cards, Card{CardID: "tl", CardType: CTTimeline, Title: "Timeline",
			Summary: fmt.Sprintf("%d events", inputs.EventCount), Priority: 10})
	}
	return cards
}

func (c *Composer) diversificationScore(allocations []AllocEntry) float64 {
	if len(allocations) == 0 { return 0 }
	sumSq := 0.0
	for _, a := range allocations {
		p := a.Percent / 100.0
		sumSq += p * p
	}
	if sumSq == 0 { return 0 }
	return 1 / sumSq
}
