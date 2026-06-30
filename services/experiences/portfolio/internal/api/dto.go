package api

import "github.com/horizon/core/services/experiences/portfolio/internal/engine"

type DashboardResponse struct {
	Success  bool                     `json:"success"`
	Data     *engine.PortfolioDashboard `json:"data,omitempty"`
	Metadata *Metadata                `json:"metadata,omitempty"`
}

type AllocationResponse struct {
	Success  bool                 `json:"success"`
	Data     *engine.AllocationView `json:"data,omitempty"`
	Metadata *Metadata            `json:"metadata,omitempty"`
}

type PerformanceResponse struct {
	Success  bool                   `json:"success"`
	Data     *engine.PerformanceView `json:"data,omitempty"`
	Metadata *Metadata              `json:"metadata,omitempty"`
}

type RiskResponse struct {
	Success  bool            `json:"success"`
	Data     *engine.RiskView `json:"data,omitempty"`
	Metadata *Metadata       `json:"metadata,omitempty"`
}

type ProjectionResponse struct {
	Success  bool                 `json:"success"`
	Data     *engine.ProjectionView `json:"data,omitempty"`
	Metadata *Metadata            `json:"metadata,omitempty"`
}

type RecsResponse struct {
	Success  bool             `json:"success"`
	Data     *engine.CardView  `json:"data,omitempty"`
	Metadata *Metadata        `json:"metadata,omitempty"`
}

type OptResponse struct {
	Success  bool             `json:"success"`
	Data     *engine.CardView  `json:"data,omitempty"`
	Metadata *Metadata        `json:"metadata,omitempty"`
}

type SimResponse struct {
	Success  bool             `json:"success"`
	Data     *engine.SimView  `json:"data,omitempty"`
	Metadata *Metadata        `json:"metadata,omitempty"`
}

type TimelineResponse struct {
	Success  bool             `json:"success"`
	Data     *engine.CardView  `json:"data,omitempty"`
	Metadata *Metadata        `json:"metadata,omitempty"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
