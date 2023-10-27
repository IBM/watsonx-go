package foundation_models

/*
https://ibm.github.io/watson-machine-learning-sdk/_modules/ibm_watson_machine_learning/foundation_models/model.html#Model
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.ibm.com/robby-ibm/go-watsonx/src/foundation_models/client"
	"github.ibm.com/robby-ibm/go-watsonx/src/utils"
)

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
		credentials.setDefaultUrl()
	}
	if projectId == "" && spaceId == "" {
		return nil, errors.New("one of these parameters is required: ['project_id', 'space_id']")
	}

	println(credentials.ApiKey)

	// Get Bearer token from IAM
	iamToken, err := utils.GetIAMToken(credentials.ApiKey)
	if err != nil {
		fmt.Println("Error getting IAM token: ", err)
		return nil, err
	}

	credentials.ApiKey = iamToken

	model := &Model{
		ModelId:     modelId,
		Params:      *NewGenParams(nil),
		Credentials: credentials,
		ProjectId:   projectId,
		SpaceId:     spaceId,
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
		return GenerateResponse{}, errors.New("prompt cannot be empty")
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
		return GenerateResponse{}, fmt.Errorf(fmt.Sprintf("Request failed with: %s (%d)", response.Status, response.StatusCode))
	}

	return response, nil

}

// _makeGenerateRequest sends the generate request and handles the response using the http package.
func (model *Model) _makeGenerateRequest(payload map[string]interface{}) (GenerateResponse, error) {
	var response GenerateResponse

	// Construct the URL for the generate endpoint
	generateTextURL := fmt.Sprintf("%s/%s/text?version=%s", model.Credentials.Url, TEXT_GENERATION_ENDPOINT, VERSION_PARAM)

	// Send the HTTP POST request using PostRequest
	resp, err := client.PostRequest(generateTextURL, payload, model.Credentials.ApiKey)
	if err != nil {
		return response, err
	}

	// Check for HTTP status code errors
	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		return response, fmt.Errorf(fmt.Sprintf("Request failed with status code: %d", resp.StatusCode))
	}

	// Parse the response JSON into the GenerateResponse struct
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}

func (m Model) String() string {
	return fmt.Sprintf("ModelId: %s\nParams: %s\nCredentials: %s\nProjectId: %s\nSpaceId: %s",
		m.ModelId, m.Params, m.Credentials, m.ProjectId, m.SpaceId)
}
