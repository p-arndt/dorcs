// Package server provides the HTTP handler for serving markdown documentation.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SearchResponse represents the JSON response for search API.
type SearchResponse struct {
	Results []SearchResultItem `json:"results"`
	Query   string             `json:"query"`
}

// SearchResultItem represents a single search result in the API response.
type SearchResultItem struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	Snippet     string `json:"snippet"`
	Score       int    `json:"score"`
	HeadingID   string `json:"heading_id,omitempty"`
	HeadingText string `json:"heading_text,omitempty"`
}

// ServeSearch handles the search API endpoint.
func (h *Handler) ServeSearch(w http.ResponseWriter, r *http.Request) {
	h.handleSearch(w, r)
}

// handleSearch handles the search API endpoint.
func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		// Return empty results for empty query
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResponse{
			Results: []SearchResultItem{},
			Query:   "",
		})
		return
	}

	h.mu.RLock()
	site := h.cfg.Site
	hideDraft := h.cfg.HideDraft
	basePath := h.cfg.BasePath
	h.mu.RUnlock()

	if site == nil {
		http.Error(w, "server not configured", http.StatusInternalServerError)
		return
	}

	// Perform search
	searchResults := site.SearchDocs(query, !hideDraft, 100) // Max 100 results

	// Convert to API response format
	results := make([]SearchResultItem, 0, len(searchResults))
	for _, sr := range searchResults {
		// Build full path with base path
		path := sr.Path
		if basePath != "" {
			path = basePath + path
		}

		results = append(results, SearchResultItem{
			Key:         sr.Key,
			Title:       sr.Title,
			Path:        path,
			Snippet:     sr.Snippet,
			Score:       sr.Score,
			HeadingID:   sr.HeadingID,
			HeadingText: sr.HeadingText,
		})
	}

	response := SearchResponse{
		Results: results,
		Query:   query,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ServeSitemap handles the sitemap.xml endpoint.
func (h *Handler) ServeSitemap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.mu.RLock()
	site := h.cfg.Site
	hideDraft := h.cfg.HideDraft
	basePath := h.cfg.BasePath
	h.mu.RUnlock()

	if site == nil {
		http.Error(w, "server not configured", http.StatusInternalServerError)
		return
	}

	// Get all non-draft documents
	docs := site.ListDocs(!hideDraft)

	// Build base URL from request
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	host := r.Host
	if host == "" {
		host = "localhost"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	if basePath != "" {
		baseURL += basePath
	}

	// Generate sitemap XML
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour

	// Write XML header and root element
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n"))
	w.Write([]byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"))

	// Add each document to the sitemap
	var urlBuilder strings.Builder
	for _, doc := range docs {
		// Build URL path
		urlBuilder.Reset()
		urlBuilder.WriteString(baseURL)
		if doc.Key == "" {
			// Root index page
			if !strings.HasSuffix(baseURL, "/") {
				urlBuilder.WriteByte('/')
			}
		} else {
			// Ensure single slash between basePath and key
			if !strings.HasSuffix(baseURL, "/") {
				urlBuilder.WriteByte('/')
			}
			// URL-encode each path segment
			parts := strings.Split(doc.Key, "/")
			for i, part := range parts {
				if i > 0 {
					urlBuilder.WriteByte('/')
				}
				urlBuilder.WriteString(url.PathEscape(part))
			}
		}
		urlPath := urlBuilder.String()

		// Format last modified date (W3C datetime format)
		lastmod := doc.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z")

		// Determine priority based on document type
		priority := "0.7" // Default priority
		if doc.Key == "" {
			priority = "1.0" // Homepage gets highest priority
		} else if !strings.Contains(doc.Key, "/") {
			priority = "0.9" // Top-level pages
		}

		// Determine changefreq (how often the page is likely to change)
		changefreq := "monthly" // Default
		if doc.Key == "" {
			changefreq = "weekly" // Homepage changes more frequently
		}

		// Write URL entry
		w.Write([]byte("  <url>\n"))
		fmt.Fprintf(w, "    <loc>%s</loc>\n", escapeXML(urlPath))
		fmt.Fprintf(w, "    <lastmod>%s</lastmod>\n", lastmod)
		fmt.Fprintf(w, "    <changefreq>%s</changefreq>\n", changefreq)
		fmt.Fprintf(w, "    <priority>%s</priority>\n", priority)
		w.Write([]byte("  </url>\n"))
	}

	// Close root element
	w.Write([]byte("</urlset>\n"))
}
