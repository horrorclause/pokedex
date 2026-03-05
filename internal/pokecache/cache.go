package pokecache

import (
	"sync"
	"time"
)

// Struct for Mutex, and cacheEntries, the key for the map will be URLs
type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

// Struct for new Cache Entries, val is the response bytes
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {

	// Creates a Cache and returns a pointer to it
	c := &Cache{
		entries: make(map[string]cacheEntry), //Initializes map
	}

	// Concurrently strats the reapLoop go routine
	go c.reapLoop(interval)

	return c
}

// Add method, where a new cache entry is created.
// MU locks the Mutex so only one Go Routine can access
// the map
func (c *Cache) Add(key string, val []byte) {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

}

// Get Method for entry values. Checks if the key exists in the entries map
// if not return false, if so, return the entry value and true
func (c *Cache) Get(key string) ([]byte, bool) {

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	// Creates an infinite ticking loop
	for range ticker.C {
		c.mu.Lock()

		// Iterates over entries in Cache
		for key, entry := range c.entries {

			// Measures the time since the entry was created
			// If it is older than the interval set, delete it.
			if time.Since(entry.createdAt) > interval {
				delete(c.entries, key)
			}
		}

		c.mu.Unlock()

	}

}
