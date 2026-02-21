package site

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStartWatcherReload(t *testing.T) {
	// Create temp docs dir
	docs := t.TempDir()

	// Create site and build initial index
	s, err := New(docs, "github", "", "")
	if err != nil {
		t.Fatalf("New site: %v", err)
	}
	if err := s.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex: %v", err)
	}

	b := NewReloadBroadcaster()
	cleanup, err := s.StartWatcher(docs, b, nil, "", nil)
	if err != nil {
		t.Fatalf("StartWatcher: %v", err)
	}
	defer cleanup()

	// Subscribe to events
	ch := b.Subscribe()
	defer b.Unsubscribe(ch)

	// Create a new markdown file that should trigger a reload
	p := filepath.Join(docs, "newpage.md")
	if err := os.WriteFile(p, []byte("# New Page\n\nContent\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Wait for a reload event (allow some time for debounce + index rebuild)
	select {
	case msg := <-ch:
		if msg != "reload" {
			t.Fatalf("unexpected reload message: %q", msg)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for reload event")
	}

	// Verify the index now contains the new doc
	if _, ok := s.GetDoc("newpage"); !ok {
		t.Fatal("newpage not found in index after reload")
	}
}
