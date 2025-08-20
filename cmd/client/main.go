package main

import (
	"fmt"

	"github.com/khabib-developer/chat-application/internal/client"
	"github.com/khabib-developer/chat-application/internal/user"
)

func main() {
	client.Draw()
	user := user.User{}

	state := make(chan string)
	

	if err := client.Run(&user); err != nil {

	}

	if err := client.Auth(&user); err != nil {
		fmt.Println("err", err)
	}

	client.Connect(&user)

	go client.Listen(&user, state)

	client.ReceiveCommands(&user, state)
}
