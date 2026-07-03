package engine

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// detectState
// ---------------------------------------------------------------------------

func TestDetectState_FirstTimeUser(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{IsFirstTimeUser: true, HasGoals: false})
	if state != DSFirstTimeUser {
		t.Fatalf("expected DSFirstTimeUser, got %v", state)
	}
}

func TestDetectState_FirstTimeUser_OverridesEvenWithHighHealth(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{
		IsFirstTimeUser: true,
		HasGoals:        false,
		HealthScore:     95,
		RiskScore:       10,
	})
	if state != DSFirstTimeUser {
		t.Fatalf("expected DSFirstTimeUser (highest priority), got %v", state)
	}
}

func TestDetectState_NoData(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{HasAccounts: true, HasEvents: false})
	if state != DSNoData {
		t.Fatalf("expected DSNoData, got %v", state)
	}
}

func TestDetectState_CriticalByHealth(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{HealthScore: 39, RiskScore: 50})
	if state != DSCritical {
		t.Fatalf("expected DSCritical (health<40), got %v", state)
	}
}

func TestDetectState_CriticalByRisk(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{HealthScore: 70, RiskScore: 80})
	if state != DSCritical {
		t.Fatalf("expected DSCritical (risk>=80), got %v", state)
	}
}

func TestDetectState_GoalAtRisk(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{
		GoalsAtRisk: 2,
		HealthScore: 50,
		RiskScore:   60,
	})
	if state != DSGoalAtRisk {
		t.Fatalf("expected DSGoalAtRisk, got %v", state)
	}
}

func TestDetectState_Healthy(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{HealthScore: 60, RiskScore: 30})
	if state != DSHealthy {
		t.Fatalf("expected DSHealthy, got %v", state)
	}
}

func TestDetectState_AttentionNeeded(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{
		HealthScore: 40,
		RiskScore:   60,
	})
	if state != DSAttentionNeeded {
		t.Fatalf("expected DSAttentionNeeded, got %v", state)
	}
}

func TestDetectState_HealthyBoundary(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{HealthScore: 60, RiskScore: 0})
	if state != DSHealthy {
		t.Fatalf("expected DSHealthy at boundary 60, got %v", state)
	}
}

func TestDetectState_AttentionNeededBoundary(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{HealthScore: 59, RiskScore: 79})
	if state != DSAttentionNeeded {
		t.Fatalf("expected DSAttentionNeeded at boundaries, got %v", state)
	}
}

func TestDetectState_AllZerosFallsToCritical(t *testing.T) {
	c := &Composer{}
	state := c.detectState(Inputs{})
	// IsFirstTimeUser false, HasAccounts false, HealthScore=0 < 40
	if state != DSCritical {
		t.Fatalf("expected DSCritical (HealthScore=0 < 40), got %v", state)
	}
}

// ---------------------------------------------------------------------------
// buildWidgets – unconditional widgets
// ---------------------------------------------------------------------------

func TestBuildWidgets_AlwaysContainsHealthScore(t *testing.T) {
	c := &Composer{}
	t1, _, _, _ := c.buildWidgets(Inputs{}, DSHealthy, DMOverview)
	found := false
	for _, w := range t1 {
		if w.WidgetID == "health" {
			found = true
			if w.WidgetType != "health_score" {
				t.Errorf("expected health_score type, got %s", w.WidgetType)
			}
			if w.Confidence != "High" {
				t.Errorf("expected High confidence, got %s", w.Confidence)
			}
			break
		}
	}
	if !found {
		t.Fatal("health widget missing from tier1")
	}
}

func TestBuildWidgets_AlwaysContainsNetWorth(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{NetWorth: 50000}, DSHealthy, DMOverview)
	found := false
	for _, w := range t2 {
		if w.WidgetID == "nw" {
			found = true
			if w.WidgetType != "net_worth" {
				t.Errorf("expected net_worth type, got %s", w.WidgetType)
			}
			break
		}
	}
	if !found {
		t.Fatal("net worth widget missing from tier2")
	}
}

func TestBuildWidgets_AlwaysContainsCashPosition(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{}, DSHealthy, DMOverview)
	found := false
	for _, w := range t2 {
		if w.WidgetID == "cash" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("cash position widget missing from tier2")
	}
}

func TestBuildWidgets_AlwaysContainsRiskSummary(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{RiskScore: 45, RiskLevel: "Moderate"}, DSHealthy, DMOverview)
	found := false
	for _, w := range t2 {
		if w.WidgetID == "risk" {
			found = true
			if w.WidgetType != "risk_summary" {
				t.Errorf("expected risk_summary type, got %s", w.WidgetType)
			}
			break
		}
	}
	if !found {
		t.Fatal("risk summary widget missing from tier2")
	}
}

func TestBuildWidgets_AlwaysContainsCashFlow(t *testing.T) {
	c := &Composer{}
	_, _, t3, _ := c.buildWidgets(Inputs{MonthlyIncome: 100000, MonthlyExpenses: 50000}, DSHealthy, DMOverview)
	found := false
	for _, w := range t3 {
		if w.WidgetID == "cf" {
			found = true
			if w.WidgetType != "cash_flow" {
				t.Errorf("expected cash_flow type, got %s", w.WidgetType)
			}
			break
		}
	}
	if !found {
		t.Fatal("cash flow widget missing from tier3")
	}
}

func TestBuildWidgets_DebtSummaryVisibleWhenLiabilitiesPositive(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{TotalLiabilities: 200000}, DSHealthy, DMOverview)
	found := false
	for _, w := range t2 {
		if w.WidgetID == "debt" {
			found = true
			if !w.Visible {
				t.Error("debt summary should be visible when TotalLiabilities > 0")
			}
			break
		}
	}
	if !found {
		t.Fatal("debt summary widget missing from tier2")
	}
}

func TestBuildWidgets_DebtSummaryHiddenWhenLiabilitiesZero(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{TotalLiabilities: 0}, DSHealthy, DMOverview)
	for _, w := range t2 {
		if w.WidgetID == "debt" && w.Visible {
			t.Error("debt summary should not be visible when TotalLiabilities is 0")
		}
	}
}

// ---------------------------------------------------------------------------
// buildWidgets – conditional widgets
// ---------------------------------------------------------------------------

func TestBuildWidgets_GoalProgressShownWhenGoalsExist(t *testing.T) {
	c := &Composer{}
	t1, _, _, _ := c.buildWidgets(Inputs{GoalCount: 5, GoalsOnTrack: 3}, DSHealthy, DMOverview)
	found := false
	for _, w := range t1 {
		if w.WidgetID == "goals" {
			found = true
			if w.WidgetType != "goal_progress" {
				t.Errorf("expected goal_progress type, got %s", w.WidgetType)
			}
			if !strings.Contains(w.Title, "3/5") {
				t.Errorf("title should contain 3/5, got %q", w.Title)
			}
			break
		}
	}
	if !found {
		t.Fatal("goal progress widget missing when GoalCount > 0")
	}
}

func TestBuildWidgets_GoalProgressHiddenWhenNoGoals(t *testing.T) {
	c := &Composer{}
	t1, _, _, _ := c.buildWidgets(Inputs{GoalCount: 0}, DSHealthy, DMOverview)
	for _, w := range t1 {
		if w.WidgetID == "goals" {
			t.Fatal("goal progress widget should not appear when GoalCount == 0")
		}
	}
}

func TestBuildWidgets_RecommendationShownWhenAvailable(t *testing.T) {
	c := &Composer{}
	t1, _, _, _ := c.buildWidgets(Inputs{HasRecommendation: true}, DSHealthy, DMOverview)
	found := false
	for _, w := range t1 {
		if w.WidgetID == "rec" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("recommendation widget missing when HasRecommendation is true")
	}
}

func TestBuildWidgets_RecommendationHiddenWhenNotAvailable(t *testing.T) {
	c := &Composer{}
	t1, _, _, _ := c.buildWidgets(Inputs{HasRecommendation: false}, DSHealthy, DMOverview)
	for _, w := range t1 {
		if w.WidgetID == "rec" {
			t.Fatal("recommendation widget should not appear when HasRecommendation is false")
		}
	}
}

func TestBuildWidgets_PortfolioSnapshotShownWhenValuePositive(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{PortfolioValue: 1000000}, DSHealthy, DMOverview)
	found := false
	for _, w := range t2 {
		if w.WidgetID == "pf" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("portfolio snapshot missing when PortfolioValue > 0")
	}
}

func TestBuildWidgets_PortfolioSnapshotHiddenWhenZero(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{PortfolioValue: 0}, DSHealthy, DMOverview)
	for _, w := range t2 {
		if w.WidgetID == "pf" {
			t.Fatal("portfolio snapshot should not appear when PortfolioValue == 0")
		}
	}
}

func TestBuildWidgets_UpcomingBillsShownWhenDue(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{BillsDue: 15000}, DSHealthy, DMOverview)
	found := false
	for _, w := range t2 {
		if w.WidgetID == "bills" {
			found = true
			if w.WidgetType != "upcoming_bills" {
				t.Errorf("expected upcoming_bills type, got %s", w.WidgetType)
			}
			break
		}
	}
	if !found {
		t.Fatal("upcoming bills missing when BillsDue > 0")
	}
}

func TestBuildWidgets_UpcomingBillsHiddenWhenNoneDue(t *testing.T) {
	c := &Composer{}
	_, t2, _, _ := c.buildWidgets(Inputs{BillsDue: 0}, DSHealthy, DMOverview)
	for _, w := range t2 {
		if w.WidgetID == "bills" {
			t.Fatal("upcoming bills should not appear when BillsDue == 0")
		}
	}
}

func TestBuildWidgets_RecentEventsShownWhenPresent(t *testing.T) {
	c := &Composer{}
	_, _, t3, _ := c.buildWidgets(Inputs{RecentEvents: 3}, DSHealthy, DMOverview)
	found := false
	for _, w := range t3 {
		if w.WidgetID == "events" {
			found = true
			if w.WidgetType != "recent_events" {
				t.Errorf("expected recent_events type, got %s", w.WidgetType)
			}
			break
		}
	}
	if !found {
		t.Fatal("recent events missing when RecentEvents > 0")
	}
}

func TestBuildWidgets_RecentEventsHiddenWhenNone(t *testing.T) {
	c := &Composer{}
	_, _, t3, _ := c.buildWidgets(Inputs{RecentEvents: 0}, DSHealthy, DMOverview)
	for _, w := range t3 {
		if w.WidgetID == "events" {
			t.Fatal("recent events should not appear when RecentEvents == 0")
		}
	}
}

func TestBuildWidgets_AchievementsShownWhenPresent(t *testing.T) {
	c := &Composer{}
	_, _, t3, _ := c.buildWidgets(Inputs{Achievements: 5}, DSHealthy, DMOverview)
	found := false
	for _, w := range t3 {
		if w.WidgetID == "achievements" {
			found = true
			if w.WidgetType != "achievements" {
				t.Errorf("expected achievements type, got %s", w.WidgetType)
			}
			break
		}
	}
	if !found {
		t.Fatal("achievements missing when Achievements > 0")
	}
}

func TestBuildWidgets_AchievementsHiddenWhenZero(t *testing.T) {
	c := &Composer{}
	_, _, t3, _ := c.buildWidgets(Inputs{Achievements: 0}, DSHealthy, DMOverview)
	for _, w := range t3 {
		if w.WidgetID == "achievements" {
			t.Fatal("achievements should not appear when Achievements == 0")
		}
	}
}

// ---------------------------------------------------------------------------
// buildWidgets – critical alert
// ---------------------------------------------------------------------------

func TestBuildWidgets_CriticalAlertWhenStateCritical(t *testing.T) {
	c := &Composer{}
	_, _, _, alert := c.buildWidgets(Inputs{}, DSCritical, DMOverview)
	if alert == nil {
		t.Fatal("expected critical alert when state is DSCritical")
	}
	if alert.WidgetType != "critical_alert" {
		t.Errorf("expected critical_alert type, got %s", alert.WidgetType)
	}
}

func TestBuildWidgets_CriticalAlertWhenRiskHigh(t *testing.T) {
	c := &Composer{}
	_, _, _, alert := c.buildWidgets(Inputs{RiskScore: 85}, DSHealthy, DMOverview)
	if alert == nil {
		t.Fatal("expected critical alert when RiskScore >= 80 even in healthy state")
	}
}

func TestBuildWidgets_NoCriticalAlertWhenLowRiskAndHealthy(t *testing.T) {
	c := &Composer{}
	_, _, _, alert := c.buildWidgets(Inputs{RiskScore: 50}, DSHealthy, DMOverview)
	if alert != nil {
		t.Fatal("expected no critical alert when state is healthy and risk < 80")
	}
}

// ---------------------------------------------------------------------------
// buildWidgets – tier and priority assignment
// ---------------------------------------------------------------------------

func TestBuildWidgets_PrioritiesWithinTier1(t *testing.T) {
	c := &Composer{}
	inputs := Inputs{HealthScore: 75, HealthGrade: "B", GoalCount: 3, GoalsOnTrack: 2, HasRecommendation: true}
	t1, _, _, _ := c.buildWidgets(inputs, DSHealthy, DMOverview)
	// health=1, goals=2, rec=3
	for _, w := range t1 {
		switch w.WidgetID {
		case "health":
			if w.Priority != 1 {
				t.Errorf("health priority: expected 1, got %d", w.Priority)
			}
		case "goals":
			if w.Priority != 2 {
				t.Errorf("goals priority: expected 2, got %d", w.Priority)
			}
		case "rec":
			if w.Priority != 3 {
				t.Errorf("rec priority: expected 3, got %d", w.Priority)
			}
		}
	}
}

func TestBuildWidgets_TierAssignment(t *testing.T) {
	c := &Composer{}
	inputs := Inputs{
		HealthScore: 80, HealthGrade: "A", RiskScore: 30, RiskLevel: "Low",
		GoalCount: 3, GoalsOnTrack: 2, HasRecommendation: true,
		NetWorth: 500000, CashBalance: 100000, PortfolioValue: 200000,
		BillsDue: 5000, TotalLiabilities: 300000,
		MonthlyIncome: 100000, MonthlyExpenses: 70000,
		RecentEvents: 2, Achievements: 1,
	}
	t1, t2, t3, _ := c.buildWidgets(inputs, DSHealthy, DMOverview)

	for _, w := range t1 {
		if w.Tier != Tier1 {
			t.Errorf("widget %s in tier1 slice should have Tier1, has %d", w.WidgetID, w.Tier)
		}
	}
	for _, w := range t2 {
		if w.Tier != Tier2 {
			t.Errorf("widget %s in tier2 slice should have Tier2, has %d", w.WidgetID, w.Tier)
		}
	}
	for _, w := range t3 {
		if w.Tier != Tier3 {
			t.Errorf("widget %s in tier3 slice should have Tier3, has %d", w.WidgetID, w.Tier)
		}
	}
}

// ---------------------------------------------------------------------------
// Compose – integration
// ---------------------------------------------------------------------------

func TestCompose_BasicFlow(t *testing.T) {
	c := NewComposer()
	out := c.Compose(Inputs{
		UserID: "user-1", HealthScore: 75, HealthGrade: "B+",
		RiskScore: 35, RiskLevel: "Low",
		NetWorth: 1000000, CashBalance: 250000,
		GoalCount: 4, GoalsOnTrack: 3,
	}, DMOverview)

	if out == nil {
		t.Fatal("Compose returned nil")
	}
	if out.DashboardID != "dash-user-1" {
		t.Errorf("expected dash-user-1, got %s", out.DashboardID)
	}
	if out.Mode != DMOverview {
		t.Errorf("expected DMOverview mode, got %s", out.Mode)
	}
	if out.State != DSHealthy {
		t.Errorf("expected DSHealthy state, got %s", out.State)
	}
	if out.LastRefreshTime == "" {
		t.Error("LastRefreshTime should not be empty")
	}
	if len(out.Tier1) == 0 {
		t.Error("Tier1 should not be empty")
	}
}

func TestCompose_DefaultMode(t *testing.T) {
	c := NewComposer()
	out := c.Compose(Inputs{HealthScore: 70, HealthGrade: "B", RiskScore: 30, RiskLevel: "Low"}, "")
	if out.Mode != DMOverview {
		t.Errorf("expected default mode DMOverview, got %s", out.Mode)
	}
}

// ---------------------------------------------------------------------------
// Compose – mode passthrough (all 6 modes)
// ---------------------------------------------------------------------------

func TestCompose_AllModes(t *testing.T) {
	c := NewComposer()
	base := Inputs{
		UserID: "u1", HealthScore: 80, HealthGrade: "A",
		RiskScore: 25, RiskLevel: "Low",
		NetWorth: 100000, CashBalance: 50000,
	}
	modes := []DashboardMode{DMOverview, DMDecision, DMGoal, DMPortfolio, DMPlanning, DMHistorical}
	for _, mode := range modes {
		out := c.Compose(base, mode)
		if out.Mode != mode {
			t.Errorf("expected mode %s, got %s", mode, out.Mode)
		}
		if out.DashboardID != "dash-u1" {
			t.Errorf("expected dash-u1, got %s", out.DashboardID)
		}
	}
}

// ---------------------------------------------------------------------------
// BuildSummary
// ---------------------------------------------------------------------------

func TestBuildSummary(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		NetWorth:         1500000,
		HealthScore:      82,
		HealthGrade:      "A",
		RiskScore:        30,
		RiskLevel:        "Low",
		GoalsOnTrack:     4,
		GoalCount:        5,
		CashBalance:      500000,
		PortfolioValue:   800000,
		TotalLiabilities: 300000,
	}
	s := c.BuildSummary(inputs)
	if s == nil {
		t.Fatal("BuildSummary returned nil")
	}
	if s.NetWorth != 1500000 {
		t.Errorf("NetWorth: expected 1500000, got %d", s.NetWorth)
	}
	if s.HealthScore != 82 {
		t.Errorf("HealthScore: expected 82, got %d", s.HealthScore)
	}
	if s.HealthGrade != "A" {
		t.Errorf("HealthGrade: expected A, got %s", s.HealthGrade)
	}
	if s.RiskScore != 30 {
		t.Errorf("RiskScore: expected 30, got %d", s.RiskScore)
	}
	if s.RiskLevel != "Low" {
		t.Errorf("RiskLevel: expected Low, got %s", s.RiskLevel)
	}
	if s.GoalsOnTrack != 4 {
		t.Errorf("GoalsOnTrack: expected 4, got %d", s.GoalsOnTrack)
	}
	if s.TotalGoals != 5 {
		t.Errorf("TotalGoals: expected 5, got %d", s.TotalGoals)
	}
	if s.CashBalance != 500000 {
		t.Errorf("CashBalance: expected 500000, got %d", s.CashBalance)
	}
	if s.PortfolioValue != 800000 {
		t.Errorf("PortfolioValue: expected 800000, got %d", s.PortfolioValue)
	}
	if s.TotalDebt != 300000 {
		t.Errorf("TotalDebt: expected 300000, got %d", s.TotalDebt)
	}
}

func TestBuildSummary_ZeroInputs(t *testing.T) {
	c := NewComposer()
	s := c.BuildSummary(Inputs{})
	if s == nil {
		t.Fatal("BuildSummary returned nil")
	}
	if s.NetWorth != 0 || s.HealthScore != 0 || s.RiskScore != 0 {
		t.Error("all numeric fields should default to zero")
	}
}

// ---------------------------------------------------------------------------
// fmtHealth
// ---------------------------------------------------------------------------

func TestFmtHealth(t *testing.T) {
	tests := []struct {
		score int
		grade string
		want  string
	}{
		{80, "A", "Health: 80 (A)"},
		{0, "F", "Health: 0 (F)"},
		{100, "A+", "Health: 100 (A+)"},
	}
	for _, tt := range tests {
		got := fmtHealth(tt.score, tt.grade)
		if got != tt.want {
			t.Errorf("fmtHealth(%d, %q) = %q; want %q", tt.score, tt.grade, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// fmtMoney
// ---------------------------------------------------------------------------

func TestFmtMoney(t *testing.T) {
	tests := []struct {
		v    int64
		want string
	}{
		{0, "₹0"},
		{500, "₹500"},
		{100000, "₹100000"},
		{10000000, "₹10000000"},
		{999999999, "₹999999999"},
		{-1000, "₹-1000"},
	}
	for _, tt := range tests {
		got := fmtMoney(tt.v)
		if got != tt.want {
			t.Errorf("fmtMoney(%d) = %q; want %q", tt.v, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Widget type inventory – verify all 13 widget types are exercised
// ---------------------------------------------------------------------------

func TestAllWidgetTypesCovered(t *testing.T) {
	c := &Composer{}
	inputs := Inputs{
		HealthScore: 80, HealthGrade: "A",
		RiskScore: 45, RiskLevel: "Moderate",
		GoalCount: 3, GoalsOnTrack: 2,
		HasRecommendation: true,
		NetWorth: 500000, CashBalance: 100000,
		PortfolioValue: 200000, BillsDue: 5000,
		TotalLiabilities: 300000,
		MonthlyIncome: 100000, MonthlyExpenses: 70000,
		RecentEvents: 4, Achievements: 2,
	}
	types := make(map[string]int)
	collect := func(list []Widget) {
		for _, w := range list {
			types[w.WidgetType]++
		}
	}
	t1, t2, t3, alert2 := c.buildWidgets(inputs, DSCritical, DMOverview)
	collect(t1)
	collect(t2)
	collect(t3)
	if alert2 != nil {
		types[alert2.WidgetType]++
	}

	expected := []string{
		"critical_alert", "health_score", "goal_progress", "top_recommendation",
		"net_worth", "cash_position", "portfolio_snapshot", "risk_summary",
		"upcoming_bills", "debt_summary", "cash_flow", "recent_events", "achievements",
	}
	for _, wt := range expected {
		if types[wt] == 0 {
			t.Errorf("widget type %q not produced by buildWidgets", wt)
		}
	}
	if len(types) < len(expected) {
		t.Errorf("expected at least %d widget types, got %d", len(expected), len(types))
	}
}

// ---------------------------------------------------------------------------
// NewComposer
// ---------------------------------------------------------------------------

func TestNewComposer(t *testing.T) {
	c := NewComposer()
	if c == nil {
		t.Fatal("NewComposer returned nil")
	}
}

// ---------------------------------------------------------------------------
// Table-driven compound scenarios
// ---------------------------------------------------------------------------

func TestDetectStateScenarios(t *testing.T) {
	tests := []struct {
		name     string
		inputs   Inputs
		expected DashboardState
	}{
		{
			name:     "first time user without goals",
			inputs:   Inputs{IsFirstTimeUser: true, HasGoals: false},
			expected: DSFirstTimeUser,
		},
		{
			name:     "first time user with goals – falls through to healthy",
			inputs:   Inputs{IsFirstTimeUser: true, HasGoals: true, HealthScore: 70},
			expected: DSHealthy,
		},
		{
			name:     "has accounts but no events",
			inputs:   Inputs{HasAccounts: true, HasEvents: false},
			expected: DSNoData,
		},
		{
			name:     "critical via health below 40",
			inputs:   Inputs{HealthScore: 39},
			expected: DSCritical,
		},
		{
			name:     "critical via risk >= 80",
			inputs:   Inputs{RiskScore: 80, HealthScore: 70},
			expected: DSCritical,
		},
		{
			name:     "goals at risk overrides attention needed",
			inputs:   Inputs{GoalsAtRisk: 1, HealthScore: 50, RiskScore: 50},
			expected: DSGoalAtRisk,
		},
		{
			name:     "healthy",
			inputs:   Inputs{HealthScore: 60, RiskScore: 50},
			expected: DSHealthy,
		},
		{
			name:     "attention needed default",
			inputs:   Inputs{HealthScore: 50, RiskScore: 70},
			expected: DSAttentionNeeded,
		},
	}

	c := &Composer{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.detectState(tt.inputs)
			if got != tt.expected {
				t.Errorf("detectState() = %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestComposeEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		inputs Inputs
		mode   DashboardMode
		checks func(t *testing.T, out *DashboardOutput)
	}{
		{
			name:   "empty inputs still produce dashboard",
			inputs: Inputs{UserID: "empty"},
			mode:   DMOverview,
			checks: func(t *testing.T, out *DashboardOutput) {
				if out == nil {
					t.Fatal("output is nil")
				}
				if len(out.Tier1) == 0 {
					t.Error("expected at least one widget in tier1")
				}
			},
		},
		{
			name: "maximal inputs produce all widgets",
			inputs: Inputs{
				UserID: "max", HealthScore: 90, HealthGrade: "A",
				RiskScore: 20, RiskLevel: "Low",
				NetWorth: 10000000, CashBalance: 2000000,
				PortfolioValue: 5000000, BillsDue: 10000,
				TotalLiabilities: 1000000,
				GoalCount: 10, GoalsOnTrack: 8,
				HasRecommendation: true,
				MonthlyIncome: 500000, MonthlyExpenses: 300000,
				RecentEvents: 20, Achievements: 15,
			},
			mode: DMGoal,
			checks: func(t *testing.T, out *DashboardOutput) {
				if len(out.Tier1) != 3 {
					t.Errorf("expected 3 tier1 widgets, got %d", len(out.Tier1))
				}
				if len(out.Tier2) != 6 {
					t.Errorf("expected 6 tier2 widgets, got %d", len(out.Tier2))
				}
				if len(out.Tier3) != 3 {
					t.Errorf("expected 3 tier3 widgets, got %d", len(out.Tier3))
				}
				if out.CriticalAlert != nil {
					t.Error("should not have critical alert with healthy inputs")
				}
			},
		},
		{
			name: "critical state produces alert",
			inputs: Inputs{
				UserID: "crit", HealthScore: 30, HealthGrade: "D",
				RiskScore: 85, RiskLevel: "High",
			},
			mode: DMDecision,
			checks: func(t *testing.T, out *DashboardOutput) {
				if out.CriticalAlert == nil {
					t.Fatal("expected critical alert")
				}
				if out.State != DSCritical {
					t.Errorf("expected DSCritical state, got %s", out.State)
				}
			},
		},
	}
	c := NewComposer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := c.Compose(tt.inputs, tt.mode)
			tt.checks(t, out)
		})
	}
}

// ---------------------------------------------------------------------------
// Confidence fields
// ---------------------------------------------------------------------------

func TestWidgetConfidence(t *testing.T) {
	c := &Composer{}
	t1, _, _, alert := c.buildWidgets(Inputs{HealthScore: 75, HealthGrade: "B", RiskScore: 90, HasRecommendation: true}, DSCritical, DMOverview)
	for _, w := range t1 {
		if w.Confidence == "" {
			t.Errorf("widget %s has empty confidence", w.WidgetID)
		}
	}
	if alert != nil && alert.Confidence != "High" {
		t.Errorf("critical alert confidence should be High, got %s", alert.Confidence)
	}
}
