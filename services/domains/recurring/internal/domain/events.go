package domain

import (
	"time"
)

type EventType string

const (
	EventRecurringCreated              EventType = "recurring.created"
	EventRecurringActivated            EventType = "recurring.activated"
	EventRecurringPaused               EventType = "recurring.paused"
	EventRecurringCompleted            EventType = "recurring.completed"
	EventRecurringCancelled            EventType = "recurring.cancelled"
	EventRecurringArchived             EventType = "recurring.archived"
	EventRecurringOccurrenceGenerated  EventType = "recurring.occurrence.generated"
	EventRecurringSkipped              EventType = "recurring.skipped"
)

type RecurringCreated struct {
	RecurringID string
	UserID      string
	Timestamp   time.Time
}

type RecurringActivated struct {
	RecurringID string
	Timestamp   time.Time
}

type RecurringPaused struct {
	RecurringID string
	Timestamp   time.Time
}

type RecurringCompleted struct {
	RecurringID string
	Timestamp   time.Time
}

type RecurringCancelled struct {
	RecurringID string
	Timestamp   time.Time
}

type RecurringArchived struct {
	RecurringID string
	Timestamp   time.Time
}

type RecurringOccurrenceGenerated struct {
	RecurringID string
	GeneratedID string
	Timestamp   time.Time
}

type RecurringSkipped struct {
	RecurringID string
	Reason      string
	Timestamp   time.Time
}
