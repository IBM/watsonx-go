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

type generatePayload struct {
	projectID  string           `json:"project_id"`
	model      string           `json:"model_id"`
	prompt     string           `json:"input"`
	parameters *GenerateOptions `json:"parameters,omitempty"`
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

	payload := generatePayload{
		projectID:  m.projectID,
		model:      m.modelType,
		prompt:     prompt,
		parameters: opts,
	}

	response, err := m.generateTextRequest(payload)
	if err != nil {
		return "", err
	}

	statusCode := response.StatusCode

	if statusCode < 200 || statusCode >= 300 {
		return "", fmt.Errorf(fmt.Sprintf("Request failed with: %s (%d)", response.Status, statusCode))
	}

	result := response.Results[0].GeneratedText

	return result, nil
}

// generateTextRequest sends the generate request and handles the response using the http package.
func (m *Model) generateTextRequest(payload generatePayload) (generateResponse, error) {
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

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return generateResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		body, _ := io.ReadAll(resp.Body)
		return generateResponse{}, fmt.Errorf(fmt.Sprintf("Request failed with status code %d and error %s", resp.StatusCode, body))
	}

	var response generateResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return generateResponse{}, err
	}

	return response, nil
}
