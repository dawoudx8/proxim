package session

import (
	"server/types"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	SessionID   string
	UUID        string
	Status      string
	DeviceInfo  types.DeviceInfo
	DesktopConn *websocket.Conn
	MobileConn  *websocket.Conn
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

var (
	sessions = make(map[string]Session)
	mu       sync.RWMutex
)

// Create creates a new session with sessionID key
func Create(sessionID string, session Session) {
	mu.Lock()
	defer mu.Unlock()
	sessions[sessionID] = session
}

// Get returns session by ID and a bool indicating if it exists
func Get(sessionID string) (Session, bool) {
	mu.RLock()
	defer mu.RUnlock()
	s, ok := sessions[sessionID]
	return s, ok
}

// Update updates a session by ID
func Update(sessionID string, updated Session) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := sessions[sessionID]
	if !exists {
		return false
	}
	sessions[sessionID] = updated
	return true
}

// Delete removes a session
func Delete(sessionID string) {
	mu.Lock()
	defer mu.Unlock()
	delete(sessions, sessionID)
}

// Cleanup clears expired sessions
func Cleanup() {
	mu.Lock()
	defer mu.Unlock()
	now := time.Now()
	for id, s := range sessions {
		if !s.ExpiresAt.IsZero() && s.ExpiresAt.Before(now) {
			delete(sessions, id)
		}
	}
}
