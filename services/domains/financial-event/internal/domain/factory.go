package domain

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type FinancialEventFactory struct {
	sequenceCounter int64
}

func NewFinancialEventFactory() *FinancialEventFactory {
	return &FinancialEventFactory{}
}

func NewBaseEvent(eventID, userID string, eventType EventType, amount int64, currency string,
	eventDate, effectiveDate time.Time, description string, origin EventOrigin,
	confidence EventConfidence, createdBy CreatedBy, state EventState,
) *FinancialEvent {
	now := time.Now().UTC()
	return &FinancialEvent{
		eventID:       eventID,
		userID:        userID,
		eventType:     eventType,
		amount:        amount,
		currency:      currency,
		eventDate:     eventDate.UTC(),
		effectiveDate: effectiveDate.UTC(),
		description:   description,
		origin:        origin,
		confidence:    confidence,
		createdBy:     createdBy,
		state:         state,
		createdAt:     now,
		updatedAt:     now,
	}
}

func (f *FinancialEventFactory) nextSequence() int64 {
	return atomic.AddInt64(&f.sequenceCounter, 1)
}

func (f *FinancialEventFactory) CreateDraft(
	eventID, userID string,
	eventType EventType,
	amount int64,
	currency string,
	eventDate, effectiveDate time.Time,
	description string,
	origin EventOrigin,
	confidence EventConfidence,
	createdBy CreatedBy,
	source, destination, reference, notes, householdID, eventSubType string,
	attachments []string,
) (*FinancialEvent, error) {

	if err := ValidateEventType(string(eventType)); err != nil {
		return nil, err
	}
	if err := ValidateAmount(amount); err != nil {
		return nil, err
	}
	if err := ValidateCurrency(currency); err != nil {
		return nil, err
	}
	if err := ValidateEventDate(eventDate); err != nil {
		return nil, err
	}
	if err := ValidateEffectiveDate(effectiveDate, eventDate, 90); err != nil {
		return nil, err
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	if !AllOrigins[origin] {
		return nil, fmt.Errorf("invalid origin: %s", origin)
	}
	if !AllCreatedBy[createdBy] {
		return nil, fmt.Errorf("invalid createdBy: %s", createdBy)
	}

	now := time.Now().UTC()
	return &FinancialEvent{
		eventID:       eventID,
		userID:        userID,
		householdID:   householdID,
		eventType:     eventType,
		eventSubType:  eventSubType,
		eventDate:     eventDate.UTC(),
		effectiveDate: effectiveDate.UTC(),
		currency:      currency,
		amount:        amount,
		source:        source,
		destination:   destination,
		description:   description,
		reference:     reference,
		attachments:   attachments,
		origin:        origin,
		confidence:    confidence,
		createdBy:     createdBy,
		notes:         notes,
		state:         StateDraft,
		createdAt:     now,
		updatedAt:     now,
	}, nil
}
