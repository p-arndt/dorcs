package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/p-arndt/dorcs/internal/config"
	"github.com/p-arndt/dorcs/internal/site"
)

func TestRebuildSitesWithExplicitNavRollsBackOnError(t *testing.T) {
	firstDir := t.TempDir()
	secondDir := t.TempDir()

	writeDoc := func(root, relPath, content string) {
		fullPath := filepath.Join(root, filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("failed to create dir for %q: %v", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %q: %v", fullPath, err)
		}
	}

	writeDoc(firstDir, "index.md", "# Home")
	writeDoc(firstDir, "guide.md", "# Guide")
	writeDoc(secondDir, "index.md", "# Home")

	firstSite, err := site.New(firstDir, "github", "", "")
	if err != nil {
		t.Fatalf("New(firstSite) failed: %v", err)
	}
	secondSite, err := site.New(secondDir, "github", "", "")
	if err != nil {
		t.Fatalf("New(secondSite) failed: %v", err)
	}

	if err := firstSite.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex(firstSite) failed: %v", err)
	}
	if err := secondSite.BuildIndex(); err != nil {
		t.Fatalf("BuildIndex(secondSite) failed: %v", err)
	}

	newNav := config.NavItems{{Label: "Guide", Page: "guide.md"}}
	if err := rebuildSitesWithExplicitNav([]*site.Site{firstSite, secondSite}, newNav); err == nil {
		t.Fatal("expected rebuildSitesWithExplicitNav to fail when one site is missing an explicit nav doc")
	}

	if err := os.Remove(filepath.Join(firstDir, "guide.md")); err != nil {
		t.Fatalf("failed to remove guide.md: %v", err)
	}

	if err := firstSite.BuildIndex(); err != nil {
		t.Fatalf("expected rollback to restore previous explicit nav, BuildIndex() failed: %v", err)
	}
}
