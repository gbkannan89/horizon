package cache

import (
	"context"
	"sync"
	"time"

	"github.com/horizon/core/services/experiences/timeline/internal/engine"
)

// InMemoryCache is a read-through TTL-based cache for timeline items.
type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]*cacheEntry
	ttl  time.Duration
}

type cacheEntry struct {
	items   []engine.RawItem
	expires time.Time
}

func NewInMemory(ttl time.Duration) *InMemoryCache {
	return &InMemoryCache{data: make(map[string]*cacheEntry), ttl: ttl}
}

func (c *InMemoryCache) Get(_ context.Context, key string) ([]engine.RawItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.data[key]
	if !ok || time.Now().After(entry.expires) { return nil, false }
	return entry.items, true
}

func (c *InMemoryCache) Set(_ context.Context, key string, items []engine.RawItem) {
	c.mu.Lock()
	c.data[key] = &cacheEntry{items: items, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *InMemoryCache) Invalidate(_ context.Context, key string) {
	c.mu.Lock(); delete(c.data, key); c.mu.Unlock()
}
