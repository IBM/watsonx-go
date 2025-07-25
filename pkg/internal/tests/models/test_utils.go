package test

import (
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
