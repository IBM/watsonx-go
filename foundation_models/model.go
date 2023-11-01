package foundation_models

/*
 *  https://ibm.github.io/watson-machine-learning-sdk/_modules/ibm_watson_machine_learning/foundation_models/model.html#Model
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.ibm.com/robby-ibm/go-watsonx/foundation_models/client"
	"github.ibm.com/robby-ibm/go-watsonx/utils"
)

const (
	VERSION_PARAM            string = "2023-05-02"
	TEXT_GENERATION_ENDPOINT string = "/ml/v1-beta/generation"
)

type Model struct {
	ModelId     ModelType
	Params      *GenParams
	Credentials Credentials
	ProjectId   string
	SpaceId     string
}

func NewModel(modelId ModelType, credentials Credentials, params *GenParams, projectId, spaceId string) (*Model, error) {
	if modelId == "" {
		modelId = GetDefaultModelType()
	}
	if credentials.ApiKey == "" {
		return nil, errors.New("API key must be provided")
	}
	if credentials.Url == "" {
		credentials.setDefaultUrl()
	}
	if params == nil {
		params = DefaultGenParams()
	}
	if projectId == "" && spaceId == "" {
		return nil, errors.New("one of these parameters is required: ['project_id', 'space_id']")
	}

	accessToken, err := utils.GetIAMToken(credentials.ApiKey)

	if err != nil {
		return nil, errors.New("error getting IAM token")
	} else {
		credentials.AccessToken = accessToken
	}

	model := &Model{
		ModelId:     modelId,
		Params:      params,
		Credentials: credentials,
		ProjectId:   projectId,
		SpaceId:     spaceId,
	}

	return model, nil
}

type GenerateResult struct {
	GeneratedText string `json:"generated_text"`
	StopReason    string `json:"stop_reason"`
}

type generateResponse struct {
	Status     string           `json:"status"`
	StatusCode int              `json:"status_code"`
	Results    []GenerateResult `json:"results"`
}

// GenerateText generates completion text based on a given prompt and parameters (nil to use model params).
func (model *Model) GenerateText(prompt string, params *GenParams) (string, error) {
	// Validate input parameters
	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	if params == nil {
		// Use model params
		params = model.Params
	}

	payload := map[string]interface{}{
		"model_id":   model.ModelId,
		"input":      prompt,
		"parameters": *params,
	}

	if model.ProjectId != "" {
		payload["project_id"] = model.ProjectId
	} else if model.SpaceId != "" {
		payload["space_id"] = model.SpaceId
	}

	var (
		response generateResponse
		err      error
	)
	// Retry request up to 3 times on certain status codes
	for retries := 0; retries < 3; retries++ {
		response, err = model.makeGenerateRequest(payload)
		if err != nil {
			return "", err
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
		return "", fmt.Errorf(fmt.Sprintf("Request failed with: %s (%d)", response.Status, response.StatusCode))
	}

	result := response.Results[0].GeneratedText

	return result, nil
}

// makeGenerateRequest sends the generate request and handles the response using the http package.
func (model *Model) makeGenerateRequest(payload map[string]interface{}) (generateResponse, error) {
	var response generateResponse

	generateTextURL := fmt.Sprintf("%s/%s/text?version=%s", model.Credentials.Url, TEXT_GENERATION_ENDPOINT, VERSION_PARAM)

	resp, err := client.PostRequest(generateTextURL, payload, model.Credentials.AccessToken)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		body, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf(fmt.Sprintf("Request failed with status code %d and error %s", resp.StatusCode, body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}

func (m Model) String() string {
	return fmt.Sprintf(
		"\nModelId:\t%s\nParams:\t%s\nCredentials:\t%s\nProjectId:\t%s\nSpaceId: %s",
		m.ModelId, m.Params, m.Credentials, m.ProjectId, m.SpaceId)
}
