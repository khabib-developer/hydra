package client

import (
	"encoding/json"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)


func create(u *user.User, channel string) {
	if u == nil {
		return
	}

	jsonChannel, err := json.Marshal(channel)
	if err != nil {
		return
	}

	send(u, dto.MessageTypeCreate, jsonChannel)
}

func join(u *user.User, channel string) {
	jsonChannel, err := json.Marshal(channel)
	if err != nil {
		return
	}

	send(u, dto.MessageTypeJoin, jsonChannel)
}

func broadcast(u *user.User, message string) {
	if u == nil {
		return
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return
	}

	send(u, dto.MessageTypeBroadcast, jsonMessage)
}

func destroyChannel(u *user.User, channelName string) {
	jsonMessage, err := json.Marshal(channelName)
	if err != nil {
		return
	}

	send(u, dto.MessageTypeDestroy, jsonMessage)
}