package outbox

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/horizon/core/services/infra/internal/eventbus"
	"github.com/horizon/core/services/infra/internal/persistence"
)

// Config holds the outbox dispatcher configuration.
type Config struct {
	PollInterval  time.Duration
	BatchSize     int
	MaxRetries    int
	ShutdownGrace time.Duration
}

// DefaultConfig returns sensible defaults for the outbox dispatcher.
func DefaultConfig() Config {
	return Config{
		PollInterval:  100 * time.Millisecond,
		BatchSize:     100,
		MaxRetries:    5,
		ShutdownGrace: 5 * time.Second,
	}
}

// Dispatcher polls the outbox table and publishes events through NATS.
type Dispatcher struct {
	repo  *persistence.OutboxRepository
	bus   *eventbus.Bus
	cfg   Config
	stop  chan struct{}
	done  chan struct{}
}

// NewDispatcher creates a new outbox dispatcher.
func NewDispatcher(repo *persistence.OutboxRepository, bus *eventbus.Bus, cfg Config) *Dispatcher {
	return &Dispatcher{
		repo: repo,
		bus:  bus,
		cfg:  cfg,
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
}

// Start begins polling the outbox table in a background goroutine.
func (d *Dispatcher) Start(ctx context.Context) {
	go d.loop(ctx)
	log.Printf("[outbox] dispatcher started (poll: %v, batch: %d)", d.cfg.PollInterval, d.cfg.BatchSize)
}

// Stop signals the dispatcher to shut down gracefully.
func (d *Dispatcher) Stop() {
	close(d.stop)
	select {
	case <-d.done:
	case <-time.After(d.cfg.ShutdownGrace):
		log.Print("[outbox] dispatcher shutdown timed out")
	}
	log.Print("[outbox] dispatcher stopped")
}

func (d *Dispatcher) loop(ctx context.Context) {
	defer close(d.done)

	ticker := time.NewTicker(d.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-d.stop:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.processBatch(ctx)
		}
	}
}

func (d *Dispatcher) processBatch(ctx context.Context) {
	rows, err := d.repo.FetchPending(ctx, d.cfg.BatchSize)
	if err != nil {
		log.Printf("[outbox] fetch error: %v", err)
		return
	}

	if len(rows) == 0 {
		return
	}

	for _, row := range rows {
		select {
		case <-d.stop:
			return
		default:
		}
		d.processRow(ctx, row)
	}
}

func (d *Dispatcher) processRow(ctx context.Context, row persistence.OutboxRow) {
	// Parse the payload to determine the subject
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(row.EventPayload, &payloadMap); err != nil {
		log.Printf("[outbox] unmarshal error for %s: %v", row.ID, err)
		d.repo.MarkCompleted(ctx, row.ID) // can't recover — remove from outbox
		return
	}

	aggregateType, _ := payloadMap["aggregate_type"].(string)
	eventType := row.EventType
	if aggregateType == "" {
		aggregateType = "unknown"
	}

	subject := eventbus.SubjectFor(eventType, aggregateType)

	if err := d.bus.Publish(subject, row.EventPayload); err != nil {
		log.Printf("[outbox] publish error for %s: %v", row.ID, err)
		d.repo.MarkFailed(ctx, row.ID)
		return
	}

	d.repo.MarkCompleted(ctx, row.ID)
}

// OutboxWriter writes events to the outbox table within the same transaction.
type OutboxWriter struct {
	repo *persistence.OutboxRepository
}

// NewOutboxWriter creates a new outbox writer.
func NewOutboxWriter(repo *persistence.OutboxRepository) *OutboxWriter {
	return &OutboxWriter{repo: repo}
}

// Write inserts an event into the outbox.
func (w *OutboxWriter) Write(ctx context.Context, eventID, eventType, idempotencyKey, correlationID, partitionKey string, payload []byte) error {
	return w.repo.Insert(ctx, &persistence.OutboxRow{
		ID:             eventID,
		EventType:      eventType,
		EventPayload:   payload,
		IdempotencyKey: idempotencyKey,
		CorrelationID:  correlationID,
		PartitionKey:   partitionKey,
	})
}
