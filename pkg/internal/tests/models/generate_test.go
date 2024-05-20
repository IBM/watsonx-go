package test

import (
	"os"
	"testing"

	wx "github.com/h0rv/go-watsonx/pkg/models"
)

func getModel(t *testing.T) *wx.Model {
	apiKey := os.Getenv(wx.WatsonxAPIKeyEnvVarName)
	projectID := os.Getenv(wx.WatsonxProjectIDEnvVarName)
	if apiKey == "" {
		t.Fatal("No watsonx API key provided")
	}
	if projectID == "" {
		t.Fatal("No watsonx project ID provided")
	}

	model, err := wx.NewModel(
		wx.WithWatsonxAPIKey(apiKey),
		wx.WithWatsonxProjectID(projectID),
	)
	if err != nil {
		t.Fatalf("Failed to create model for testing. Error: %v", err)
	}

	return model
}

func TestEmptyPromptError(t *testing.T) {
	model := getModel(t)

	_, err := model.GenerateText(
		"dumby model",
		"",
	)
	if err == nil {
		t.Fatalf("Expected error for an empty prompt, but got nil")
	}
}

func TestNilOptions(t *testing.T) {
	model := getModel(t)

	_, err := model.GenerateText(
		"meta-llama/llama-3-70b-instruct",
		"What day is it?",
		nil,
	)
	if err != nil {
		t.Fatalf("Expected no error for nil options, but got %v", err)
	}
}

func TestValidPrompt(t *testing.T) {
	model := getModel(t)

	_, err := model.GenerateText(
		"meta-llama/llama-3-70b-instruct",
		"Test prompt",
	)
	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}
}

func TestGenerateText(t *testing.T) {
	model := getModel(t)

	result, err := model.GenerateText(
		"meta-llama/llama-3-70b-instruct",
		"Hi, who are you?",
		wx.WithTemperature(0.9),
		wx.WithTopP(.5),
		wx.WithTopK(10),
		wx.WithMaxNewTokens(512),
	)
	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}
	if result.Text == "" {
		t.Fatal("Expected a result, but got an empty string")
	}
}

func TestGenerateTextWithNilOptions(t *testing.T) {
	model := getModel(t)

	result, err := model.GenerateText(
		"meta-llama/llama-3-70b-instruct",
		"Who are you?",
		nil,
	)
	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}
	if result.Text == "" {
		t.Fatal("Expected a result, but got an empty string")
	}
}
