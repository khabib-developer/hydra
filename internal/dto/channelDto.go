package dto

type ChannelDto struct {
	Name    string `json:"name"`
	Members int    `json:"members"`
	Mine    bool   `json:"mine"`
}