package test

import (
	"os"
	"testing"

	"github.com/IBM/watsonx-go/pkg/constants"
	wx "github.com/IBM/watsonx-go/pkg/models"
)

func eitherProjectOrSpace(t *testing.T, projectId, spaceId string) wx.ClientOption {
	t.Helper()

	if projectId != "" {
		return wx.WithWatsonxProjectID(projectId)
	} else if spaceId != "" {
		return wx.WithWatsonxSpaceID(spaceId)
	}
	t.Fatal("No watsonx project Id or space ID provided")

	// this satisfies go linter, but fatal should stop execution before we get here
	return func(co *wx.ClientOptions) {}
}

func getClient(t *testing.T) *wx.Client {
	apiKey := os.Getenv(constants.WatsonxAPIKeyEnvVarName)

	projectID, spaceID := os.Getenv(constants.WatsonxProjectIDEnvVarName), os.Getenv(constants.WatsonxSpaceIDEnvVarName)
	option := eitherProjectOrSpace(t, projectID, spaceID)

	if apiKey == "" {
		t.Fatal("No watsonx API key provided")
	}

	client, err := wx.NewClient(
		wx.WithWatsonxAPIKey(apiKey),
		option,
	)
	if err != nil {
		t.Fatalf("Failed to create client for testing. Error: %v", err)
	}

	return client
}
