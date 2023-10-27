package foundation_models

/*
https://ibm.github.io/watson-machine-learning-sdk/_modules/ibm_watson_machine_learning/foundation_models/model.html#Model
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Credentials struct {
	ApiKey string
	Url    string
}

func getDefaultCredentialUrl() string {
	return "https://us-south.ml.cloud.ibm.com"
}

const (
	VERSION_PARAM            string = "23-05-02"
	TEXT_GENERATION_ENDPOINT string = "/ml/v1-beta/generation"
)

// Model represents the model interface.
type Model struct {
	ModelId     ModelType
	Params      GenParams
	Credentials Credentials
	ProjectId   string
	SpaceId     string
	Client      *resty.Client
}

// NewModel initializes a new Model instance.
func NewModel(modelId ModelType, credentials Credentials, params GenParams, projectId, spaceId string) (*Model, error) {
	if modelId == "" {
		modelId = GetDefaultModelType()
	}
	if credentials.ApiKey == "" {
		return nil, errors.New("API key must be provided")
	}
	if credentials.Url == "" {
		credentials.Url = getDefaultCredentialUrl()
	}
	if projectId == "" && spaceId == "" {
		return nil, errors.New("One of these parameters is required: ['project_id', 'space_id']")
	}

	model := &Model{
		ModelId:     modelId,
		Params:      *NewGenParams(nil),
		Credentials: credentials,
		ProjectId:   projectId,
		SpaceId:     spaceId,
		Client:      resty.New(),
	}

	return model, nil
}

type GenerateResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Results    []struct {
		GeneratedText string `json:"generated_text"`
		StopReason    string `json:"stop_reason"`
	} `json:"results"`
}

// Generate generates completion text based on a given prompt and parameters.
func (model *Model) Generate(prompt string, params map[string]interface{}) (GenerateResponse, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	// Validate input parameters
	if prompt == "" {
		return GenerateResponse{}, errors.New("Prompt cannot be empty")
	}

	payload := map[string]interface{}{
		"model_id": model.ModelId,
		"input":    prompt,
	}

	if params != nil {
		payload["parameters"] = params
	} else {
		payload["parameters"] = model.Params
	}

	// Handle decoding method if it's an enum
	if decodingMethod, ok := payload["parameters"].(map[string]interface{})["DecodingMethod"]; ok {
		if dm, ok := decodingMethod.(DecodingMethods); ok {
			payload["parameters"].(map[string]interface{})["DecodingMethod"] = dm
		}
	}

	if model.ProjectId != "" {
		payload["project_id"] = model.ProjectId
	} else if model.SpaceId != "" {
		payload["space_id"] = model.SpaceId
	}

	// Check return options
	if returnOptions, ok := payload["parameters"].(map[string]interface{})["return_options"].(map[string]interface{}); ok {
		if !returnOptions["input_text"].(bool) && !returnOptions["input_tokens"].(bool) {
			return GenerateResponse{
				Results: []struct {
					GeneratedText string `json:"generated_text"`
					StopReason    string `json:"stop_reason"`
				}{
					{
						GeneratedText: "Response failed with error 'fm_required_parameters_not_provided' on prompt",
						StopReason:    "ERROR",
					},
				},
			}, nil
		}
	}

	var response GenerateResponse

	// Retry request up to 3 times on certain status codes
	for retries := 0; retries < 3; retries++ {
		response, err := model._makeGenerateRequest(payload)
		if err != nil {
			return GenerateResponse{}, err
		}

		statusCode := response.StatusCode

		// Check if the status code indicates a retryable error
		if statusCode == 429 || statusCode == 503 || statusCode == 504 || statusCode == 520 {
			// Sleep for an exponentially increasing duration
			sleepDuration := time.Duration(1<<retries) * time.Second
			time.Sleep(sleepDuration)
		} else {
			// No need to retry, break out of the loop
			break
		}
	}

	if response.StatusCode >= 400 && response.StatusCode <= 599 {
		return GenerateResponse{}, errors.New(fmt.Sprintf("Request failed with: %s (%d)", response.Status, response.StatusCode))
	}

	return response, nil

}

// _makeGenerateRequest sends the generate request and handles the response using Resty.
func (model *Model) _makeGenerateRequest(payload map[string]interface{}) (GenerateResponse, error) {
	var response GenerateResponse

	// Convert the payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return response, err
	}

	// Construct the URL for the generate endpoint
	generateTextURL := fmt.Sprintf("%s/text?version=%s", TEXT_GENERATION_ENDPOINT, VERSION_PARAM)

	// Send the HTTP POST request
	resp, err := model.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+model.Credentials.ApiKey). // Replace with your actual access token
		SetBody(payloadJSON).
		Post(generateTextURL)

	if err != nil {
		return response, err
	}

	// Check for HTTP status code errors
	if resp.StatusCode() >= 400 && resp.StatusCode() <= 599 {
		return response, errors.New(fmt.Sprintf("Request failed with status code: %d", resp.StatusCode()))
	}

	// Parse the response JSON into the GenerateResponse struct
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return response, err
	}

	return response, nil
}
