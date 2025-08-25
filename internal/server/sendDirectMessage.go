package server

import (
	"encoding/json"
	"strings"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)


func(server *Server) sendDirectMessage(payload json.RawMessage, sender *user.User) {
	var messagePayload dto.SendMessageDto

	if err := json.Unmarshal(payload, &messagePayload); err != nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "invalid message payload")
		return
	}

	var receiver *user.User

	for _, user := range server.Users {
		if messagePayload.Receiver == user.Username {
			receiver = user
		}
	}

	if receiver == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Username not found")
		return
	}

	if len(receiver.Password) != 0 {
		permanentData := user.PermanentData{
			Expect:   receiver.Password,
			Data:     messagePayload.Message,
			Receiver: receiver,
		}
		sender.PermanentData = &permanentData
		server.sendRawMessage(sender.SafeConn, dto.MessageTypePassword, "Password of user: ")
		return
	}

	server.sendActualMessage(sender, receiver, messagePayload.Message)
}


func(server *Server) handlePassword(payload json.RawMessage, sender *user.User) {
	var password string
    if err := json.Unmarshal(payload, &password); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Unsupperted type of password")
		return
    }
	
	if sender.PermanentData == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "User did not expect password")
		return
	}

	if sender.PermanentData.Expect != password {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Wrong password")
		return
	}

	server.sendActualMessage(sender, sender.PermanentData.Receiver, sender.PermanentData.Data)

	sender.PermanentData = nil
}


func(server *Server) sendRawMessage(safeConn *user.SafeConn, messageType dto.MessageType, message string) {
	messageJson, err := json.Marshal(strings.TrimSpace(message))
	if err != nil {
		return
	}
	server.sendMessage(safeConn, messageType,  messageJson)

}

func(server *Server) sendActualMessage(sender *user.User, receiver *user.User, message string) {
	responsePayloadBytes, error := json.Marshal(dto.ReceiveMessageDto{
		Sender: sender.Username,
		Message: message,
	})

	if error != nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "invalid message payload")
		return
	}

	server.sendMessage(receiver.SafeConn, dto.MessageTypeMessage,  responsePayloadBytes)

	server.sendRawMessage(sender.SafeConn, dto.MessageTypeInfo, "Your message successfully sent")
}