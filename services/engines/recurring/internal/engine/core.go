package engine

import (
	"context"
	"time"

	"github.com/horizon/core/packages/events"
)

type Engine struct {
	provider RecurringProvider
	creator  FinancialEventCreator
	eventBus events.Publisher
}

func NewEngine(p RecurringProvider, c FinancialEventCreator, bus events.Publisher) *Engine {
	return &Engine{
		provider: p,
		creator:  c,
		eventBus: bus,
	}
}

func (e *Engine) ProcessDueTransactions(ctx context.Context, asOf time.Time) error {
	dueTransactions, err := e.provider.GetDueByDate(ctx, asOf)
	if err != nil {
		return err
	}

	for _, rt := range dueTransactions {
		err := e.processTransaction(ctx, rt, asOf)
		if err != nil {
			// Log error, but continue with others
			continue
		}
	}
	return nil
}

func (e *Engine) processTransaction(ctx context.Context, rt *RecurringTransaction, asOf time.Time) error {
	occurrenceDate := asOf
	if rt.NextOccurrence != nil {
		occurrenceDate = *rt.NextOccurrence
	}

	// 1. Create Financial Event
	if err := e.creator.CreateEvent(ctx, rt, occurrenceDate); err != nil {
		return err
	}

	// 2. Advance the occurrence schedule
	// (Simulated for engine side - in reality calls API)

	// 3. Save the updated RecurringTransaction
	if err := e.provider.Save(ctx, rt); err != nil {
		return err
	}

	// 4. Publish Event
	data := []byte(`{"recurring_id":"` + rt.ID + `","generated_id":"TODO","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`)
	_ = e.eventBus.Publish(ctx, events.NewEnvelope("recurring.occurrence.generated", 1, data))

	return nil
}
