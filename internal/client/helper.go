package client

import "fmt"
func helper() {
	fmt.Println("\n📜 Available Commands:")
	fmt.Println("  /help                → Show this help menu")
	fmt.Println("  /create <name>       → Create a channel")
	fmt.Println("  /broadcast <message> → Send a message to all users")
	fmt.Println("  /join <name>         → Join a channel")
	fmt.Println("  /msg <user> <msg>    → Send a private message")
	fmt.Println("  /exit                → Exit the application")
	fmt.Println("  /users               → List of users")
	fmt.Println("  /channels            → List of channels")
	fmt.Println("  /destroy <name>      → Destroy a channel")
	fmt.Println("  /members <name>      → List members of a channel")
	fmt.Println("  /current             → Show current channel")
	fmt.Println()
}