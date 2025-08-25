package server

import (
	"encoding/json"
	"fmt"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)

func (server *Server) createChannel(payload json.RawMessage, sender *user.User) {
	var channelName string
    if err := json.Unmarshal(payload, &channelName); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Unsupported type of channel name")
		return
    }

	for _, ch := range server.Channels {
		if ch.Name == channelName {
			server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Channel already exists")
			return
		}
	}

	if sender.JoinedChannel != nil {
		for i, u := range sender.JoinedChannel.Users {
			if u == sender {
				sender.JoinedChannel.Users = append(sender.JoinedChannel.Users[:i], sender.JoinedChannel.Users[i+1:]...)
				break
			}
		}
		sender.JoinedChannel = nil
	}

	channel := &user.Channel{
		Name:   channelName,
		Owner:  sender,
		Users:  []*user.User{sender},
	}

	sender.JoinedChannel = channel

	server.Channels = append(server.Channels, channel)
		
	server.sendRawMessage(sender.SafeConn, dto.MessageTypeInfo, "Channel created successfully")
}

func (server *Server) joinChannel(payload json.RawMessage, sender *user.User) {
	var channelName string
    if err := json.Unmarshal(payload, &channelName); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Unsupperted type of channel name")
		return
    }

	var channelToJoin *user.Channel
	for _, ch := range server.Channels {
		if ch.Name == channelName {
			channelToJoin = ch
			break
		}
	}

	if channelToJoin == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Channel not found")
		return
	}

	if sender.JoinedChannel != nil {
		for i, u := range sender.JoinedChannel.Users {
			if u == sender {
				fmt.Println("Removing user from previous channel")
				sender.JoinedChannel.Users = append(sender.JoinedChannel.Users[:i], sender.JoinedChannel.Users[i+1:]...)
				break
			}
		}
		sender.JoinedChannel = nil
	}

	sender.JoinedChannel = channelToJoin
	channelToJoin.Users = append(channelToJoin.Users, sender)
	server.sendRawMessage(sender.SafeConn, dto.MessageTypeInfo, "Successfully joined to channel") 
}

func (server *Server) broadcastMessage(payload json.RawMessage, sender *user.User) {
	var message string
	if err := json.Unmarshal(payload, &message); err != nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Invalid message payload")
		return
	}

	if sender.JoinedChannel == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "You are not in a channel")
		return
	}

	channelMessage := dto.ChannelMessageDto{
		Channel: sender.JoinedChannel.Name,
		Sender:  sender.Username,
		Message: message,
	}

	jsonMessage, err := json.Marshal(channelMessage)
	if err != nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Failed to marshal channel message")
		return
	}

	for _, user := range sender.JoinedChannel.Users {
		server.sendMessage(user.SafeConn, dto.MessageTypeBroadcast, jsonMessage)
	}
}

func (server *Server) destroyChannel(payload json.RawMessage, sender *user.User) {
	var channelName string

	if err := json.Unmarshal(payload, &channelName); err != nil {
		println(err)
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Unsupported type of channel name")
		return
	}

	var channelToDestroy *user.Channel
	var index int = -1

	for i, ch := range server.Channels {
		if channelName == ch.Name {
			channelToDestroy = ch
			index = i
			break
		}
	}

	if channelToDestroy == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Channel not found")
		return
	}

	if channelToDestroy.Owner != sender {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "You are not the owner of this channel")
		return
	}

	for _, user := range channelToDestroy.Users {
		user.JoinedChannel = nil
		server.sendRawMessage(user.SafeConn, dto.MessageTypeInfo, fmt.Sprintf("Channel '%s' has been destroyed", channelToDestroy.Name))
	}

	if index != -1 {
		server.Channels = append(server.Channels[:index], server.Channels[index+1:]...)
	}
}

func(server *Server) destroyUserChannels(u *user.User) {
	for _, channel := range server.Channels {
		if channel.Owner == u {
			
			channelNameJson, err := json.Marshal(channel.Name)

			if err != nil {
				return
			}

			server.destroyChannel(channelNameJson, u)

		}
	}
}

	
