package client

import (
	"encoding/json"
	"fmt"

	"github.com/khabib-developer/chat-application/internal/user"
)


func onReceiveInfo(_ *user.User, payload json.RawMessage, _ chan string) error {
	var info string
	if err := json.Unmarshal(payload, &info); err != nil {
		return fmt.Errorf("failed to unmarshal info payload: %w", err)
	}

	// ANSI color codes
	const (
		green  = "\033[32m"
		yellow = "\033[33m"
		red    = "\033[31m"
		reset  = "\033[0m"
	)

	// Example: print green
	fmt.Println(":",string(green), info, string(reset))

	return nil
}
