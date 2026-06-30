package events

import (
	"encoding/json"
	"time"
)

type Envelope struct {
	EventID       string          `json:"event_id"`
	EventType     string          `json:"event_type"`
	SchemaVersion int             `json:"version"`
	Timestamp     time.Time       `json:"timestamp"`
	CorrelationID string          `json:"correlation_id"`
	TraceContext  json.RawMessage `json:"trace_context,omitempty"`
	UserID        string          `json:"user_id,omitempty"`
	Source        string          `json:"source"`
	Data          json.RawMessage `json:"data"`
}

func NewEnvelope(eventType string, schemaVersion int, data json.RawMessage) Envelope {
	return Envelope{
		EventType:     eventType,
		SchemaVersion: schemaVersion,
		Timestamp:     time.Now().UTC(),
		Data:          data,
	}
}

func (e Envelope) WithCorrelationID(id string) Envelope {
	e.CorrelationID = id
	return e
}

func (e Envelope) WithUserID(id string) Envelope {
	e.UserID = id
	return e
}

func (e Envelope) WithSource(source string) Envelope {
	e.Source = source
	return e
}
