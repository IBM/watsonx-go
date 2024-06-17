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

type StopReason = string

const (
	NotFinished        StopReason = "NOT_FINISHED"  // Possibly more tokens to be streamed
	MaxTokens          StopReason = "MAX_TOKENS"    // Maximum requested tokens reached
	EndOfSequenceToken StopReason = "EOS_TOKEN"     // End of sequence token encountered
	Cancelled          StopReason = "CANCELLED"     // Request canceled by the client
	TimeLimit          StopReason = "TIME_LIMIT"    // Time limit reached
	StopSequence       StopReason = "STOP_SEQUENCE" // Stop sequence encountered
	TokenLimit         StopReason = "TOKEN_LIMIT"   // Token limit reached
	Error              StopReason = "ERROR"         // Error encountered
)

type GenerateTextResult struct {
	Text                string     `json:"generated_text"`
	GeneratedTokenCount int        `json:"generated_token_count"`
	InputTokenCount     int        `json:"input_token_count"`
	StopReason          StopReason `json:"stop_reason"`
}

type GenerateTextPayload struct {
	ProjectID  string           `json:"project_id"`
	Model      string           `json:"model_id"`
	Prompt     string           `json:"input"`
	Parameters *GenerateOptions `json:"parameters,omitempty"`
}

type generateTextResponse struct {
	Status     string               `json:"status"`
	StatusCode int                  `json:"status_code"`
	Results    []GenerateTextResult `json:"results"`
}

// GenerateText generates completion text based on a given prompt and parameters
func (m *Client) GenerateText(model, prompt string, options ...GenerateOption) (GenerateTextResult, error) {
	m.CheckAndRefreshToken()

	if prompt == "" {
		return GenerateTextResult{}, errors.New("prompt cannot be empty")
	}

	opts := &GenerateOptions{}
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	payload := GenerateTextPayload{
		ProjectID:  m.projectID,
		Model:      model,
		Prompt:     prompt,
		Parameters: opts,
	}

	response, err := m.generateTextRequest(payload)
	if err != nil {
		return GenerateTextResult{}, err
	}

	if len(response.Results) == 0 {
		return GenerateTextResult{}, errors.New("no result recieved")
	}

	result := response.Results[0]

	return result, nil
}

// generateTextRequest sends the generate request and handles the response using the http package.
// Returns error on non-2XX response
func (m *Client) generateTextRequest(payload GenerateTextPayload) (generateTextResponse, error) {
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
		return generateTextResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, generateTextURL.String(), bytes.NewBuffer(payloadJSON))
	if err != nil {
		return generateTextResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token.value)

	res, err := m.httpClient.Do(req)
	if err != nil {
		return generateTextResponse{}, err
	}

	statusCode := res.StatusCode

	if statusCode < 200 || statusCode >= 300 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return generateTextResponse{}, fmt.Errorf("request failed with status code %d", statusCode)
		}
		return generateTextResponse{}, fmt.Errorf("request failed with status code %d and error %s", statusCode, body)
	}
	defer res.Body.Close()

	var generateRes generateTextResponse

	if err := json.NewDecoder(res.Body).Decode(&generateRes); err != nil {
		return generateTextResponse{}, err
	}

	return generateRes, nil
}
