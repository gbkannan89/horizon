package api

import "github.com/horizon/core/services/experiences/planning/internal/engine"

type DashboardResponse struct {
	Success  bool                   `json:"success"`
	Data     *engine.PlanningDashboard `json:"data,omitempty"`
	Metadata *Metadata              `json:"metadata,omitempty"`
}

type ProjectionResponse struct {
	Success  bool                  `json:"success"`
	Data     *engine.ProjectionData `json:"data,omitempty"`
	Metadata *Metadata             `json:"metadata,omitempty"`
}

type ScenarioListResponse struct {
	Success  bool              `json:"success"`
	Data     *ScenarioList     `json:"data,omitempty"`
	Metadata *Metadata         `json:"metadata,omitempty"`
}

type ScenarioList struct {
	Scenarios []engine.Scenario `json:"scenarios"`
	Total     int               `json:"total"`
}

type ComparisonResponse struct {
	Success  bool                  `json:"success"`
	Data     *engine.PlanComparison `json:"data,omitempty"`
	Metadata *Metadata             `json:"metadata,omitempty"`
}

type CardListResponse struct {
	Success  bool        `json:"success"`
	Data     *CardList   `json:"data,omitempty"`
	Metadata *Metadata   `json:"metadata,omitempty"`
}

type CardList struct {
	Cards []engine.Card `json:"cards"`
	Count int           `json:"count"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
