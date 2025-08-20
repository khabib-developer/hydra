package client

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/khabib-developer/chat-application/internal/dto"
)

func send(ws *websocket.Conn, messageType dto.MessageType, payload json.RawMessage)  {
	data := dto.WebsocketDto{
		MessageType: messageType,
		Payload:     payload,
	}

	message, err := json.Marshal(data)


	if err != nil {
		fmt.Println("marshal error:", err)
		ws.Close()
	}

	if err := ws.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		fmt.Println("write error:", err)
		ws.Close()
	}
}