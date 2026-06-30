package eventbus

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Config holds NATS event bus configuration.
type Config struct {
	URL          string
	StreamName   string
	MaxAge       time.Duration
	MaxMsgs      int64
	Replicas     int
	RetryBase    time.Duration
	RetryMax     time.Duration
}

// DefaultConfig returns a sensible default configuration.
func DefaultConfig() Config {
	return Config{
		URL:        "nats://localhost:4222",
		StreamName: "horizon",
		MaxAge:     7 * 24 * time.Hour,
		MaxMsgs:    1000000,
		Replicas:   1,
		RetryBase:  100 * time.Millisecond,
		RetryMax:   10 * time.Second,
	}
}

// Bus manages NATS JetStream connections, streams, and subscriptions.
type Bus struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	cfg    Config
}

// New creates a new NATS event bus connection and initialises streams.
func New(cfg Config) (*Bus, error) {
	conn, err := nats.Connect(cfg.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(5*time.Second),
		nats.Name("horizon-eventbus"),
	)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("jetstream: %w", err)
	}

	bus := &Bus{conn: conn, js: js, cfg: cfg}
	if err := bus.initStreams(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("init streams: %w", err)
	}

	log.Printf("[eventbus] connected to %s, stream: %s", cfg.URL, cfg.StreamName)
	return bus, nil
}

// SubjectFor returns the NATS subject for a given event type and aggregate.
const SubjectPrefix = "horizon.event"

func SubjectFor(eventType, aggregateType string) string {
	return fmt.Sprintf("%s.%s.%s", SubjectPrefix, aggregateType, eventType)
}

// SubjectForDLQ returns the dead-letter subject for a given aggregate.
func SubjectForDLQ(aggregateType string) string {
	return fmt.Sprintf("%s.%s.dlq", SubjectPrefix, aggregateType)
}

// Publish sends an event to NATS JetStream.
func (b *Bus) Publish(subject string, data []byte) error {
	_, err := b.js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("publish %s: %w", subject, err)
	}
	return nil
}

// Subscribe creates a durable consumer for a subject.
func (b *Bus) Subscribe(subject, consumerName string, handler func(msg []byte) error) (*nats.Subscription, error) {
	sub, err := b.js.Subscribe(subject, func(m *nats.Msg) {
		if err := handler(m.Data); err != nil {
			log.Printf("[eventbus] handler error for %s: %v", subject, err)
			m.Nak()
			return
		}
		m.Ack()
	},
		nats.Durable(consumerName),
		nats.DeliverNew(),
		nats.AckExplicit(),
		nats.MaxDeliver(3),
		nats.BackOff(b.backoffSchedule()),
	)
	if err != nil {
		return nil, fmt.Errorf("subscribe %s: %w", subject, err)
	}
	return sub, nil
}

// SubscribeDLQ subscribes to the dead-letter queue for monitoring.
func (b *Bus) SubscribeDLQ(aggregateType string, handler func(msg []byte) error) (*nats.Subscription, error) {
	return b.Subscribe(SubjectForDLQ(aggregateType), aggregateType+"-dlq", handler)
}

// Close gracefully shuts down the NATS connection.
func (b *Bus) Close() {
	b.conn.Drain()
	b.conn.Close()
}

// Health returns nil if the connection is healthy.
func (b *Bus) Health() error {
	if !b.conn.IsConnected() {
		return fmt.Errorf("nats not connected")
	}
	return nil
}

func (b *Bus) initStreams() error {
	cfg := &nats.StreamConfig{
		Name:       b.cfg.StreamName,
		Subjects:   []string{SubjectPrefix + ".>"},
		MaxAge:     b.cfg.MaxAge,
		MaxMsgs:    b.cfg.MaxMsgs,
		Storage:    nats.FileStorage,
		Replicas:   b.cfg.Replicas,
		Retention:  nats.LimitsPolicy,
		Discard:    nats.DiscardOld,
		MaxMsgSize: 1024 * 1024, // 1 MB
	}

	_, err := b.js.AddStream(cfg)
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return fmt.Errorf("add stream: %w", err)
	}

	// DLQ stream
	dlqCfg := &nats.StreamConfig{
		Name:       b.cfg.StreamName + "-dlq",
		Subjects:   []string{SubjectForDLQ("*")},
		MaxAge:     30 * 24 * time.Hour,
		MaxMsgs:    100000,
		Storage:    nats.FileStorage,
		Replicas:   b.cfg.Replicas,
		Retention:  nats.LimitsPolicy,
		Discard:    nats.DiscardOld,
	}
	_, err = b.js.AddStream(dlqCfg)
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return fmt.Errorf("add dlq stream: %w", err)
	}

	return nil
}

func (b *Bus) backoffSchedule() []time.Duration {
	return []time.Duration{
		b.cfg.RetryBase,
		b.cfg.RetryBase * 2,
		b.cfg.RetryBase * 5,
	}
}

func (b *Bus) nextBackoff(seq uint64) time.Duration {
	d := b.cfg.RetryBase
	for i := uint64(0); i < seq%3; i++ {
		d *= 2
	}
	if d > b.cfg.RetryMax {
		return b.cfg.RetryMax
	}
	return d
}
