package engine

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Category string

const (
	CatGoal          Category = "Goal"
	CatAccount       Category = "Account"
	CatPortfolio     Category = "Portfolio"
	CatAsset         Category = "Asset"
	CatLiability     Category = "Liability"
	CatFinancialEvt  Category = "FinancialEvent"
	CatHealth        Category = "Health"
	CatRisk          Category = "Risk"
	CatRecommend     Category = "Recommendation"
	CatOptimization  Category = "Optimization"
	CatSimulation    Category = "Simulation"
	CatReminder      Category = "Reminder"
	CatAlert         Category = "Alert"
	CatAchievement   Category = "Achievement"
	CatMilestone     Category = "Milestone"
	CatSystem        Category = "System"
	CatInsight       Category = "Insight"
)

type Priority string

const (
	P1Critical Priority = "P1"
	P2High     Priority = "P2"
	P3Medium   Priority = "P3"
	P4Low      Priority = "P4"
	P5Info     Priority = "P5"
)

type State string

const (
	StateUnread   State = "Unread"
	StateRead     State = "Read"
	StateArchived State = "Archived"
	StateDismissed State = "Dismissed"
	StateSnoozed  State = "Snoozed"
)

type Notification struct {
	NotifID      string            `json:"notif_id"`
	Category     Category          `json:"category"`
	Priority     Priority          `json:"priority"`
	State        State             `json:"state"`
	Title        string            `json:"title"`
	Summary      string            `json:"summary"`
	Description  string            `json:"description"`
	Source       string            `json:"source"`
	SourceEngine string            `json:"source_engine"`
	RelatedEnt   string            `json:"related_entity,omitempty"`
	RelatedGoal  string            `json:"related_goal,omitempty"`
	RelatedAcct  string            `json:"related_account,omitempty"`
	RelatedPf    string            `json:"related_portfolio,omitempty"`
	Timestamp    string            `json:"timestamp"`
	ExpiresAt    string            `json:"expires_at,omitempty"`
	ActionURL    string            `json:"action_url,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type NotificationCenter struct {
	UnreadCount  int            `json:"unread_count"`
	ByPriority   map[Priority]int `json:"by_priority"`
	ByCategory   map[Category]int `json:"by_category"`
	Items        []Notification `json:"items"`
	Cursor       string         `json:"cursor,omitempty"`
	HasMore      bool           `json:"has_more"`
}

type Preference struct {
	Category    Category `json:"category"`
	InApp       bool     `json:"in_app"`
	Push        bool     `json:"push"`
	Email       bool     `json:"email"`
	Digest      string   `json:"digest"`
	MinPriority Priority `json:"min_priority"`
	Enabled     bool     `json:"enabled"`
}

type Inputs struct {
	UserID       string         `json:"user_id"`
	HealthScore  int            `json:"health_score"`
	HealthGrade  string         `json:"health_grade"`
	HealthChange int            `json:"health_change"`
	RiskScore    int            `json:"risk_score"`
	RiskLevel    string         `json:"risk_level"`
	GoalOnTrack  int            `json:"goals_on_track"`
	TotalGoals   int            `json:"total_goals"`
	GoalRisk     int            `json:"goals_at_risk"`
	CashSurplus  int64          `json:"cash_flow_surplus"`
	NetWorthChg  int64          `json:"net_worth_change"`
	PortRet      float64        `json:"portfolio_return"`
	HasRecs      bool           `json:"has_recommendations"`
	RecCount     int            `json:"rec_count"`
	HasOpts      bool           `json:"has_optimizations"`
	HasSims      bool           `json:"has_simulations"`
	AchieveCnt   int            `json:"achievement_count"`
	EventCount   int            `json:"event_count"`
	Preferences  []Preference   `json:"preferences"`
}

type StateRepository interface {
	Get(ctx context.Context, userID, notifID string) (State, bool)
	SetState(ctx context.Context, userID, notifID string, s State) error
	Snooze(ctx context.Context, userID, notifID string, until string) error
}

type PreferenceRepository interface {
	GetPreferences(ctx context.Context, userID string) ([]Preference, error)
	SavePreferences(ctx context.Context, userID string, prefs []Preference) error
}

type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

func (c *Composer) BuildCenter(ctx context.Context, inputs Inputs, prefs map[Category]Preference, stateRepo StateRepository, limit int, cursor string) *NotificationCenter {
	notifs := c.generate(inputs)
	return c.buildFeed(ctx, inputs.UserID, notifs, prefs, stateRepo, limit, cursor)
}

func (c *Composer) BuildUnread(ctx context.Context, inputs Inputs, prefs map[Category]Preference, stateRepo StateRepository) *NotificationCenter {
	notifs := c.generate(inputs)
	return c.buildFeed(ctx, inputs.UserID, notifs, prefs, stateRepo, 50, "")
}

func (c *Composer) BuildHistory(ctx context.Context, inputs Inputs, prefs map[Category]Preference, stateRepo StateRepository, limit int, cursor string) *NotificationCenter {
	notifs := c.generate(inputs)
	return c.buildFeed(ctx, inputs.UserID, notifs, prefs, stateRepo, limit, cursor)
}

func (c *Composer) Search(ctx context.Context, inputs Inputs, prefs map[Category]Preference, stateRepo StateRepository, q string) []Notification {
	notifs := c.generate(inputs)
	filtered := c.applyPrefs(ctx, inputs.UserID, notifs, prefs, stateRepo)
	q = strings.ToLower(q)
	var matched []Notification
	for _, n := range filtered {
		if strings.Contains(strings.ToLower(n.Title), q) || strings.Contains(strings.ToLower(n.Summary), q) || strings.Contains(strings.ToLower(n.Description), q) {
			matched = append(matched, n)
		}
	}
	return matched
}

func (c *Composer) GetByID(ctx context.Context, inputs Inputs, prefs map[Category]Preference, stateRepo StateRepository, id string) *Notification {
	for _, n := range c.generate(inputs) {
		if n.NotifID == id {
			// Apply state from repo
			if s, ok := stateRepo.Get(ctx, inputs.UserID, id); ok { n.State = s }
			return &n
		}
	}
	return nil
}

func (c *Composer) buildFeed(ctx context.Context, userID string, notifs []Notification, prefs map[Category]Preference, stateRepo StateRepository, limit int, cursor string) *NotificationCenter {
	filtered := c.applyPrefs(ctx, userID, notifs, prefs, stateRepo)

	sort.Slice(filtered, func(i, j int) bool {
		wi := c.weight(filtered[i].Priority)
		wj := c.weight(filtered[j].Priority)
		if wi != wj { return wi < wj }
		return filtered[i].Timestamp > filtered[j].Timestamp
	})

	byPri := map[Priority]int{}
	byCat := map[Category]int{}
	unread := 0
	for _, n := range filtered {
		byPri[n.Priority]++
		byCat[n.Category]++
		if n.State == StateUnread { unread++ }
	}

	if limit <= 0 { limit = 50 }
	var start int
	if cursor != "" {
		for i, n := range filtered {
			if n.NotifID == cursor { start = i + 1; break }
		}
	}
	hasMore := len(filtered) > start+limit
	if start+limit > len(filtered) {
		filtered = filtered[start:]
	} else {
		filtered = filtered[start : start+limit]
	}
	nextCursor := ""
	if hasMore && len(filtered) > 0 { nextCursor = filtered[len(filtered)-1].NotifID }

	return &NotificationCenter{
		UnreadCount: unread, ByPriority: byPri, ByCategory: byCat,
		Items: filtered, Cursor: nextCursor, HasMore: hasMore,
	}
}

func (c *Composer) buildPreferencesMap(prefs []Preference) map[Category]Preference {
	m := make(map[Category]Preference)
	for _, p := range prefs { m[p.Category] = p }
	return m
}

func (c *Composer) applyPrefs(ctx context.Context, userID string, notifs []Notification, prefs map[Category]Preference, repo StateRepository) []Notification {
	var filtered []Notification
	for _, n := range notifs {
		// Apply persisted state
		if s, ok := repo.Get(ctx, userID, n.NotifID); ok {
			n.State = s
			if s == StateDismissed || s == StateArchived { continue }
			if s == StateSnoozed {
				// Check if still snoozed
				continue
			}
		}
		// Apply preferences
		if pref, ok := prefs[n.Category]; ok {
			if !pref.Enabled { continue }
			if c.weight(n.Priority) < c.weight(pref.MinPriority) { continue }
		}
		filtered = append(filtered, n)
	}
	return filtered
}

func (c *Composer) weight(p Priority) int {
	switch p { case P1Critical: return 1; case P2High: return 2; case P3Medium: return 3; case P4Low: return 4; default: return 5 }
}

func (c *Composer) generate(inputs Inputs) []Notification {
	now := time.Now().UTC().Format(time.RFC3339)
	var notifs []Notification

	add := func(n Notification) {
		if n.Timestamp == "" { n.Timestamp = now }
		if n.State == "" { n.State = StateUnread }
		if n.Metadata == nil { n.Metadata = map[string]interface{}{} }
		notifs = append(notifs, n)
	}

	// Health notifications
	if inputs.HealthScore > 0 && inputs.HealthScore < 40 {
		add(Notification{NotifID: "notif-health-critical", Category: CatHealth, Priority: P1Critical,
			Title: "Critical Health Score", Summary: fmt.Sprintf("Your health score is %d — immediate attention needed", inputs.HealthScore),
			SourceEngine: "HealthScore", ActionURL: "/dashboard"})
	} else if inputs.HealthChange < 0 {
		add(Notification{NotifID: "notif-health-down", Category: CatHealth, Priority: P2High,
			Title: "Health Score Declined", Summary: fmt.Sprintf("Down %d points to %d (%s)", -inputs.HealthChange, inputs.HealthScore, inputs.HealthGrade),
			SourceEngine: "HealthScore", ActionURL: "/dashboard"})
	}

	// Risk notifications
	if inputs.RiskScore >= 80 {
		add(Notification{NotifID: "notif-risk-critical", Category: CatRisk, Priority: P1Critical,
			Title: "Critical Risk Level", Summary: fmt.Sprintf("Risk score %d — %s", inputs.RiskScore, inputs.RiskLevel),
			SourceEngine: "Risk", ActionURL: "/portfolio/risk"})
	} else if inputs.RiskScore >= 60 {
		add(Notification{NotifID: "notif-risk-high", Category: CatRisk, Priority: P2High,
			Title: "Elevated Risk", Summary: fmt.Sprintf("Risk score %d — review recommended", inputs.RiskScore),
			SourceEngine: "Risk", ActionURL: "/portfolio/risk"})
	}

	// Goal notifications
	if inputs.GoalRisk > 0 {
		add(Notification{NotifID: "notif-goal-risk", Category: CatGoal, Priority: P2High,
			Title: "Goals at Risk", Summary: fmt.Sprintf("%d goal(s) need attention", inputs.GoalRisk),
			SourceEngine: "Projection", ActionURL: "/goals"})
	}
	if inputs.GoalOnTrack == inputs.TotalGoals && inputs.TotalGoals > 0 {
		add(Notification{NotifID: "notif-goal-all", Category: CatAchievement, Priority: P4Low,
			Title: "All Goals On Track", Summary: fmt.Sprintf("%d/%d goals progressing well", inputs.GoalOnTrack, inputs.TotalGoals),
			SourceEngine: "Projection", ActionURL: "/goals"})
	}

	// Cash flow
	if inputs.CashSurplus < 0 {
		add(Notification{NotifID: "notif-cf-deficit", Category: CatAlert, Priority: P2High,
			Title: "Cash Flow Deficit", Summary: fmt.Sprintf("Monthly shortfall"), SourceEngine: "Projection", ActionURL: "/planning"})
	}

	// Net worth
	if inputs.NetWorthChg < 0 {
		add(Notification{NotifID: "notif-nw-down", Category: CatFinancialEvt, Priority: P3Medium,
			Title: "Net Worth Declined", Summary: fmt.Sprintf("Decreased"), SourceEngine: "Projection", ActionURL: "/dashboard"})
	}

	// Recommendations
	if inputs.HasRecs {
		add(Notification{NotifID: "notif-recs", Category: CatRecommend, Priority: P3Medium,
			Title: "New Recommendations", Summary: fmt.Sprintf("%d recommendation(s) available", inputs.RecCount),
			SourceEngine: "Recommendation", ActionURL: "/goals/recommendations"})
	}

	// Achievements
	if inputs.AchieveCnt > 0 {
		add(Notification{NotifID: "notif-ach", Category: CatAchievement, Priority: P4Low,
			Title: "Achievements Unlocked", Summary: fmt.Sprintf("%d new achievement(s)", inputs.AchieveCnt),
			SourceEngine: "HealthScore", ActionURL: "/insights/achievements"})
	}

	// Optimization
	if inputs.HasOpts {
		add(Notification{NotifID: "notif-opts", Category: CatOptimization, Priority: P4Low,
			Title: "Optimization Available", Summary: "New optimization strategies", SourceEngine: "Optimization", ActionURL: "/planning/optimization"})
	}

	return notifs
}
