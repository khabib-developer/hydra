package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
	"github.com/khabib-developer/chat-application/internal/utils"
)

const ChunkSize /*bytes*/ = 1024 * 1024 // 1mb

const dir = "./files"


func fileTransfer(u *user.User,username string, path string) {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	fileStats, err := f.Stat()

	if err != nil {
		fmt.Println("Incorrect file")
		return
	}

	totalChunks := (fileStats.Size() + ChunkSize - 1) / ChunkSize

	id := uuid.NewString()

	fileMetaData := user.FileMetadata{
		ID: id,
		Receiver: username,
		Filename: fileStats.Name(),
		Total: totalChunks,
		Size: fileStats.Size(),
	}

	fileMetaDataJson, err := json.Marshal(fileMetaData)

	if err != nil {
		fmt.Println("Error while marshaling json")
		return
	}

	
	send(u, dto.MessageTypeFile, fileMetaDataJson) //17034654

	// go fileChunkTransfer(f, id, u)
	
}

// func fileChunkTransfer(f *os.File, id string, u *user.User) {
// 	buf := make([]byte, ChunkSize)
// 	index := 0
// 	defer f.Close()
// 	for {
// 		n, err := f.Read(buf)

// 		if n > 0 {
// 			fmt.Println(n)
// 			chunkDto := dto.FileChunkDto{
// 				ID:    id,
// 				Body:  buf[:n], 
// 				Index: int64(index),
// 			}
// 			chunkDtoJson, err := json.Marshal(chunkDto)
// 			if err != nil {
// 				onError(u)
// 				return
// 			}
// 			send(u, dto.MessageTypeFileChunk, chunkDtoJson)
// 		}

// 		if err == io.EOF {
// 			fmt.Println("File has been sent")
// 			break
// 		}

// 		if err != nil {
// 			onError(u)
// 			return
// 		}

// 		index++

// 	}

// }

// func onError(u *user.User) {
// 	message := "Something went wrong"
// 	messageJson, err := json.Marshal(message)
// 	if err != nil {
// 		return
// 	}
// 	send(u, dto.MessageTypeCancel, messageJson)
// }

func onReceiveFileMetadata(u *user.User, payload json.RawMessage, _ chan string) error {
	
	if u.PermanentFile != nil {
		return fmt.Errorf("someone trying to send you file but you cannot receive several files in the same time")
	}
	
	permanentFileData := user.PermanentFileData{}
	if err := json.Unmarshal(payload, &permanentFileData); err != nil {
		return err
	}

	u.PermanentFile = &permanentFileData

	os.MkdirAll(dir, 0755)

	free, _ := utils.GetFreeSpace(dir)

	fmt.Printf("%v", free)

	if free < uint64(permanentFileData.Size) {
		return fmt.Errorf("not enough disk space: need %d, have %d", permanentFileData.Size, free)
	}

	fname := fmt.Sprintf("%s_%d", permanentFileData.Filename, time.Now().Unix())

	path := filepath.Join(dir,fname)


	dst, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	} 

	permanentFileData.File = dst
	permanentFileData.Filename = fname

	return nil
}

func onReceiveFileChunk(u *user.User, payload json.RawMessage, _ chan string) error {

	

	return nil
}


func onReceiveCancelTransfer(u *user.User, payload json.RawMessage, _ chan string) error {
	
	return nil
}
