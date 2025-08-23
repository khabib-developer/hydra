package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/khabib-developer/chat-application/internal/dto"
	"github.com/khabib-developer/chat-application/internal/user"
	"golang.org/x/term"
)



func ReceiveCommands(u *user.User, state chan string) error {
    if u == nil {
        return fmt.Errorf("nil user pointer")
    }

    reader := bufio.NewReader(os.Stdin)
    currentState := StateNormal

    fmt.Println("💬 Type /help to see available commands:")

    for {
        select {
        case s := <-state:
			currentState = s

        default:

			if currentState == StatePassword {
				fmt.Print("Password of user: ")
				bytePwd, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Println() // move to next line after password input
				if err != nil {
					return fmt.Errorf("read password: %w", err)
				}
				pwdJSON, err := json.Marshal(strings.TrimSpace(string(bytePwd)))
				if err != nil {
					return fmt.Errorf("marshal password: %w", err)
				}
                send(u, dto.MessageTypePassword, pwdJSON)
                currentState = StateNormal
                continue
            }


            fmt.Print("> ")
            input, err := reader.ReadString('\n')
            if err != nil {
                return fmt.Errorf("read command: %w", err)
            }
            input = strings.TrimSpace(input)
            if input == "" {
                continue
            }


            // otherwise normal command handling
            args := strings.SplitN(input, " ", 3)
            cmd := dto.Command(args[0])

            if !isCommand(cmd) {
                fmt.Println("❌ Unknown command. Type /help to see available commands.")
                continue
            }

            switch cmd {
			case dto.CmdHelp:
				helper()

			case dto.CmdCreate:
				if len(args) < 2 {
					fmt.Println("⚠ Usage: /create <channel_name>")
					continue
				}
				channelName := args[1]
				fmt.Printf("✅ Creating channel: %s\n", channelName)
				create(u, channelName)

			case dto.CmdBroadcast:
				if len(args) < 2 {
					fmt.Println("⚠ Usage: /broadcast <message>")
					continue
				}
				message := strings.Join(args[1:], " ")
				broadcast(u, message)

			case dto.CmdJoin:
				if len(args) < 2 {
					fmt.Println("⚠ Usage: /join <channel_name>")
					continue
				}
				channelName := args[1]
				join(u, channelName)

			case dto.CmdMessage:
				if len(args) < 3 {
					fmt.Println("⚠ Usage: /msg <username> <message>")
					continue
				}
				targetUser := args[1]
				message := args[2]
				payload := dto.SendMessageDto{
					Receiver: targetUser,
					Message: message,
				}
				payloadBytes, err := json.Marshal(payload)
				if err != nil {
					return err
				}
				send(u, dto.MessageTypeMessage, payloadBytes)
				time.Sleep(time.Millisecond * 500)

			case dto.CmdChannels: 
				getChannels(u)

			case dto.CmdUsers: 
				getActiveUsers(u)

			case dto.CmdCurrent: 
				getCurrentChannel(u)

			case dto.CmdMembers: 
				getChannelMembers(u)

			case dto.CmdDestroy: 
			if len(args) < 2 {
					fmt.Println("⚠ Usage: /destroy <channel_name>")
					continue
				}
				destroyChannel(u, args[1])

			case dto.CmdFile:
				if len(args) < 3 {
					fmt.Println("⚠ Usage: /file <username> <path>")
					continue
				}
				fileTransfer(u, args[1], args[2])
				

			case dto.CmdInfo: 
				getProfileInfo(u)

			case dto.CmdExit:
				fmt.Println("👋 Exiting...")
				return nil
			}
        }
    }
}
func isCommand(cmd dto.Command) bool {
	return slices.Contains(dto.AllCommands, cmd)
}
