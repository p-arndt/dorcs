package site

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ConfigReloadCallback is a function that gets called when the config file changes.
// It should reload the config and return an error if reloading fails.
type ConfigReloadCallback func() error

// RebuildCallback is a function that rebuilds site index(es) when files change.
// If nil, the default site's BuildIndex is used.
type RebuildCallback func() error

// StartWatcher starts a file watcher that automatically rebuilds the index
// when markdown files change. watchRoot is the top-level directory to watch
// (e.g. the docs directory); it should include all language/version subdirs
// so changes in any of them are detected. If rebuildAll is provided, it is
// called instead of s.BuildIndex() to rebuild all affected sites.
// If broadcaster is provided, it will notify connected clients for live reload.
// If configReload is provided, it will be called when config files change.
// configFilePath is an optional path to a config file outside the docs directory to watch.
// It returns a cleanup function that should be called to stop the watcher.
func (s *Site) StartWatcher(watchRoot string, broadcaster *ReloadBroadcaster, configReload ConfigReloadCallback, configFilePath string, rebuildAll RebuildCallback) (func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Use watchRoot if provided, otherwise fall back to site's RootDir
	rootToWatch := watchRoot
	if rootToWatch == "" {
		rootToWatch = s.RootDir
	}
	absWatchRoot, err := filepath.Abs(rootToWatch)
	if err != nil {
		watcher.Close()
		return nil, err
	}

	// Add the root directory and all subdirectories
	if err := s.addWatchRecursive(watcher, absWatchRoot); err != nil {
		watcher.Close()
		return nil, err
	}

	// If config file is outside the watched directory, watch its directory too
	if configFilePath != "" {
		configDir := filepath.Dir(configFilePath)
		absConfigDir, err := filepath.Abs(configDir)
		if err == nil {
			// Only add if it's different from the watch root
			if absConfigDir != absWatchRoot {
				if err := watcher.Add(absConfigDir); err != nil {
					log.Printf("dorcs: warning: failed to watch config directory %s: %v", absConfigDir, err)
				}
			}
		}
	}

	// Debounce timers to avoid rebuilding multiple times for rapid changes
	// 500ms gives enough time for users who are actively typing/saving
	var debounceTimer *time.Timer
	var configDebounceTimer *time.Timer
	debounceDuration := 500 * time.Millisecond
	// Track if we saw Create/Remove/Rename (index must be rebuilt); Write/Chmod only need reload
	var structuralChange atomic.Bool
	// Config files get longer debounce since they're changed less frequently
	// and full page reloads are more expensive
	configDebounceDuration := 1000 * time.Millisecond

	stopChan := make(chan struct{})
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)
		defer watcher.Close()
		// Ensure debounce timers are stopped when the goroutine exits to avoid races
		defer func() {
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			if configDebounceTimer != nil {
				configDebounceTimer.Stop()
			}
		}()

		for {
			select {
			case <-stopChan:
				return

			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Check if this is a config file change
				// Only react to Write events to avoid reacting to Chmod, Create, etc.
				if s.isConfigFile(event.Name) && event.Has(fsnotify.Write) {
					if configReload != nil {
						// Debounce: reset timer on each event
						// Use longer debounce for config files to prevent rapid reloads
						if configDebounceTimer != nil {
							configDebounceTimer.Stop()
						}

						configDebounceTimer = time.AfterFunc(configDebounceDuration, func() {
							log.Printf("dorcs: detected config file changes, reloading config...")
							if err := configReload(); err != nil {
								log.Printf("dorcs: error reloading config: %v", err)
							} else {
								log.Printf("dorcs: config reloaded successfully")
								// Notify connected browsers to do a full page reload
								// Config changes (theme, title, etc.) require full reload
								if broadcaster != nil {
									broadcaster.Notify("fullreload")
								}
							}
						})
					}
					continue
				}

				// Only process markdown files and directory operations
				if !s.shouldReloadForEvent(event) {
					continue
				}

				// Create/Remove/Rename affect index (nav, doc list); Write/Chmod are content-only
				if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					structuralChange.Store(true)
				}

				// If a new directory was created, add it to the watcher
				if event.Has(fsnotify.Create) {
					if info, err := filepath.Abs(event.Name); err == nil {
						if stat, err := os.Lstat(info); err == nil && stat.IsDir() {
							_ = s.addWatchRecursive(watcher, info)
						}
					}
				}

				// Debounce: reset timer on each event
				if debounceTimer != nil {
					debounceTimer.Stop()
				}

				debounceTimer = time.AfterFunc(debounceDuration, func() {
					needRebuild := structuralChange.Swap(false)
					if needRebuild {
						log.Printf("dorcs: detected structural changes, rebuilding index...")
						var err error
						if rebuildAll != nil {
							err = rebuildAll()
						} else {
							err = s.BuildIndex()
						}
						if err != nil {
							log.Printf("dorcs: error rebuilding index: %v", err)
						} else {
							log.Printf("dorcs: index rebuilt successfully")
						}
					}
					// Always notify - browser fetches fresh content (RenderDoc reads from disk)
					if broadcaster != nil {
						broadcaster.Notify("reload")
					}
				})

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("dorcs: watcher error: %v", err)
			}
		}
	}()

	// Return cleanup function
	cleanup := func() {
		close(stopChan)
		<-doneChan // Wait for goroutine to finish
	}

	return cleanup, nil
}

// addWatchRecursive adds a directory and all its subdirectories to the watcher.
func (s *Site) addWatchRecursive(watcher *fsnotify.Watcher, dir string) error {
	if err := watcher.Add(dir); err != nil {
		return err
	}

	entries, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := os.Lstat(entry)
		if err != nil {
			continue
		}

		// Skip hidden directories and files
		base := filepath.Base(entry)
		if strings.HasPrefix(base, ".") {
			continue
		}

		if info.IsDir() {
			if err := s.addWatchRecursive(watcher, entry); err != nil {
				log.Printf("dorcs: warning: failed to watch %s: %v", entry, err)
			}
		}
	}

	return nil
}

// shouldReloadForEvent determines if a file system event should trigger a reload.
func (s *Site) shouldReloadForEvent(event fsnotify.Event) bool {
	// We care about: Write, Create, Remove, Rename, Chmod
	// Chmod is included because some editors on Windows/macOS may emit it when saving
	if !event.Has(fsnotify.Write) &&
		!event.Has(fsnotify.Create) &&
		!event.Has(fsnotify.Remove) &&
		!event.Has(fsnotify.Rename) &&
		!event.Has(fsnotify.Chmod) {
		return false
	}

	name := filepath.Base(event.Name)

	// Skip hidden files
	if strings.HasPrefix(name, ".") {
		return false
	}

	// Skip temporary files (common editor patterns)
	if strings.HasSuffix(name, "~") ||
		strings.HasSuffix(name, ".swp") ||
		strings.HasSuffix(name, ".tmp") ||
		strings.HasPrefix(name, "#") {
		return false
	}

	// For file operations, only care about .md files
	if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Chmod) {
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".md" && ext != ".markdown" {
			// Unless it's a directory operation
			if info, err := os.Lstat(event.Name); err != nil || !info.IsDir() {
				return false
			}
		}
	}

	return true
}

// isConfigFile checks if the given file path is a config file (dorcs.yaml, dorcs.yml, or dorcs.json).
// It also checks if the file path matches a specific config file path.
func (s *Site) isConfigFile(filePath string) bool {
	name := filepath.Base(filePath)
	return name == "dorcs.yaml" || name == "dorcs.yml" || name == "dorcs.json"
}
