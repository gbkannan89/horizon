package engine

// EventCategory represents the source category of a timeline item.
type EventCategory string

const (
	ECFinancial    EventCategory = "FinancialEvent"
	ECGoal         EventCategory = "GoalEvent"
	ECAccount      EventCategory = "AccountEvent"
	ECAsset        EventCategory = "AssetEvent"
	ECLiability    EventCategory = "LiabilityEvent"
	ECPortfolio    EventCategory = "PortfolioEvent"
	ECHealth       EventCategory = "HealthEvent"
	ECRisk         EventCategory = "RiskEvent"
	ECRecommend    EventCategory = "RecommendationEvent"
	ECSimulation   EventCategory = "SimulationEvent"
	ECOptimization EventCategory = "OptimizationEvent"
	ECAchievement  EventCategory = "Achievement"
	ECMilestone    EventCategory = "Milestone"
	ECUser         EventCategory = "UserEvent"
)

// Severity level for a timeline item.
type Severity string

const (
	SevInfo     Severity = "info"
	SevWarning  Severity = "warning"
	SevCritical Severity = "critical"
	SevSuccess  Severity = "success"
	SevMilestone Severity = "milestone"
)

// TimelineItem is the full model for a single event on the timeline.
type TimelineItem struct {
	TimelineID    string        `json:"timeline_id"`
	Timestamp     string        `json:"timestamp"`
	EventType     string        `json:"event_type"`
	Category      EventCategory `json:"category"`
	Title         string        `json:"title"`
	Summary       string        `json:"summary"`
	Description   string        `json:"description"`
	Severity      Severity      `json:"severity"`
	RelatedEntity string        `json:"related_entity,omitempty"`
	RelatedAgg    string        `json:"related_aggregate,omitempty"`
	RelatedGoal   string        `json:"related_goal,omitempty"`
	RelatedAcct   string        `json:"related_account,omitempty"`
	RelatedAsset  string        `json:"related_asset,omitempty"`
	Amount        int64         `json:"amount,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Icon          string        `json:"icon"`
	Color         string        `json:"color"`
}

// TimelineFeed is the paginated feed response.
type TimelineFeed struct {
	Items   []TimelineItem `json:"items"`
	Cursor  string         `json:"cursor,omitempty"`
	HasMore bool           `json:"has_more"`
	Total   int            `json:"total"`
}

// TimelineView represents one of the timeline views.
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

// FilterState represents the active filters on the timeline.
type FilterState struct {
	Categories []EventCategory `json:"categories,omitempty"`
	EntityID   string          `json:"entity_id,omitempty"`
	StartDate  string          `json:"start_date,omitempty"`
	EndDate    string          `json:"end_date,omitempty"`
	SearchText string          `json:"search_text,omitempty"`
	Severity   Severity        `json:"severity,omitempty"`
	Goal       string          `json:"goal,omitempty"`
	Account    string          `json:"account,omitempty"`
	Asset      string          `json:"asset,omitempty"`
	Liability  string          `json:"liability,omitempty"`
	Portfolio  string          `json:"portfolio,omitempty"`
	Rec        string          `json:"recommendation,omitempty"`
	GroupBy    string          `json:"group_by,omitempty"`
}

// TimelineOutput is the complete timeline response.
type TimelineOutput struct {
	View    TimelineView `json:"view"`
	Feed    TimelineFeed `json:"feed"`
	Filters FilterState  `json:"filters"`
	Total   int          `json:"total"`
}

// Inputs for the timeline experience.
type Inputs struct {
	Items   []RawItem    `json:"items"`
	View    TimelineView `json:"view"`
	Filters FilterState  `json:"filters"`
	UserID  string       `json:"user_id"`
	Cursor  string       `json:"cursor,omitempty"`
	Limit   int          `json:"limit,omitempty"`
}

// RawItem represents raw event data from a provider.
type RawItem struct {
	TimelineID    string        `json:"timeline_id"`
	Timestamp     string        `json:"timestamp"`
	EventType     string        `json:"event_type"`
	Category      EventCategory `json:"category"`
	Title         string        `json:"title"`
	Summary       string        `json:"summary"`
	Description   string        `json:"description"`
	Severity      Severity      `json:"severity"`
	RelatedEntity string        `json:"related_entity,omitempty"`
	RelatedAgg    string        `json:"related_aggregate,omitempty"`
	RelatedGoal   string        `json:"related_goal,omitempty"`
	RelatedAcct   string        `json:"related_account,omitempty"`
	RelatedAsset  string        `json:"related_asset,omitempty"`
	Amount        int64         `json:"amount,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Icon          string        `json:"icon"`
	Color         string        `json:"color"`
}
