package main

import (
	"net/http"
	"server/utils"
	"server/ws"
)

func main() {

	log := utils.New(true)

	// Set up WebSocket handler
	http.HandleFunc("/ws", ws.HandleWebSocket) // Upgrade connection inside handler
	log.Success("Signaling server listening on: 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Errorf("ListenAndServe error: %v", err)
	}

}
