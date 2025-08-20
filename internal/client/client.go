package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/khabib-developer/chat-application/internal/user"
	"golang.org/x/term"
)

// Run asks for username and password and fills the provided user struct.
func Run(u *user.User) error {
	if u == nil {
		return fmt.Errorf("nil user pointer")
	}

	reader := bufio.NewReader(os.Stdin)

	// Ask username until it's non-empty
	var username string
	for {
		fmt.Print("Username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read username: %w", err)
		}
		username = strings.TrimSpace(input)

		if username == "" {
			fmt.Println("Username cannot be empty. Please try again.")
			continue
		}
		break
	}

	// Ask password (hidden)
	fmt.Print("Password (leave empty to skip): ")
	bytePwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // move to next line after password input
	if err != nil {
		return fmt.Errorf("read password: %w", err)
	}
	password := strings.TrimSpace(string(bytePwd))

	// Fill user struct
	u.Username = username
	u.Password = password

	return nil
}
