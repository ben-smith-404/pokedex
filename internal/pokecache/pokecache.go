package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu           sync.Mutex
	cacheEntries map[string]CacheEntry
}

func (cache Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	cache.cacheEntries[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	cache.mu.Unlock()
}

func (cache Cache) Get(key string) (val []byte, exists bool) {
	cache.mu.Lock()
	entry, ok := cache.cacheEntries[key]
	cache.mu.Unlock()
	if ok {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (cache Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for tick := range ticker.C {
		cache.mu.Lock()
		for key, entry := range cache.cacheEntries {
			if entry.createdAt.Sub(tick)*-1 > interval {
				delete(cache.cacheEntries, key)
			}
		}
		cache.mu.Unlock()
	}
}

func NewCache(interval time.Duration) Cache {
	var cache = Cache{
		mu:           sync.Mutex{},
		cacheEntries: map[string]CacheEntry{},
	}
	go cache.reapLoop(interval)
	return cache
}
