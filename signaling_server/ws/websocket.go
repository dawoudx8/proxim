package ws

import (
	"github.com/dawoudx8/minilog"
	"github.com/gorilla/websocket"
	"net/http"
)

var log = minilog.New(true)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Upgrade error: %s", err)
		return
	}
	defer ws.Close()
}
