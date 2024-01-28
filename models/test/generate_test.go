package test

import (
	"os"
	"testing"

	wx "github.com/h0rv/go-watsonx/models"
)

func getModel(t *testing.T) *wx.Model {
	apiKey := os.Getenv(wx.IBMCloudAPIKeyEnvVarName)
	projectID := os.Getenv(wx.WatsonxProjectIDEnvVarName)
	if apiKey == "" {
        t.Fatal("No IBM Cloud API key provided")
	}
	if projectID == "" {
        t.Fatal("No watsonx project ID provided")
	}

	model, err := wx.NewModel(
		apiKey,
		projectID,
		wx.WithModel(wx.FLAN_UL2),
	)
	if err != nil {
		t.Fatalf("Failed to create model for testing. Error: %v", err)
	}

	return model
}

func TestEmptyPromptError(t *testing.T) {
	model := getModel(t)

	_, err := model.GenerateText("")
	if err == nil {
		t.Fatalf("Expected error for an empty prompt, but got nil")
	}
}

func TestNilOptions(t *testing.T) {
	model := getModel(t)

	_, err := model.GenerateText("What day is it?", nil)
	if err != nil {
		t.Fatalf("Expected no error for nil options, but got %v", err)
	}
}

func TestValidPrompt(t *testing.T) {
	model := getModel(t)

	prompt := "Test prompt"
	_, err := model.GenerateText(prompt)
	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}
}

func TestGenerateText(t *testing.T) {
	model := getModel(t)

	prompt := "Hi, who are you?"
	result, err := model.GenerateText(
		prompt,
		wx.WithTemperature(0.9),
		wx.WithTopP(.5),
		wx.WithTopK(10),
		wx.WithMaxNewTokens(512),
		wx.WithDecodingMethod(wx.Greedy),
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

	prompt := "Who are you?"
	result, err := model.GenerateText(
		prompt,
		nil,
	)
	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}
	if result.Text == "" {
		t.Fatal("Expected a result, but got an empty string")
	}
}
