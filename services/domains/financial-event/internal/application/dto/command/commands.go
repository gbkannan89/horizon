package command

import "time"

type CreateDraftCommand struct {
	UserID        string    `json:"user_id" validate:"required"`
	HouseholdID   string    `json:"household_id,omitempty"`
	EventType     string    `json:"event_type" validate:"required"`
	EventSubType  string    `json:"event_sub_type,omitempty"`
	Amount        int64     `json:"amount" validate:"required,ne=0"`
	Currency      string    `json:"currency" validate:"required,len=3"`
	EventDate     time.Time `json:"event_date" validate:"required"`
	EffectiveDate time.Time `json:"effective_date" validate:"required"`
	Source        string    `json:"source,omitempty"`
	Destination   string    `json:"destination,omitempty"`
	Description   string    `json:"description" validate:"required"`
	Notes         string    `json:"notes,omitempty"`
	Reference     string    `json:"reference,omitempty"`
	Origin        string    `json:"origin" validate:"required"`
	Confidence    string    `json:"confidence" validate:"required"`
	CreatedBy     string    `json:"created_by" validate:"required"`
	ImportedFrom  string    `json:"imported_from,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
}

type SubmitCommand struct {
	EventID string `json:"event_id" validate:"required"`
}

type ConfirmCommand struct {
	EventID string `json:"event_id" validate:"required"`
}

type PostCommand struct {
	EventID string `json:"event_id" validate:"required"`
}

type ReverseCommand struct {
	EventID       string    `json:"event_id" validate:"required"`
	UserID        string    `json:"user_id" validate:"required"`
	Reason        string    `json:"reason,omitempty"`
	EventDate     time.Time `json:"event_date"`
	EffectiveDate time.Time `json:"effective_date"`
	Description   string    `json:"description"`
}

type CancelCommand struct {
	EventID string `json:"event_id" validate:"required"`
	Reason  string `json:"reason,omitempty"`
}

type ArchiveCommand struct {
	EventID string `json:"event_id" validate:"required"`
}

type ImportCommand struct {
	UserID        string    `json:"user_id" validate:"required"`
	EventType     string    `json:"event_type" validate:"required"`
	Amount        int64     `json:"amount" validate:"required,ne=0"`
	Currency      string    `json:"currency" validate:"required,len=3"`
	EventDate     time.Time `json:"event_date" validate:"required"`
	EffectiveDate time.Time `json:"effective_date" validate:"required"`
	Source        string    `json:"source,omitempty"`
	Destination   string    `json:"destination,omitempty"`
	Description   string    `json:"description" validate:"required"`
	Reference     string    `json:"reference,omitempty"`
	ImportedFrom  string    `json:"imported_from" validate:"required"`
	SourceOfTruth string    `json:"source_of_truth" validate:"required"`
}

type DetectDuplicateCommand struct {
	EventID string `json:"event_id" validate:"required"`
}

type CreateDraftResult struct {
	EventID string `json:"event_id"`
	State   string `json:"state"`
}

type CommandResult struct {
	Success bool   `json:"success"`
	EventID string `json:"event_id,omitempty"`
	State   string `json:"state,omitempty"`
}

type DetectDuplicateResult struct {
	IsDuplicate bool    `json:"is_duplicate"`
	MatchedID   string  `json:"matched_id,omitempty"`
	Confidence  float64 `json:"confidence"`
}
