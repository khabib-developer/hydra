package client

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/khabib-developer/chat-application/internal/user"
)

var wsurl string

func init() {
	_ = godotenv.Load()
	wsurl = os.Getenv("WS_URL")
	if wsurl == "" {
		wsurl = "ws://13.48.44.95:8080"
	}
}

func Connect(u *user.User) {

	headers := http.Header{}
	headers.Set("connID", u.ID)

	// Dial the websocket server
	conn, resp, err := websocket.DefaultDialer.Dial(wsurl+"/connect", headers)
	if err != nil {
		if resp != nil {
			log.Printf("Handshake failed with status %d", resp.StatusCode)
		}
		log.Fatalf("Failed to connect: %v", err)
	}

	log.Println("Connected to server with connID:", u.ID)
	safeConn := user.SafeConn{
		Conn: conn,
		Mutex: sync.Mutex{},
	}
	u.SafeConn = &safeConn
}