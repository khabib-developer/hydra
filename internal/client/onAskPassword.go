package client

import (
	"encoding/json"

	"github.com/khabib-developer/chat-application/internal/user"
)

func onAskPassword(u *user.User, _ json.RawMessage, state chan string) error {
	state <- StatePassword
    return nil

}