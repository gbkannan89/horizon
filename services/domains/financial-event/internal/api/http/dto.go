package http

import "time"

type CreateEventRequest struct {
	EventType     string    `json:"event_type"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	EventDate     time.Time `json:"event_date"`
	EffectiveDate time.Time `json:"effective_date"`
	Source        string    `json:"source"`
	Destination   string    `json:"destination"`
	Description   string    `json:"description"`
	Reference     string    `json:"reference"`
	Notes         string    `json:"notes"`
	Origin        string    `json:"origin"`
	Confidence    string    `json:"confidence"`
	CreatedBy     string    `json:"created_by"`
}

type EventResponse struct {
	EventID       string    `json:"event_id"`
	EventType     string    `json:"event_type"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	Description   string    `json:"description"`
	State         string    `json:"state"`
	EventDate     time.Time `json:"event_date"`
	EffectiveDate time.Time `json:"effective_date"`
	CreatedAt     time.Time `json:"created_at"`
}
