package api

import (
	"github.com/horizon/core/services/experiences/timeline/internal/engine"
)

type TimelineResponse struct {
	Success  bool                `json:"success"`
	Data     *engine.TimelineOutput `json:"data,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type ItemResponse struct {
	Success  bool               `json:"success"`
	Data     *engine.TimelineItem `json:"data,omitempty"`
	Metadata *Metadata          `json:"metadata,omitempty"`
}

type ItemsResponse struct {
	Success  bool                `json:"success"`
	Data     *ItemList           `json:"data,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type ItemList struct {
	Items []engine.TimelineItem `json:"items"`
	Count int                   `json:"count"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
