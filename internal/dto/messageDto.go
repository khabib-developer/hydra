package dto

type SendMessageDto struct {
	Receiver string `json:"receiver"`
	Message  string	`json:"message"`
}


type ReceiveMessageDto struct {
	Sender string `json:"sender"`
	Message  string	`json:"message"`
}

type ChannelMessageDto struct {
	Channel string `json:"channel"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}