package engine

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// mock state repository
// ---------------------------------------------------------------------------

type mockStateRepo struct {
	states map[string]State
}

func (m *mockStateRepo) Get(_ context.Context, _, notifID string) (State, bool) {
	if m == nil || m.states == nil {
		return "", false
	}
	s, ok := m.states[notifID]
	return s, ok
}

func (m *mockStateRepo) SetState(_ context.Context, _, notifID string, s State) error {
	if m.states == nil {
		m.states = make(map[string]State)
	}
	m.states[notifID] = s
	return nil
}

func (m *mockStateRepo) Snooze(_ context.Context, _, notifID string, _ string) error {
	if m.states == nil {
		m.states = make(map[string]State)
	}
	m.states[notifID] = StateSnoozed
	return nil
}

func emptyRepo() *mockStateRepo {
	return &mockStateRepo{states: map[string]State{}}
}

func ctx() context.Context {
	return context.Background()
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
// weight
// ---------------------------------------------------------------------------

func TestWeight(t *testing.T) {
	c := &Composer{}
	tests := []struct {
		p    Priority
		want int
	}{
		{P1Critical, 1},
		{P2High, 2},
		{P3Medium, 3},
		{P4Low, 4},
		{P5Info, 5},
	}
	for _, tt := range tests {
		got := c.weight(tt.p)
		if got != tt.want {
			t.Errorf("weight(%s) = %d; want %d", tt.p, got, tt.want)
		}
	}
}

func TestWeight_DefaultIsP5(t *testing.T) {
	c := &Composer{}
	got := c.weight(Priority("unknown"))
	if got != 5 {
		t.Errorf("weight(unknown) = %d; want 5", got)
	}
}

// ---------------------------------------------------------------------------
// generate – health notifications
// ---------------------------------------------------------------------------

func TestGenerate_HealthCritical(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HealthScore: 39, HealthGrade: "D", HealthChange: -10})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-health-critical" {
			found = true
			if n.Category != CatHealth {
				t.Errorf("category: expected %s, got %s", CatHealth, n.Category)
			}
			if n.Priority != P1Critical {
				t.Errorf("priority: expected %s, got %s", P1Critical, n.Priority)
			}
			if n.SourceEngine != "HealthScore" {
				t.Errorf("source: expected HealthScore, got %s", n.SourceEngine)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-health-critical not generated when HealthScore=39")
	}
}

func TestGenerate_HealthCritical_PrecedesDecline(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HealthScore: 39, HealthGrade: "D", HealthChange: -10})
	for _, n := range notifs {
		if n.NotifID == "notif-health-down" {
			t.Fatal("notif-health-down should not appear when health is critical")
		}
	}
}

func TestGenerate_HealthDeclined(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HealthScore: 55, HealthGrade: "C", HealthChange: -5})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-health-down" {
			found = true
			if n.Category != CatHealth {
				t.Errorf("category: expected %s, got %s", CatHealth, n.Category)
			}
			if n.Priority != P2High {
				t.Errorf("priority: expected %s, got %s", P2High, n.Priority)
			}
			if n.SourceEngine != "HealthScore" {
				t.Errorf("source: expected HealthScore, got %s", n.SourceEngine)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-health-down not generated when HealthChange=-5")
	}
}

func TestGenerate_NoHealthNotifications(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HealthScore: 75, HealthGrade: "B", HealthChange: 0})
	for _, n := range notifs {
		if n.Category == CatHealth {
			t.Fatal("no health notification expected when score>=40 and change>=0")
		}
	}
}

// ---------------------------------------------------------------------------
// generate – risk notifications
// ---------------------------------------------------------------------------

func TestGenerate_RiskCritical(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{RiskScore: 80, RiskLevel: "High"})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-risk-critical" {
			found = true
			if n.Category != CatRisk {
				t.Errorf("category: expected %s, got %s", CatRisk, n.Category)
			}
			if n.Priority != P1Critical {
				t.Errorf("priority: expected %s, got %s", P1Critical, n.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-risk-critical not generated when RiskScore=80")
	}
}

func TestGenerate_RiskHigh(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{RiskScore: 70, RiskLevel: "Elevated"})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-risk-high" {
			found = true
			if n.Category != CatRisk {
				t.Errorf("category: expected %s, got %s", CatRisk, n.Category)
			}
			if n.Priority != P2High {
				t.Errorf("priority: expected %s, got %s", P2High, n.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-risk-high not generated when RiskScore=70")
	}
}

func TestGenerate_RiskHighBoundary(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{RiskScore: 60, RiskLevel: "Moderate"})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-risk-high" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("notif-risk-high should be generated at boundary RiskScore=60")
	}
}

func TestGenerate_NoRiskNotifications(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{RiskScore: 59, RiskLevel: "Moderate"})
	for _, n := range notifs {
		if n.Category == CatRisk {
			t.Fatal("no risk notification expected when RiskScore < 60")
		}
	}
}

// ---------------------------------------------------------------------------
// generate – goal notifications
// ---------------------------------------------------------------------------

func TestGenerate_GoalRisk(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{GoalRisk: 2})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-goal-risk" {
			found = true
			if n.Category != CatGoal {
				t.Errorf("category: expected %s, got %s", CatGoal, n.Category)
			}
			if n.Priority != P2High {
				t.Errorf("priority: expected %s, got %s", P2High, n.Priority)
			}
			if n.SourceEngine != "Projection" {
				t.Errorf("source: expected Projection, got %s", n.SourceEngine)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-goal-risk not generated when GoalRisk=2")
	}
}

func TestGenerate_NoGoalRisk(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{GoalRisk: 0, GoalOnTrack: 3, TotalGoals: 5})
	for _, n := range notifs {
		if n.NotifID == "notif-goal-risk" {
			t.Fatal("notif-goal-risk should not appear when GoalRisk=0")
		}
	}
}

func TestGenerate_AllGoalsOnTrack(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{GoalOnTrack: 5, TotalGoals: 5})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-goal-all" {
			found = true
			if n.Category != CatAchievement {
				t.Errorf("category: expected %s, got %s", CatAchievement, n.Category)
			}
			if n.Priority != P4Low {
				t.Errorf("priority: expected %s, got %s", P4Low, n.Priority)
			}
			if !strings.Contains(n.Summary, "5/5") {
				t.Errorf("summary should contain 5/5, got %q", n.Summary)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-goal-all not generated when all goals on track")
	}
}

func TestGenerate_AllGoalsOnTrack_NotWhenTotalGoalsZero(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{GoalOnTrack: 0, TotalGoals: 0})
	for _, n := range notifs {
		if n.NotifID == "notif-goal-all" {
			t.Fatal("notif-goal-all should not appear when TotalGoals=0")
		}
	}
}

// ---------------------------------------------------------------------------
// generate – cash flow & net worth
// ---------------------------------------------------------------------------

func TestGenerate_CashFlowDeficit(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{CashSurplus: -1000})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-cf-deficit" {
			found = true
			if n.Category != CatAlert {
				t.Errorf("category: expected %s, got %s", CatAlert, n.Category)
			}
			if n.Priority != P2High {
				t.Errorf("priority: expected %s, got %s", P2High, n.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-cf-deficit not generated when CashSurplus=-1000")
	}
}

func TestGenerate_NoCashFlowDeficit(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{CashSurplus: 0})
	for _, n := range notifs {
		if n.NotifID == "notif-cf-deficit" {
			t.Fatal("notif-cf-deficit should not appear when CashSurplus >= 0")
		}
	}
}

func TestGenerate_NetWorthDeclined(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{NetWorthChg: -5000})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-nw-down" {
			found = true
			if n.Category != CatFinancialEvt {
				t.Errorf("category: expected %s, got %s", CatFinancialEvt, n.Category)
			}
			if n.Priority != P3Medium {
				t.Errorf("priority: expected %s, got %s", P3Medium, n.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-nw-down not generated when NetWorthChg=-5000")
	}
}

func TestGenerate_NoNetWorthDeclined(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{NetWorthChg: 1000})
	for _, n := range notifs {
		if n.NotifID == "notif-nw-down" {
			t.Fatal("notif-nw-down should not appear when NetWorthChg >= 0")
		}
	}
}

// ---------------------------------------------------------------------------
// generate – recommendations
// ---------------------------------------------------------------------------

func TestGenerate_Recommendations(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true, RecCount: 3})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-recs" {
			found = true
			if n.Category != CatRecommend {
				t.Errorf("category: expected %s, got %s", CatRecommend, n.Category)
			}
			if n.Priority != P3Medium {
				t.Errorf("priority: expected %s, got %s", P3Medium, n.Priority)
			}
			if !strings.Contains(n.Summary, "3") {
				t.Errorf("summary should contain rec count, got %q", n.Summary)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-recs not generated when HasRecs=true")
	}
}

func TestGenerate_NoRecommendations(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: false})
	for _, n := range notifs {
		if n.NotifID == "notif-recs" {
			t.Fatal("notif-recs should not appear when HasRecs=false")
		}
	}
}

// ---------------------------------------------------------------------------
// generate – achievements
// ---------------------------------------------------------------------------

func TestGenerate_Achievements(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{AchieveCnt: 2})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-ach" {
			found = true
			if n.Category != CatAchievement {
				t.Errorf("category: expected %s, got %s", CatAchievement, n.Category)
			}
			if n.Priority != P4Low {
				t.Errorf("priority: expected %s, got %s", P4Low, n.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-ach not generated when AchieveCnt=2")
	}
}

func TestGenerate_NoAchievements(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{AchieveCnt: 0})
	for _, n := range notifs {
		if n.NotifID == "notif-ach" {
			t.Fatal("notif-ach should not appear when AchieveCnt=0")
		}
	}
}

// ---------------------------------------------------------------------------
// generate – optimization
// ---------------------------------------------------------------------------

func TestGenerate_Optimization(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasOpts: true})
	found := false
	for _, n := range notifs {
		if n.NotifID == "notif-opts" {
			found = true
			if n.Category != CatOptimization {
				t.Errorf("category: expected %s, got %s", CatOptimization, n.Category)
			}
			if n.Priority != P4Low {
				t.Errorf("priority: expected %s, got %s", P4Low, n.Priority)
			}
			break
		}
	}
	if !found {
		t.Fatal("notif-opts not generated when HasOpts=true")
	}
}

func TestGenerate_NoOptimization(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasOpts: false})
	for _, n := range notifs {
		if n.NotifID == "notif-opts" {
			t.Fatal("notif-opts should not appear when HasOpts=false")
		}
	}
}

// ---------------------------------------------------------------------------
// generate – default fields (Timestamp, State, Metadata)
// ---------------------------------------------------------------------------

func TestGenerate_DefaultFields(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasOpts: true, HasRecs: true})
	for _, n := range notifs {
		if n.NotifID == "notif-opts" || n.NotifID == "notif-recs" {
			if n.Timestamp == "" {
				t.Errorf("%s: Timestamp should be set", n.NotifID)
			}
			if n.State != StateUnread {
				t.Errorf("%s: State should default to Unread, got %s", n.NotifID, n.State)
			}
			if n.Metadata == nil {
				t.Errorf("%s: Metadata should be non-nil", n.NotifID)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// generate – all notifications simultaneous
// ---------------------------------------------------------------------------

func TestGenerate_AllNotifications(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{
		HealthScore: 39, HealthGrade: "D", HealthChange: -5,
		RiskScore: 85, RiskLevel: "Critical",
		GoalRisk: 3, GoalOnTrack: 5, TotalGoals: 5,
		CashSurplus: -1, NetWorthChg: -100,
		HasRecs: true, RecCount: 4,
		AchieveCnt: 1,
		HasOpts:    true,
	})
	ids := make(map[string]bool)
	for _, n := range notifs {
		ids[n.NotifID] = true
	}
	expected := []string{
		"notif-health-critical",
		"notif-risk-critical",
		"notif-goal-risk",
		"notif-goal-all",
		"notif-cf-deficit",
		"notif-nw-down",
		"notif-recs",
		"notif-ach",
		"notif-opts",
	}
	for _, id := range expected {
		if !ids[id] {
			t.Errorf("expected %s in generated notifications but not found", id)
		}
	}
}

// ---------------------------------------------------------------------------
// generate – zero inputs (no notifications)
// ---------------------------------------------------------------------------

func TestGenerate_ZeroInputs(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{})
	if len(notifs) != 0 {
		t.Errorf("expected 0 notifications for zero inputs, got %d", len(notifs))
	}
}

// ---------------------------------------------------------------------------
// buildPreferencesMap
// ---------------------------------------------------------------------------

func TestBuildPreferencesMap(t *testing.T) {
	c := &Composer{}
	prefs := []Preference{
		{Category: CatHealth, Enabled: true, MinPriority: P3Medium},
		{Category: CatRisk, Enabled: false, MinPriority: P1Critical},
	}
	m := c.buildPreferencesMap(prefs)
	if len(m) != 2 {
		t.Fatalf("expected map of length 2, got %d", len(m))
	}
	if p, ok := m[CatHealth]; !ok {
		t.Error("CatHealth missing from map")
	} else {
		if p.Enabled != true {
			t.Error("CatHealth.Enabled should be true")
		}
		if p.MinPriority != P3Medium {
			t.Error("CatHealth.MinPriority should be P3Medium")
		}
	}
	if p, ok := m[CatRisk]; !ok {
		t.Error("CatRisk missing from map")
	} else {
		if p.Enabled != false {
			t.Error("CatRisk.Enabled should be false")
		}
	}
}

func TestBuildPreferencesMap_Empty(t *testing.T) {
	c := &Composer{}
	m := c.buildPreferencesMap(nil)
	if m == nil {
		t.Fatal("buildPreferencesMap(nil) returned nil")
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

// ---------------------------------------------------------------------------
// applyPrefs – preference filtering
// ---------------------------------------------------------------------------

func TestApplyPrefs_DisabledCategory(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true, HasOpts: true})
	prefs := map[Category]Preference{
		CatRecommend:    {Category: CatRecommend, Enabled: false, MinPriority: P5Info},
		CatOptimization: {Category: CatOptimization, Enabled: true, MinPriority: P5Info},
	}
	filtered := c.applyPrefs(ctx(), "u1", notifs, prefs, emptyRepo())
	for _, n := range filtered {
		if n.NotifID == "notif-recs" {
			t.Fatal("notif-recs should be filtered out when CatRecommend is disabled")
		}
		if n.NotifID == "notif-opts" {
			// should remain
		}
	}
}

func TestApplyPrefs_MinPriorityFilters(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true})
	// P3Medium min priority → notifications with weight >= 3 (P3, P4, P5) are kept.
	// notif-recs has P3Medium weight=3 → 3 < 3 is false, so it's kept.
	prefs := map[Category]Preference{
		CatRecommend: {Category: CatRecommend, Enabled: true, MinPriority: P3Medium},
	}
	filtered := c.applyPrefs(ctx(), "u1", notifs, prefs, emptyRepo())
	found := false
	for _, n := range filtered {
		if n.NotifID == "notif-recs" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("notif-recs (P3) should pass when min is P3")
	}
}

func TestApplyPrefs_MinPriorityExcludesHigher(t *testing.T) {
	c := &Composer{}
	// Add health-critical notification (P1) alongside recs (P3) and opts (P4)
	notifs := c.generate(Inputs{HasRecs: true, HasOpts: true, HealthScore: 39, HealthGrade: "D", HealthChange: -10})
	// P2High min priority → only weight >= 2 (P2, P3, P4, P5) are kept; P1 (weight 1) is excluded.
	prefs := map[Category]Preference{
		CatHealth:       {Category: CatHealth, Enabled: true, MinPriority: P2High},
		CatRecommend:    {Category: CatRecommend, Enabled: true, MinPriority: P2High},
		CatOptimization: {Category: CatOptimization, Enabled: true, MinPriority: P2High},
	}
	filtered := c.applyPrefs(ctx(), "u1", notifs, prefs, emptyRepo())
	for _, n := range filtered {
		if n.NotifID == "notif-health-critical" {
			t.Error("notif-health-critical (P1) should be filtered when min is P2")
		}
	}
}

func TestApplyPrefs_NoPrefs(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true})
	filtered := c.applyPrefs(ctx(), "u1", notifs, nil, emptyRepo())
	if len(filtered) != 1 {
		t.Fatalf("expected 1 notification with nil prefs, got %d", len(filtered))
	}
}

func TestApplyPrefs_EmptyPrefs(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true})
	filtered := c.applyPrefs(ctx(), "u1", notifs, map[Category]Preference{}, emptyRepo())
	if len(filtered) != 1 {
		t.Fatalf("expected 1 notification with empty prefs, got %d", len(filtered))
	}
}

// ---------------------------------------------------------------------------
// applyPrefs – state management
// ---------------------------------------------------------------------------

func TestApplyPrefs_DismissedState(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true})
	repo := &mockStateRepo{states: map[string]State{"notif-recs": StateDismissed}}
	filtered := c.applyPrefs(ctx(), "u1", notifs, nil, repo)
	for _, n := range filtered {
		if n.NotifID == "notif-recs" {
			t.Fatal("dismissed notification should be filtered out")
		}
	}
}

func TestApplyPrefs_ArchivedState(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true})
	repo := &mockStateRepo{states: map[string]State{"notif-recs": StateArchived}}
	filtered := c.applyPrefs(ctx(), "u1", notifs, nil, repo)
	for _, n := range filtered {
		if n.NotifID == "notif-recs" {
			t.Fatal("archived notification should be filtered out")
		}
	}
}

func TestApplyPrefs_SnoozedState(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true})
	repo := &mockStateRepo{states: map[string]State{"notif-recs": StateSnoozed}}
	filtered := c.applyPrefs(ctx(), "u1", notifs, nil, repo)
	for _, n := range filtered {
		if n.NotifID == "notif-recs" {
			t.Fatal("snoozed notification should be filtered out")
		}
	}
}

func TestApplyPrefs_ReadState(t *testing.T) {
	c := &Composer{}
	notifs := c.generate(Inputs{HasRecs: true})
	repo := &mockStateRepo{states: map[string]State{"notif-recs": StateRead}}
	filtered := c.applyPrefs(ctx(), "u1", notifs, nil, repo)
	found := false
	for _, n := range filtered {
		if n.NotifID == "notif-recs" {
			found = true
			if n.State != StateRead {
				t.Errorf("state should be Read, got %s", n.State)
			}
			break
		}
	}
	if !found {
		t.Fatal("read notification should remain in feed")
	}
}

// ---------------------------------------------------------------------------
// buildFeed – sorting and aggregation
// ---------------------------------------------------------------------------

func TestBuildFeed_Sorting(t *testing.T) {
	c := &Composer{}
	notifs := []Notification{
		{NotifID: "a", Priority: P4Low, Timestamp: "2025-01-01T00:00:00Z"},
		{NotifID: "b", Priority: P1Critical, Timestamp: "2025-01-03T00:00:00Z"},
		{NotifID: "c", Priority: P3Medium, Timestamp: "2025-01-02T00:00:00Z"},
		{NotifID: "d", Priority: P1Critical, Timestamp: "2025-01-01T00:00:00Z"},
	}
	center := c.buildFeed(ctx(), "u1", notifs, nil, emptyRepo(), 10, "")
	items := center.Items
	// Expected order: P1 first (b then d by timestamp), then P3 (c), then P4 (a)
	if len(items) != 4 {
		t.Fatalf("expected 4 items, got %d", len(items))
	}
	if items[0].NotifID != "b" {
		t.Errorf("first should be b (P1, latest), got %s", items[0].NotifID)
	}
	if items[1].NotifID != "d" {
		t.Errorf("second should be d (P1, earlier), got %s", items[1].NotifID)
	}
	if items[2].NotifID != "c" {
		t.Errorf("third should be c (P3), got %s", items[2].NotifID)
	}
	if items[3].NotifID != "a" {
		t.Errorf("fourth should be a (P4), got %s", items[3].NotifID)
	}
}

func TestBuildFeed_AggregationCounts(t *testing.T) {
	c := &Composer{}
	notifs := []Notification{
		{NotifID: "n1", Priority: P1Critical, Category: CatHealth, State: StateUnread},
		{NotifID: "n2", Priority: P2High, Category: CatRisk, State: StateRead},
		{NotifID: "n3", Priority: P1Critical, Category: CatHealth, State: StateUnread},
		{NotifID: "n4", Priority: P3Medium, Category: CatGoal, State: StateRead},
	}
	center := c.buildFeed(ctx(), "u1", notifs, nil, emptyRepo(), 10, "")
	if center.UnreadCount != 2 {
		t.Errorf("UnreadCount: expected 2, got %d", center.UnreadCount)
	}
	if center.ByPriority[P1Critical] != 2 {
		t.Errorf("ByPriority[P1]: expected 2, got %d", center.ByPriority[P1Critical])
	}
	if center.ByCategory[CatHealth] != 2 {
		t.Errorf("ByCategory[Health]: expected 2, got %d", center.ByCategory[CatHealth])
	}
}

// ---------------------------------------------------------------------------
// buildFeed – pagination
// ---------------------------------------------------------------------------

func TestBuildFeed_Pagination(t *testing.T) {
	c := &Composer{}
	notifs := make([]Notification, 10)
	for i := 0; i < 10; i++ {
		notifs[i] = Notification{
			NotifID:   fmtID(i),
			Priority:  P5Info,
			Timestamp: fmtTS(i),
		}
	}
	center := c.buildFeed(ctx(), "u1", notifs, nil, emptyRepo(), 3, "")
	if len(center.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(center.Items))
	}
	if !center.HasMore {
		t.Error("HasMore should be true")
	}
	if center.Cursor == "" {
		t.Error("Cursor should be non-empty when HasMore")
	}
}

func TestBuildFeed_PaginationWithCursor(t *testing.T) {
	c := &Composer{}
	notifs := make([]Notification, 10)
	for i := 0; i < 10; i++ {
		notifs[i] = Notification{
			NotifID:   fmtID(i),
			Priority:  P5Info,
			Timestamp: fmtTS(i),
		}
	}
	// Items are sorted by timestamp descending (all P5), so order is:
	// n-9, n-8, n-7, n-6, n-5, n-4, n-3, n-2, n-1, n-0
	// cursor "n-4" → start = index 5+1 = 6 → items n-3, n-2, n-1
	center := c.buildFeed(ctx(), "u1", notifs, nil, emptyRepo(), 3, fmtID(4))
	if len(center.Items) != 3 {
		t.Fatalf("expected 3 items after cursor, got %d", len(center.Items))
	}
	if center.Items[0].NotifID != fmtID(3) {
		t.Errorf("first after cursor should be %s, got %s", fmtID(3), center.Items[0].NotifID)
	}
	if !center.HasMore {
		t.Error("HasMore should be true after cursor with remaining items")
	}
}

func TestBuildFeed_PaginationLastPage(t *testing.T) {
	c := &Composer{}
	notifs := make([]Notification, 5)
	for i := 0; i < 5; i++ {
		notifs[i] = Notification{
			NotifID:   fmtID(i),
			Priority:  P5Info,
			Timestamp: fmtTS(i),
		}
	}
	center := c.buildFeed(ctx(), "u1", notifs, nil, emptyRepo(), 5, "")
	if center.HasMore {
		t.Error("HasMore should be false on last page")
	}
	if center.Cursor != "" {
		t.Errorf("Cursor should be empty on last page, got %s", center.Cursor)
	}
}

func TestBuildFeed_DefaultLimit(t *testing.T) {
	c := &Composer{}
	notifs := make([]Notification, 60)
	for i := 0; i < 60; i++ {
		notifs[i] = Notification{
			NotifID:   fmtID(i),
			Priority:  P5Info,
			Timestamp: fmtTS(i),
		}
	}
	center := c.buildFeed(ctx(), "u1", notifs, nil, emptyRepo(), 0, "")
	if len(center.Items) != 50 {
		t.Fatalf("expected default limit 50, got %d", len(center.Items))
	}
	if !center.HasMore {
		t.Error("HasMore should be true with 60 notifs and limit 50")
	}
}

func TestBuildFeed_Empty(t *testing.T) {
	c := &Composer{}
	center := c.buildFeed(ctx(), "u1", nil, nil, emptyRepo(), 10, "")
	if center == nil {
		t.Fatal("buildFeed should return non-nil center")
	}
	if len(center.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(center.Items))
	}
	if center.UnreadCount != 0 {
		t.Errorf("UnreadCount should be 0, got %d", center.UnreadCount)
	}
	if center.HasMore {
		t.Error("HasMore should be false")
	}
}

// ---------------------------------------------------------------------------
// BuildCenter – integration
// ---------------------------------------------------------------------------

func TestBuildCenter(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID: "u1",
		HealthScore: 39, HealthGrade: "D", HealthChange: -5,
		RiskScore: 85, RiskLevel: "High",
		GoalRisk: 2, GoalOnTrack: 3, TotalGoals: 5,
		CashSurplus: -500, NetWorthChg: -1000,
		HasRecs: true, RecCount: 3,
		AchieveCnt: 2,
		HasOpts:    true,
	}
	center := c.BuildCenter(ctx(), inputs, nil, emptyRepo(), 10, "")
	if center == nil {
		t.Fatal("BuildCenter returned nil")
	}
	if len(center.Items) == 0 {
		t.Fatal("BuildCenter returned empty items")
	}
	if center.UnreadCount <= 0 {
		t.Errorf("expected unread count > 0, got %d", center.UnreadCount)
	}
	if len(center.ByPriority) == 0 {
		t.Error("ByPriority should be non-empty")
	}
	if len(center.ByCategory) == 0 {
		t.Error("ByCategory should be non-empty")
	}
}

func TestBuildCenter_PreferenceFiltering(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:  "u1",
		HasRecs: true, RecCount: 3,
	}
	// Disable recommendations
	prefs := map[Category]Preference{
		CatRecommend: {Category: CatRecommend, Enabled: false, MinPriority: P5Info},
	}
	center := c.BuildCenter(ctx(), inputs, prefs, emptyRepo(), 10, "")
	for _, n := range center.Items {
		if n.NotifID == "notif-recs" {
			t.Fatal("notif-recs should be filtered out by preferences")
		}
	}
}

func TestBuildCenter_StateFiltering(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:  "u1",
		HasRecs: true,
		HasOpts: true,
	}
	repo := &mockStateRepo{states: map[string]State{"notif-recs": StateDismissed}}
	center := c.BuildCenter(ctx(), inputs, nil, repo, 10, "")
	for _, n := range center.Items {
		if n.NotifID == "notif-recs" {
			t.Fatal("dismissed notif-recs should be filtered out")
		}
	}
}

// ---------------------------------------------------------------------------
// BuildUnread
// ---------------------------------------------------------------------------

func TestBuildUnread(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:     "u1",
		HasRecs:    true,
		AchieveCnt: 1,
	}
	center := c.BuildUnread(ctx(), inputs, nil, emptyRepo())
	if center == nil {
		t.Fatal("BuildUnread returned nil")
	}
	if len(center.Items) == 0 {
		t.Fatal("BuildUnread returned empty items")
	}
}

// ---------------------------------------------------------------------------
// BuildHistory
// ---------------------------------------------------------------------------

func TestBuildHistory(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:     "u1",
		HasRecs:    true,
		AchieveCnt: 1,
	}
	center := c.BuildHistory(ctx(), inputs, nil, emptyRepo(), 5, "")
	if center == nil {
		t.Fatal("BuildHistory returned nil")
	}
	if len(center.Items) == 0 {
		t.Fatal("BuildHistory returned empty items")
	}
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

func TestSearch_ByTitle(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:  "u1",
		HasRecs: true,
		HasOpts: true,
	}
	results := c.Search(ctx(), inputs, nil, emptyRepo(), "Recommendation")
	if len(results) == 0 {
		t.Fatal("expected match for 'Recommendation' in title")
	}
}

func TestSearch_BySummary(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:    "u1",
		HasRecs:   true,
		RecCount:  5,
	}
	results := c.Search(ctx(), inputs, nil, emptyRepo(), "5 recommendation")
	if len(results) == 0 {
		t.Fatal("expected match for '5 recommendation' in summary")
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:  "u1",
		HasRecs: true,
		HasOpts: true,
	}
	resultsUpper := c.Search(ctx(), inputs, nil, emptyRepo(), "OPTIMIZATION")
	resultsLower := c.Search(ctx(), inputs, nil, emptyRepo(), "optimization")
	if len(resultsUpper) != len(resultsLower) {
		t.Errorf("case-insensitive search gave different counts: upper=%d, lower=%d", len(resultsUpper), len(resultsLower))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:  "u1",
		HasRecs: true,
	}
	results := c.Search(ctx(), inputs, nil, emptyRepo(), "zzz_nonexistent_zzz")
	if len(results) != 0 {
		t.Errorf("expected 0 results for non-matching query, got %d", len(results))
	}
}

func TestSearch_FiltersByPreferences(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:  "u1",
		HasRecs: true,
		HasOpts: true,
	}
	prefs := map[Category]Preference{
		CatOptimization: {Category: CatOptimization, Enabled: false, MinPriority: P5Info},
	}
	results := c.Search(ctx(), inputs, prefs, emptyRepo(), "Optimization")
	if len(results) != 0 {
		t.Error("should return no results when optimization is disabled in prefs")
	}
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

func TestGetByID_Found(t *testing.T) {
	c := NewComposer()
	inputs := Inputs{
		UserID:  "u1",
		HasOpts: true,
	}
	n := c.GetByID(ctx(), inputs, nil, emptyRepo(), "notif-opts")
	if n == nil {
		t.Fatal("GetByID should find notif-opts")
	}
	if n.NotifID != "notif-opts" {
		t.Errorf("NotifID: expected notif-opts, got %s", n.NotifID)
	}
	if n.Category != CatOptimization {
		t.Errorf("Category: expected %s, got %s", CatOptimization, n.Category)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	c := &Composer{}
	inputs := Inputs{
		UserID:  "u1",
		HasOpts: true,
	}
	n := c.GetByID(ctx(), inputs, nil, emptyRepo(), "non-existent-id")
	if n != nil {
		t.Fatal("GetByID should return nil for unknown ID")
	}
}

func TestGetByID_AppliesState(t *testing.T) {
	c := &Composer{}
	inputs := Inputs{
		UserID:  "u1",
		HasRecs: true,
	}
	repo := &mockStateRepo{states: map[string]State{"notif-recs": StateRead}}
	n := c.GetByID(ctx(), inputs, nil, repo, "notif-recs")
	if n == nil {
		t.Fatal("GetByID should find notif-recs")
	}
	if n.State != StateRead {
		t.Errorf("State should be Read from repo, got %s", n.State)
	}
}

// ---------------------------------------------------------------------------
// All 17 notification category constants
// ---------------------------------------------------------------------------

func TestAllCategoriesDefined(t *testing.T) {
	categories := []Category{
		CatGoal, CatAccount, CatPortfolio, CatAsset, CatLiability,
		CatFinancialEvt, CatHealth, CatRisk, CatRecommend, CatOptimization,
		CatSimulation, CatReminder, CatAlert, CatAchievement, CatMilestone,
		CatSystem, CatInsight,
	}
	seen := make(map[Category]bool)
	for _, cat := range categories {
		if seen[cat] {
			t.Errorf("duplicate category constant: %s", cat)
		}
		seen[cat] = true
	}
	if len(seen) != 17 {
		t.Errorf("expected 17 unique categories, got %d", len(seen))
	}
}

// ---------------------------------------------------------------------------
// All 5 priority constants
// ---------------------------------------------------------------------------

func TestAllPrioritiesDefined(t *testing.T) {
	priorities := []Priority{P1Critical, P2High, P3Medium, P4Low, P5Info}
	seen := make(map[Priority]bool)
	for _, p := range priorities {
		if seen[p] {
			t.Errorf("duplicate priority constant: %s", p)
		}
		seen[p] = true
	}
	if len(seen) != 5 {
		t.Errorf("expected 5 unique priorities, got %d", len(seen))
	}
}

// ---------------------------------------------------------------------------
// All 5 state constants
// ---------------------------------------------------------------------------

func TestAllStatesDefined(t *testing.T) {
	states := []State{StateUnread, StateRead, StateArchived, StateDismissed, StateSnoozed}
	seen := make(map[State]bool)
	for _, s := range states {
		if seen[s] {
			t.Errorf("duplicate state constant: %s", s)
		}
		seen[s] = true
	}
	if len(seen) != 5 {
		t.Errorf("expected 5 unique states, got %d", len(seen))
	}
}

// ---------------------------------------------------------------------------
// Table-driven generate scenarios
// ---------------------------------------------------------------------------

func TestGenerateScenarios(t *testing.T) {
	tests := []struct {
		name     string
		inputs   Inputs
		expected []string // notification IDs that should be present
		not      []string // notification IDs that should NOT be present
	}{
		{
			name:     "no triggers",
			inputs:   Inputs{},
			expected: nil,
		},
		{
			name:     "health only critical",
			inputs:   Inputs{HealthScore: 25, HealthGrade: "F", HealthChange: -10},
			expected: []string{"notif-health-critical"},
			not:      []string{"notif-health-down"},
		},
		{
			name:     "health only decline",
			inputs:   Inputs{HealthScore: 60, HealthGrade: "B", HealthChange: -5},
			expected: []string{"notif-health-down"},
			not:      []string{"notif-health-critical"},
		},
		{
			name:     "risk critical",
			inputs:   Inputs{RiskScore: 90, RiskLevel: "Critical"},
			expected: []string{"notif-risk-critical"},
			not:      []string{"notif-risk-high"},
		},
		{
			name:     "risk high only",
			inputs:   Inputs{RiskScore: 65, RiskLevel: "Elevated"},
			expected: []string{"notif-risk-high"},
			not:      []string{"notif-risk-critical"},
		},
		{
			name:     "goal risk on track partial",
			inputs:   Inputs{GoalRisk: 1, GoalOnTrack: 3, TotalGoals: 5},
			expected: []string{"notif-goal-risk"},
			not:      []string{"notif-goal-all"},
		},
		{
			name:     "all goals on track no risk",
			inputs:   Inputs{GoalRisk: 0, GoalOnTrack: 5, TotalGoals: 5},
			expected: []string{"notif-goal-all"},
		},
		{
			name:     "cash flow deficit",
			inputs:   Inputs{CashSurplus: -1},
			expected: []string{"notif-cf-deficit"},
		},
		{
			name:     "net worth declined",
			inputs:   Inputs{NetWorthChg: -1},
			expected: []string{"notif-nw-down"},
		},
		{
			name:     "recommendations only",
			inputs:   Inputs{HasRecs: true, RecCount: 2},
			expected: []string{"notif-recs"},
		},
		{
			name:     "achievements only",
			inputs:   Inputs{AchieveCnt: 3},
			expected: []string{"notif-ach"},
		},
		{
			name:     "optimization only",
			inputs:   Inputs{HasOpts: true},
			expected: []string{"notif-opts"},
		},
	}

	c := &Composer{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifs := c.generate(tt.inputs)
			ids := make(map[string]bool)
			for _, n := range notifs {
				ids[n.NotifID] = true
			}
			for _, id := range tt.expected {
				if !ids[id] {
					t.Errorf("expected notification %q not found", id)
				}
			}
			for _, id := range tt.not {
				if ids[id] {
					t.Errorf("unexpected notification %q found", id)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Table-driven priority assignment for generated notifications
// ---------------------------------------------------------------------------

func TestGeneratedNotificationPriorities(t *testing.T) {
	tests := []struct {
		name     string
		inputs   Inputs
		checks   map[string]Priority // notifID -> expected priority
	}{
		{
			name:   "P1 critical notifications",
			inputs: Inputs{HealthScore: 30, HealthGrade: "F", RiskScore: 85, RiskLevel: "High"},
			checks: map[string]Priority{
				"notif-health-critical": P1Critical,
				"notif-risk-critical":   P1Critical,
			},
		},
		{
			name: "P2 high notifications",
			inputs: Inputs{
				HealthScore: 60, HealthGrade: "C", HealthChange: -10,
				RiskScore: 70, RiskLevel: "Elevated",
				GoalRisk: 2,
				CashSurplus: -100,
			},
			checks: map[string]Priority{
				"notif-health-down": P2High,
				"notif-risk-high":   P2High,
				"notif-goal-risk":   P2High,
				"notif-cf-deficit":  P2High,
			},
		},
		{
			name: "P3 medium notifications",
			inputs: Inputs{
				NetWorthChg: -500,
				HasRecs:     true,
			},
			checks: map[string]Priority{
				"notif-nw-down": P3Medium,
				"notif-recs":    P3Medium,
			},
		},
		{
			name: "P4 low notifications",
			inputs: Inputs{
				GoalOnTrack: 5, TotalGoals: 5,
				AchieveCnt: 2,
				HasOpts:    true,
			},
			checks: map[string]Priority{
				"notif-goal-all": P4Low,
				"notif-ach":      P4Low,
				"notif-opts":     P4Low,
			},
		},
	}

	c := &Composer{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifs := c.generate(tt.inputs)
			byID := make(map[string]Notification)
			for _, n := range notifs {
				byID[n.NotifID] = n
			}
			for id, wantPri := range tt.checks {
				n, ok := byID[id]
				if !ok {
					t.Errorf("notification %q not generated", id)
					continue
				}
				if n.Priority != wantPri {
					t.Errorf("%s: expected priority %s, got %s", id, wantPri, n.Priority)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func fmtID(i int) string {
	return fmt.Sprintf("n-%d", i)
}

func fmtTS(i int) string {
	return time.Date(2025, 1, 1, 0, 0, i, 0, time.UTC).Format(time.RFC3339)
}
