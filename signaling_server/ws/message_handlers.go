package ws

import (
	"github.com/gorilla/websocket"
	"server/session"
	"server/types"
)

// handlePairRequest processes an incoming pair_request from a mobile client.
func handlePairRequest(msg map[string]interface{}, role string, mobileConn, desktopConn *websocket.Conn) {

	if role != "mobile" {
		log.Errorf("[%s] Only mobile clients can send pair_request", role)
		return
	}

	sessionID, ok := msg["sessionID"].(string)
	if !ok || sessionID == "" {
		log.Error("Missing or invalid 'sessionID' in pair_request")
		return
	}

	deviceInfoRaw, ok := msg["deviceInfo"].(map[string]interface{})
	if !ok {
		log.Error("Missing or invalid 'deviceInfo' in pair_request")
		return
	}

	name, _ := deviceInfoRaw["name"].(string)
	id, _ := deviceInfoRaw["id"].(string)

	log.Infof("Pairing request received from device: %s (%s)", name, id)

	if desktopConn == nil {
		log.Error("No desktop connected — cannot forward pair_request")
		return
	}

	// Store session
	session.Create(sessionID, session.Session{
		SessionID:  sessionID,
		MobileConn: mobileConn,
		DeviceInfo: types.DeviceInfo{
			Name: name,
			ID:   id,
		},
	})

	// Send prompt to desktop
	prompt := map[string]interface{}{
		"type":       "pair_prompt",
		"sessionID":  sessionID,
		"deviceInfo": deviceInfoRaw,
	}

	err := desktopConn.WriteJSON(prompt)
	if err != nil {
		log.Errorf("Failed to send pair_prompt to desktop: %v", err)
	}
}

// handlePairResponse processes the desktop's approval or rejection of the pair_request.
func handlePairResponse(msg map[string]interface{}, role string) {
	if role != "desktop" {
		log.Errorf("[%s] Only desktop clients can send pair_response", role)
		return
	}

	sessionID, ok := msg["sessionID"].(string)
	if !ok || sessionID == "" {
		log.Error("Missing or invalid 'sessionID' in pair_response")
		return
	}

	approved, ok := msg["approved"].(bool)
	if !ok {
		log.Error("Missing or invalid 'approved' field in pair_response")
		return
	}

	sess, exists := session.Get(sessionID)
	if !exists || sess.MobileConn == nil {
		log.Error("Session not found or no mobile connection to respond to")
		return
	}

	if approved {
		log.Info("Desktop approved the connection request ✅")
	} else {
		log.Info("Desktop rejected the connection request ❌")
	}

	result := map[string]interface{}{
		"type":      "pair_result",
		"sessionID": sessionID,
		"status":    "rejected",
	}
	if approved {
		result["status"] = "approved"
	}

	err := sess.MobileConn.WriteJSON(result)
	if err != nil {
		log.Errorf("Failed to send pair_result to mobile: %v", err)
	}
}

// handleSignalMessage forwards a WebRTC signaling message to the other peer.
func handleSignalMessage(msg map[string]interface{}, senderConn *websocket.Conn, senderRole string) {
	sessionID, ok := msg["sessionID"].(string)
	if !ok || sessionID == "" {
		log.Error("Missing or invalid 'sessionID' in signal message")
		return
	}

	sess, exists := session.Get(sessionID)
	if !exists {
		log.Errorf("[%s] Session not found for signaling: %s", senderRole, sessionID)
		return
	}

	var targetConn *websocket.Conn
	if senderRole == "mobile" {
		targetConn = sess.DesktopConn
	} else if senderRole == "desktop" {
		targetConn = sess.MobileConn
	}

	if targetConn == nil {
		log.Errorf("[%s] Target connection not found for signaling", senderRole)
		return
	}

	err := targetConn.WriteJSON(msg)
	if err != nil {
		log.Errorf("[%s] Failed to forward signaling message: %v", senderRole, err)
	}
}
