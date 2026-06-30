package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	pool *redis.Pool
}

func NewCache(addr string) *Cache {
	pool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &Cache{pool: pool}
}

func (c *Cache) Get(key string) ([]byte, error) {
	conn := c.pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil { return nil, nil }
	return data, err
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", key, int(ttl.Seconds()), value)
	return err
}

func (c *Cache) Delete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

func (c *Cache) Ping() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}

func (c *Cache) Close() error { return c.pool.Close() }
