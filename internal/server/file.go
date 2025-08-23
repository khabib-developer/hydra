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
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "Unsupported type of file metadata")
		return
	}

	var receiver *user.User

	for _, user := range server.Users {
		if fileMetadata.Receiver == user.Username {
			receiver = user
		}
	}


	if receiver == nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "Username not found")
		return
	}

	if sender.PermanentFile != nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "You can send only one file in the same time")
		return
	}

	if receiver.PermanentFile != nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "Receiver can accept only one file in the same time.")
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


	fileDtoJson, err := json.Marshal(permanentFileData)
	if err != nil {
		fmt.Println(err)
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "Unsupported type of file metadata")
		return
	}

	server.sendMessage(sender.Conn, dto.MessageTypeFile, fileDtoJson)
}

func(server *Server) TransferFileChunk(payload json.RawMessage, sender *user.User) {

	fileChunkDto := dto.FileChunkDto{}

	err := json.Unmarshal(payload, &fileChunkDto)

	if err != nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "Unsupported type of file")
		return
	}

	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "User is not expecting file")
		return
	}

	if sender.PermanentFile.Receiver == nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "User is not expecting file")
		return
	}

	if sender.PermanentFile.ID != fileChunkDto.ID {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "You can send only one file in the same time")
		return
	}

	if sender.PermanentFile.Total >= sender.PermanentFile.Index {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "File has been sent completely")
		return
	}

	if sender.PermanentFile.Index == 0 {
		server.sendRawMessage(sender.Conn, dto.MessageTypeInfo, "Your file has been started to transfer. To see progress type /progress. To cancel file transfer type /cancel")
	}

	sender.PermanentFile.Index++

	server.sendMessage(sender.PermanentData.Receiver.Conn, dto.MessageTypeFileChunk, payload)

	if sender.PermanentFile.Total == sender.PermanentFile.Index {
		sender.PermanentFile = nil
		sender.PermanentFile.Receiver.PermanentFile = nil
	}
}	

func(server *Server) ProgressFileTransfer(_ json.RawMessage, sender *user.User) {
	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "You are not sending file")
		return
	}
	progress := (sender.PermanentFile.Index * 100) / sender.PermanentFile.Total
	server.sendRawMessage(sender.Conn, dto.MessageTypeInfo, fmt.Sprintf("Progress %v", progress))
}


func(server *Server) CancelFileTransfer(_ json.RawMessage, sender *user.User) {
	if sender.PermanentFile == nil {
		server.sendRawMessage(sender.Conn, dto.MessageTypeError, "You are not transfering file")
		return
	}

	sender.PermanentFile = nil
	sender.PermanentFile.Receiver.PermanentFile = nil

	if sender == sender.PermanentFile.Sender {
		server.sendRawMessage(sender.PermanentFile.Receiver.Conn, dto.MessageTypeCancel, sender.PermanentFile.ID)
	} else {
		server.sendRawMessage(sender.PermanentFile.Sender.Conn, dto.MessageTypeCancel, sender.PermanentFile.ID)
	}

	server.sendRawMessage(sender.PermanentFile.Sender.Conn, dto.MessageTypeInfo, "File transfer has been canceled")
	server.sendRawMessage(sender.PermanentFile.Receiver.Conn, dto.MessageTypeInfo, "File transfer has been canceled")
}