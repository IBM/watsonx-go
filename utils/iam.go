package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetIAMToken(apiKey string) (string, error) {
	// Replace these with your IBM Cloud API key and IAM API endpoint.
	iamEndpoint := "https://iam.cloud.ibm.com/identity/token"

	payload := strings.NewReader(fmt.Sprintf("grant_type=urn%%3Aibm%%3Aparams%%3Aoauth%%3Agrant-type%%3Aapikey&apikey=%s", apiKey))

	// Create an HTTP request to the IAM endpoint.
	req, err := http.NewRequest("POST", iamEndpoint, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set the request headers.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse the response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// Parse the JSON response.
	var tokenResponse TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return "", err
	}

	// Access the IAM token.
	accessToken := tokenResponse.AccessToken

	return accessToken, nil
}
