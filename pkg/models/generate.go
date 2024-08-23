package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	GenerationEndpoint         string = "/ml/v1/text"
	GenerateTextEndpoint       string = GenerationEndpoint + "/generation"
	GenerateTextStreamEndpoint string = GenerationEndpoint + "/generation_stream"
)

type StopReason = string

const (
	NotFinished        StopReason = "not_finished"  // Possibly more tokens to be streamed
	MaxTokens          StopReason = "max_tokens"    // Maximum requested tokens reached
	EndOfSequenceToken StopReason = "eos_token"     // End of sequence token encountered
	Cancelled          StopReason = "cancelled"     // Request canceled by the client
	TimeLimit          StopReason = "time_limit"    // Time limit reached
	StopSequence       StopReason = "stop_sequence" // Stop sequence encountered
	TokenLimit         StopReason = "token_limit"   // Token limit reached
	Error              StopReason = "error"         // Error encountered
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
	textUrl := m.generateUrlFromEndpoint(GenerateTextEndpoint)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return generateTextResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, textUrl, bytes.NewBuffer(payloadJSON))
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

// GenerateTextStream generates completion text channel (stream) based on a given prompt and parameters
func (m *Client) GenerateTextStream(model, prompt string, options ...GenerateOption) (<-chan GenerateTextResult, error) {
	dataChan := make(chan GenerateTextResult)

	if prompt == "" {
		close(dataChan)
		return dataChan, errors.New("prompt cannot be empty")
	}

	go func() {
		defer close(dataChan)

		m.CheckAndRefreshToken()

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

		responseChan, _ := m.generateTextStreamRequest(payload)

		for data := range responseChan {
			for _, result := range data.Results {
				dataChan <- result
			}
		}
	}()

	return dataChan, nil
}

// generateTextStreamRequest sends the generate request and handles the response using the http package.
// Closes the channel on non-200 response
// If any error happens during the streaming, it will be logged and the channel will be closed
func (m *Client) generateTextStreamRequest(payload GenerateTextPayload) (<-chan generateTextResponse, error) {
	dataChan := make(chan generateTextResponse)

	go func() {
		defer close(dataChan)

		streamUrl := m.generateUrlFromEndpoint(GenerateTextStreamEndpoint)

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling payload: ", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, streamUrl, bytes.NewBuffer(payloadJSON))
		if err != nil {
			log.Println("error creating request: ", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.token.value)
		req.Header.Set("Accept", "text/event-stream")

		res, err := m.httpClient.Do(req)
		if err != nil {
			log.Println("error making request: ", err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("request failed with status code %d", res.StatusCode)
			} else {
				log.Printf("request failed with status code %d and error %s", res.StatusCode, body)
			}
			return
		}

		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := line[6:]
			var generation generateTextResponse

			if err := json.Unmarshal([]byte(data), &generation); err != nil {
				log.Println("error unmarshalling data: ", err)
				return
			}
			dataChan <- generation
		}
	}()

	return dataChan, nil
}
