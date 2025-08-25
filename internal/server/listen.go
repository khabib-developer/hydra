package server

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/khabib-developer/chat-application/internal/dto"
)

func (server *Server) listen(connID string) {
	sender := server.Users[connID]
	defer func() {
		sender.SafeConn.Conn.Close()
		delete(server.Users, connID)

		server.destroyUserChannels(sender)
		
		fmt.Println("connection closed for", connID)
	}()

	for {
		messageType, msg, err := sender.SafeConn.Conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			break
		}

		if messageType == websocket.TextMessage {

			var payload dto.WebsocketDto

			if err := json.Unmarshal(msg, &payload); err != nil {
				server.sendMessage(sender.SafeConn, dto.MessageTypeError, []byte(`"wrong type of command"`))
				return
			}

			if handler, ok := server.handlers[payload.MessageType]; ok {
				handler(payload.Payload, sender)
			} else {
				server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "unknown command")
			}
			
		}
	}
}

