package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/p-arndt/dorcs/internal/config"
)

func writeDocs(t testing.TB, dir string, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%03d-test.md", i)
		p := filepath.Join(dir, name)
		content := fmt.Sprintf("# Page %d\n\nThis is page %d.\n", i, i)
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			t.Fatalf("write doc: %v", err)
		}
	}
}

func TestBuildLanguageConcurrent(t *testing.T) {
	// Create temp docs and output dirs
	docsDir := t.TempDir()
	outDir := t.TempDir()

	// Write a moderate number of docs to exercise concurrency
	writeDocs(t, docsDir, 20)

	// Minimal document template that defines "doc" so server uses it
	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))

	b := New(Config{
		DocsDir:      docsDir,
		RootDir:      docsDir,
		OutputDir:    outDir,
		BasePath:     "",
		SiteConfig:   config.Default(),
		DocumentTmpl: tmpl,
		Parallelism:  4,
	})

	if err := b.buildLanguage("", true, "", false); err != nil {
		t.Fatalf("buildLanguage failed: %v", err)
	}

	// ensure at least one generated page exists (any .html file)
	found := false
	_ = filepath.Walk(outDir, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() && filepath.Ext(p) == ".html" {
			found = true
			return filepath.SkipDir
		}
		return nil
	})
	if !found {
		t.Fatalf("no generated HTML files found in output dir")
	}
}

func BenchmarkBuildLanguage(bm *testing.B) {
	// Benchmark smaller number of docs to keep test time reasonable
	docsDir := bm.TempDir()
	outDir := bm.TempDir()
	writeDocs(bm, docsDir, 100)

	tmpl := template.Must(template.New("doc").Parse(`{{define "doc"}}{{.HTML}}{{end}}`))

	builder := New(Config{
		DocsDir:      docsDir,
		RootDir:      docsDir,
		OutputDir:    outDir,
		BasePath:     "",
		SiteConfig:   config.Default(),
		DocumentTmpl: tmpl,
		Parallelism:  0, // use auto
	})

	bm.ResetTimer()
	for i := 0; i < bm.N; i++ {
		if err := builder.buildLanguage("", true, "", false); err != nil {
			bm.Fatalf("buildLanguage failed: %v", err)
		}
	}
}
