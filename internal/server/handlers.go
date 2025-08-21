package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)

type AuthRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func(server *Server) Connect(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }

	connID := r.Header.Get("connID");

    if  connID == "" {
		server.closeConnection(ws, "connID is not exist")
	}

	err = server.add(connID, ws)

	if err != nil {
		server.closeConnection(ws, err.Error())
	}
	
	go server.listen(connID, ws)
}

// get Iformation

func(server *Server) GetActiveUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]user.UserDTO, 0, len(server.Users))
	for _, userItem := range server.Users {
		users = append(users, user.UserDTO{
			Username: userItem.Username,
			Private:  len(userItem.Password) > 0,
		})
	}

	msg, err := json.Marshal(users)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
    w.Write([]byte(msg))
}

func(server *Server) GetChannels(w http.ResponseWriter, r *http.Request) {
	channels := make([]dto.ChannelDto, len(server.Channels))
	i := 0
	connID := r.Header.Get("connID")
	user := server.Users[connID]


	for _, channel := range server.Channels {
		channels[i] = dto.ChannelDto{
			Name:    channel.Name,
			Members: len(channel.Users),
			Mine:    user == channel.Owner,
		}
		i++
	}
	msg, err := json.Marshal(channels)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func(server *Server) GetCurrentChannel(w http.ResponseWriter, r *http.Request) {
	connID := r.Header.Get("connID")
	user := server.Users[connID]	
	
	if user.JoinedChannel == nil {
		http.Error(w, "User is not in a channel", http.StatusBadRequest)
		return
	}

	channelName, err := json.Marshal(user.JoinedChannel.Name)
	if err != nil {
		http.Error(w, "Failed to marshal channel name", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(channelName)
}

func(server *Server) GetChannelMembers(w http.ResponseWriter, r *http.Request) {
	users := make([]user.UserDTO, 0, len(server.Users))
	for _, userItem := range server.Users {
		users = append(users, user.UserDTO{
			Username: userItem.Username,
			Private:  len(userItem.Password) > 0,
		})
	}

	msg, err := json.Marshal(users)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
    w.Write([]byte(msg))
	

}

func(server *Server) sendMessage(ws *websocket.Conn, messageType dto.MessageType, payload json.RawMessage ) {
	data := dto.WebsocketDto{
		MessageType: messageType,
		Payload:     payload,
	}
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal error:", err)
		ws.Close()
		return
	}

	if err = ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		fmt.Println("write error:", err)
		ws.Close()
	}
}

func(server *Server) closeConnection(ws *websocket.Conn, msg string) {
	server.sendMessage(ws, dto.MessageTypeClose, []byte(msg))
	ws.Close()
}


func(server *Server) add(connID string, ws *websocket.Conn) error {
	if server.Users == nil {
		server.Users = make(map[string]*user.User)
	}

	u, ok := server.Users[connID]
	if !ok {
		return fmt.Errorf("user with connID %q not found", connID)
	}

	// attach connection to the copy and save it back
	u.Conn = ws
	server.Users[connID] = u

	return nil
}


func(server *Server) usernameExists(username string) bool {
	for _, user := range server.Users {
        if strings.EqualFold(user.Username, username) {
            return true
        }
    }
    return false

}



func(server *Server) Auth(w http.ResponseWriter, r *http.Request) {
	req := AuthRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

	if candidate := server.usernameExists(req.Username); candidate {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	connID := uuid.New().String()

	server.Users[connID] = &user.User{
		ID: connID,
		Username: req.Username,
		Password: req.Password,
	}


    w.WriteHeader(http.StatusOK)
    w.Write([]byte(connID))
}

