package cache

import (
	"context"
	"sync"
	"time"

	"github.com/horizon/core/services/experiences/dashboard/internal/engine"
)

// InMemoryCache is a read-through, TTL-based cache for dashboard inputs.
type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]*cacheEntry
	ttl  time.Duration
}

type cacheEntry struct {
	inputs  *engine.Inputs
	expires time.Time
}

func NewInMemory(ttl time.Duration) *InMemoryCache {
	return &InMemoryCache{data: make(map[string]*cacheEntry), ttl: ttl}
}

func (c *InMemoryCache) Get(_ context.Context, key string) (*engine.Inputs, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	if !ok || time.Now().After(entry.expires) { return nil, false }
	return entry.inputs, true
}

func (c *InMemoryCache) Set(_ context.Context, key string, inputs *engine.Inputs) {
	if inputs == nil {
		c.mu.Lock(); delete(c.data, key); c.mu.Unlock()
		return
	}
	c.mu.Lock()
	c.data[key] = &cacheEntry{inputs: inputs, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *InMemoryCache) Invalidate(_ context.Context, key string) {
	c.mu.Lock(); delete(c.data, key); c.mu.Unlock()
}

func (c *InMemoryCache) Clear() {
	c.mu.Lock(); c.data = make(map[string]*cacheEntry); c.mu.Unlock()
}
