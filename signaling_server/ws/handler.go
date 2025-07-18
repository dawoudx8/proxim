package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"server/types"
	"server/utils"
)

// Accept all origins temporarily during development.
// TODO: Restrict origins before production deployment.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var log = utils.New(true)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Log connection attempt
	log.Info("New WebSocket connection attempt from: " + r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("WebSocket upgrade failed: %v", err)
		return
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	// Read the first message from client (expected to be an "identify" message)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Errorf("Failed to read initial message: %v", err)
		return
	}

	// Parse the JSON message
	var identify types.IdentifyMessage
	if err := json.Unmarshal(msg, &identify); err != nil {
		log.Errorf("Invalid identify message format: %v", err)
		return
	}

	// Validate message type
	if identify.Type != "identify" || (identify.Role != "mobile" && identify.Role != "desktop") || identify.SessionID == "" {
		log.Error("Missing or invalid fields in identify message")
		return
	}

	// Create new Client object
	client := &Client{
		ID:   identify.SessionID,
		Conn: conn,
		Role: Role(identify.Role),
		Send: make(chan []byte),
	}

	// Register client in session
	if paired := RegisterClient(client); paired {
		log.Infof("Session %s is now paired (mobile + desktop connected)", client.ID)
	} else {
		log.Infof("Client connected as %s for session %s", client.Role, client.ID)
	}

	// Block until connection is closed
	select {}
}
