package dto

import (
	"encoding/json"
	"fmt"
	"slices"
)

type Command     string
type MessageType string

const (
	CmdHelp      Command = "/help"
	CmdCreate    Command = "/create"
	CmdBroadcast Command = "/broadcast"
	CmdJoin      Command = "/join"
	CmdMessage   Command = "/msg"
	CmdExit      Command = "/exit"
	CmdUsers     Command = "/users"
	CmdChannels  Command = "/channels"
)


const (
	MessageTypeError    MessageType = "error"
	MessageTypePassword MessageType = "password"
	MessageTypeInfo     MessageType = "info"
	MessageTypeMessage  MessageType = "message"
	MessageTypeClose    MessageType = "close"
	MessageTypeJoin     MessageType = "join"
	MessageTypeCreate   MessageType = "create"
	MessageTypeDestroy  MessageType = "destroy"
	MessageTypeFile     MessageType = "file"
)

var AllCommands = []Command{
	CmdHelp,
	CmdCreate,
	CmdBroadcast,
	CmdJoin,
	CmdMessage,
	CmdExit,
	CmdUsers,
	CmdChannels,
}

var AllMessageTypes = []MessageType{
	MessageTypeError,
	MessageTypeInfo,
	MessageTypeMessage,
	MessageTypePassword,
}

type WebsocketDto struct {
	MessageType MessageType        `json:"command"`
	Payload json.RawMessage `json:"payload"` 
}

func (c MessageType) IsValid() bool {
	return slices.Contains(AllMessageTypes, c)
}


func (w WebsocketDto) Validate() error {
	if !w.MessageType.IsValid() {
		return fmt.Errorf("invalid command: %s", w.MessageType)
	}
	return nil
}
