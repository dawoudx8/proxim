package ws

import (
	"github.com/gorilla/websocket"
)

func HandleConnection(role string, ws *websocket.Conn, desktopConn *websocket.Conn, mobileConn *websocket.Conn) {

	for {
		var raw map[string]interface{}

		err := ws.ReadJSON(&raw)
		if err != nil {
			log.Errorf("[%s] Failed to read message: %v", role, err)
			return
		}

		msgType, ok := raw["type"].(string)
		if !ok {
			log.Errorf("[%s] Message missing type field", role)
			continue
		}

		log.Infof("[%s] Received message of type: %s", role, msgType)

		DispatchMessage(role, raw, ws, desktopConn, mobileConn)
	}
}
