package ws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"server/session"
	"server/types"
)

func handlePairRequest(raw map[string]interface{}, role string, mobileConn, desktopConn *websocket.Conn) {
	if role != "mobile" {
		log.Errorf("[%s] Only mobile clients can send pair_request", role)
		return
	}

	var msg types.PairRequest
	if err := mapToStruct(raw, &msg); err != nil {
		log.Errorf("Failed to decode pair_request: %v", err)
		return
	}

	if msg.SessionID == "" {
		log.Error("Missing sessionID in pair_request")
		return
	}

	log.Infof("Pairing request received from device: %s (%s)", msg.DeviceInfo.Name, msg.DeviceInfo.ID)

	if desktopConn == nil {
		log.Error("No desktop connected — cannot forward pair_request")
		return
	}

	session.Create(msg.SessionID, session.Session{
		SessionID:  msg.SessionID,
		MobileConn: mobileConn,
		DeviceInfo: msg.DeviceInfo,
	})

	prompt := types.PairPrompt{
		Type:       "pair_prompt",
		SessionID:  msg.SessionID,
		DeviceInfo: msg.DeviceInfo,
	}

	if err := desktopConn.WriteJSON(prompt); err != nil {
		log.Errorf("Failed to send pair_prompt to desktop: %v", err)
	}
}

func handlePairResponse(raw map[string]interface{}, role string) {
	if role != "desktop" {
		log.Errorf("[%s] Only desktop clients can send pair_response", role)
		return
	}

	var msg types.PairResponse
	if err := mapToStruct(raw, &msg); err != nil {
		log.Errorf("Failed to decode pair_response: %v", err)
		return
	}

	if msg.SessionID == "" {
		log.Error("Missing sessionID in pair_response")
		return
	}

	sess, exists := session.Get(msg.SessionID)
	if !exists || sess.MobileConn == nil {
		log.Error("Session not found or no mobile connection to respond to")
		return
	}

	if msg.Approved {
		log.Info("Desktop approved the connection request ✅")
	} else {
		log.Info("Desktop rejected the connection request ❌")
	}

	result := types.PairResult{
		Type:      "pair_result",
		SessionID: msg.SessionID,
		Status:    "rejected",
	}
	if msg.Approved {
		result.Status = "approved"
	}

	if err := sess.MobileConn.WriteJSON(result); err != nil {
		log.Errorf("Failed to send pair_result to mobile: %v", err)
	}
}

func handleSignalMessage(raw map[string]interface{}, senderConn *websocket.Conn, senderRole string) {
	var msg types.SignalMessage
	if err := mapToStruct(raw, &msg); err != nil {
		log.Errorf("Failed to decode signal message: %v", err)
		return
	}

	if msg.SessionID == "" {
		log.Error("Missing sessionID in signal message")
		return
	}

	sess, exists := session.Get(msg.SessionID)
	if !exists {
		log.Errorf("[%s] Session not found for signaling: %s", senderRole, msg.SessionID)
		return
	}

	var targetConn *websocket.Conn
	if senderRole == "mobile" {
		targetConn = sess.DesktopConn
	} else {
		targetConn = sess.MobileConn
	}

	if targetConn == nil {
		log.Errorf("[%s] Target connection not found for signaling", senderRole)
		return
	}

	if err := targetConn.WriteJSON(msg); err != nil {
		log.Errorf("[%s] Failed to forward signaling message: %v", senderRole, err)
	}
}
func mapToStruct(m map[string]interface{}, out interface{}) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}
