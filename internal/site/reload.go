package site

import (
	"fmt"
	"net/http"
	"sync"
)

// ReloadBroadcaster manages Server-Sent Events (SSE) connections for live reload.
type ReloadBroadcaster struct {
	mu      sync.RWMutex
	clients map[chan string]struct{}
}

// NewReloadBroadcaster creates a new broadcaster for live reload events.
func NewReloadBroadcaster() *ReloadBroadcaster {
	return &ReloadBroadcaster{
		clients: make(map[chan string]struct{}),
	}
}

// Subscribe adds a new SSE client connection.
func (b *ReloadBroadcaster) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan string, 10)
	b.clients[ch] = struct{}{}
	return ch
}

// Unsubscribe removes an SSE client connection.
func (b *ReloadBroadcaster) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.clients, ch)
	close(ch)
}

// Notify sends a reload event to all connected clients.
func (b *ReloadBroadcaster) Notify(message string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.clients {
		select {
		case ch <- message:
		default:
			// Client is slow or disconnected, skip
		}
	}
}

// ClientCount returns the number of connected clients.
func (b *ReloadBroadcaster) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}

// ServeHTTP handles SSE connections for live reload.
func (b *ReloadBroadcaster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if response writer supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Set SSE headers before writing anything
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering

	// Write status code explicitly
	w.WriteHeader(http.StatusOK)

	// Subscribe to events
	messageChan := b.Subscribe()
	defer b.Unsubscribe(messageChan)

	// Send initial connection message
	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	// Listen for events or client disconnect
	for {
		select {
		case <-r.Context().Done():
			// Client disconnected (normal close)
			return
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			// Send reload event
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
