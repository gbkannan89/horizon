package api

import "github.com/horizon/core/services/experiences/advisor/internal/engine"

type WorkspaceResponse struct {
	Success  bool                   `json:"success"`
	Data     *engine.AdvisorWorkspace `json:"data,omitempty"`
	Metadata *Metadata              `json:"metadata,omitempty"`
}

type ContextResponse struct {
	Success  bool                `json:"success"`
	Data     *engine.AdvisorContext `json:"data,omitempty"`
	Metadata *Metadata           `json:"metadata,omitempty"`
}

type SummaryResponse struct {
	Success  bool                      `json:"success"`
	Data     *engine.FinancialSummaryCtx `json:"data,omitempty"`
	Metadata *Metadata                 `json:"metadata,omitempty"`
}

type CardListResponse struct {
	Success  bool      `json:"success"`
	Data     *CardList `json:"data,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type CardList struct {
	Cards []engine.Card `json:"cards"`
	Count int           `json:"count"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
