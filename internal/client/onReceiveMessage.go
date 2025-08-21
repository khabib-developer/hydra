package client

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)


func onReceiveMessage(_ *user.User, payload json.RawMessage, _ chan string) error {
	var messageDto dto.ReceiveMessageDto
	if err := json.Unmarshal(payload, &messageDto); err != nil {
		log.Printf("⚠ Payload error: %v", err)
		return err
	}

	// Nice formatted output
	fmt.Printf("%s[%s]%s %s%s%s: %s\n",
		cyan, time.Now().Format("15:04"), reset,
		green, messageDto.Sender, reset,
		messageDto.Message,
	)
	return nil
}

func onReceiveMessageFromChannel(_ *user.User, payload json.RawMessage, _ chan string) error {
	var channelMessageDto dto.ChannelMessageDto
	if err := json.Unmarshal(payload, &channelMessageDto); err != nil {
		log.Printf("⚠ Payload error: %v", err)
		return err
	}

	// Example: [15:04] #general @Alice: Hello world
	fmt.Printf("%s[%s]%s %s#%s%s %s@%s%s: %s\n",
		cyan, time.Now().Format("15:04"), reset,      // timestamp
		magenta, channelMessageDto.Channel, reset,   // channel name
		green, channelMessageDto.Sender, reset,      // sender
		channelMessageDto.Message,                   // message
	)

	return nil
}