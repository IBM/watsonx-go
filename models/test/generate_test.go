package test

import (
	"os"
	"testing"

	"github.ibm.com/robby-ibm/go-watsonx/models"
)

const (
	dumbyProjectId = "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	dumbySpaceId   = "c9d1e1b2-936c-4e70-9d54-21ed4049e131"
)

func TestGenerateText(t *testing.T) {
	// ENV Variables For Testing
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		// If the environment variable is not set, use a dummy value
		apiKey = "your_dummy_api_key"
	}

	// Create a test model with dummy data
	model, err := models.NewModel(
		"",
		models.Credentials{ApiKey: apiKey, Url: ""},
		dumbyProjectId,
		dumbySpaceId,
		nil,
	)
	if model == nil {
		t.Error("Expected proper creation of model. Error: ", err)
		return
	}

	// Test case 1: Empty prompt should return an error
	_, err = model.GenerateText("", nil)
	if err == nil {
		t.Error("Expected error for an empty prompt, but got nil")
	}

	// Test case 2: Valid prompt with no additional parameters
	prompt := "Test prompt"
	_, err = model.GenerateText(prompt, nil)
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	// Add more test cases as needed
}
