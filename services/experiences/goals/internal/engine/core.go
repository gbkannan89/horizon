package engine

import "fmt"

// CardType represents goal dashboard card types.
type CardType string

const (
	CTOverview       CardType = "Overview"
	CTProgress       CardType = "Progress"
	CTFunding        CardType = "Funding"
	CTProjection     CardType = "Projection"
	CTRecommendation CardType = "Recommendation"
	CTOptimization   CardType = "Optimization"
	CTTimeline       CardType = "Timeline"
	CTMilestone      CardType = "Milestone"
)

// Card represents a goal dashboard card.
type Card struct {
	CardID   string   `json:"card_id"`
	CardType CardType `json:"card_type"`
	Title    string   `json:"title"`
	Summary  string   `json:"summary"`
	Data     interface{} `json:"data,omitempty"`
	Priority int      `json:"priority"`
}

// GoalDashboard is the full dashboard for a single goal.
type GoalDashboard struct {
	GoalID      string `json:"goal_id"`
	GoalName    string `json:"goal_name"`
	Progress    int    `json:"progress"`
	Status      string `json:"status"`
	Importance  string `json:"importance"`
	GoalType    string `json:"goal_type"`
	TargetAmount float64 `json:"target_amount"`
	CurrentValue float64 `json:"current_value"`
	FundingGap   float64 `json:"funding_gap"`
	TargetDate   string `json:"target_date,omitempty"`
	Age         int    `json:"age_days,omitempty"`
	Cards       []Card `json:"cards"`
}

// GoalProgress tracks progress details.
type GoalProgress struct {
	GoalID            string  `json:"goal_id"`
	GoalName          string  `json:"goal_name"`
	ProgressPct       int     `json:"progress_pct"`
	CurrentValue      float64 `json:"current_value"`
	TargetAmount      float64 `json:"target_amount"`
	RemainingAmount   float64 `json:"remaining_amount"`
	Status            string  `json:"status"`
	MonthlyContribution float64 `json:"monthly_contribution,omitempty"`
	MonthsToTarget    int     `json:"months_to_target,omitempty"`
}

// GoalProjectionView wraps projection engine data.
type GoalProjectionView struct {
	GoalID              string  `json:"goal_id"`
	ProjectedDate       string  `json:"projected_date,omitempty"`
	ProjectedValue      float64 `json:"projected_value,omitempty"`
	Confidence          string  `json:"confidence,omitempty"`
	OnTrack             bool    `json:"on_track"`
	MonthlyContribution float64 `json:"monthly_contribution,omitempty"`
}

// GoalRecommendationView wraps recommendation engine data.
type GoalRecommendationView struct {
	GoalID      string        `json:"goal_id"`
	Recs        []RecItem     `json:"recommendations"`
	TotalCount  int           `json:"total_count"`
}

type RecItem struct {
	RecID    string `json:"rec_id"`
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	Priority int    `json:"priority"`
	Impact   string `json:"impact"`
}

// GoalOptimizationView wraps optimization engine data.
type GoalOptimizationView struct {
	GoalID            string      `json:"goal_id"`
	Opts              []OptItem   `json:"optimizations"`
	TotalCount        int         `json:"total_count"`
}

type OptItem struct {
	OptID    string  `json:"opt_id"`
	Strategy string  `json:"strategy"`
	Score    float64 `json:"score"`
	Summary  string  `json:"summary"`
}

// GoalTimelineView wraps recent goal timeline events.
type GoalTimelineView struct {
	GoalID     string       `json:"goal_id"`
	Events     []GoalEvent  `json:"events"`
	TotalCount int          `json:"total_count"`
}

type GoalEvent struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
}

// GoalMilestoneView wraps milestone data.
type GoalMilestoneView struct {
	GoalID     string       `json:"goal_id"`
	Milestones []Milestone  `json:"milestones"`
	TotalCount int          `json:"total_count"`
}

type Milestone struct {
	MilestoneID string `json:"milestone_id"`
	Title       string `json:"title"`
	ProgressPct int    `json:"progress_pct"`
	Reached     bool   `json:"reached"`
}

// Inputs aggregates all data for the goals experience.
type Inputs struct {
	UserID string     `json:"user_id"`
	Goals  []GoalData `json:"goals"`
}

// GoalData wraps one goal with all engine outputs.
type GoalData struct {
	GoalID            string           `json:"goal_id"`
	Name              string           `json:"name"`
	TargetAmount      float64          `json:"target_amount"`
	CurrentValue      float64          `json:"current_value"`
	FundingGap        float64          `json:"funding_gap"`
	Importance        string           `json:"importance"`
	GoalType          string           `json:"goal_type"`
	Priority          int              `json:"priority"`
	Status            string           `json:"status"`
	TargetDate        string           `json:"target_date,omitempty"`
	CreatedAt         string           `json:"created_at,omitempty"`
	MonthlyContribution float64        `json:"monthly_contribution,omitempty"`

	ProjectedDate       string  `json:"projected_date,omitempty"`
	ProjectedValue      float64 `json:"projected_value,omitempty"`
	OnTrack             bool    `json:"on_track"`
	HasRecommendation   bool    `json:"has_recommendation"`
	HasOptimization     bool    `json:"has_optimization"`
	HasRisk             bool    `json:"has_risk"`
	RecCount            int     `json:"rec_count"`
	OptCount            int     `json:"opt_count"`
	MilestoneCount      int     `json:"milestone_count"`
	EventCount          int     `json:"event_count"`
}

// Composer builds goal experience view models.
type Composer struct{}

func NewComposer() *Composer { return &Composer{} }

// BuildList returns goal summaries.
func (c *Composer) BuildList(inputs Inputs) *GoalSummaryList {
	summaries := make([]GoalSummary, len(inputs.Goals))
	for i, g := range inputs.Goals {
		summaries[i] = c.toSummary(g)
	}
	return &GoalSummaryList{Goals: summaries}
}

// BuildDashboard returns the full dashboard for a single goal.
func (c *Composer) BuildDashboard(goal GoalData) *GoalDashboard {
	progress := c.calcProgress(goal)
	cards := c.buildCards(goal, progress)
	return &GoalDashboard{
		GoalID: goal.GoalID, GoalName: goal.Name,
		Progress: progress, Status: goal.Status,
		Importance: goal.Importance, GoalType: goal.GoalType,
		TargetAmount: goal.TargetAmount, CurrentValue: goal.CurrentValue,
		FundingGap: goal.FundingGap, TargetDate: goal.TargetDate,
		Cards: cards,
	}
}

// BuildProgress returns progress details for a goal.
func (c *Composer) BuildProgress(goal GoalData) *GoalProgress {
	pct := c.calcProgress(goal)
	remaining := goal.TargetAmount - goal.CurrentValue
	if remaining < 0 { remaining = 0 }

	months := 0
	if goal.MonthlyContribution > 0 && remaining > 0 {
		months = int(remaining / goal.MonthlyContribution)
	}

	return &GoalProgress{
		GoalID: goal.GoalID, GoalName: goal.Name,
		ProgressPct: pct, CurrentValue: goal.CurrentValue,
		TargetAmount: goal.TargetAmount, RemainingAmount: remaining,
		Status: goal.Status,
		MonthlyContribution: goal.MonthlyContribution,
		MonthsToTarget: months,
	}
}

// BuildProjection returns projection view for a goal.
func (c *Composer) BuildProjection(goal GoalData) *GoalProjectionView {
	return &GoalProjectionView{
		GoalID: goal.GoalID,
		ProjectedDate: goal.ProjectedDate,
		ProjectedValue: goal.ProjectedValue,
		OnTrack: goal.OnTrack,
		MonthlyContribution: goal.MonthlyContribution,
	}
}

// BuildRecommendations returns rec view for a goal.
func (c *Composer) BuildRecommendations(goal GoalData, recs []RecItem) *GoalRecommendationView {
	return &GoalRecommendationView{
		GoalID: goal.GoalID, Recs: recs, TotalCount: len(recs),
	}
}

// BuildOptimization returns optimization view for a goal.
func (c *Composer) BuildOptimization(goal GoalData, opts []OptItem) *GoalOptimizationView {
	return &GoalOptimizationView{
		GoalID: goal.GoalID, Opts: opts, TotalCount: len(opts),
	}
}

// BuildTimeline returns timeline events for a goal.
func (c *Composer) BuildTimeline(goal GoalData, evts []GoalEvent) *GoalTimelineView {
	return &GoalTimelineView{
		GoalID: goal.GoalID, Events: evts, TotalCount: len(evts),
	}
}

// BuildMilestones returns milestones for a goal.
func (c *Composer) BuildMilestones(goal GoalData, milestones []Milestone) *GoalMilestoneView {
	return &GoalMilestoneView{
		GoalID: goal.GoalID, Milestones: milestones, TotalCount: len(milestones),
	}
}

func (c *Composer) calcProgress(goal GoalData) int {
	if goal.TargetAmount <= 0 { return 0 }
	p := int(goal.CurrentValue / goal.TargetAmount * 100)
	if p > 100 { p = 100 }
	return p
}

func (c *Composer) toSummary(g GoalData) GoalSummary {
	return GoalSummary{
		GoalID: g.GoalID, Name: g.Name, Progress: c.calcProgress(g),
		Status: g.Status, Importance: g.Importance, Priority: g.Priority,
		HasRecommendation: g.HasRecommendation,
	}
}

func (c *Composer) buildCards(goal GoalData, progress int) []Card {
	cards := []Card{
		{CardID: "ov", CardType: CTOverview, Title: goal.Name,
			Summary:  fmt.Sprintf("%d%% funded — %s", progress, goal.Status),
			Data:     map[string]interface{}{"target": goal.TargetAmount, "current": goal.CurrentValue},
			Priority: 1},
		{CardID: "prog", CardType: CTProgress, Title: "Progress",
			Summary:  fmt.Sprintf("₹%.0f / ₹%.0f", goal.CurrentValue, goal.TargetAmount),
			Data:     map[string]interface{}{"pct": progress, "remaining": goal.FundingGap},
			Priority: 2},
		{CardID: "fund", CardType: CTFunding, Title: "Funding Status",
			Summary:  fmt.Sprintf("Gap: ₹%.0f", goal.FundingGap),
			Data:     map[string]interface{}{"gap": goal.FundingGap, "monthly": goal.MonthlyContribution},
			Priority: 3},
		{CardID: "proj", CardType: CTProjection, Title: "Projection",
			Summary:  fmt.Sprintf("On track: %t", goal.OnTrack),
			Data:     map[string]interface{}{"projected_date": goal.ProjectedDate, "projected_value": goal.ProjectedValue},
			Priority: 4},
	}
	if goal.HasRecommendation {
		cards = append(cards, Card{CardID: "rec", CardType: CTRecommendation,
			Title: "Recommendations", Summary: fmt.Sprintf("%d available", goal.RecCount),
			Data: map[string]interface{}{"count": goal.RecCount}, Priority: 5})
	}
	if goal.HasOptimization {
		cards = append(cards, Card{CardID: "opt", CardType: CTOptimization,
			Title: "Optimization", Summary: fmt.Sprintf("%d strategies", goal.OptCount),
			Data: map[string]interface{}{"count": goal.OptCount}, Priority: 6})
	}
	if goal.EventCount > 0 {
		cards = append(cards, Card{CardID: "tl", CardType: CTTimeline,
			Title: "Timeline", Summary: fmt.Sprintf("%d events", goal.EventCount),
			Data: map[string]interface{}{"count": goal.EventCount}, Priority: 7})
	}
	if goal.MilestoneCount > 0 {
		cards = append(cards, Card{CardID: "ms", CardType: CTMilestone,
			Title: "Milestones", Summary: fmt.Sprintf("%d milestones", goal.MilestoneCount),
			Data: map[string]interface{}{"count": goal.MilestoneCount}, Priority: 8})
	}
	return cards
}

type GoalSummaryList struct {
	Goals []GoalSummary `json:"goals"`
}

type GoalSummary struct {
	GoalID           string `json:"goal_id"`
	Name             string `json:"name"`
	Progress         int    `json:"progress"`
	Status           string `json:"status"`
	Importance       string `json:"importance"`
	Priority         int    `json:"priority"`
	HasRecommendation bool  `json:"has_recommendation"`
}
