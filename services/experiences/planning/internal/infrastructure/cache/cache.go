package cache

import (
	"context"
	"sync"
	"time"
)

type InMemoryCache struct {
	mu   sync.RWMutex
	data map[string]*entry
	ttl  time.Duration
}

type entry struct {
	val     interface{}
	expires time.Time
}

func NewInMemory(ttl time.Duration) *InMemoryCache {
	return &InMemoryCache{data: make(map[string]*entry), ttl: ttl}
}

func (c *InMemoryCache) Get(_ context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.data[key]
	if !ok || time.Now().After(e.expires) { return nil, false }
	return e.val, true
}

func (c *InMemoryCache) Set(_ context.Context, key string, val interface{}) {
	c.mu.Lock()
	c.data[key] = &entry{val: val, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *InMemoryCache) Invalidate(_ context.Context, key string) {
	c.mu.Lock(); delete(c.data, key); c.mu.Unlock()
}
