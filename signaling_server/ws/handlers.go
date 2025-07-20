package ws

import (
	"github.com/gorilla/websocket"
	"server/types"
)

func handleInitMessage(ws *websocket.Conn) (string, error) {

	var initMsg types.InitMessage

	err := ws.ReadJSON(&initMsg)
	if err != nil {
		return "", err
	}

	if initMsg.Role == "desktop" {
		desktopConn = ws
	} else if initMsg.Role == "mobile" {
		log.Info("Waiting for pair_request from mobile...")
	}

	return initMsg.Role, nil
}

func handleMobileInit(ws *websocket.Conn) {
	var pairMsg types.PairRequest

	err := ws.ReadJSON(&pairMsg)
	if err != nil {
		log.Errorf("Failed to read pair_request: %v", err)
		return
	}

	if pairMsg.Type != "pair_request" {
		log.Errorf("Expected pair_request, got: %s", pairMsg.Type)
		return
	}

	log.Infof("Pairing request received from device: %s (%s)", pairMsg.DeviceInfo.Name, pairMsg.DeviceInfo.ID)

	if desktopConn == nil {
		log.Error("No desktop connected to forward pair request")
		return
	}

	prompt := types.PairPrompt{
		Type:       "pair_prompt",
		SessionID:  pairMsg.SessionID,
		DeviceInfo: pairMsg.DeviceInfo,
	}

	err = desktopConn.WriteJSON(prompt)
	if err != nil {
		log.Errorf("Failed to send pair_prompt to desktop: %v", err)
		return
	}

	var response types.PairResponse
	err = desktopConn.ReadJSON(&response)
	if err != nil {
		log.Errorf("Failed to read pair_response from desktop: %v", err)
		return
	}

	if response.Type != "pair_response" {
		log.Errorf("Expected pair_response, got: %s", response.Type)
		return
	}

	if response.Approved {
		log.Info("Desktop approved the connection request ✅")
	} else {
		log.Info("Desktop rejected the connection request ❌")
	}

	result := types.PairResult{
		Type:      "pair_result",
		SessionID: response.SessionID,
	}

	if response.Approved {
		result.Status = "approved"
	} else {
		result.Status = "rejected"
	}

	err = ws.WriteJSON(result)
	if err != nil {
		log.Errorf("Failed to send pair_result to mobile: %v", err)
		return
	}
}
