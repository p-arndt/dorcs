// Package server provides edit mode handlers for authenticated users.
package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dorcs-v2/internal/auth"
)

// EditHandlers provides handlers for edit mode operations.
type EditHandlers struct {
	authManager *auth.AuthManager
	docsDir     string
	basePath    string
}

// NewEditHandlers creates a new edit handlers instance.
func NewEditHandlers(authManager *auth.AuthManager, docsDir, basePath string) *EditHandlers {
	return &EditHandlers{
		authManager: authManager,
		docsDir:     docsDir,
		basePath:    basePath,
	}
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response.
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// HandleLogin handles login requests.
func (h *EditHandlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	session, err := h.authManager.Login(req.Username, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	auth.SetSessionCookie(w, session.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
	})
}

// HandleLogout handles logout requests.
func (h *EditHandlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, ok := auth.GetSessionFromRequest(r, h.authManager.GetStore())
	if ok {
		h.authManager.Logout(session.ID)
	}

	auth.ClearSessionCookie(w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
	})
}

// HandleCheckAuth checks if the user is authenticated.
func (h *EditHandlers) HandleCheckAuth(w http.ResponseWriter, r *http.Request) {
	isAuth := h.authManager.IsAuthenticated(r)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": isAuth,
	})
}

// FileInfo represents file information.
type FileInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	IsDir     bool      `json:"is_dir"`
	Size      int64     `json:"size"`
	Modified  time.Time `json:"modified"`
	Extension string    `json:"extension,omitempty"`
}

// HandleListFiles lists files in the docs directory.
func (h *EditHandlers) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get path parameter (optional, defaults to root)
	pathParam := r.URL.Query().Get("path")
	if pathParam == "" {
		pathParam = "."
	}

	// Clean and validate path
	relPath := filepath.Clean(pathParam)
	if strings.HasPrefix(relPath, "..") || strings.HasPrefix(relPath, "/") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.docsDir, relPath)

	// Ensure path is within docs directory
	absDocsDir, err := filepath.Abs(h.docsDir)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(absFilePath, absDocsDir+string(filepath.Separator)) && absFilePath != absDocsDir {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Read directory
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		filePath := filepath.Join(relPath, entry.Name())
		if relPath == "." {
			filePath = entry.Name()
		}

		fileInfo := FileInfo{
			Name:     entry.Name(),
			Path:     filePath,
			IsDir:    entry.IsDir(),
			Size:     info.Size(),
			Modified: info.ModTime(),
		}

		if !entry.IsDir() {
			fileInfo.Extension = strings.ToLower(filepath.Ext(entry.Name()))
		}

		files = append(files, fileInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"files": files,
		"path":  relPath,
	})
}

// HandleReadFile reads a file's content.
func (h *EditHandlers) HandleReadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get path parameter
	pathParam := r.URL.Query().Get("path")
	if pathParam == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}

	// Clean and validate path
	relPath := filepath.Clean(pathParam)
	if strings.HasPrefix(relPath, "..") || strings.HasPrefix(relPath, "/") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.docsDir, relPath)

	// Ensure path is within docs directory
	absDocsDir, err := filepath.Abs(h.docsDir)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(absFilePath, absDocsDir+string(filepath.Separator)) && absFilePath != absDocsDir {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Read file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "read error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"path":    relPath,
		"content": string(content),
	})
}

// SaveFileRequest represents a save file request.
type SaveFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// HandleSaveFile saves file content.
func (h *EditHandlers) HandleSaveFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SaveFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Clean and validate path
	relPath := filepath.Clean(req.Path)
	if strings.HasPrefix(relPath, "..") || strings.HasPrefix(relPath, "/") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.docsDir, relPath)

	// Ensure path is within docs directory
	absDocsDir, err := filepath.Abs(h.docsDir)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(absFilePath, absDocsDir+string(filepath.Separator)) && absFilePath != absDocsDir {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, "failed to create directory", http.StatusInternalServerError)
		return
	}

	// Write file
	if err := os.WriteFile(fullPath, []byte(req.Content), 0644); err != nil {
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"path":    relPath,
	})
}

// CreateFileRequest represents a create file request.
type CreateFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content,omitempty"`
	IsDir   bool   `json:"is_dir,omitempty"`
}

// HandleCreateFile creates a new file or directory.
func (h *EditHandlers) HandleCreateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Clean and validate path
	relPath := filepath.Clean(req.Path)
	if strings.HasPrefix(relPath, "..") || strings.HasPrefix(relPath, "/") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.docsDir, relPath)

	// Ensure path is within docs directory
	absDocsDir, err := filepath.Abs(h.docsDir)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(absFilePath, absDocsDir+string(filepath.Separator)) && absFilePath != absDocsDir {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Check if file/dir already exists
	if _, err := os.Stat(fullPath); err == nil {
		http.Error(w, "file or directory already exists", http.StatusConflict)
		return
	}

	if req.IsDir {
		// Create directory
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			http.Error(w, "failed to create directory", http.StatusInternalServerError)
			return
		}
	} else {
		// Ensure parent directory exists
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			http.Error(w, "failed to create directory", http.StatusInternalServerError)
			return
		}

		// Create file with content
		content := req.Content
		if content == "" {
			content = ""
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			http.Error(w, "failed to create file", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"path":    relPath,
	})
}

// DeleteFileRequest represents a delete file request.
type DeleteFileRequest struct {
	Path string `json:"path"`
}

// HandleDeleteFile deletes a file or directory.
func (h *EditHandlers) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Clean and validate path
	relPath := filepath.Clean(req.Path)
	if strings.HasPrefix(relPath, "..") || strings.HasPrefix(relPath, "/") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.docsDir, relPath)

	// Ensure path is within docs directory
	absDocsDir, err := filepath.Abs(h.docsDir)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	absFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(absFilePath, absDocsDir+string(filepath.Separator)) && absFilePath != absDocsDir {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// Prevent deleting the root docs directory
	if absFilePath == absDocsDir {
		http.Error(w, "cannot delete root directory", http.StatusBadRequest)
		return
	}

	// Delete file or directory
	if err := os.RemoveAll(fullPath); err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

