package cache

import (
	"fmt"
	"time"
)

func NewCache(interval time.Duration) Cache {
	cache := &SimpleCache{
		store: make(map[string]CacheEntry),
		interval: interval,
	}

	go cache.reapLoop()
	return cache
}

func (c* SimpleCache) Set(key string, data []byte)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = CacheEntry{
		createdAt: time.Now(),
		value: data,
	}
}

func (c* SimpleCache) Get(key string) ([]byte, bool)  {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.store[key]
	if !ok {
		return  nil, false
	}

	return value.value, true
}

func (c* SimpleCache) reapLoop() {
	for {
		c.mu.Lock()
		for key, value := range c.store {
			if time.Since(value.createdAt) > c.interval {
				fmt.Printf("%s is elapsed, deleting", key)
				delete(c.store, key)			
			}
		}
		c.mu.Unlock()
		time.Sleep(c.interval)
	}
}