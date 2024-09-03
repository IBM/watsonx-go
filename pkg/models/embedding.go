package models

import (
	"bytes"
	"encoding/json"
	"errors"
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

// EmbedDocuments embeds the given texts using the specified model.
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

// EmbedQuery embeds the given text using the specified model.
func (m *Client) EmbedQuery(model string, text string, options ...EmbeddingOption) (EmbeddingResponse, error) {
	return m.EmbedDocuments(model, []string{text}, options...)
}

// generateEmbeddingRequest sends a request to the embedding endpoint with the given payload.
// return the response from the server if and only if the request is successful, code 200.
func (m *Client) generateEmbeddingRequest(payload EmbeddingPayload) (embeddingResponse, error) {
	embeddingUrl := m.generateUrlFromEndpoint(EmbeddingEndpoint)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return embeddingResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, embeddingUrl, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return embeddingResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.token.value)

	res, err := m.httpClient.DoWithRetry(req)
	if err != nil {
		return embeddingResponse{}, err
	}
	defer res.Body.Close()

	var embeddingRes embeddingResponse

	if err := json.NewDecoder(res.Body).Decode(&embeddingRes); err != nil {
		return embeddingResponse{}, err
	}

	return embeddingRes, nil
}
