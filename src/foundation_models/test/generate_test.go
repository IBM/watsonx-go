package test

import (
	"os"
	"testing"

	"github.ibm.com/robby-ibm/go-watsonx/src/foundation_models"
)

func TestGenerate(t *testing.T) {
	// ENV Variables For Testing
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		// If the environment variable is not set, use a dummy value
		apiKey = "your_dummy_api_key"
	}

	// Create a test model with dummy data
	model, err := foundation_models.NewModel("", foundation_models.Credentials{ApiKey: apiKey, Url: ""}, *foundation_models.NewGenParams(nil), "dummyProjectId", "dummySpaceId")
	if model == nil {
		t.Error("Expected proper creation of model. Error: ", err)
		return
	}

	t.Log("Model:\n\n", model)

	// Test case 1: Empty prompt should return an error
	_, err = model.Generate("", nil)
	if err == nil {
		t.Error("Expected error for an empty prompt, but got nil")
	}

	// Test case 2: Valid prompt with no additional parameters
	prompt := "Test prompt"
	response, err := model.Generate(prompt, nil)
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}
	if response.Status != "OK" {
		t.Errorf("Expected 'OK' status, but got %s", response.Status)
	}

	// Add more test cases as needed
}
