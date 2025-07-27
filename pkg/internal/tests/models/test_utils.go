package test

import (
	"encoding/json"
	"os"
	"testing"

	wx "github.com/IBM/watsonx-go/pkg/models"
)

func getClient(t *testing.T) *wx.Client {
	apiKey, projectID := os.Getenv(wx.WatsonxAPIKeyEnvVarName), os.Getenv(wx.WatsonxProjectIDEnvVarName)

	if apiKey == "" {
		t.Fatal("No watsonx API key provided")
	}
	if projectID == "" {
		t.Fatal("No watsonx project ID provided")
	}

	client, err := wx.NewClient(
		wx.WithClientRetryConfig(wx.NewRetryConfig(
			wx.WithReturnHTTPStatusAsErr(false)),
		),
		wx.WithWatsonxAPIKey(apiKey),
		wx.WithWatsonxProjectID(projectID),
	)
	if err != nil {
		t.Fatalf("Failed to create client for testing. Error: %v", err)
	}

	return client
}

func getClientWithLegacyRetry(t *testing.T) *wx.Client {
	apiKey, projectID := os.Getenv(wx.WatsonxAPIKeyEnvVarName), os.Getenv(wx.WatsonxProjectIDEnvVarName)

	if apiKey == "" {
		t.Fatal("No watsonx API key provided")
	}
	if projectID == "" {
		t.Fatal("No watsonx project ID provided")
	}

	client, err := wx.NewClient(
		wx.WithWatsonxAPIKey(apiKey),
		wx.WithWatsonxProjectID(projectID),
	)
	if err != nil {
		t.Fatalf("Failed to create client for testing. Error: %v", err)
	}

	return client
}

type WatsonxError struct {
	Errors []struct {
		Code     string `json:"code"`
		Message  string `json:"message"`
		MoreInfo string `json:"more_info"`
	} `json:"errors"`
	Trace      string `json:"trace"`
	StatusCode int    `json:"status_code"`
}

func parseResponseErrMessage(t *testing.T, err error) *WatsonxError {
	if err == nil {
		return nil
	}
	msg := err.Error()
	var wxErr WatsonxError
	if err := json.Unmarshal([]byte(msg), &wxErr); err != nil {
		t.Logf("Failed to unmarshal error JSON: %v", err)
		return nil
	}
	return &wxErr
}
