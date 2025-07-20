package ws

import (
	"github.com/dawoudx8/minilog"
	"github.com/gorilla/websocket"
	"net/http"
)

var log = minilog.New(true)

var desktopConn *websocket.Conn

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

	role, err := handleInitMessage(ws)
	if err != nil {
		log.Errorf("Failed to read init message: %v", err)
		return
	}

	log.Infof("Client connected with role: %s", role)

	if role == "mobile" {
		handleMobileInit(ws)
	}

}
