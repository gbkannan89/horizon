package api

import (
	"github.com/horizon/core/services/experiences/dashboard/internal/engine"
)

type DashboardResponse struct {
	Success  bool                  `json:"success"`
	Data     *engine.DashboardOutput `json:"data,omitempty"`
	Metadata *Metadata             `json:"metadata,omitempty"`
}

type SummaryResponse struct {
	Success  bool                  `json:"success"`
	Data     *engine.SummaryOutput `json:"data,omitempty"`
	Metadata *Metadata             `json:"metadata,omitempty"`
}

type WidgetsResponse struct {
	Success  bool            `json:"success"`
	Data     *WidgetList     `json:"data,omitempty"`
	Metadata *Metadata       `json:"metadata,omitempty"`
}

type WidgetList struct {
	Widgets []engine.Widget `json:"widgets"`
	Count   int             `json:"count"`
}

type HealthResponse struct {
	Success bool          `json:"success"`
	Data    *engine.Widget `json:"data,omitempty"`
	Metadata *Metadata    `json:"metadata,omitempty"`
}

type MilestonesResponse struct {
	Success  bool            `json:"success"`
	Data     *MilestoneList  `json:"data,omitempty"`
	Metadata *Metadata       `json:"metadata,omitempty"`
}

type MilestoneList struct {
	Milestones []engine.Widget `json:"milestones"`
	Count      int             `json:"count"`
}

type RecommendationsResponse struct {
	Success  bool            `json:"success"`
	Data     *RecList        `json:"data,omitempty"`
	Metadata *Metadata       `json:"metadata,omitempty"`
}

type RecList struct {
	Recommendations []engine.Widget `json:"recommendations"`
	Count           int             `json:"count"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

type ErrorPayload struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
