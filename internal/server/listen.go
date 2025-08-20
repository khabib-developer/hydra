package server

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/khabib-developer/chat-application/internal/dto"
)

func (server *Server) listen(connID string, ws *websocket.Conn) {
	defer func() {
		ws.Close()
		delete(server.Users, connID)
		fmt.Println("connection closed for", connID)
	}()

	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}

		if messageType == websocket.TextMessage {


			var payload dto.WebsocketDto

			if err := json.Unmarshal(msg, &payload); err != nil {
				server.sendMessage(ws, dto.MessageTypeError, []byte(`"wrong type of command"`))
				return
			}

			sender := server.Users[connID]

			switch payload.MessageType {
			case dto.MessageTypeMessage:
				server.sendDirectMessage(payload.Payload, sender)
			case dto.MessageTypeJoin:
			case dto.MessageTypeCreate:
			case dto.MessageTypePassword:

				server.handlePassword(payload.Payload, sender)


			}
			
		}
	}
}