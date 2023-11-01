package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// PostRequest sends an HTTP POST request to the specified URL with the given payload and headers, and returns the response as an HTTP response object.
func PostRequest(url string, payload map[string]interface{}, access_token string) (*http.Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+access_token) // Replace with your actual access token

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
