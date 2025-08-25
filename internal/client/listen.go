package client

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)

const (
	green   = "\033[32m"
	cyan    = "\033[36m"
	magenta = "\033[35m"
	reset   = "\033[0m"
)


// Listen starts listening for messages from the server
func Listen(u *user.User, state chan string) {
	if u.SafeConn == nil {
		log.Println("❌ No active websocket connection")
		return
	}

	for {
		messageType, msg, err := u.SafeConn.Conn.ReadMessage()
		if err != nil {
			log.Printf("🔌 Connection closed: %v", err)
			_ = u.SafeConn.Conn.Close()
			os.Exit(1) // Or notify caller to reconnect
		}

		if messageType != websocket.TextMessage {
			continue // ignore binary/ping
		}

		var websocketDto dto.WebsocketDto
		if err := json.Unmarshal(msg, &websocketDto); err != nil {
			log.Printf("⚠ Invalid JSON: %v", err)
		}

		// Dispatch to the right handler
		if handler, ok := commandHandlers[websocketDto.MessageType]; ok {
			if err := handler(u, websocketDto.Payload, state); err != nil {
				log.Printf("⚠ Handler error: %v", err)
			}
		} else {
			log.Printf("⚠ Unknown command: %s", websocketDto.MessageType)
		}
		
	}
}

var commandHandlers = map[dto.MessageType]func(*user.User, json.RawMessage, chan string) error {
	dto.MessageTypeMessage:   onReceiveMessage,
	dto.MessageTypePassword:  onAskPassword,
	dto.MessageTypeInfo:      onReceiveInfo,
	dto.MessageTypeError:     onReceiveError,
	dto.MessageTypeBroadcast: onReceiveMessageFromChannel,
	dto.MessageTypeFile:      onReceiveFileMetadata,
	dto.MessageTypeFileChunk: onReceiveFileChunk,
	dto.MessageTypeCancel:    onReceiveCancelTransfer,
}
