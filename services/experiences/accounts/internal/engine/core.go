package engine

import "fmt"

type CardType string

const (
	CTOverview       CardType = "Overview"
	CTCashBalance    CardType = "CashBalance"
	CTCreditUtil     CardType = "CreditUtilization"
	CTAccountHealth  CardType = "AccountHealth"
	CTInstitution    CardType = "Institution"
	CTCashFlow       CardType = "CashFlow"
	CTProjection     CardType = "Projection"
	CTRisk           CardType = "Risk"
	CTRecommendation CardType = "Recommendation"
	CTTimeline       CardType = "Timeline"
)

type Card struct {
	CardID   string      `json:"card_id"`
	CardType CardType    `json:"card_type"`
	Title    string      `json:"title"`
	Summary  string      `json:"summary"`
	Data     interface{} `json:"data,omitempty"`
	Priority int         `json:"priority"`
}

type AccountsDashboard struct {
	TotalBalance     int64  `json:"total_balance"`
	TotalAccounts    int    `json:"total_accounts"`
	AccountsByType   map[string]int `json:"accounts_by_type"`
	Cards            []Card `json:"cards"`
}

type BalanceSummary struct {
	TotalBalance    int64            `json:"total_balance"`
	TotalAvailable  int64            `json:"total_available"`
	TotalSpendable  int64            `json:"total_spendable"`
	CreditUtilized  float64          `json:"credit_utilization_pct"`
	AccountBalances []AccountBalItem `json:"account_balances"`
}

type AccountBalItem struct {
	AccountID   string `json:"account_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Balance     int64  `json:"balance"`
	Available   int64  `json:"available"`
}

type CashFlowSummary struct {
	PeriodInflow  int64 `json:"period_inflow"`
	PeriodOutflow int64 `json:"period_outflow"`
	NetFlow       int64 `json:"net_flow"`
	Projected     int64 `json:"projected_flow"`
}

type AccountHealthSummary struct {
	AccountsByHealth map[string]int `json:"accounts_by_health"`
	HealthyPct       float64        `json:"healthy_pct"`
}

// Extended Inputs — add engine data
type Inputs struct {
	Accounts       []AccountInput `json:"accounts"`
	Account        *AccountInput  `json:"account,omitempty"`
	Transactions   []TransInput   `json:"transactions,omitempty"`
	UserID         string         `json:"user_id"`

	CashBalance    int64  `json:"cash_balance"`
	TotalInflow    int64  `json:"total_inflow"`
	TotalOutflow   int64  `json:"total_outflow"`
	ProjectedFlow  int64  `json:"projected_flow"`
	RiskScore      int    `json:"risk_score"`
	RiskLevel      string `json:"risk_level"`
	HealthScore    int    `json:"health_score"`
	HealthGrade    string `json:"health_grade"`
	HasRecs        bool   `json:"has_recommendations"`
	RecCount       int    `json:"rec_count"`
	HasOpts        bool   `json:"has_optimizations"`
	HasSims        bool   `json:"has_simulations"`
	EventCount     int    `json:"event_count"`
}

type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

// Existing methods remain from composer.go — this adds dashboard composition
func (c *Composer) BuildDashboard(inputs Inputs) *AccountsDashboard {
	cards := c.buildCards(inputs)

	byType := map[string]int{}
	totalBal := int64(0)
	for _, a := range inputs.Accounts {
		byType[a.AccountType]++
		totalBal += a.CurrentBalance
	}

	return &AccountsDashboard{
		TotalBalance:   totalBal,
		TotalAccounts:  len(inputs.Accounts),
		AccountsByType: byType,
		Cards:          cards,
	}
}

func (c *Composer) BuildBalanceSummary(inputs Inputs) *BalanceSummary {
	totalAvail := int64(0)
	totalSpend := int64(0)
	var items []AccountBalItem
	for _, a := range inputs.Accounts {
		totalAvail += a.AvailableBalance
		totalSpend += a.SpendableBalance
		items = append(items, AccountBalItem{
			AccountID: a.AccountID, Name: a.AccountName,
			Type: a.AccountType, Balance: a.CurrentBalance,
			Available: a.AvailableBalance,
		})
	}
	return &BalanceSummary{
		TotalBalance:   inputs.CashBalance,
		TotalAvailable: totalAvail,
		TotalSpendable: totalSpend,
		AccountBalances: items,
	}
}

func (c *Composer) BuildCashFlowSummary(inputs Inputs) *CashFlowSummary {
	return &CashFlowSummary{
		PeriodInflow: inputs.TotalInflow,
		PeriodOutflow: inputs.TotalOutflow,
		NetFlow:      inputs.TotalInflow - inputs.TotalOutflow,
		Projected:    inputs.ProjectedFlow,
	}
}

func (c *Composer) BuildHealthSummary(inputs Inputs) *AccountHealthSummary {
	byHealth := map[string]int{}
	total := 0
	for _, a := range inputs.Accounts {
		byHealth[a.AccountHealth]++
		total++
	}
	healthyPct := 0.0
	if total > 0 {
		healthyPct = float64(byHealth["healthy"]+byHealth["good"]) / float64(total) * 100
	}
	return &AccountHealthSummary{
		AccountsByHealth: byHealth,
		HealthyPct:       healthyPct,
	}
}

func (c *Composer) buildCards(inputs Inputs) []Card {
	cards := []Card{
		{CardID: "ov", CardType: CTOverview, Title: "Accounts Overview",
			Summary:  fmt.Sprintf("%d accounts — %s total", len(inputs.Accounts), fmtMoney(inputs.CashBalance)),
			Data:     map[string]interface{}{"count": len(inputs.Accounts), "balance": inputs.CashBalance},
			Priority: 1},
		{CardID: "cash", CardType: CTCashBalance, Title: "Cash Balance",
			Summary:  fmtMoney(inputs.CashBalance),
			Priority: 2},
		{CardID: "health", CardType: CTAccountHealth, Title: "Account Health",
			Summary:  fmt.Sprintf("Score: %d (%s)", inputs.HealthScore, inputs.HealthGrade),
			Data:     map[string]interface{}{"score": inputs.HealthScore},
			Priority: 3},
		{CardID: "cf", CardType: CTCashFlow, Title: "Cash Flow",
			Summary:  fmt.Sprintf("In: %s | Out: %s", fmtMoney(inputs.TotalInflow), fmtMoney(inputs.TotalOutflow)),
			Priority: 4},
		{CardID: "risk", CardType: CTRisk, Title: "Risk Assessment",
			Summary:  fmt.Sprintf("Score: %d — %s", inputs.RiskScore, inputs.RiskLevel),
			Data:     map[string]interface{}{"score": inputs.RiskScore},
			Priority: 5},
		{CardID: "proj", CardType: CTProjection, Title: "Projection",
			Summary:  fmt.Sprintf("Projected flow: %s", fmtMoney(inputs.ProjectedFlow)),
			Priority: 6},
	}
	if inputs.HasRecs {
		cards = append(cards, Card{CardID: "rec", CardType: CTRecommendation,
			Title: "Recommendations", Summary: fmt.Sprintf("%d available", inputs.RecCount), Priority: 7})
	}
	if inputs.EventCount > 0 {
		cards = append(cards, Card{CardID: "tl", CardType: CTTimeline,
			Title: "Recent Activity", Summary: fmt.Sprintf("%d events", inputs.EventCount), Priority: 8})
	}
	return cards
}

func SummaryCount(n int) string { if n == 0 { return "No items" }; return fmt.Sprintf("%d items", n) }
func SummaryMoney(v int64) string { return fmt.Sprintf("₹%d projected", v) }
func fmtMoney(v int64) string {
	if v >= 10000000 { return fmt.Sprintf("₹%.2fCr", float64(v)/10000000) }
	if v >= 100000 { return fmt.Sprintf("₹%.2fL", float64(v)/100000) }
	return fmt.Sprintf("₹%d", v)
}
