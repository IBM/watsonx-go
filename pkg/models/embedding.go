package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	EmbeddingEndpoint string = "/ml/v1/text/embeddings"
)

type EmbeddingPayload struct {
	ProjectID  string            `json:"project_id"`
	Model      string            `json:"model_id"`
	Inputs     []string          `json:"inputs"`
	Parameters *EmbeddingOptions `json:"parameters,omitempty"`
}

type EmbeddingResponse struct {
	Model           string            `json:"model_id"`
	Results         []EmbeddingResult `json:"results"`
	CreatedAt       time.Time         `json:"created_at"`
	InputTokenCount int               `json:"input_token_count"`
}

type EmbeddingResult struct {
	Embedding []float64 `json:"embedding"`
	Input     string    `json:"input,omitempty"`
}

type embeddingResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	EmbeddingResponse
}

func (m *Client) EmbedDocuments(model string, texts []string, options ...EmbeddingOption) (EmbeddingResponse, error) {
	m.CheckAndRefreshToken()

	opts := &EmbeddingOptions{}
	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	payload := EmbeddingPayload{
		ProjectID:  m.projectID,
		Model:      model,
		Inputs:     texts,
		Parameters: opts,
	}

	response, err := m.generateEmbeddingRequest(payload)
	if err != nil {
		return EmbeddingResponse{}, err
	}

	if len(response.Results) == 0 {
		return EmbeddingResponse{}, errors.New("no result received")
	}

	return response.EmbeddingResponse, nil
}

func (m *Client) EmbedQuery(model string, text string, options ...EmbeddingOption) (EmbeddingResponse, error) {
	return m.EmbedDocuments(model, []string{text}, options...)
}

func (m *Client) generateEmbeddingRequest(payload EmbeddingPayload) (embeddingResponse, error) {
	embeddingUrl := m.generateUrlFromEndpoint(EmbeddingEndpoint)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return embeddingResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, embeddingUrl, bytes.NewBuffer(payloadJSON))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token.value)

	res, err := m.httpClient.Do(req)
	if err != nil {
		return embeddingResponse{}, err
	}

	statusCode := res.StatusCode
	if statusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return embeddingResponse{}, fmt.Errorf("request failed with status code %d", statusCode)
		}
		return embeddingResponse{}, fmt.Errorf("request failed with status code %d and error %s", statusCode, body)
	}
	defer res.Body.Close()

	var embeddingRes embeddingResponse

	if err := json.NewDecoder(res.Body).Decode(&embeddingRes); err != nil {
		return embeddingResponse{}, err
	}

	return embeddingRes, nil
}
