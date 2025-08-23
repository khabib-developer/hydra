package user

import "os"


type FileMetadata struct {
	ID		 string  `json:"id"`
	Receiver string  `json:"receiver"`
	Filename string  `json:"filename"`
	Total    int64   `json:"total"`
	Size     int64   `json:"size"`
}


type PermanentFileData struct {
	ID  	 string
	Receiver *User
	Sender   *User
	Index    int64
	Total    int64
	Filename string
	Size     int64 
	File 	 *os.File
}