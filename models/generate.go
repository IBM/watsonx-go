package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	GenerationEndpoint   string = "/ml/v1-beta/generation"
	GenerateTextEndpoint string = GenerationEndpoint + "/text"
)

type GenerateResult struct {
	GeneratedText string `json:"generated_text"`
	StopReason    string `json:"stop_reason"`
}

type GeneratePayload struct {
	ProjectID  string           `json:"project_id"`
	Model      string           `json:"model_id"`
	Prompt     string           `json:"input"`
	Parameters *GenerateOptions `json:"parameters,omitempty"`
}

type generateResponse struct {
	Status     string           `json:"status"`
	StatusCode int              `json:"status_code"`
	Results    []GenerateResult `json:"results"`
}

// GenerateText generates completion text based on a given prompt and parameters
func (m *Model) GenerateText(prompt string, options ...GenerateOption) (string, error) {
	m.CheckAndRefreshToken()

	if prompt == "" {
		return "", errors.New("prompt cannot be empty")
	}

	opts := &GenerateOptions{}
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	payload := GeneratePayload{
		ProjectID:  m.projectID,
		Model:      m.modelType,
		Prompt:     prompt,
		Parameters: opts,
	}

	response, err := m.generateTextRequest(payload)
	if err != nil {
		return "", err
	}

	result := response.Results[0].GeneratedText

	return result, nil
}

// generateTextRequest sends the generate request and handles the response using the http package.
// Returns error on non-2XX response
func (m *Model) generateTextRequest(payload GeneratePayload) (generateResponse, error) {
	params := url.Values{
		"version": {m.apiVersion},
	}

	generateTextURL := url.URL{
		Scheme:   "https",
		Host:     m.url,
		Path:     GenerateTextEndpoint,
		RawQuery: params.Encode(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return generateResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, generateTextURL.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return generateResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token.value)

	res, err := m.httpClient.Do(req)
	if err != nil {
		return generateResponse{}, err
	}

	statusCode := res.StatusCode

	if statusCode < 200 || statusCode >= 300 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return generateResponse{}, fmt.Errorf("request failed with status code %d", statusCode)
		}
		return generateResponse{}, fmt.Errorf("request failed with status code %d and error %s", statusCode, body)
	}
	defer res.Body.Close()

	var generateRes generateResponse

	if err := json.NewDecoder(res.Body).Decode(&generateRes); err != nil {
		return generateResponse{}, err
	}

	return generateRes, nil
}
