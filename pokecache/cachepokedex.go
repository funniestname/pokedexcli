package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	interval time.Duration
	entries  map[string]cacheEntry
	mu       sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(dur time.Duration) *Cache {
	c := &Cache{
		interval: dur,
		entries:  make(map[string]cacheEntry),
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.Sub(entry.createdAt) > c.interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}
