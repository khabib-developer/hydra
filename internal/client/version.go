package client

import (
	"fmt"

	"github.com/khabib-developer/chat-application/internal/version"
)

func GetVersion() error {

	fmt.Println("Checking version...")
	
	respBody, err := httpSender("GET", "/version", nil, nil)

	if err != nil {
		return err
	}

	serverVersion := string(respBody)
	clientVersion := version.GetVersion()

	if serverVersion != clientVersion {
		return fmt.Errorf(
			"❌ version mismatch: client is %s, server is %s\n👉 Please update your app to match.",
			clientVersion, serverVersion,
		)
	}

	fmt.Printf("Version check passed (v%s)\n", clientVersion)

	return nil
}