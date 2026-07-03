package events

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewEnvelope(t *testing.T) {
	data := json.RawMessage(`{"key":"value"}`)
	e := NewEnvelope("test.event", 1, data)

	if e.EventType != "test.event" {
		t.Errorf("expected event type test.event, got %s", e.EventType)
	}
	if e.SchemaVersion != 1 {
		t.Errorf("expected version 1, got %d", e.SchemaVersion)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if string(e.Data) != `{"key":"value"}` {
		t.Errorf("unexpected data: %s", string(e.Data))
	}
}

func TestWithCorrelationID(t *testing.T) {
	e := NewEnvelope("test", 1, json.RawMessage(`{}`))
	e2 := e.WithCorrelationID("corr-123")

	if e2.CorrelationID != "corr-123" {
		t.Errorf("expected corr-123, got %s", e2.CorrelationID)
	}
	if e.CorrelationID != "" {
		t.Error("original envelope should be unchanged")
	}
}

func TestWithUserID(t *testing.T) {
	e := NewEnvelope("test", 1, json.RawMessage(`{}`))
	e2 := e.WithUserID("user-1")

	if e2.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", e2.UserID)
	}
}

func TestWithSource(t *testing.T) {
	e := NewEnvelope("test", 1, json.RawMessage(`{}`))
	e2 := e.WithSource("test-service")

	if e2.Source != "test-service" {
		t.Errorf("expected test-service, got %s", e2.Source)
	}
}

func TestEnvelopeJSON(t *testing.T) {
	data := json.RawMessage(`{"amount":100}`)
	e := NewEnvelope("financial.event", 2, data).
		WithCorrelationID("corr-abc").
		WithUserID("uid-42").
		WithSource("domain-service")

	raw, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Envelope
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.EventType != "financial.event" {
		t.Errorf("expected financial.event, got %s", decoded.EventType)
	}
	if decoded.SchemaVersion != 2 {
		t.Errorf("expected version 2, got %d", decoded.SchemaVersion)
	}
	if decoded.CorrelationID != "corr-abc" {
		t.Errorf("expected corr-abc, got %s", decoded.CorrelationID)
	}
	if decoded.UserID != "uid-42" {
		t.Errorf("expected uid-42, got %s", decoded.UserID)
	}
	if decoded.Source != "domain-service" {
		t.Errorf("expected domain-service, got %s", decoded.Source)
	}
}

func TestEnvelopeTimestampIsUTC(t *testing.T) {
	data := json.RawMessage(`{}`)
	e := NewEnvelope("test", 1, data)

	if e.Timestamp.Location() != time.UTC {
		t.Error("timestamp should be in UTC")
	}
}
