package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	val map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.val == nil {
		c.val = make(map[string]cacheEntry)
	}
	c.val[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	val, ok := c.val[key]
	if !ok {
		return nil, false
	}
	return val.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for t := range ticker.C {
		c.mu.Lock()
		for k, v := range c.val {
			if t.Sub(v.createdAt) > interval {
				delete(c.val, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		val: make(map[string]cacheEntry),
		mu: sync.Mutex{}, 
	}
	go c.reapLoop(interval)
	return c
}