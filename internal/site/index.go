package site

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"dorcs-v2/internal/markdown"
)

// BuildIndex scans the RootDir recursively for ".md" files and builds an in-memory index.
// It reads front matter for metadata and uses the filename as a fallback title.
func (s *Site) BuildIndex() error {
	type found struct {
		key string
		doc *Doc
	}

	var docs []found

	err := filepath.WalkDir(s.RootDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			// Skip common folders that shouldn't be indexed.
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == ".idea" || name == ".vscode" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		rel, err := filepath.Rel(s.RootDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		key := keyFromRel(rel) // drops .md, applies index.md rules (root index -> "")
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		// Parse front matter (best effort); do not fail the whole index if a file has bad metadata.
		meta, contentHash, pmErr := markdown.ReadFrontMatterAndHash(path)
		if pmErr != nil {
			// Keep file indexed with fallbacks; include hash if possible.
			meta = markdown.FrontMatter{}
		}

		doc := &Doc{
			Key:         key,
			FilePath:    path,
			RelPath:     rel,
			DirKey:      dirKeyFromKey(key),
			Title:       strings.TrimSpace(meta.Title),
			Description: strings.TrimSpace(meta.Description),
			Tags:        append([]string(nil), meta.Tags...),
			Draft:       meta.Draft,
			Order:       meta.Order,
			Author:      strings.TrimSpace(meta.Author),
			UpdatedAt:   stat.ModTime(),
			ContentHash: contentHash,
		}

		// Parse date (optional).
		if ds := strings.TrimSpace(meta.Date); ds != "" {
			if t, ok := parseDate(ds); ok {
				doc.Date = t
			}
		}

		// Title fallback:
		// - index pages should show the folder (or site) title, not "index"
		if doc.Title == "" {
			if isIndexRel(rel) {
				doc.Title = titleFromIndexRel(rel)
			} else {
				doc.Title = titleFromKey(doc.Key)
			}
		}

		docs = append(docs, found{key: key, doc: doc})
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk root: %w", err)
	}

	// Atomic swap of index + nav tree.
	s.mu.Lock()
	defer s.mu.Unlock()

	newIndex := make(map[string]*Doc, len(docs))
	for _, f := range docs {
		// Important: allow empty key (root index) to be indexed.
		newIndex[f.key] = f.doc
	}
	s.index = newIndex
	s.nav = buildNavTree(newIndex)
	return nil
}

// ListDocs returns indexed documents sorted by directory then title.
// Draft documents are included if includeDraft is true.
func (s *Site) ListDocs(includeDraft bool) []*Doc {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*Doc, 0, len(s.index))
	for _, d := range s.index {
		if !includeDraft && d.Draft {
			continue
		}
		out = append(out, d)
	}

	sort.Slice(out, func(i, j int) bool {
		di, dj := out[i], out[j]

		// If both have order set (non-zero), sort by order first
		if di.Order != 0 && dj.Order != 0 {
			if di.Order != dj.Order {
				return di.Order < dj.Order
			}
		} else if di.Order != 0 {
			// di has order, dj doesn't - di comes first
			return true
		} else if dj.Order != 0 {
			// dj has order, di doesn't - dj comes first
			return false
		}

		// Check numeric prefixes from filenames
		diPrefix := extractNumericPrefix(di.RelPath)
		djPrefix := extractNumericPrefix(dj.RelPath)
		if diPrefix != 0 && djPrefix != 0 {
			if diPrefix != djPrefix {
				return diPrefix < djPrefix
			}
		} else if diPrefix != 0 {
			// di has prefix, dj doesn't - di comes first
			return true
		} else if djPrefix != 0 {
			// dj has prefix, di doesn't - dj comes first
			return false
		}

		// Prefer date descending if present in both.
		if !di.Date.IsZero() && !dj.Date.IsZero() && !di.Date.Equal(dj.Date) {
			return di.Date.After(dj.Date)
		}

		if di.DirKey != dj.DirKey {
			return di.DirKey < dj.DirKey
		}
		// Title then key.
		if di.Title != dj.Title {
			return strings.ToLower(di.Title) < strings.ToLower(dj.Title)
		}
		return di.Key < dj.Key
	})
	return out
}

// GetDoc returns a document by its URL key (no extension).
// Note: empty key ("") is valid and represents the root index page.
func (s *Site) GetDoc(key string) (*Doc, bool) {
	nk := normalizeKey(key)

	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.index[nk]
	return d, ok
}
