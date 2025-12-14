// Package auth provides authentication and session management for dorcs.
package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	// Session cookie name
	sessionCookieName = "dorcs_session"
	// Session duration
	sessionDuration = 24 * time.Hour
	// Session file name
	sessionFile = ".dorcs_sessions.json"
)

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	// Username for login
	Username string
	// Password hash (argon2id)
	PasswordHash string
	// Session secret for signing cookies
	SessionSecret string
	// Sessions file path
	SessionsPath string
}

// Session represents a user session.
type Session struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionStore manages active sessions.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	filePath string
}

// NewSessionStore creates a new session store.
func NewSessionStore(filePath string) (*SessionStore, error) {
	store := &SessionStore{
		sessions: make(map[string]*Session),
		filePath: filePath,
	}

	// Load existing sessions from file
	if err := store.loadSessions(); err != nil {
		// If file doesn't exist, that's okay - start with empty store
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("load sessions: %w", err)
		}
	}

	// Clean up expired sessions periodically
	go store.cleanupExpired()

	return store, nil
}

// loadSessions loads sessions from disk.
func (s *SessionStore) loadSessions() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var sessions []*Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sess := range sessions {
		if sess.ExpiresAt.After(time.Now()) {
			s.sessions[sess.ID] = sess
		}
	}

	return nil
}

// saveSessions saves sessions to disk.
func (s *SessionStore) saveSessions() error {
	s.mu.RLock()
	sessions := make([]*Session, 0, len(s.sessions))
	for _, sess := range s.sessions {
		sessions = append(sessions, sess)
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first, then rename (atomic write)
	tmpPath := s.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.filePath)
}

// cleanupExpired periodically removes expired sessions.
func (s *SessionStore) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for id, sess := range s.sessions {
			if sess.ExpiresAt.Before(now) {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()

		// Save after cleanup
		s.saveSessions()
	}
}

// CreateSession creates a new session for a user.
func (s *SessionStore) CreateSession(username string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	sess := &Session{
		ID:        sessionID,
		Username:  username,
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	s.mu.Lock()
	s.sessions[sessionID] = sess
	s.mu.Unlock()

	// Save to disk
	if err := s.saveSessions(); err != nil {
		// Log error but don't fail - session is still valid in memory
		fmt.Printf("warning: failed to save session: %v\n", err)
	}

	return sess, nil
}

// GetSession retrieves a session by ID.
func (s *SessionStore) GetSession(sessionID string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[sessionID]
	if !ok {
		return nil, false
	}
	if sess.ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	return sess, true
}

// DeleteSession removes a session.
func (s *SessionStore) DeleteSession(sessionID string) {
	s.mu.Lock()
	delete(s.sessions, sessionID)
	s.mu.Unlock()
	s.saveSessions()
}

// generateSessionID generates a random session ID.
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// HashPassword hashes a password using argon2id.
func HashPassword(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash password with argon2id
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Encode: salt + hash
	encoded := base64.RawStdEncoding.EncodeToString(append(salt, hash...))
	return encoded, nil
}

// VerifyPassword verifies a password against a hash.
func VerifyPassword(password, hash string) (bool, error) {
	// Decode hash
	decoded, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	if len(decoded) < 16 {
		return false, fmt.Errorf("invalid hash format")
	}

	// Extract salt and hash
	salt := decoded[:16]
	expectedHash := decoded[16:]

	// Compute hash
	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Constant-time comparison
	return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1, nil
}

// GetSessionFromRequest extracts the session from a request.
func GetSessionFromRequest(r *http.Request, store *SessionStore) (*Session, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, false
	}

	sess, ok := store.GetSession(cookie.Value)
	return sess, ok
}

// SetSessionCookie sets the session cookie on the response.
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(sessionDuration.Seconds()),
	})
}

// ClearSessionCookie clears the session cookie.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// AuthManager manages authentication.
type AuthManager struct {
	config *AuthConfig
	store  *SessionStore
}

// NewAuthManager creates a new authentication manager.
func NewAuthManager(config *AuthConfig) (*AuthManager, error) {
	// Ensure sessions directory exists
	sessionsDir := filepath.Dir(config.SessionsPath)
	if err := os.MkdirAll(sessionsDir, 0700); err != nil {
		return nil, fmt.Errorf("create sessions dir: %w", err)
	}

	store, err := NewSessionStore(config.SessionsPath)
	if err != nil {
		return nil, fmt.Errorf("create session store: %w", err)
	}

	return &AuthManager{
		config: config,
		store:  store,
	}, nil
}

// Login authenticates a user and creates a session.
func (a *AuthManager) Login(username, password string) (*Session, error) {
	// Verify username
	if username != a.config.Username {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	valid, err := VerifyPassword(password, a.config.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("verify password: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Create session
	return a.store.CreateSession(username)
}

// Logout removes a session.
func (a *AuthManager) Logout(sessionID string) {
	a.store.DeleteSession(sessionID)
}

// IsAuthenticated checks if a request is authenticated.
func (a *AuthManager) IsAuthenticated(r *http.Request) bool {
	_, ok := GetSessionFromRequest(r, a.store)
	return ok
}

// RequireAuth is middleware that requires authentication.
func (a *AuthManager) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.IsAuthenticated(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetStore returns the session store (for internal use).
func (a *AuthManager) GetStore() *SessionStore {
	return a.store
}

