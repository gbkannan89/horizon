package domain

import (
	"time"
)

func NewRecurringTransaction(
	id string, userID string, householdID string, name string, desc string,
	amount int64, currency string, freq Frequency, interval int,
	startDate time.Time, endDate *time.Time,
	eventTemplate EventTemplate, skipHols bool, skipWknds bool,
	tags []string, metadata map[string]string,
) (*RecurringTransaction, error) {
	rt := &RecurringTransaction{
		id:             id,
		userID:         userID,
		householdID:    householdID,
		name:           name,
		description:    desc,
		amount:         amount,
		currency:       currency,
		frequency:      freq,
		interval:       interval,
		startDate:      startDate,
		endDate:        endDate,
		status:         StatusActive,
		eventTemplate:  eventTemplate,
		skipHolidays:   skipHols,
		skipWeekends:   skipWknds,
		tags:           tags,
		metadata:       metadata,
		createdAt:      time.Now().UTC(),
		updatedAt:      time.Now().UTC(),
	}

	if err := rt.Validate(); err != nil {
		return nil, err
	}

	nextOcc, err := rt.CalculateNextOccurrence(rt.startDate)
	if err != nil {
		return nil, err
	}
	rt.nextOccurrence = &nextOcc

	return rt, nil
}
