package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

// ParseErrorResponse attempts to parse an HTTP error response and return a formatted error
func ParseErrorResponse(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("nil response")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("watsonx error (%d): failed to read error response", resp.StatusCode)
	}

	// Restore the body so it can be read again if needed
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	// If body is empty, return status-based error
	if len(body) == 0 {
		return fmt.Errorf("watsonx error (%d): %s", resp.StatusCode, resp.Status)
	}

	// Try to parse as WatsonxErrorResponse
	var wxErr WatsonxErrorResponse
	if err := json.Unmarshal(body, &wxErr); err != nil {
		// If parsing fails, return generic error with status code and body
		return fmt.Errorf("watsonx error (%d): %s", resp.StatusCode, string(body))
	}

	// Format: watsonx error (404): model_not_supported - Model 'xyz' was not found
	if len(wxErr.Errors) > 0 {
		firstError := wxErr.Errors[0]
		return fmt.Errorf("watsonx error (%d): %s - %s", resp.StatusCode, firstError.Code, firstError.Message)
	}

	// Fallback if no errors array
	return fmt.Errorf("watsonx error (%d): %s", resp.StatusCode, resp.Status)
}
