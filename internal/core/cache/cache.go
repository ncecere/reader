package cache

import (
	"sync"
	"time"
)

// Entry represents a cached item with its metadata
type Entry struct {
	Content   string
	Timestamp time.Time
}

// Cache provides thread-safe caching functionality
type Cache struct {
	store    map[string]*Entry
	mu       sync.RWMutex
	maxAge   time.Duration
	maxItems int
}

// Options configures the cache behavior
type Options struct {
	MaxAge   time.Duration
	MaxItems int
}

// DefaultOptions returns the default cache configuration
func DefaultOptions() *Options {
	return &Options{
		MaxAge:   24 * time.Hour,
		MaxItems: 1000,
	}
}

// New creates a new cache instance
func New(opts *Options) *Cache {
	if opts == nil {
		opts = DefaultOptions()
	}

	return &Cache{
		store:    make(map[string]*Entry),
		maxAge:   opts.MaxAge,
		maxItems: opts.MaxItems,
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.store[key]
	if !exists {
		return "", false
	}

	// Check if entry has expired
	if time.Since(entry.Timestamp) > c.maxAge {
		go c.Delete(key) // Clean up expired entry asynchronously
		return "", false
	}

	return entry.Content, true
}

// Set stores a value in the cache
func (c *Cache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict items
	if len(c.store) >= c.maxItems {
		c.evictOldest()
	}

	c.store[key] = &Entry{
		Content:   value,
		Timestamp: time.Now(),
	}
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]*Entry)
}

// evictOldest removes the oldest entries when cache is full
func (c *Cache) evictOldest() {
	var oldest time.Time
	var oldestKey string
	first := true

	for key, entry := range c.store {
		if first || entry.Timestamp.Before(oldest) {
			oldest = entry.Timestamp
			oldestKey = key
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.store, oldestKey)
	}
}

// Count returns the number of items in the cache
func (c *Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}

// Stats returns cache statistics
type Stats struct {
	ItemCount int
	HitCount  int64
	MissCount int64
}

// Stats returns current cache statistics
func (c *Cache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Stats{
		ItemCount: len(c.store),
	}
}
