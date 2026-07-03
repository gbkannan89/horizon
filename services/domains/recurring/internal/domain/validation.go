package domain

import (
	"errors"
)

var (
	ErrInvalidAmount   = errors.New("amount must be greater than zero")
	ErrInvalidInterval = errors.New("interval must be greater than zero")
	ErrInvalidDates    = errors.New("start date cannot be after end date")
	ErrInvalidName     = errors.New("name cannot be empty")
)

func (r *RecurringTransaction) Validate() error {
	if r.name == "" {
		return ErrInvalidName
	}
	if r.amount <= 0 {
		return ErrInvalidAmount
	}
	if r.interval <= 0 {
		return ErrInvalidInterval
	}
	if r.endDate != nil && r.startDate.After(*r.endDate) {
		return ErrInvalidDates
	}
	return nil
}
