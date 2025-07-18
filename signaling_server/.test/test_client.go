package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"time"
)

// IdentifyMessage represents the initial message sent to the signaling server
type IdentifyMessage struct {
	Type      string `json:"type"`
	Role      string `json:"role"`       // "mobile" or "desktop"
	SessionID string `json:"session_id"` // UUID shared between both sides
}

func main() {
	// Connect to the signaling server
	url := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to", url)

	// Construct the identify message
	role := "desktop"               // Change to "mobile" to simulate the mobile side
	sessionID := "test-session-123" // Use same ID to simulate pairing

	msg := IdentifyMessage{
		Type:      "identify",
		Role:      role,
		SessionID: sessionID,
	}

	// Encode to JSON and send
	data, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("JSON encode error:", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Fatal("Failed to send identify message:", err)
	}

	fmt.Println("Identify message sent")

	// Wait for incoming messages (optional, or can exit)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Connection closed:", err)
				os.Exit(0)
			}
			fmt.Println("Received:", string(message))
		}
	}()

	// Keep the client running
	for {
		time.Sleep(5 * time.Second)
	}
}
