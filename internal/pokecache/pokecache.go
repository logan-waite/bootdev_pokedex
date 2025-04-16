package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	mu           sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{cacheEntries: make(map[string]cacheEntry)}
	go cache.reapLoop(interval)
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.cacheEntries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.cacheEntries[key]
	if ok {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			cache.mu.Lock()
			for key, entry := range cache.cacheEntries {
				if time.Since(entry.createdAt) > interval {
					delete(cache.cacheEntries, key)
				}
			}
			cache.mu.Unlock()
		}
	}
}
