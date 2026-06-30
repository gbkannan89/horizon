package engine

// EventCategory represents the source category of a timeline event.
type EventCategory string
const (
	ECFinancial   EventCategory = "FinancialEvent"
	ECGoal        EventCategory = "GoalEvent"
	ECAsset       EventCategory = "AssetEvent"
	ECLiability   EventCategory = "LiabilityEvent"
	ECPortfolio   EventCategory = "PortfolioEvent"
	ECHealth      EventCategory = "HealthEvent"
	ECRisk        EventCategory = "RiskEvent"
	ECRecommend   EventCategory = "RecommendationEvent"
	ECSimulation  EventCategory = "SimulationEvent"
	ECOptimization EventCategory = "OptimizationEvent"
	ECAchievement EventCategory = "Achievement"
	ECLifeEvent   EventCategory = "LifeEvent"
	ECMilestone   EventCategory = "Milestone"
)

// TimelineEvent represents a single event in the timeline feed.
type TimelineEvent struct {
	EventID        string        `json:"event_id"`
	SourceType     EventCategory `json:"source_type"`
	EventType      string        `json:"event_type"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Amount         int64         `json:"amount,omitempty"`
	EventDate      string        `json:"event_date"`
	IsCelebration  bool          `json:"is_celebration"`
	IsWarning      bool          `json:"is_warning"`
	IsHistoricalOnly bool        `json:"is_historical_only"`
	SourceEntityID string        `json:"source_entity_id,omitempty"`
	SourceRoute    string        `json:"source_route,omitempty"`
}

// TimelineFeed represents the paginated timeline feed.
type TimelineFeed struct {
	Events    []TimelineEvent `json:"events"`
	Cursor    string          `json:"cursor,omitempty"`
	HasMore   bool            `json:"has_more"`
	Total     int             `json:"total"`
}

// DailySummary represents the aggregated summary for a single day.
type DailySummary struct {
	Date               string `json:"date"`
	EventCount         int    `json:"event_count"`
	NetFinancialImpact int64  `json:"net_financial_impact"`
	MostSignificant    string `json:"most_significant,omitempty"`
}

// TimelineView represents one of the 9 timeline views.
type TimelineView string
const (
	TVToday          TimelineView = "Today"
	TVThisWeek       TimelineView = "ThisWeek"
	TVThisMonth      TimelineView = "ThisMonth"
	TVThisYear       TimelineView = "ThisYear"
	TVCustomRange    TimelineView = "CustomRange"
	TVLifeJourney    TimelineView = "LifeJourney"
	TVGoalJourney    TimelineView = "GoalJourney"
	TVPortfolioJourney TimelineView = "PortfolioJourney"
)

// FilterState represents the active filters applied to the timeline.
type FilterState struct {
	Categories []EventCategory `json:"categories,omitempty"`
	EntityID   string          `json:"entity_id,omitempty"`
	StartDate  string          `json:"start_date,omitempty"`
	EndDate    string          `json:"end_date,omitempty"`
	SearchText string          `json:"search_text,omitempty"`
	GroupBy    string          `json:"group_by,omitempty"` // day, week, month, year, category
}

// TimelineOutput is the complete timeline experience view.
type TimelineOutput struct {
	View          TimelineView   `json:"view"`
	Feed          TimelineFeed   `json:"feed"`
	DailySummaries []DailySummary `json:"daily_summaries,omitempty"`
	Filters       FilterState    `json:"filters"`
	TotalEvents   int            `json:"total_events"`
}

// Inputs for the timeline experience.
type Inputs struct {
	Events   []EventInput `json:"events"`
	View     TimelineView `json:"view"`
	Filters  FilterState  `json:"filters"`
	UserID   string       `json:"user_id"`
	Cursor   string       `json:"cursor,omitempty"`
	Limit    int          `json:"limit,omitempty"`
}

// EventInput represents raw event data from domains and engines.
type EventInput struct {
	EventID        string        `json:"event_id"`
	SourceType     EventCategory `json:"source_type"`
	EventType      string        `json:"event_type"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Amount         int64         `json:"amount,omitempty"`
	EventDate      string        `json:"event_date"`
	IsCelebration  bool          `json:"is_celebration"`
	IsWarning      bool          `json:"is_warning"`
	IsHistorical   bool          `json:"is_historical"`
	SourceEntityID string        `json:"source_entity_id,omitempty"`
	SourceRoute    string        `json:"source_route,omitempty"`
}
