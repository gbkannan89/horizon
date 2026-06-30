package envelope

import (
	"encoding/json"
	"time"
)

// Envelope wraps every event published through the event bus.
type Envelope struct {
	EventID         string          `json:"event_id"`
	EventType       string          `json:"event_type"`
	AggregateID     string          `json:"aggregate_id"`
	AggregateType   string          `json:"aggregate_type"`
	AggregateVersion int            `json:"aggregate_version"`
	EventVersion    int             `json:"event_version"`
	CorrelationID   string          `json:"correlation_id"`
	CausationID     string          `json:"causation_id"`
	Timestamp       time.Time       `json:"timestamp"`
	SourceService   string          `json:"source_service"`
	Payload         json.RawMessage `json:"payload"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// New creates a new event envelope.
func New(eventType, aggregateID, aggregateType string, aggregateVersion int, payload json.RawMessage) Envelope {
	return Envelope{
		EventID:         newID(),
		EventType:       eventType,
		AggregateID:     aggregateID,
		AggregateType:   aggregateType,
		AggregateVersion: aggregateVersion,
		EventVersion:    1,
		CorrelationID:   "",
		CausationID:     "",
		Timestamp:       time.Now().UTC(),
		SourceService:   "",
		Payload:         payload,
		Metadata:        make(map[string]string),
	}
}

// WithCorrelationID sets the correlation ID.
func (e Envelope) WithCorrelationID(cid string) Envelope {
	e.CorrelationID = cid
	return e
}

// WithCausationID sets the causation ID.
func (e Envelope) WithCausationID(cid string) Envelope {
	e.CausationID = cid
	return e
}

// WithSource sets the source service.
func (e Envelope) WithSource(s string) Envelope {
	e.SourceService = s
	return e
}

// Marshal serializes the envelope to JSON.
func (e Envelope) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// Unmarshal deserializes JSON into an envelope.
func Unmarshal(data []byte) (Envelope, error) {
	var e Envelope
	err := json.Unmarshal(data, &e)
	return e, err
}

// ToOutboxPayload converts the envelope to an outbox JSON payload.
func (e Envelope) ToOutboxPayload() ([]byte, error) {
	return e.Marshal()
}

func newID() string {
	b := make([]byte, 16)
	now := time.Now().UnixMilli()
	b[0] = byte(now >> 40)
	b[1] = byte(now >> 32)
	b[2] = byte(now >> 24)
	b[3] = byte(now >> 16)
	b[4] = byte(now >> 8)
	b[5] = byte(now)
	b[6] = (b[6] & 0x0f) | 0x70
	b[8] = (b[8] & 0x3f) | 0x80
	return string(b)
}
