package domain

import "time"

type BaseDomainEvent struct {
	Type string    `json:"event_type"`
	ID   string    `json:"event_id"`
	Time time.Time `json:"timestamp"`
}

func NewBaseDomainEvent(eventType, aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		Type: eventType,
		ID:   aggregateID,
		Time: time.Now().UTC(),
	}
}

func (e BaseDomainEvent) EventName() string  { return e.Type }
func (e BaseDomainEvent) EntityID() string   { return e.ID }

type FinancialEventCreated struct {
	BaseDomainEvent
	UserID        string    `json:"user_id"`
	EvtType       string    `json:"event_type"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	EffectiveDate time.Time `json:"effective_date"`
	Description   string    `json:"description"`
}

func NewFinancialEventCreated(e *FinancialEvent) FinancialEventCreated {
	return FinancialEventCreated{
		BaseDomainEvent: NewBaseDomainEvent("FinancialEventCreated", e.EventID()),
		UserID:          e.UserID(),
		EvtType:         string(e.EventType()),
		Amount:          e.Amount(),
		Currency:        e.Currency(),
		EffectiveDate:   e.EffectiveDate(),
		Description:     e.Description(),
	}
}

type FinancialEventConfirmed struct {
	BaseDomainEvent
	UserID        string    `json:"user_id"`
	EvtType       string    `json:"event_type"`
	Amount        int64     `json:"amount"`
	EffectiveDate time.Time `json:"effective_date"`
}

func NewFinancialEventConfirmed(e *FinancialEvent) FinancialEventConfirmed {
	return FinancialEventConfirmed{
		BaseDomainEvent: NewBaseDomainEvent("FinancialEventConfirmed", e.EventID()),
		UserID:          e.UserID(),
		EvtType:         string(e.EventType()),
		Amount:          e.Amount(),
		EffectiveDate:   e.EffectiveDate(),
	}
}

type FinancialEventPosted struct {
	BaseDomainEvent
	UserID        string    `json:"user_id"`
	EvtType       string    `json:"event_type"`
	Amount        int64     `json:"amount"`
	EffectiveDate time.Time `json:"effective_date"`
	Source        string    `json:"source,omitempty"`
	Destination   string    `json:"destination,omitempty"`
}

func NewFinancialEventPosted(e *FinancialEvent) FinancialEventPosted {
	return FinancialEventPosted{
		BaseDomainEvent: NewBaseDomainEvent("FinancialEventPosted", e.EventID()),
		UserID:          e.UserID(),
		EvtType:         string(e.EventType()),
		Amount:          e.Amount(),
		EffectiveDate:   e.EffectiveDate(),
		Source:          e.Source(),
		Destination:     e.Destination(),
	}
}

type FinancialEventReversed struct {
	BaseDomainEvent
	UserID          string    `json:"user_id"`
	ReversalEventID string    `json:"reversal_event_id"`
	Amount          int64     `json:"amount"`
	EffectiveDate   time.Time `json:"effective_date"`
}

func NewFinancialEventReversed(original, reversal *FinancialEvent) FinancialEventReversed {
	return FinancialEventReversed{
		BaseDomainEvent:  NewBaseDomainEvent("FinancialEventReversed", original.EventID()),
		UserID:           original.UserID(),
		ReversalEventID:  reversal.EventID(),
		Amount:           reversal.Amount(),
		EffectiveDate:    reversal.EffectiveDate(),
	}
}

type FinancialEventCancelled struct {
	BaseDomainEvent
	UserID string `json:"user_id"`
	Reason string `json:"reason,omitempty"`
}

func NewFinancialEventCancelled(e *FinancialEvent, reason string) FinancialEventCancelled {
	return FinancialEventCancelled{
		BaseDomainEvent: NewBaseDomainEvent("FinancialEventCancelled", e.EventID()),
		UserID:          e.UserID(),
		Reason:          reason,
	}
}

type FinancialEventArchived struct {
	BaseDomainEvent
	UserID string `json:"user_id"`
}

func NewFinancialEventArchived(e *FinancialEvent) FinancialEventArchived {
	return FinancialEventArchived{
		BaseDomainEvent: NewBaseDomainEvent("FinancialEventArchived", e.EventID()),
		UserID:          e.UserID(),
	}
}

type FinancialEventImported struct {
	BaseDomainEvent
	UserID       string `json:"user_id"`
	ImportedFrom string `json:"imported_from"`
	Reference    string `json:"reference,omitempty"`
}

func NewFinancialEventImported(e *FinancialEvent) FinancialEventImported {
	return FinancialEventImported{
		BaseDomainEvent: NewBaseDomainEvent("FinancialEventImported", e.EventID()),
		UserID:          e.UserID(),
		ImportedFrom:    e.ImportedFrom(),
		Reference:       e.Reference(),
	}
}

type FinancialEventDuplicateDetected struct {
	BaseDomainEvent
	MatchedEventID string  `json:"matched_event_id"`
	Confidence     float64 `json:"confidence"`
}

func NewFinancialEventDuplicateDetected(eventID, matchedID string, confidence float64) FinancialEventDuplicateDetected {
	return FinancialEventDuplicateDetected{
		BaseDomainEvent: NewBaseDomainEvent("FinancialEventDuplicateDetected", eventID),
		MatchedEventID:  matchedID,
		Confidence:      confidence,
	}
}
