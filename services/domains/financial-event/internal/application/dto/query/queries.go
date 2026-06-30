package query

import "time"

type GetEventQuery struct {
	EventID string `json:"event_id" validate:"required"`
}

type ListEventsByUserQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type ListEventsByAccountQuery struct {
	AccountID string `json:"account_id" validate:"required"`
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

type ListEventsByDateRangeQuery struct {
	UserID string    `json:"user_id" validate:"required"`
	Start  time.Time `json:"start" validate:"required"`
	End    time.Time `json:"end" validate:"required"`
	Cursor string    `json:"cursor,omitempty"`
	Limit  int       `json:"limit,omitempty"`
}

type ListEventsByTypeQuery struct {
	UserID    string `json:"user_id" validate:"required"`
	EventType string `json:"event_type" validate:"required"`
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

type GetTimelineQuery struct {
	UserID string `json:"user_id" validate:"required"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type EventResult struct {
	EventID         string    `json:"event_id"`
	UserID          string    `json:"user_id"`
	EventType       string    `json:"event_type"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	EventDate       time.Time `json:"event_date"`
	EffectiveDate   time.Time `json:"effective_date"`
	Description     string    `json:"description"`
	State           string    `json:"state"`
	Origin          string    `json:"origin"`
	Confidence      string    `json:"confidence"`
	CreatedBy       string    `json:"created_by"`
	Source          string    `json:"source,omitempty"`
	Destination     string    `json:"destination,omitempty"`
	Reference       string    `json:"reference,omitempty"`
	ImportedFrom    string    `json:"imported_from,omitempty"`
	ReversalOfID    string    `json:"reversal_of_event_id,omitempty"`
	CorrelationID   string    `json:"correlation_id,omitempty"`
	OrderIndex      string    `json:"order_index,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PaginatedResult struct {
	Events     []EventResult `json:"events"`
	NextCursor string        `json:"next_cursor,omitempty"`
	HasMore    bool          `json:"has_more"`
}
