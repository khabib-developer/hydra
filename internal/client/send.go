package client

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)

func send(u *user.User, messageType dto.MessageType, payload json.RawMessage)  {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	data := dto.WebsocketDto{
		MessageType: messageType,
		Payload:     payload,
	}

	message, err := json.Marshal(data)


	if err != nil {
		fmt.Println("marshal error:", err)
		u.Conn.Close()
	}

	if err := u.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		fmt.Println("write error:", err)
		u.Conn.Close()
	}
}