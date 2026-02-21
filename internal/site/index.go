package site

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/p-arndt/dorcs/internal/markdown"
)

// BuildIndex scans the RootDir recursively for ".md" files and builds an in-memory index.
// It reads front matter for metadata and uses the filename as a fallback title.
// If GitHub integration is enabled, local files are skipped and only GitHub files are indexed.
func (s *Site) BuildIndex() error {
	type found struct {
		key string
		doc *Doc
	}

	var docs []found

	// Skip local file indexing if GitHub is enabled
	s.mu.RLock()
	githubEnabled := s.githubClient != nil
	s.mu.RUnlock()

	// Only walk local directory if GitHub is not enabled
	if !githubEnabled {
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
				Key:          key,
				FilePath:     path,
				RelPath:      rel,
				DirKey:       dirKeyFromKey(key),
				Title:        strings.TrimSpace(meta.Title),
				Description:  strings.TrimSpace(meta.Description),
				Tags:         append([]string(nil), meta.Tags...),
				Draft:        meta.Draft,
				Order:        meta.Order,
				Author:       strings.TrimSpace(meta.Author),
				After:        strings.TrimSpace(meta.After),
				Presentation: meta.Presentation,
				UpdatedAt:    stat.ModTime(),
				ContentHash:  contentHash,
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
	}

	// Index GitHub files if configured
	if s.githubClient != nil {
		githubFiles, err := s.githubClient.DiscoverMarkdownFiles(s.githubOwner, s.githubRepo, s.githubBranch, s.githubPath)
		if err != nil {
			// If GitHub is enabled, we must have GitHub files - fail if we can't get them
			return fmt.Errorf("failed to discover GitHub files from %s/%s@%s/%s: %w", s.githubOwner, s.githubRepo, s.githubBranch, s.githubPath, err)
		}
		if len(githubFiles) > 0 {
			fmt.Printf("dorcs: discovered %d markdown files from GitHub\n", len(githubFiles))
			for _, filePath := range githubFiles {

				// Generate key from relative path (keyFromRel expects .md extension)
				key := keyFromRel(filePath)

				// Create a relative path without .md extension for RelPath field
				relPath := filePath
				if strings.HasSuffix(strings.ToLower(relPath), ".md") {
					relPath = relPath[:len(relPath)-3]
				}

				// Create cache key for GitHub content
				fullGitHubPath := filePath
				if s.githubPath != "" {
					fullGitHubPath = s.githubPath + "/" + filePath
				}
				cacheKey := fmt.Sprintf("%s/%s/%s/%s", s.githubOwner, s.githubRepo, s.githubBranch, fullGitHubPath)

				// Create virtual Doc entry
				doc := &Doc{
					Key:            key,
					FilePath:       "", // No local file path for GitHub docs
					RelPath:        relPath + ".md",
					DirKey:         dirKeyFromKey(key),
					Title:          titleFromKey(key),
					IsGitHub:       true,
					GitHubPath:     fullGitHubPath,
					GitHubCacheKey: cacheKey,
					UpdatedAt:      time.Now(), // Use current time as fallback
				}

				// Try to fetch and parse front matter for metadata
				// This is best effort - if it fails, we still index the file with fallback title
				content, err := s.githubClient.FetchMarkdown(s.githubOwner, s.githubRepo, s.githubBranch, fullGitHubPath)
				if err != nil {
					// Log warning but continue - file will be indexed with fallback metadata
					fmt.Printf("dorcs: warning: failed to fetch GitHub file %s for metadata: %v\n", fullGitHubPath, err)
				} else {
					// Parse front matter from content
					meta, _, parseErr := markdown.ParseFrontMatterFromContent(content)
					if parseErr != nil {
						// Log but continue - front matter parsing is best effort
						fmt.Printf("dorcs: warning: failed to parse front matter for %s: %v\n", fullGitHubPath, parseErr)
					} else if meta != nil {
						doc.Title = strings.TrimSpace(meta.Title)
						doc.Description = strings.TrimSpace(meta.Description)
						doc.Tags = append([]string(nil), meta.Tags...)
						doc.Draft = meta.Draft
						doc.Order = meta.Order
						doc.Author = strings.TrimSpace(meta.Author)
						doc.After = strings.TrimSpace(meta.After)
						doc.Presentation = meta.Presentation

						// Parse date
						if ds := strings.TrimSpace(meta.Date); ds != "" {
							if t, ok := parseDate(ds); ok {
								doc.Date = t
							}
						}
					}
				}

				// Title fallback
				if doc.Title == "" {
					if isIndexRel(relPath + ".md") {
						doc.Title = titleFromIndexRel(relPath + ".md")
					} else {
						doc.Title = titleFromKey(doc.Key)
					}
				}

				docs = append(docs, found{key: key, doc: doc})
			}
		}
	}

	// Atomic swap of index + nav tree.
	s.mu.Lock()
	newIndex := make(map[string]*Doc, len(docs))
	for _, f := range docs {
		// Important: allow empty key (root index) to be indexed.
		newIndex[f.key] = f.doc
	}
	s.index = newIndex
	s.nav = buildNavTree(newIndex)
	s.mu.Unlock()

	// Check for broken links after building the index
	s.ReportBrokenLinks()

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
