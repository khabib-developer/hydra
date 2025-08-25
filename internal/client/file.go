package client

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
	"github.com/khabib-developer/chat-application/internal/utils"
)

const ChunkSize /*bytes*/ = 1024 * 256  // 1mb


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

	fmt.Println("📥 Sending file:")
	fmt.Printf(" 📝 Name:   %s\n", fileMetaData.Filename)
	fmt.Printf(" 📦 Size:   %.2f MB (%d bytes)\n", float64(fileMetaData.Size)/(1024*1024), fileMetaData.Size)
	fmt.Println("--------------------------------------------------")
	fmt.Println("❌ To cancel file transfer, type: /cancel")
	fmt.Println("--------------------------------------------------")
	
	send(u, dto.MessageTypeFile, fileMetaDataJson)

	go fileChunkTransfer(f, id, u, totalChunks)
	
}

func fileChunkTransfer(f *os.File, id string, u *user.User, totalChunks int64) {
	buf := make([]byte, ChunkSize)
	index := 0
	
	defer f.Close()
	for {
		n, err := f.Read(buf)

		if n > 0 {
			chunkDto := dto.FileChunkDto{
				ID:    id,
				Body:  buf[:n], 
				Index: int64(index),
			}
			chunkDtoJson, err := json.Marshal(chunkDto)
			if err != nil {
				onError(u)
				return
			}
			send(u, dto.MessageTypeFileChunk, chunkDtoJson)

			progress := int(((index + 1) * 100) / int(totalChunks))
            visualizeProgress(progress)
		}

		if err == io.EOF {
			fmt.Println()
			fmt.Printf("🎉 File transfer complete!\n")
			fmt.Println()
			fmt.Println("-==================================================-")
			fmt.Println()
			fmt.Print(">")
			break
		}

		if err != nil {
			onError(u)
			return
		}

		index++
	}
}


func visualizeProgress(progress int) {
	barWidth := 50 
	filled := (progress * barWidth) / 100
	fmt.Printf("\r[%-*s] %3d%%", barWidth, strings.Repeat("=", filled), progress)

	if progress == 100 {
        fmt.Println() // move to new line when finished
    }

}

func onError(u *user.User) {
	message := "Something went wrong"
	messageJson, err := json.Marshal(message)
	if err != nil {
		return
	}
	send(u, dto.MessageTypeCancel, messageJson)
}

func onReceiveFileMetadata(u *user.User, payload json.RawMessage, _ chan string) error {
	if u.PermanentFile != nil {
		return fmt.Errorf("someone trying to send you file but you cannot receive several files in the same time")
	}
	
	fileDto := dto.FileDto{}
	if err := json.Unmarshal(payload, &fileDto); err != nil {
		return err
	}

	u.PermanentFile = &user.PermanentFileData{
		ID: fileDto.ID,
		Index: 0,
		Total: fileDto.Total,
		Size: fileDto.Size,
	}

	dir, err := getSaveDir()

	if err != nil {return err}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", dir, err)
	}

	free, err := utils.GetFreeSpace(dir)

	if err != nil {
		return err
	}

	if free < uint64(fileDto.Size) {
		return fmt.Errorf("not enough disk space: need %d, have %d", fileDto.Size, free)
	}

	fname := fmt.Sprintf("%d_%s",  time.Now().Unix(), fileDto.Filename)

	path := filepath.Join(dir,fname)

	dst, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	} 

	u.PermanentFile.File = dst
	u.PermanentFile.Filename = fname

	fmt.Println("📥 Receiving file:")
	fmt.Printf("  📝 Name:   %s\n", fileDto.Filename)
	fmt.Printf("  👤 Sender: %s\n", fileDto.Sender)
	fmt.Printf("  📂 Save:   %s\n", path)
	fmt.Printf("  📦 Size:   %.2f MB (%d bytes)\n", float64(fileDto.Size)/(1024*1024), fileDto.Size)
	fmt.Printf("  💾 Free:   %.2f GB available\n", float64(free)/(1024*1024*1024))
	fmt.Println("--------------------------------------------------")
	fmt.Println("❌ To cancel file transfer, type: /cancel")
	fmt.Println("--------------------------------------------------")

	return nil
}

func onReceiveFileChunk(u *user.User, payload json.RawMessage, _ chan string) error {

	if u.PermanentFile == nil {
		return fmt.Errorf("you are not receiving file")
	}

	fileChunkDto := dto.FileChunkDto{}
	if err := json.Unmarshal(payload, &fileChunkDto); err != nil {
		return fmt.Errorf("failed to parse file chunk: %w", err)
	}

	if u.PermanentFile.ID != fileChunkDto.ID {
		return fmt.Errorf("filechunk id is not correct")
	}


	if u.PermanentFile.Index == u.PermanentFile.Total {
		return nil
	}

	if _, err := u.PermanentFile.File.Write(fileChunkDto.Body); err != nil {
		return fmt.Errorf("failed to write chunk to file: %w", err)
	}

	u.PermanentFile.Index++

	

	if fileChunkDto.Index+1 == u.PermanentFile.Total {
		u.PermanentFile.File.Close()

		fmt.Println()
		fmt.Printf("🎉 File transfer complete!\n")
		fmt.Println()
		fmt.Println("-==================================================-")
		fmt.Println()
		fmt.Print(">")

		u.PermanentFile = nil
	} else {
		progress := int(((u.PermanentFile.Index + 1) * 100) / (u.PermanentFile.Total))

		visualizeProgress(progress)
	}
	

	return nil
}


func onReceiveCancelTransfer(u *user.User, payload json.RawMessage, _ chan string) error {
	
	return nil
}


func getSaveDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot detect home dir: %w", err)
	}

	dir := filepath.Join(home, "hydra-files")

	// Create the folder if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create save dir: %w", err)
	}

	return dir, nil
}


//	 /Users/habib/Downloads/jivotnoe.mp3