package pokecache

import (
	"sync"
	"time"
)

// cacheEntry holds the raw data and creation timestamp for a cache entry.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache is the main structure for the caching system.
// It uses a map to store data and a Mutex for concurrent safety.
type Cache struct {
	data map[string]cacheEntry
	mu   *sync.Mutex // Mutex to protect the map from race conditions
}

// NewCache creates a new Cache instance and starts the reaping background goroutine.
func NewCache(interval time.Duration) Cache {
	c := Cache{
		data: make(map[string]cacheEntry),
		mu:   &sync.Mutex{},
	}

	// Start the background goroutine to clean out old entries
	go c.reapLoop(interval)

	return c
}

// Add inserts a new key-value pair into the cache with the current timestamp.
func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves a value from the cache. The returned boolean indicates success.
func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

// reapLoop removes entries older than the given interval from the cache.
func (c Cache) reapLoop(interval time.Duration) {
	// Create a ticker that fires every 'interval' duration.
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // Ensure the ticker stops when the function returns

	// Loop indefinitely, executing the cleanup logic when the ticker fires.
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		for key, entry := range c.data {
			// If the entry is older than the interval, delete it.
			if now.Sub(entry.createdAt) > interval {
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
