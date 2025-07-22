package ws

import (
	"encoding/json"
	"github.com/dawoudx8/minilog"
	"github.com/gorilla/websocket"
	"net/http"
	"server/types"
)

var log = minilog.New(true)

// Global single-user connection holders
var desktopConn *websocket.Conn
var mobileConn *websocket.Conn

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	clientIP := r.RemoteAddr
	log.Infof("New WebSocket client connected from: %s", clientIP)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Upgrade error: %s", err)
		return
	}
	defer ws.Close()

	// Step 1: Read and parse the init message
	_, msgBytes, err := ws.ReadMessage()
	if err != nil {
		log.Errorf("Failed to read init message: %v", err)
		return
	}

	var initMsg types.InitMessage
	err = json.Unmarshal(msgBytes, &initMsg)
	if err != nil {
		log.Errorf("Invalid init message JSON: %v", err)
		return
	}

	if initMsg.Type != "init" || (initMsg.Role != "mobile" && initMsg.Role != "desktop") {
		log.Errorf("Invalid init message contents: %+v", initMsg)
		return
	}

	log.Infof("Client connected with role: %s", initMsg.Role)

	// Step 2: Store connection based on role
	if initMsg.Role == "desktop" {
		desktopConn = ws
		log.Info("Waiting for pairing request...")
	} else if initMsg.Role == "mobile" {
		mobileConn = ws
	}

	// Step 3: Begin message handling loop
	HandleConnection(initMsg.Role, ws, desktopConn, mobileConn)
}
