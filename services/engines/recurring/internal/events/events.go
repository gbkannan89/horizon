package events

import (
	"context"
	"log"

	"github.com/horizon/core/packages/events"
)

type RecurringEventListener struct {
	// engine or other dependencies could go here
}

func NewRecurringEventListener() *RecurringEventListener {
	return &RecurringEventListener{}
}

func (l *RecurringEventListener) Handle(ctx context.Context, env events.Envelope) error {
	log.Printf("Received recurring event")
	return nil
}
