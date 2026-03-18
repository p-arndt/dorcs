package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDocsDirExistingDirectory(t *testing.T) {
	dir := t.TempDir()

	resolved, cleanup, err := ResolveDocsDir(dir, false)
	if err != nil {
		t.Fatalf("ResolveDocsDir() error = %v", err)
	}
	defer cleanup()

	if resolved != dir {
		t.Fatalf("ResolveDocsDir() resolved = %q, want %q", resolved, dir)
	}
}

func TestResolveDocsDirMissingWithoutGitHubFails(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "docs")

	_, cleanup, err := ResolveDocsDir(missing, false)
	if cleanup != nil {
		defer cleanup()
	}
	if err == nil {
		t.Fatal("ResolveDocsDir() expected error for missing local docs dir")
	}
}

func TestResolveDocsDirMissingWithGitHubCreatesInternalDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "docs")

	resolved, cleanup, err := ResolveDocsDir(missing, true)
	if err != nil {
		t.Fatalf("ResolveDocsDir() error = %v", err)
	}

	info, statErr := os.Stat(resolved)
	if statErr != nil {
		t.Fatalf("stat resolved dir: %v", statErr)
	}
	if !info.IsDir() {
		t.Fatalf("resolved path is not a directory: %s", resolved)
	}
	if resolved == missing {
		t.Fatalf("ResolveDocsDir() should use an internal directory, got original missing path %q", resolved)
	}

	cleanup()
	if _, err := os.Stat(resolved); !os.IsNotExist(err) {
		t.Fatalf("cleanup did not remove internal dir %q", resolved)
	}
}
