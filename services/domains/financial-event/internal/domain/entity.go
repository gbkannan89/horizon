package domain

import (
	"fmt"
	"time"
)

type FinancialEvent struct {
	eventID           string
	userID            string
	householdID       string
	eventType         EventType
	eventSubType      string
	eventDate         time.Time
	effectiveDate     time.Time
	currency          string
	amount            int64
	source            string
	destination       string
	description       string
	notes             string
	reference         string
	attachments       []string
	origin            EventOrigin
	confidence        EventConfidence
	sourceOfTruth     SourceOfTruth
	createdBy         CreatedBy
	importedFrom      string
	correlationID     string
	reversalOfEventID string
	tags              []string
	state             EventState
	sequenceNumber    int64
	orderIndex        string
	createdAt         time.Time
	updatedAt         time.Time
}

func (e *FinancialEvent) EventID() string              { return e.eventID }
func (e *FinancialEvent) UserID() string               { return e.userID }
func (e *FinancialEvent) HouseholdID() string          { return e.householdID }
func (e *FinancialEvent) EventType() EventType         { return e.eventType }
func (e *FinancialEvent) EventSubType() string         { return e.eventSubType }
func (e *FinancialEvent) EventDate() time.Time         { return e.eventDate }
func (e *FinancialEvent) EffectiveDate() time.Time     { return e.effectiveDate }
func (e *FinancialEvent) Currency() string             { return e.currency }
func (e *FinancialEvent) Amount() int64                { return e.amount }
func (e *FinancialEvent) Source() string               { return e.source }
func (e *FinancialEvent) Destination() string          { return e.destination }
func (e *FinancialEvent) Description() string          { return e.description }
func (e *FinancialEvent) Notes() string                { return e.notes }
func (e *FinancialEvent) Reference() string            { return e.reference }
func (e *FinancialEvent) Attachments() []string        { return e.attachments }
func (e *FinancialEvent) Origin() EventOrigin          { return e.origin }
func (e *FinancialEvent) Confidence() EventConfidence  { return e.confidence }
func (e *FinancialEvent) SourceOfTruth() SourceOfTruth { return e.sourceOfTruth }
func (e *FinancialEvent) CreatedBy() CreatedBy         { return e.createdBy }
func (e *FinancialEvent) ImportedFrom() string         { return e.importedFrom }
func (e *FinancialEvent) CorrelationID() string        { return e.correlationID }
func (e *FinancialEvent) ReversalOfEventID() string    { return e.reversalOfEventID }
func (e *FinancialEvent) Tags() []string               { return e.tags }
func (e *FinancialEvent) State() EventState            { return e.state }
func (e *FinancialEvent) SequenceNumber() int64        { return e.sequenceNumber }
func (e *FinancialEvent) OrderIndex() string           { return e.orderIndex }
func (e *FinancialEvent) CreatedAt() time.Time         { return e.createdAt }
func (e *FinancialEvent) UpdatedAt() time.Time         { return e.updatedAt }

func (e *FinancialEvent) IsImport() bool  { return e.createdBy == CreatedByImport }
func (e *FinancialEvent) IsReversal() bool { return e.reversalOfEventID != "" }

func (e *FinancialEvent) CanTransitionTo(target EventState) error {
	transitions := map[EventState]map[EventState]bool{
		StateDraft:     {StatePending: true, StateCancelled: true},
		StatePending:   {StateConfirmed: true, StateCancelled: true},
		StateConfirmed: {StatePosted: true, StateReversed: true},
		StatePosted:    {StateReversed: true, StateArchived: true},
		StateReversed:  {StateArchived: true},
		StateCancelled: {StateArchived: true},
	}
	if allowed, ok := transitions[e.state]; ok {
		if allowed[target] {
			return nil
		}
	}
	if e.state == target {
		return fmt.Errorf("event is already in %s state", target)
	}
	return fmt.Errorf("cannot transition from %s to %s", e.state, target)
}

func (e *FinancialEvent) SetSequenceNumber(seq int64)  { e.sequenceNumber = seq }
func (e *FinancialEvent) SetOrderIndex(idx string)       { e.orderIndex = idx }
func (e *FinancialEvent) SetCreatedAt(t time.Time)       { e.createdAt = t.UTC() }
func (e *FinancialEvent) SetUpdatedAt(t time.Time)       { e.updatedAt = t.UTC() }


// ReconstructFromDB creates a FinancialEvent from persistent storage without validation.
func ReconstructFromDB(id, userID, householdID string, eventType EventType, eventSubType string,
	eventDate, effectiveDate time.Time, currency string, amount int64, source, destination, description, notes, reference string,
	origin EventOrigin, confidence EventConfidence, sourceOfTruth SourceOfTruth, createdBy CreatedBy,
	importedFrom, correlationID, reversalOfEventID string, tags []string, state EventState,
	sequenceNumber int64, orderIndex string, createdAt, updatedAt time.Time) *FinancialEvent {
	return &FinancialEvent{
		eventID: id, userID: userID, householdID: householdID, eventType: eventType, eventSubType: eventSubType,
		eventDate: eventDate, effectiveDate: effectiveDate, currency: currency, amount: amount,
		source: source, destination: destination, description: description, notes: notes, reference: reference,
		origin: origin, confidence: confidence, sourceOfTruth: sourceOfTruth, createdBy: createdBy,
		importedFrom: importedFrom, correlationID: correlationID, reversalOfEventID: reversalOfEventID,
		tags: tags, state: state, sequenceNumber: sequenceNumber, orderIndex: orderIndex,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}