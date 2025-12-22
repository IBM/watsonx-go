package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WatsonxError represents a structured WatsonX API error
type WatsonxError struct {
	StatusCode int
	Errors     []ErrorDetail
	Trace      string
}

// Error implements the error interface
func (e *WatsonxError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf(
			"watsonx error (%d): %s - %s",
			e.StatusCode,
			e.Errors[0].Code,
			e.Errors[0].Message,
		)
	}
	return fmt.Sprintf("watsonx error (%d)", e.StatusCode)
}

// WatsonxErrorResponse represents the error response structure from Watson X API
type WatsonxErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`
	Trace  string        `json:"trace"`
}

// ErrorDetail represents individual error information
type ErrorDetail struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
}

// DecodeWatsonxError attempts to parse an HTTP error response and return a structured watsonx error
func DecodeWatsonxError(resp *http.Response) error {
	if resp == nil {
		return &WatsonxError{}
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &WatsonxError{
			StatusCode: resp.StatusCode,
		}
	}

	// Restore body so it can be read again
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	wxErr := &WatsonxError{
		StatusCode: resp.StatusCode,
	}

	// Empty body → status-only error
	if len(body) == 0 {
		return wxErr
	}

	// Try to parse WatsonX error schema
	var apiErr WatsonxErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// Unknown / non-JSON error
		return wxErr
	}

	wxErr.Errors = apiErr.Errors
	wxErr.Trace = apiErr.Trace

	return wxErr
}
