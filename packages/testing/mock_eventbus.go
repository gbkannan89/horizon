package testing

import (
	"context"
	"sync"

	"github.com/horizon/core/packages/events"
)

type MockEventBus struct {
	mu       sync.Mutex
	published []events.Envelope
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{}
}

func (b *MockEventBus) Publish(ctx context.Context, envelope events.Envelope) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.published = append(b.published, envelope)
	return nil
}

func (b *MockEventBus) Published() []events.Envelope {
	b.mu.Lock()
	defer b.mu.Unlock()
	result := make([]events.Envelope, len(b.published))
	copy(result, b.published)
	return result
}

func (b *MockEventBus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.published = nil
}
