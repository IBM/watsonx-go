package test

import (
	"fmt"
	"os"
	"testing"

	wx "github.com/IBM/watsonx-go/pkg/models"
)

const (
	modelLlama3  = "meta-llama/llama-3-3-70b-instruct"
	modelFlanUL2 = "google/flan-ul2"
	modelInvalid = "invalid-model-test"
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
		wx.WithClientRetryConfig(wx.NewRetryConfig(
			wx.WithReturnHTTPStatusAsErr(false),
		)),
		wx.WithWatsonxAPIKey(apiKey),
		wx.WithWatsonxProjectID(projectID),
	)

	if err != nil {
		t.Fatalf("Expected no error for creating client with passing secrets, but got %v", err)
	}
}

func TestEmptyPromptError(t *testing.T) {
	client := getClient(t)

	_, err := client.GenerateText(
		modelLlama3,
		"",
	)
	if err == nil {
		t.Fatal("Expected error for an empty prompt, but got nil")
	}
}

func TestNilOptions(t *testing.T) {
	client := getClient(t)

	_, err := client.GenerateText(
		modelLlama3,
		"What day is it?",
		nil,
	)
	if err != nil {
		t.Fatalf("Expected no error for nil options, but got %v", err)
	}
}

func TestInvalidParameter(t *testing.T) {
	client := getClient(t)

	result, err := client.GenerateText(
		modelLlama3,
		"What day is it?",
		wx.WithTemperature(100), // Invalid temperature
	)
	if err == nil {
		t.Fatal("Expected error for invalid parameter but got nil")
	}
	if result.Text != "" {
		t.Fatal("Expected empty response but got:", result.Text)
	}

	errorMessage := parseResponseErrMessage(t, err)
	if errorMessage == nil {
		t.Fatalf("Expected JSON error message, but got: %v", err)
	}

	if errorMessage.StatusCode != 400 {
		t.Fatalf("Expected status code 400, but got: %d", errorMessage.StatusCode)
	}
	expectedErrText := "Json document validation error: parameters.temperature should not be greater than 2.0"
	if len(errorMessage.Errors) == 0 || errorMessage.Errors[0].Message != expectedErrText {
		t.Fatalf("Expected error message to be %s, but got: %v", expectedErrText, errorMessage.Errors[0].Message)
	}
}

func TestInvalidModel(t *testing.T) {
	client := getClient(t)

	result, err := client.GenerateText(
		modelInvalid,
		"What day is it?",
	)

	if err == nil {
		t.Fatal("Expected error for invalid model but got nil")
	}
	if result.Text != "" {
		t.Fatal("Expected empty response but got:", result.Text)
	}

	errorMessage := parseResponseErrMessage(t, err)
	if errorMessage == nil {
		t.Fatalf("Expected JSON error message, but got: %v", err)
	}

	if errorMessage.StatusCode != 404 {
		t.Fatalf("Expected status code 404, but got: %d", errorMessage.StatusCode)
	}
	expectedErrText := fmt.Sprintf("Model '%s' is not supported", modelInvalid)
	if len(errorMessage.Errors) == 0 || errorMessage.Errors[0].Message != expectedErrText {
		t.Fatalf("Expected error message to be %s, but got: %v", expectedErrText, errorMessage.Errors[0].Message)
	}
}

func TestInvalidModelWithLegacyRetry(t *testing.T) {
	client := getClientWithLegacyRetry(t)

	result, err := client.GenerateText(
		modelInvalid,
		"What day is it?",
	)

	if err == nil {
		t.Fatal("Expected error for invalid model but got nil")
	}
	if result.Text != "" {
		t.Fatal("Expected empty response but got:", result.Text)
	}
	expectedErrText := "404 Not Found"
	if err.Error() != expectedErrText {
		t.Fatalf("Expected error to be '%s', but got %v", expectedErrText, err)
	}
}

func TestValidPrompt(t *testing.T) {
	client := getClient(t)

	result, err := client.GenerateText(
		modelLlama3,
		"Test prompt",
	)

	if err != nil {
		t.Fatalf("Expected no error, but got an error: %v", err)
	}
	if result.Text == "" {
		t.Fatal("Expected a result, but got an empty string")
	}
}

func TestGenerateText(t *testing.T) {
	client := getClient(t)

	result, err := client.GenerateText(
		modelLlama3,
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
		modelFlanUL2,
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
		modelFlanUL2,
		"",
		wx.WithTemperature(0.9),
		wx.WithTopP(.5),
		wx.WithTopK(10),
		wx.WithMinNewTokens(10),
		wx.WithMaxNewTokens(10),
		wx.WithRandomSeed(1),
	)

	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}

	if err.Error() != "prompt cannot be empty" {
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
		modelLlama3,
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
