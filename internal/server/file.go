package server

import (
	"encoding/json"
	"fmt"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
)



func(server *Server) TransferFileMetadata(payload json.RawMessage, sender *user.User) {
	var fileMetadata user.FileMetadata
	
	if err := json.Unmarshal(payload, &fileMetadata); err != nil {
		// Handle error
		println(err)
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Unsupported type of file metadata")
		return
	}

	var receiver *user.User

	for _, user := range server.Users {
		if fileMetadata.Receiver == user.Username {
			receiver = user
		}
	}

	if receiver == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Username not found")
		return
	}

	if sender.PermanentFile != nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "You can send only one file in the same time")
		return
	}

	if receiver.PermanentFile != nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Receiver can accept only one file in the same time.")
		return
	}

	permanentFileData := &user.PermanentFileData{
		ID: fileMetadata.ID,
		Receiver: receiver,
		Sender: sender,
		Index: 0,
		Total: fileMetadata.Total,
		Filename: fileMetadata.Filename,
		Size: fileMetadata.Size,
	}

	sender.PermanentFile = permanentFileData
	receiver.PermanentFile = permanentFileData


	fileDto := dto.FileDto{
		ID: fileMetadata.ID,
		Receiver: receiver.Username,
		Sender: sender.Username,
		Filename: fileMetadata.Filename,
		Total: fileMetadata.Total,
		Size: fileMetadata.Size,
	}


	fileDtoJson, err := json.Marshal(fileDto)
	if err != nil {
		fmt.Println(err)
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Unsupported type of file metadata")
		return
	}

	server.sendMessage(receiver.SafeConn, dto.MessageTypeFile, fileDtoJson)
}

func(server *Server) TransferFileChunk(payload json.RawMessage, sender *user.User) {

	fileChunkDto := dto.FileChunkDto{}

	err := json.Unmarshal(payload, &fileChunkDto)

	if err != nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "Unsupported type of file")
		return
	}

	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "User is not expecting file")
		return
	}

	if sender.PermanentFile.Receiver == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "User is not expecting file")
		return
	}

	if sender.PermanentFile.ID != fileChunkDto.ID {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "You can send only one file in the same time")
		return
	}

	if sender.PermanentFile.Index >= sender.PermanentFile.Total {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "File has been sent completely")
		return
	}


	sender.PermanentFile.Index++

	server.sendMessage(sender.PermanentFile.Receiver.SafeConn, dto.MessageTypeFileChunk, payload)

	if sender.PermanentFile.Total == sender.PermanentFile.Index {
		sender.PermanentFile.Receiver.PermanentFile = nil
		sender.PermanentFile = nil
	}
}	

func(server *Server) ProgressFileTransfer(_ json.RawMessage, sender *user.User) {
	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "You are not sending file")
		return
	}
	progress := (sender.PermanentFile.Index * 100) / sender.PermanentFile.Total
	server.sendRawMessage(sender.SafeConn, dto.MessageTypeInfo, fmt.Sprintf("Progress %v", progress))
}


func(server *Server) CancelFileTransfer(_ json.RawMessage, sender *user.User) {
	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.SafeConn, dto.MessageTypeError, "You are not transfering file")
		return
	}

	sender.PermanentFile = nil
	sender.PermanentFile.Receiver.PermanentFile = nil

	if sender == sender.PermanentFile.Sender {
		server.sendRawMessage(sender.PermanentFile.Receiver.SafeConn, dto.MessageTypeCancel, sender.PermanentFile.ID)
	} else {
		server.sendRawMessage(sender.PermanentFile.Sender.SafeConn, dto.MessageTypeCancel, sender.PermanentFile.ID)
	}

	server.sendRawMessage(sender.PermanentFile.Sender.SafeConn, dto.MessageTypeInfo, "File transfer has been canceled")
	server.sendRawMessage(sender.PermanentFile.Receiver.SafeConn, dto.MessageTypeInfo, "File transfer has been canceled")
}