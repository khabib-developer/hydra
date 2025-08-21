package client

import (
	"encoding/json"
	"fmt"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)


func getActiveUsers(u *user.User) {
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

func getChannels(u *user.User) {
	headers := map[string]string{
		"connID": u.ID,
	}
	respBody, err := httpSender("GET", "/getChannels", nil, headers)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var channels []dto.ChannelDto

	if err := json.Unmarshal(respBody, &channels); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	fmt.Println("=== Channels ===")

	for _, channel := range channels {
		if channel.Mine {
			fmt.Printf("\033[34m%s (%v)\033[0m\n", channel.Name, channel.Members)
		} else {
			fmt.Printf("\033[32m%s (%v)\033[0m\n", channel.Name, channel.Members)
		}
	}
}

func getCurrentChannel(u *user.User) {
	headers := map[string]string{
		"connID": u.ID,
	}
	respBody, err := httpSender("GET", "/getCurrentChannel", nil, headers)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var channel string
	if err := json.Unmarshal(respBody, &channel); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	fmt.Println("=== Current Channel ===")
	fmt.Printf("Channel name: %s\n", channel)
}

func getChannelMembers(u *user.User) {
	headers := map[string]string{
		"connID": u.ID,
	}
	respBody, err := httpSender("GET", "/getChannelMembers", nil, headers)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var members []user.UserDTO
	if err := json.Unmarshal(respBody, &members); err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	fmt.Println("=== Channel Members ===")
	for _, member := range members {
		fmt.Printf(" - %s\n", member.Username)
	}
}

func getProfileInfo(u *user.User) {
	fmt.Printf("=== Profile Information ===\n")
	fmt.Printf("Username: %s\n", u.Username)
	fmt.Printf("Password: %s\n", u.Password)
}