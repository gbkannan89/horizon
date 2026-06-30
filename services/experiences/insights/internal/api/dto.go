package api

import "github.com/horizon/core/services/experiences/insights/internal/engine"

type DashboardResponse struct {
	Success  bool                      `json:"success"`
	Data     *engine.InsightsDashboard `json:"data,omitempty"`
	Metadata *Metadata                 `json:"metadata,omitempty"`
}

type FeedResponse struct {
	Success  bool                `json:"success"`
	Data     *engine.InsightsFeed `json:"data,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type InsightListResponse struct {
	Success  bool            `json:"success"`
	Data     *InsightList    `json:"data,omitempty"`
	Metadata *Metadata       `json:"metadata,omitempty"`
}

type InsightList struct {
	Insights []engine.Insight `json:"insights"`
	Count    int              `json:"count"`
}

type InsightResponse struct {
	Success  bool           `json:"success"`
	Data     *engine.Insight `json:"data,omitempty"`
	Metadata *Metadata      `json:"metadata,omitempty"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
