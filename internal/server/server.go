package server

import (
	"github.com/khabib-developer/chat-application/internal/channel"
	"github.com/khabib-developer/chat-application/internal/user"
)

func NewServer() *Server {
	return &Server{
		Users: make(map[string]*user.User),
	}
}

type Server struct {
	Users map[string]*user.User
	Channels []channel.Channel
}






