package tmdb

import (
	"sync"
	"time"

	"theater-showtimes/internal/models"
)

// cacheEntry holds a cached movie with expiration
type cacheEntry struct {
	movie     *models.Movie
	expiresAt time.Time
}

// Cache provides in-memory caching for TMDB data
type Cache struct {
	entries map[string]*cacheEntry
	ttl     time.Duration
	mu      sync.RWMutex
}

// NewCache creates a new cache with the given TTL
func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]*cacheEntry),
		ttl:     ttl,
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a movie from cache
func (c *Cache) Get(key string) *models.Movie {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil
	}

	// Check if expired
	if time.Now().After(entry.expiresAt) {
		return nil
	}

	return entry.movie
}

// Set stores a movie in cache
func (c *Cache) Set(key string, movie *models.Movie) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &cacheEntry{
		movie:     movie,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// cleanup periodically removes expired entries
func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.expiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// Clear removes all entries from cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*cacheEntry)
}
