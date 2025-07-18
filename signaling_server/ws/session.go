package ws

import "sync"

// Session holds both sides of a remote control session
type Session struct {
	Mobile  *Client
	Desktop *Client
}

var (
	sessions = make(map[string]*Session)
	mu       sync.Mutex
)

// RegisterClient adds a client to its session by role
func RegisterClient(client *Client) (paired bool) {
	mu.Lock()
	defer mu.Unlock()
	session, exists := sessions[client.ID]
	if !exists {
		session = &Session{}
		sessions[client.ID] = session
	}

	switch client.Role {
	case RoleDesktop:
		session.Desktop = client
	case RoleMobile:
		session.Mobile = client
	}

	// Check if both sides are connected
	if session.Desktop != nil && session.Mobile != nil {
		session.Desktop.IsPaired = true
		session.Mobile.IsPaired = true
		return true
	}

	return false
}

func RemoveClient(client *Client) {
	mu.Lock()
	defer mu.Unlock()

	session, exists := sessions[client.ID]
	if !exists {
		return
	}

	switch client.Role {
	case RoleDesktop:
		session.Desktop = nil
	case RoleMobile:
		session.Mobile = nil
	}

	// Remove entire session if both sides are gone
	if session.Desktop == nil && session.Mobile == nil {
		delete(sessions, client.ID)
	}
}
