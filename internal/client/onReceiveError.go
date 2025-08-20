package client

import (
	"encoding/json"
	"fmt"

	"github.com/khabib-developer/chat-application/internal/user"
)

func onReceiveError(_ *user.User, payload json.RawMessage, _ chan string) error {
	var errMsg string
	if err := json.Unmarshal(payload, &errMsg); err != nil {
		return fmt.Errorf("failed to unmarshal error payload: %w", err)
	}

	// ANSI escape codes
	const (
		red   = "\033[31m"
		reset = "\033[0m"
	)

	// Print error in red
	fmt.Println(string(red), "❌ ERROR:", errMsg, string(reset))

	return nil
}