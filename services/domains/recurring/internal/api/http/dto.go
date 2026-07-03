package http

import (
	"time"

	"github.com/horizon/core/services/domains/recurring/internal/domain"
)

type CreateRecurringRequest struct {
	Name         string               `json:"name"`
	Description  string               `json:"description,omitempty"`
	Amount       int64                `json:"amount"`
	Currency     string               `json:"currency"`
	Frequency    string               `json:"frequency"`
	Interval     int                  `json:"interval"`
	StartDate    time.Time            `json:"start_date"`
	EndDate      *time.Time           `json:"end_date,omitempty"`
	EventTemplate domain.EventTemplate `json:"event_template"`
	SkipHolidays bool                 `json:"skip_holidays"`
	SkipWeekends bool                 `json:"skip_weekends"`
	Tags         []string             `json:"tags,omitempty"`
	Metadata     map[string]string    `json:"metadata,omitempty"`
}

type UpdateRecurringRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	Amount       int64             `json:"amount"`
	Tags         []string          `json:"tags,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type RecurringResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Error    string      `json:"error,omitempty"`
}
