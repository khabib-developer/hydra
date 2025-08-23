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
	CmdMembers   Command = "/members"
	CmdCurrent   Command = "/current"
	CmdDestroy   Command = "/destroy"
	CmdInfo      Command = "/profile"


	CmdFile      Command = "/file"
	CmdProgress  Command = "/file_progress"
	CmdCancel    Command = "/file_cancel"
)


const (
	MessageTypeError      MessageType = "error"
	MessageTypePassword   MessageType = "password"
	MessageTypeInfo       MessageType = "info"
	MessageTypeMessage    MessageType = "message"
	MessageTypeBroadcast  MessageType = "broadcast"
	MessageTypeClose      MessageType = "close"
	MessageTypeJoin       MessageType = "join"
	MessageTypeCreate     MessageType = "create"
	MessageTypeDestroy    MessageType = "destroy"
	MessageTypeCancel     MessageType = "cancel"
	MessageTypeFile       MessageType = "file"
	MessageTypeFileChunk  MessageType = "file_chunk"
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
	CmdMembers,
	CmdCurrent,
	CmdDestroy,
	CmdInfo,
	CmdFile,

}

var AllMessageTypes = []MessageType{
	MessageTypeError,
	MessageTypeInfo,
	MessageTypeMessage,
	MessageTypePassword,
	MessageTypeBroadcast,
	MessageTypeClose,
	MessageTypeJoin,
	MessageTypeCreate,
	MessageTypeDestroy,
	MessageTypeFile,
	MessageTypeFileChunk,
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
