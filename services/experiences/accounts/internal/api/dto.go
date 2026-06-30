package api

import "github.com/horizon/core/services/experiences/accounts/internal/engine"

type DashboardResponse struct {
	Success  bool                   `json:"success"`
	Data     *engine.AccountsDashboard `json:"data,omitempty"`
	Metadata *Metadata              `json:"metadata,omitempty"`
}

type ListResponse struct {
	Success  bool                 `json:"success"`
	Data     *engine.AccountsList `json:"data,omitempty"`
	Metadata *Metadata            `json:"metadata,omitempty"`
}

type DetailResponse struct {
	Success  bool                    `json:"success"`
	Data     *engine.AccountDetailData `json:"data,omitempty"`
	Metadata *Metadata               `json:"metadata,omitempty"`
}

type BalanceResponse struct {
	Success  bool                    `json:"success"`
	Data     *engine.BalanceSummary  `json:"data,omitempty"`
	Metadata *Metadata               `json:"metadata,omitempty"`
}

type CashFlowResponse struct {
	Success  bool                    `json:"success"`
	Data     *engine.CashFlowSummary `json:"data,omitempty"`
	Metadata *Metadata               `json:"metadata,omitempty"`
}

type HealthResponse struct {
	Success  bool                       `json:"success"`
	Data     *engine.AccountHealthSummary `json:"data,omitempty"`
	Metadata *Metadata                  `json:"metadata,omitempty"`
}

type CardListResponse struct {
	Success  bool          `json:"success"`
	Data     *CardList     `json:"data,omitempty"`
	Metadata *Metadata     `json:"metadata,omitempty"`
}

type CardList struct {
	Cards []engine.Card `json:"cards"`
	Count int           `json:"count"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
