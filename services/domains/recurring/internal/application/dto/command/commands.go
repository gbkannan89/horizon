package command

import (
	"time"

	"github.com/horizon/core/services/domains/recurring/internal/domain"
)

type CreateRecurringCommand struct {
	ID            string
	UserID        string
	HouseholdID   string
	Name          string
	Description   string
	Amount        int64
	Currency      string
	Frequency     domain.Frequency
	Interval      int
	StartDate     time.Time
	EndDate       *time.Time
	EventTemplate domain.EventTemplate
	SkipHolidays  bool
	SkipWeekends  bool
	Tags          []string
	Metadata      map[string]string
}
