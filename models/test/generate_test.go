package test

import (
	"os"
	"testing"

	"github.com/h0rv/go-watsonx/models"
)

const (
	DummyProjectID = "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	DummySpaceID   = "c9d1e1b2-936c-4e70-9d54-21ed4049e131"
)

func getModel(t *testing.T) *models.Model {
	// ENV Variables For Testing
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		// If the environment variable is not set, use a dummy value
		apiKey = "your_dummy_api_key"
	}

	// Create a test model with dummy data
	model, err := models.New(
		apiKey,
		DummyProjectID,
		models.WithModel(models.FLAN_UL2),
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
		t.Errorf("Expected error for an empty prompt, but got nil")
	}
}

func TestNilOptions(t *testing.T) {
	model := getModel(t)

	_, err := model.GenerateText("What day is it?", nil)
	if err != nil {
		t.Errorf("Expected no error for nil options, but got %v", err)
	}
}

func TestValidPrompt(t *testing.T) {
	model := getModel(t)

	prompt := "Test prompt"
	_, err := model.GenerateText(prompt)
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
}
