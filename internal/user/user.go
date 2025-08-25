package user

import (
	"sync"

	"github.com/gorilla/websocket"
)

type User struct {
	ID            string
	Username      string
	Password      string
	PermanentData *PermanentData
	JoinedChannel *Channel
	PermanentFile *PermanentFileData
	SafeConn 	  *SafeConn
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

type SafeConn struct {
	Conn *websocket.Conn
	Mutex sync.Mutex
}