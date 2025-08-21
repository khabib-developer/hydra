package server

import (
	"encoding/json"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)

func NewServer() *Server {
	s := &Server{
		Users: make(map[string]*user.User),
		
	}
	s.handlers = map[dto.MessageType]func( json.RawMessage, *user.User){
			dto.MessageTypeMessage:    s.sendDirectMessage,
			dto.MessageTypePassword:   s.handlePassword,
			dto.MessageTypeJoin:       s.joinChannel,
			dto.MessageTypeCreate:     s.createChannel,
			dto.MessageTypeBroadcast:  s.broadcastMessage,
			dto.MessageTypeDestroy:    s.destroyChannel,
		}
	return s
}

type Server struct {
	Users map[string]*user.User
	Channels []*user.Channel
	handlers map[dto.MessageType]func(json.RawMessage, *user.User)
}






