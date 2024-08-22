package test

import (
	"os"
	"testing"

	wx "github.com/IBM/watsonx-go/pkg/models"
)

func TestClientCreationWithEnvVars(t *testing.T) {
	_, err := wx.NewClient()

	if err != nil {
		t.Fatalf("Expected no error for creating client with environment variables, but got %v", err)
	}
}

func TestClientCreationWithPassing(t *testing.T) {
	apiKey, projectID := os.Getenv(wx.WatsonxAPIKeyEnvVarName), os.Getenv(wx.WatsonxProjectIDEnvVarName)

	if apiKey == "" {
		t.Fatal("No watsonx API key provided")
	}
	if projectID == "" {
		t.Fatal("No watsonx project ID provided")
	}

	_, err := wx.NewClient(
		wx.WithWatsonxAPIKey(apiKey),
		wx.WithWatsonxProjectID(projectID),
	)

	if err != nil {
		t.Fatalf("Expected no error for creating client with passing secrets, but got %v", err)
	}
}

func getClient(t *testing.T) *wx.Client {
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

func TestEmptyPromptError(t *testing.T) {
	client := getClient(t)

	_, err := client.GenerateText(
		"dumby model",
		"",
	)
	if err == nil {
		t.Fatalf("Expected error for an empty prompt, but got nil")
	}
}

func TestNilOptions(t *testing.T) {
	client := getClient(t)

	_, err := client.GenerateText(
		"meta-llama/llama-3-70b-instruct",
		"What day is it?",
		nil,
	)
	if err != nil {
		t.Fatalf("Expected no error for nil options, but got %v", err)
	}
}

func TestValidPrompt(t *testing.T) {
	client := getClient(t)

	_, err := client.GenerateText(
		"meta-llama/llama-3-70b-instruct",
		"Test prompt",
	)
	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}
}

func TestGenerateText(t *testing.T) {
	client := getClient(t)

	result, err := client.GenerateText(
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

func TestGenerateTextStream(t *testing.T) {
	client := getClient(t)

	dataChan, err := client.GenerateTextStream(
		"google/flan-ul2",
		"Hi, who are you?",
		wx.WithTemperature(0.9),
		wx.WithTopP(.5),
		wx.WithTopK(10),
		wx.WithMinNewTokens(10),
		wx.WithMaxNewTokens(10),
		wx.WithRandomSeed(1),
	)

	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}

	expectedText := "I am a person. You are a"
	generatedText := ""

	for data := range dataChan {
		generatedText += data.Text
	}

	if generatedText != expectedText {
		t.Fatalf("Expected generated text to be %s, but got %s", expectedText, generatedText)
	}

}

func TestGenerateTextWithNoPrompt(t *testing.T) {
	client := getClient(t)

	dataChan, err := client.GenerateTextStream(
		"google/flan-ul2",
		"",
		wx.WithTemperature(0.9),
		wx.WithTopP(.5),
		wx.WithTopK(10),
		wx.WithMinNewTokens(10),
		wx.WithMaxNewTokens(10),
		wx.WithRandomSeed(1),
	)

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}

	if err != nil && err.Error() != "prompt cannot be empty" {
		t.Fatalf("Expected error to be 'prompt cannot be empty', but got %v", err)
	}

	generatedText := ""
	for data := range dataChan {
		generatedText += data.Text
	}

	if generatedText != "" {
		t.Fatalf("Expected generated text to be empty, but got %s", generatedText)
	}
}

func TestGenerateTextWithNilOptions(t *testing.T) {
	client := getClient(t)

	result, err := client.GenerateText(
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
