package github

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CacheEntry holds cached content with expiration time.
type CacheEntry struct {
	Content   []byte
	ExpiresAt time.Time
}

// Cache provides in-memory and file-based caching for GitHub content with TTL support.
type Cache struct {
	mu       sync.RWMutex
	items    map[string]*CacheEntry
	cacheDir string // Directory for persistent cache files
	useDisk  bool   // Whether to use disk cache
}

// NewCache creates a new cache instance.
// If cacheDir is provided, it will use persistent file-based caching.
func NewCache(cacheDir string) *Cache {
	c := &Cache{
		items:    make(map[string]*CacheEntry),
		cacheDir: cacheDir,
		useDisk:  cacheDir != "",
	}

	// Create cache directory if using disk cache
	if c.useDisk {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			// If we can't create the directory, fall back to in-memory only
			c.useDisk = false
		}
	}

	return c
}

// Get retrieves cached content if it exists and hasn't expired.
// Returns the content and true if found and valid, false otherwise.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	entry, ok := c.items[key]
	c.mu.RUnlock()

	// Check in-memory cache first
	if ok {
		// Check if expired
		if time.Now().After(entry.ExpiresAt) {
			// Don't delete here, let Set handle cleanup
			return nil, false
		}
		return entry.Content, true
	}

	// Check disk cache if enabled
	if c.useDisk {
		if content, expiresAt, ok := c.getFromDisk(key); ok {
			// Check if expired
			if time.Now().After(expiresAt) {
				// Delete expired file
				c.deleteFromDisk(key)
				return nil, false
			}
			// Load into memory cache for faster access
			c.mu.Lock()
			c.items[key] = &CacheEntry{
				Content:   content,
				ExpiresAt: expiresAt,
			}
			c.mu.Unlock()
			return content, true
		}
	}

	return nil, false
}

// Set stores content in the cache with the specified TTL.
func (c *Cache) Set(key string, content []byte, ttl time.Duration) {
	expiresAt := time.Now().Add(ttl)

	// Store in memory
	c.mu.Lock()
	// Clean up expired entries periodically
	if len(c.items) > 1000 {
		c.cleanupExpired()
	}
	c.items[key] = &CacheEntry{
		Content:   content,
		ExpiresAt: expiresAt,
	}
	c.mu.Unlock()

	// Store on disk if enabled
	if c.useDisk {
		c.saveToDisk(key, content, expiresAt)
	}
}

// Invalidate removes a specific cache entry.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all cache entries.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheEntry)
}

// cleanupExpired removes expired entries from the cache.
func (c *Cache) cleanupExpired() {
	now := time.Now()
	for key, entry := range c.items {
		if now.After(entry.ExpiresAt) {
			delete(c.items, key)
			// Also delete from disk if enabled
			if c.useDisk {
				c.deleteFromDisk(key)
			}
		}
	}
}

// getFromDisk retrieves content from disk cache.
func (c *Cache) getFromDisk(key string) ([]byte, time.Time, bool) {
	cacheFile := c.getCacheFilePath(key)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, time.Time{}, false
	}

	var entry struct {
		Content   []byte    `json:"content"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		// Invalid cache file, delete it
		os.Remove(cacheFile)
		return nil, time.Time{}, false
	}

	return entry.Content, entry.ExpiresAt, true
}

// saveToDisk saves content to disk cache.
func (c *Cache) saveToDisk(key string, content []byte, expiresAt time.Time) {
	cacheFile := c.getCacheFilePath(key)

	entry := struct {
		Content   []byte    `json:"content"`
		ExpiresAt time.Time `json:"expires_at"`
	}{
		Content:   content,
		ExpiresAt: expiresAt,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return // Silently fail - disk cache is optional
	}

	// Write to temp file first, then rename (atomic write)
	tmpFile := cacheFile + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return
	}
	os.Rename(tmpFile, cacheFile)
}

// deleteFromDisk removes a cache file from disk.
func (c *Cache) deleteFromDisk(key string) {
	cacheFile := c.getCacheFilePath(key)
	os.Remove(cacheFile)
}

// getCacheFilePath returns the file path for a cache key.
func (c *Cache) getCacheFilePath(key string) string {
	// Use SHA256 hash of key as filename to avoid filesystem issues with special characters
	hash := sha256.Sum256([]byte(key))
	filename := hex.EncodeToString(hash[:]) + ".json"
	return filepath.Join(c.cacheDir, filename)
}
