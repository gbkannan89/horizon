package engine

import "fmt"

// GoalCardType represents the type of goal workspace card.
type GoalCardType string
const (
	GCOverview       GoalCardType = "GoalOverview"
	GCFundingProgress GoalCardType = "FundingProgress"
	GCTimeline       GoalCardType = "Timeline"
	GCRecommendation GoalCardType = "Recommendation"
	GCRisk           GoalCardType = "Risk"
	GCMilestone      GoalCardType = "Milestone"
	GCAchievement    GoalCardType = "Achievement"
)

// GoalCard represents a single card in the goal workspace.
type GoalCard struct {
	CardID      string       `json:"card_id"`
	CardType    GoalCardType `json:"card_type"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	ProgressPct int          `json:"progress_pct"`
	Status      string       `json:"status"` // on_track, at_risk, behind, complete
	Priority    int          `json:"priority"`
}

// GoalWorkspace represents the full workspace for a single goal.
type GoalWorkspace struct {
	GoalID    string     `json:"goal_id"`
	GoalName  string     `json:"goal_name"`
	Progress  int        `json:"progress"`
	Status    string     `json:"status"`
	Cards     []GoalCard `json:"cards"`
}

// GoalsList is the overview of all goals.
type GoalsList struct {
	Goals []GoalSummary `json:"goals"`
}

// GoalSummary is a compact view of a goal.
type GoalSummary struct {
	GoalID           string `json:"goal_id"`
	Name             string `json:"name"`
	Progress         int    `json:"progress"`
	Status           string `json:"status"`
	Importance       string `json:"importance"`
	Priority         int    `json:"priority"`
	FundingGap       string `json:"funding_gap,omitempty"`
	HasRecommendation bool  `json:"has_recommendation"`
}

// Inputs represents data consumed by the goals experience.
type Inputs struct {
	UserID      string `json:"user_id"`
	Goals       []GoalInput `json:"goals"`
	HealthScore int    `json:"health_score"`
	RiskScore   int    `json:"risk_score"`
}

// GoalInput represents a single goal's data from the domain.
type GoalInput struct {
	GoalID          string  `json:"goal_id"`
	Name            string  `json:"name"`
	TargetAmount    int64   `json:"target_amount"`
	CurrentValue    int64   `json:"current_value"`
	Importance      string  `json:"importance"`
	Priority        int     `json:"priority"`
	Status          string  `json:"status"`
	FundingGap      int64   `json:"funding_gap"`
	HasRec          bool    `json:"has_recommendation"`
	HasRisk         bool    `json:"has_risk"`
}

// Engine composes the goals experience view models.
type Engine struct{}
func NewEngine() *Engine { return &Engine{} }

func (e *Engine) BuildList(inputs Inputs) *GoalsList {
	summaries := make([]GoalSummary, 0, len(inputs.Goals))
	for _, g := range inputs.Goals {
		progress := 0
		if g.TargetAmount > 0 { progress = int(float64(g.CurrentValue) / float64(g.TargetAmount) * 100) }
		if progress > 100 { progress = 100 }

		gap := ""
		if g.FundingGap > 0 { gap = fmt.Sprintf("₹%d", g.FundingGap) }

		summaries = append(summaries, GoalSummary{
			GoalID: g.GoalID, Name: g.Name, Progress: progress,
			Status: g.Status, Importance: g.Importance, Priority: g.Priority,
			FundingGap: gap, HasRecommendation: g.HasRec,
		})
	}
	return &GoalsList{Goals: summaries}
}

func (e *Engine) BuildWorkspace(goal GoalInput) *GoalWorkspace {
	progress := 0
	if goal.TargetAmount > 0 { progress = int(float64(goal.CurrentValue) / float64(goal.TargetAmount) * 100) }
	if progress > 100 { progress = 100 }

	cards := []GoalCard{
		{CardID: "g-ov", CardType: GCOverview, Title: goal.Name, Summary: fmt.Sprintf("%d%% funded", progress), ProgressPct: progress, Status: goal.Status, Priority: 1},
		{CardID: "g-fund", CardType: GCFundingProgress, Title: "Funding Progress", Summary: fmt.Sprintf("₹%d / ₹%d", goal.CurrentValue, goal.TargetAmount), ProgressPct: progress, Priority: 2},
	}

	if goal.HasRec {
		cards = append(cards, GoalCard{CardID: "g-rec", CardType: GCRecommendation, Title: "Recommendation", Summary: "Improve this goal", Priority: 3})
	}
	if goal.HasRisk {
		cards = append(cards, GoalCard{CardID: "g-risk", CardType: GCRisk, Title: "Risk Alert", Summary: "Goal at risk", Priority: 4})
	}

	return &GoalWorkspace{
		GoalID: goal.GoalID, GoalName: goal.Name,
		Progress: progress, Status: goal.Status, Cards: cards,
	}
}
