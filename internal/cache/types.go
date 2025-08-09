package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	value []byte
}

type Cache interface {
	Set(key string, data []byte)
	Get(key string) (data []byte, success bool)
}

type SimpleCache struct {
	store map[string]CacheEntry
	interval time.Duration
	mu sync.RWMutex
}