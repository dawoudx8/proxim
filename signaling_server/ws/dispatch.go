package ws

import (
	"github.com/gorilla/websocket"
)

// DispatchMessage routes incoming messages to the appropriate handler.
func DispatchMessage(role string, msg map[string]interface{}, conn *websocket.Conn, desktopConn *websocket.Conn, mobileConn *websocket.Conn) {

	msgTypeRaw, ok := msg["type"]
	if !ok {
		log.Errorf("[%s] Message missing 'type' field", role)
		return
	}

	msgType, ok := msgTypeRaw.(string)
	if !ok {
		log.Errorf("[%s] Invalid 'type' field (not a string)", role)
		return
	}

	switch msgType {
	case "pair_request":
		handlePairRequest(msg, role, mobileConn, desktopConn)

	case "pair_response":
		handlePairResponse(msg, role)

	case "signal":
		handleSignalMessage(msg, conn, role)

	default:
		log.Errorf("[%s] Unknown message type: %s", role, msgType)
	}
}
