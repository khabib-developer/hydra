package user

import (
	"github.com/gorilla/websocket"
)

type User struct {
	ID            string
	Username      string
	Password      string
	Conn          *websocket.Conn
	PermanantData *PermanentData
	JoinedChannel *Channel
}

type PermanentData struct {
	Expect   string
	Data     string
	Receiver *User

}

type UserDTO struct {
	Username string `json:"username"`
	Private  bool   `json:"private"`
}

type AuthDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}