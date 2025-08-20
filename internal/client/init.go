package client

import (
	"fmt"
	"time"
)


func Draw() {
	art := []string{
		"  _____ _           _   _ ",
		" / ____| |         | | | |",
		"| |    | |__   __ _| |_| |__   __ _ _ __ ___",
		"| |    | '_ \\ / _` | __| '_ \\ / _` | '__/ _ \\",
		"| |____| | | | (_| | |_| | | | (_| | | |  __/",
		" \\_____|_| |_|\\__,_|\\__|_| |_|\\__,_|_|  \\___|",
		"",
		"        Welcome to the CLI Chat — say hi!",
	}

	for _, line := range art {
		fmt.Println(line)
		time.Sleep(60 * time.Millisecond)
	}
}