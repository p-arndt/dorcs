package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// TreeEntry represents a file or directory in a GitHub tree.
type TreeEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" or "tree"
	URL  string `json:"url"`
}

// TreeResponse represents the GitHub API tree response.
type TreeResponse struct {
	Tree []TreeEntry `json:"tree"`
}

// ContentResponse represents the GitHub API content response.
type ContentResponse struct {
	Type        string `json:"type"`
	Encoding    string `json:"encoding"`
	Content     string `json:"content"`      // For files < 1MB
	DownloadURL string `json:"download_url"` // For files >= 1MB
	Path        string `json:"path"`
	Size        int    `json:"size"`
}

// RepositoryInfo holds parsed repository information.
type RepositoryInfo struct {
	Owner  string
	Repo   string
	Branch string
	Path   string
}

// Client provides GitHub API access.
type Client struct {
	httpClient *http.Client
	token      string
	cache      *Cache
	cacheTTL   time.Duration
}

// NewClient creates a new GitHub client.
func NewClient(token string, cache *Cache, cacheTTL time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token:    token,
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

// ParseRepositoryURL parses a GitHub repository tree URL.
// Supports formats:
// - https://github.com/owner/repo/tree/branch/path
// - https://github.com/owner/repo
// - owner/repo/tree/branch/path
func ParseRepositoryURL(repoURL string) (*RepositoryInfo, error) {
	// Normalize URL - add https://github.com if missing
	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		if !strings.HasPrefix(repoURL, "github.com/") {
			repoURL = "https://github.com/" + repoURL
		} else {
			repoURL = "https://" + repoURL
		}
	}

	// Parse URL
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	// Extract owner and repo from path
	// Path format: /owner/repo/tree/branch/path or /owner/repo
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repository URL: expected owner/repo")
	}

	owner := parts[0]
	repo := parts[1]
	branch := "main" // default branch
	repoPath := ""

	// Check if we have tree/branch/path structure
	if len(parts) >= 4 && parts[2] == "tree" {
		branch = parts[3]
		if len(parts) > 4 {
			repoPath = strings.Join(parts[4:], "/")
		}
	} else if len(parts) > 2 {
		// Just path without tree/branch
		repoPath = strings.Join(parts[2:], "/")
	}

	return &RepositoryInfo{
		Owner:  owner,
		Repo:   repo,
		Branch: branch,
		Path:   repoPath,
	}, nil
}

// GetDefaultBranch gets the default branch for a repository.
func (c *Client) GetDefaultBranch(owner, repo string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get repository info: %s", resp.Status)
	}

	var repoInfo struct {
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return "", err
	}

	return repoInfo.DefaultBranch, nil
}

// ListDirectory lists the contents of a directory in a GitHub repository.
func (c *Client) ListDirectory(owner, repo, branch, dirPath string) ([]TreeEntry, error) {
	// Get the SHA of the branch
	sha, err := c.getBranchSHA(owner, repo, branch)
	if err != nil {
		return nil, err
	}

	// Get the tree recursively
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", owner, repo, sha)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("directory not found: %s (404)", dirPath)
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("authentication failed: %s (check token permissions)", resp.Status)
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limit exceeded: %s (try again later)", resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list directory: %s - %s", resp.Status, string(body))
	}

	var treeResp TreeResponse
	if err := json.NewDecoder(resp.Body).Decode(&treeResp); err != nil {
		return nil, err
	}

	// Filter entries that are within the specified directory path
	var entries []TreeEntry
	dirPath = strings.Trim(dirPath, "/")
	for _, entry := range treeResp.Tree {
		if dirPath == "" {
			// Root directory - include all entries
			entries = append(entries, entry)
		} else if strings.HasPrefix(entry.Path, dirPath+"/") || entry.Path == dirPath {
			// Entry is within the directory
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// FetchMarkdown fetches a markdown file from GitHub.
func (c *Client) FetchMarkdown(owner, repo, branch, filePath string) ([]byte, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s/%s/%s/%s", owner, repo, branch, filePath)
	if c.cache != nil {
		if content, ok := c.cache.Get(cacheKey); ok {
			return content, nil
		}
	}

	// Fetch from GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, filePath, branch)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file not found: %s (404) - %s", filePath, string(body))
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("authentication failed for %s: %s - %s (check token permissions)", filePath, resp.Status, string(body))
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		// Try to read rate limit headers
		retryAfter := resp.Header.Get("Retry-After")
		rateLimitRemaining := resp.Header.Get("X-RateLimit-Remaining")
		msg := fmt.Sprintf("rate limit exceeded: %s", resp.Status)
		if retryAfter != "" {
			msg += fmt.Sprintf(" (retry after %s seconds)", retryAfter)
		}
		if rateLimitRemaining != "" {
			msg += fmt.Sprintf(" (remaining: %s)", rateLimitRemaining)
		}
		return nil, fmt.Errorf("%s - try again later or use a GitHub token to increase rate limits", msg)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch file %s: %s - %s", filePath, resp.Status, string(body))
	}

	var contentResp ContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
		return nil, err
	}

	if contentResp.Type != "file" {
		return nil, fmt.Errorf("path is not a file: %s", filePath)
	}

	var content []byte

	// Handle large files (>= 1MB) - GitHub provides download_url instead of content
	if contentResp.DownloadURL != "" {
		// Fetch content from download URL
		downloadReq, err := http.NewRequest("GET", contentResp.DownloadURL, nil)
		if err != nil {
			return nil, fmt.Errorf("create download request: %w", err)
		}

		downloadResp, err := c.httpClient.Do(downloadReq)
		if err != nil {
			return nil, fmt.Errorf("download file: %w", err)
		}
		defer downloadResp.Body.Close()

		if downloadResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(downloadResp.Body)
			return nil, fmt.Errorf("failed to download file: %s - %s", downloadResp.Status, string(body))
		}

		content, err = io.ReadAll(downloadResp.Body)
		if err != nil {
			return nil, fmt.Errorf("read downloaded content: %w", err)
		}
	} else {
		// Handle small files (< 1MB) - content is base64 encoded
		if contentResp.Encoding != "base64" {
			return nil, fmt.Errorf("unsupported encoding: %s", contentResp.Encoding)
		}

		// Decode base64 content
		var err error
		content, err = base64.StdEncoding.DecodeString(contentResp.Content)
		if err != nil {
			return nil, fmt.Errorf("decode base64: %w", err)
		}
	}

	// Cache the content
	if c.cache != nil {
		c.cache.Set(cacheKey, content, c.cacheTTL)
	}

	return content, nil
}

// DiscoverMarkdownFiles recursively discovers all .md files in a directory tree.
func (c *Client) DiscoverMarkdownFiles(owner, repo, branch, rootPath string) ([]string, error) {
	entries, err := c.ListDirectory(owner, repo, branch, rootPath)
	if err != nil {
		return nil, err
	}

	var markdownFiles []string
	rootPath = strings.Trim(rootPath, "/")
	mdRegex := regexp.MustCompile(`\.md$`)

	for _, entry := range entries {
		// Only process files (blobs), not directories (trees)
		if entry.Type != "blob" {
			continue
		}

		// Check if it's a markdown file
		if !mdRegex.MatchString(strings.ToLower(entry.Path)) {
			continue
		}

		// If rootPath is specified, remove it from the path
		relativePath := entry.Path
		if rootPath != "" && strings.HasPrefix(entry.Path, rootPath+"/") {
			relativePath = strings.TrimPrefix(entry.Path, rootPath+"/")
		} else if rootPath != "" && entry.Path == rootPath {
			// This shouldn't happen for files, but handle it
			continue
		} else if rootPath != "" && !strings.HasPrefix(entry.Path, rootPath+"/") {
			// File is not in the rootPath directory, skip it
			continue
		}

		markdownFiles = append(markdownFiles, relativePath)
	}

	return markdownFiles, nil
}

// getBranchSHA gets the SHA of a branch or tag.
func (c *Client) getBranchSHA(owner, repo, ref string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/ref/heads/%s", owner, repo, ref)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Try as a tag
		apiURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/git/ref/tags/%s", owner, repo, ref)
		req, err = http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return "", err
		}

		if c.token != "" {
			req.Header.Set("Authorization", "token "+c.token)
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get branch SHA: %s - %s", resp.Status, string(body))
	}

	var refResp struct {
		Object struct {
			SHA string `json:"sha"`
		} `json:"object"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&refResp); err != nil {
		return "", err
	}

	return refResp.Object.SHA, nil
}
