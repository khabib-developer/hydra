package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


var url string

func init() {
	_ = godotenv.Load()
	url = os.Getenv("SERVER_URL")
	if url == "" {
		url = "http://localhost:8080"
	}
}


// httpSender sends a request with method, path, payload, and headers.
// Returns response body as []byte.
func httpSender(method, path string, payload any, headers map[string]string) ([]byte, error) {
	var body io.Reader

	// Marshal payload only if it's not nil
	if payload != nil {
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal json: %w", err)
		}
		body = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, url + path, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Default content-type if payload exists
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Handle non-200
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s", string(respBody))
	}

	return respBody, nil
}
