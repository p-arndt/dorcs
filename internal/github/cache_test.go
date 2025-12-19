package github

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	t.Run("in-memory only", func(t *testing.T) {
		cache := NewCache("")
		if cache == nil {
			t.Fatal("NewCache() returned nil")
		}
		if cache.useDisk {
			t.Error("expected useDisk to be false for empty cacheDir")
		}
		if cache.cacheDir != "" {
			t.Errorf("expected cacheDir to be empty, got %q", cache.cacheDir)
		}
	})

	t.Run("with disk cache", func(t *testing.T) {
		tmpDir := t.TempDir()
		cache := NewCache(tmpDir)
		if cache == nil {
			t.Fatal("NewCache() returned nil")
		}
		if !cache.useDisk {
			t.Error("expected useDisk to be true when cacheDir is provided")
		}
		if cache.cacheDir != tmpDir {
			t.Errorf("expected cacheDir to be %q, got %q", tmpDir, cache.cacheDir)
		}
		// Verify directory was created
		if _, err := os.Stat(tmpDir); err != nil {
			t.Errorf("cache directory was not created: %v", err)
		}
	})

	t.Run("fallback to in-memory on invalid dir", func(t *testing.T) {
		// Use a path that cannot be created (on Windows, this might be a device name)
		invalidPath := "CON" // Reserved name on Windows
		cache := NewCache(invalidPath)
		if cache == nil {
			t.Fatal("NewCache() returned nil")
		}
		// Should fall back to in-memory only
		if cache.useDisk {
			t.Error("expected useDisk to be false when directory creation fails")
		}
	})
}

func TestCacheGetSet(t *testing.T) {
	t.Run("in-memory cache", func(t *testing.T) {
		cache := NewCache("")
		key := "test-key"
		content := []byte("test content")

		// Should not exist initially
		if _, ok := cache.Get(key); ok {
			t.Error("expected cache miss for new key")
		}

		// Set content
		cache.Set(key, content, 1*time.Hour)

		// Should retrieve it
		retrieved, ok := cache.Get(key)
		if !ok {
			t.Error("expected cache hit after Set")
		}
		if string(retrieved) != string(content) {
			t.Errorf("expected content %q, got %q", string(content), string(retrieved))
		}
	})

	t.Run("disk cache", func(t *testing.T) {
		tmpDir := t.TempDir()
		cache := NewCache(tmpDir)
		key := "test-key"
		content := []byte("test content for disk")

		// Set content
		cache.Set(key, content, 1*time.Hour)

		// Create new cache instance to test disk persistence
		cache2 := NewCache(tmpDir)
		retrieved, ok := cache2.Get(key)
		if !ok {
			t.Error("expected cache hit from disk after restart")
		}
		if string(retrieved) != string(content) {
			t.Errorf("expected content %q, got %q", string(content), string(retrieved))
		}
	})

	t.Run("expired entry", func(t *testing.T) {
		cache := NewCache("")
		key := "expired-key"
		content := []byte("expired content")

		// Set with very short TTL
		cache.Set(key, content, 1*time.Nanosecond)

		// Wait for expiration
		time.Sleep(10 * time.Millisecond)

		// Should not retrieve expired content
		if _, ok := cache.Get(key); ok {
			t.Error("expected cache miss for expired entry")
		}
	})

	t.Run("multiple keys", func(t *testing.T) {
		cache := NewCache("")
		keys := []string{"key1", "key2", "key3"}
		contents := [][]byte{[]byte("content1"), []byte("content2"), []byte("content3")}

		// Set all
		for i, key := range keys {
			cache.Set(key, contents[i], 1*time.Hour)
		}

		// Retrieve all
		for i, key := range keys {
			retrieved, ok := cache.Get(key)
			if !ok {
				t.Errorf("expected cache hit for key %q", key)
			}
			if string(retrieved) != string(contents[i]) {
				t.Errorf("expected content %q for key %q, got %q", string(contents[i]), key, string(retrieved))
			}
		}
	})
}

func TestCacheInvalidate(t *testing.T) {
	cache := NewCache("")
	key := "invalidate-key"
	content := []byte("content to invalidate")

	cache.Set(key, content, 1*time.Hour)

	// Should exist
	if _, ok := cache.Get(key); !ok {
		t.Error("expected cache hit before invalidation")
	}

	// Invalidate
	cache.Invalidate(key)

	// Should not exist
	if _, ok := cache.Get(key); ok {
		t.Error("expected cache miss after invalidation")
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewCache("")
	keys := []string{"key1", "key2", "key3"}

	// Set multiple entries
	for _, key := range keys {
		cache.Set(key, []byte("content"), 1*time.Hour)
	}

	// Verify they exist
	for _, key := range keys {
		if _, ok := cache.Get(key); !ok {
			t.Errorf("expected cache hit for key %q before clear", key)
		}
	}

	// Clear
	cache.Clear()

	// Verify all are gone
	for _, key := range keys {
		if _, ok := cache.Get(key); ok {
			t.Errorf("expected cache miss for key %q after clear", key)
		}
	}
}

func TestCacheDiskPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	key := "persistent-key"
	content := []byte("persistent content")

	// Set in first cache instance
	cache1 := NewCache(tmpDir)
	cache1.Set(key, content, 1*time.Hour)

	// Verify file exists on disk
	cacheFile := cache1.getCacheFilePath(key)
	if _, err := os.Stat(cacheFile); err != nil {
		t.Errorf("expected cache file to exist at %q: %v", cacheFile, err)
	}

	// Create new cache instance (simulating restart)
	cache2 := NewCache(tmpDir)
	retrieved, ok := cache2.Get(key)
	if !ok {
		t.Error("expected cache hit from disk after restart")
	}
	if string(retrieved) != string(content) {
		t.Errorf("expected content %q, got %q", string(content), string(retrieved))
	}
}

func TestCacheExpiredDiskEntry(t *testing.T) {
	tmpDir := t.TempDir()
	key := "expired-disk-key"
	content := []byte("expired disk content")

	// Set with short TTL
	cache1 := NewCache(tmpDir)
	cache1.Set(key, content, 1*time.Nanosecond)

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Create new cache instance
	cache2 := NewCache(tmpDir)
	if _, ok := cache2.Get(key); ok {
		t.Error("expected cache miss for expired disk entry")
	}

	// Verify file was deleted
	cacheFile := cache1.getCacheFilePath(key)
	if _, err := os.Stat(cacheFile); err == nil {
		t.Error("expected expired cache file to be deleted")
	}
}

func TestCacheCleanupExpired(t *testing.T) {
	cache := NewCache("")

	// Set many entries, some expired
	for i := 0; i < 100; i++ {
		ttl := 1 * time.Hour
		if i%2 == 0 {
			ttl = 1 * time.Nanosecond // Expire half of them
		}
		cache.Set(fmt.Sprintf("key-%d", i), []byte("content"), ttl)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Trigger cleanup by setting another entry (cleanup happens when len > 1000)
	// But we can also test by adding enough entries
	for i := 100; i < 1002; i++ {
		cache.Set(fmt.Sprintf("key-%d", i), []byte("content"), 1*time.Hour)
	}

	// Verify expired entries are gone
	for i := 0; i < 100; i += 2 {
		key := fmt.Sprintf("key-%d", i)
		if _, ok := cache.Get(key); ok {
			t.Errorf("expected expired entry %q to be cleaned up", key)
		}
	}

	// Verify non-expired entries still exist
	for i := 1; i < 100; i += 2 {
		key := fmt.Sprintf("key-%d", i)
		if _, ok := cache.Get(key); !ok {
			t.Errorf("expected non-expired entry %q to still exist", key)
		}
	}
}

func TestGetCacheFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	cache := NewCache(tmpDir)

	tests := []struct {
		name      string
		key       string
		wantExt   string
		wantInDir string
	}{
		{"simple key", "test-key", ".json", tmpDir},
		{"path-like key", "owner/repo/branch/file.md", ".json", tmpDir},
		{"special chars", "key with spaces & symbols!", ".json", tmpDir},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := cache.getCacheFilePath(tt.key)
			if !filepath.IsAbs(path) {
				t.Errorf("expected absolute path, got %q", path)
			}
			if !strings.HasSuffix(path, tt.wantExt) {
				t.Errorf("expected path to end with %q, got %q", tt.wantExt, path)
			}
			if !strings.HasPrefix(path, tt.wantInDir) {
				t.Errorf("expected path to be in %q, got %q", tt.wantInDir, path)
			}
			// Should be a valid filename (no path separators in filename part)
			filename := filepath.Base(path)
			if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
				t.Errorf("filename should not contain path separators: %q", filename)
			}
		})
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewCache("")
	key := "concurrent-key"
	content := []byte("concurrent content")

	// Set from one goroutine
	done := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set(key, content, 1*time.Hour)
		}
		done <- true
	}()

	// Read from another goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get(key)
		}
		done <- true
	}()

	// Wait for both
	<-done
	<-done

	// Verify final state
	retrieved, ok := cache.Get(key)
	if !ok {
		t.Error("expected cache hit after concurrent access")
	}
	if string(retrieved) != string(content) {
		t.Errorf("expected content %q, got %q", string(content), string(retrieved))
	}
}
