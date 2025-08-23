package server

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/khabib-developer/chat-application/internal/dto"
)

func (server *Server) listen(connID string, ws *websocket.Conn) {
	sender := server.Users[connID]
	defer func() {
		ws.Close()
		delete(server.Users, connID)

		server.destroyUserChannels(sender)
		
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

			fmt.Println(string(payload.Payload))


			if handler, ok := server.handlers[payload.MessageType]; ok {
				handler(payload.Payload, sender)
			} else {
				server.sendRawMessage(ws, dto.MessageTypeError, "unknown command")
			}
			
		}
	}
}

