package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidTransition = errors.New("invalid state transition")
)

func (r *RecurringTransaction) Activate() error {
	if r.status != StatusPaused {
		return ErrInvalidTransition
	}
	r.status = StatusActive
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *RecurringTransaction) Pause() error {
	if r.status != StatusActive {
		return ErrInvalidTransition
	}
	r.status = StatusPaused
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *RecurringTransaction) Cancel() error {
	if r.status == StatusCompleted || r.status == StatusArchived {
		return ErrInvalidTransition
	}
	r.status = StatusCancelled
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *RecurringTransaction) Complete() error {
	if r.status != StatusActive {
		return ErrInvalidTransition
	}
	r.status = StatusCompleted
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *RecurringTransaction) Archive() error {
	if r.status != StatusCompleted && r.status != StatusCancelled {
		return ErrInvalidTransition
	}
	r.status = StatusArchived
	r.updatedAt = time.Now().UTC()
	return nil
}

func (r *RecurringTransaction) AdvanceOccurrence(date time.Time) error {
	if r.nextOccurrence != nil {
		last := *r.nextOccurrence
		r.lastOccurrence = &last
	} else {
		r.lastOccurrence = &date
	}
	
	next, err := r.CalculateNextOccurrence(*r.lastOccurrence)
	if err != nil {
		return err
	}
	
	if r.endDate != nil && next.After(*r.endDate) {
		r.nextOccurrence = nil
		return r.Complete()
	}
	
	r.nextOccurrence = &next
	r.updatedAt = time.Now().UTC()
	return nil
}
