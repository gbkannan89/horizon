package api

import (
	"github.com/horizon/core/services/experiences/goals/internal/engine"
)

type ListResponse struct {
	Success  bool                `json:"success"`
	Data     *GoalSummaryListDTO `json:"data,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type DashboardResponse struct {
	Success  bool                `json:"success"`
	Data     *engine.GoalDashboard `json:"data,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type ProgressResponse struct {
	Success  bool               `json:"success"`
	Data     *engine.GoalProgress `json:"data,omitempty"`
	Metadata *Metadata          `json:"metadata,omitempty"`
}

type ProjectionResponse struct {
	Success  bool                   `json:"success"`
	Data     *engine.GoalProjectionView `json:"data,omitempty"`
	Metadata *Metadata              `json:"metadata,omitempty"`
}

type RecsResponse struct {
	Success  bool                      `json:"success"`
	Data     *engine.GoalRecommendationView `json:"data,omitempty"`
	Metadata *Metadata                 `json:"metadata,omitempty"`
}

type OptResponse struct {
	Success  bool                   `json:"success"`
	Data     *engine.GoalOptimizationView `json:"data,omitempty"`
	Metadata *Metadata              `json:"metadata,omitempty"`
}

type TimelineResponse struct {
	Success  bool                 `json:"success"`
	Data     *engine.GoalTimelineView `json:"data,omitempty"`
	Metadata *Metadata            `json:"metadata,omitempty"`
}

type MilestoneResponse struct {
	Success  bool                  `json:"success"`
	Data     *engine.GoalMilestoneView `json:"data,omitempty"`
	Metadata *Metadata             `json:"metadata,omitempty"`
}

type GoalSummaryListDTO struct {
	Goals []engine.GoalSummary `json:"goals"`
	Total int                  `json:"total"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
