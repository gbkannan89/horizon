package query

import (
	"time"
)

type RecurringTransactionDTO struct {
	ID             string            `json:"id"`
	UserID         string            `json:"user_id"`
	HouseholdID    string            `json:"household_id,omitempty"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	Amount         int64             `json:"amount"`
	Currency       string            `json:"currency"`
	Frequency      string            `json:"frequency"`
	Interval       int               `json:"interval"`
	StartDate      time.Time         `json:"start_date"`
	EndDate        *time.Time        `json:"end_date,omitempty"`
	NextOccurrence *time.Time        `json:"next_occurrence,omitempty"`
	LastOccurrence *time.Time        `json:"last_occurrence,omitempty"`
	Status         string            `json:"status"`
	SkipHolidays   bool              `json:"skip_holidays"`
	SkipWeekends   bool              `json:"skip_weekends"`
	Tags           []string          `json:"tags,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}
