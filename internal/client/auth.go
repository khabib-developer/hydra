package client

import (
	"fmt"

	"github.com/khabib-developer/chat-application/internal/user"
)


func Auth(u *user.User) error {
	if u == nil {
		return fmt.Errorf("nil user pointer")
	}

	payload := user.AuthDTO{
		Username: u.Username,
		Password: u.Password,
	}

	respBody, err := httpSender("POST", "/auth", &payload, nil)

	if err != nil {
		return err
	}

	u.ID = string(respBody)

	fmt.Println("Response:", string(respBody))
	return nil
}
