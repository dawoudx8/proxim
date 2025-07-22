package main

import (
	"github.com/dawoudx8/minilog"
	"net/http"
	"server/ws"
)

func main() {

	log := minilog.New(true)

	// WebSocket endpoint for desktop and mobile clients
	http.HandleFunc("/ws", ws.HandleWebSocket)
	log.Info("Server started on port 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errorf("ListenAndServe error: %s", err)
	}
}
