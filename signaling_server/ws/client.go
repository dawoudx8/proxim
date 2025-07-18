package ws

import "github.com/gorilla/websocket"

type Role string

const (
	RoleMobile  Role = "mobile"
	RoleDesktop Role = "desktop"
)

// Client represents a connected WebSocket client
type Client struct {
	ID       string
	Conn     *websocket.Conn
	Role     Role
	IsPaired bool
	Send     chan []byte
}
