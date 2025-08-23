package dto

type FileDto struct {
	ID		 string  `json:"id"`
	Sender   string  `json:"sender"`
	Filename string  `json:"filename"`
	Size     int64   `json:"size"`
}


type FileChunkDto struct {
	ID    string `json:"id"`
	Body  []byte `json:"body"`
	Index int64  `json:"index"`
}