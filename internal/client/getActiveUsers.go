package client

import (
	"encoding/json"
	"fmt"

	"github.com/khabib-developer/chat-application/internal/user"
)


func GetActiveUsers(u *user.User) {
	headers := map[string]string{
		"connID": u.ID,
	}
	respBody, err := httpSender("GET", "/getActiveUsers", nil, headers)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Decode JSON
	var users []user.UserDTO
	if err := json.Unmarshal(respBody, &users); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	// Print users beautifully
	fmt.Println("=== Active Users ===")
	for _, usr := range users {
		if usr.Username == u.Username {
			continue
		}
		if usr.Private {
			// Red for private
			fmt.Printf("\033[31m%s (private)\033[0m\n", usr.Username)
		} else {
			// Green for public
			fmt.Printf("\033[32m%s (public)\033[0m\n", usr.Username)
		}
	}
}
