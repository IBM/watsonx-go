package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wx "github.com/IBM/watsonx-go/pkg/models"
)

// TestRetryWithSuccessOnFirstRequest tests the retry mechanism with a server that always returns a 200 status code.
func TestRetryWithSuccessOnFirstRequest(t *testing.T) {
	type ResponseType struct {
		Content string `json:"content"`
		Status  int    `json:"status"`
	}

	expectedResponse := ResponseType{Content: "success"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"content":"success"}`))
	}))
	defer server.Close()

	var retryCount uint = 0
	var expectedRetries uint = 0

	sendRequest := func() (*http.Response, error) {
		return http.Get(server.URL + "/success")
	}

	resp, err := wx.Retry(
		sendRequest,
		wx.WithOnRetry(func(n uint, err error) {
			retryCount = n
			log.Printf("Retrying request after error: %v", err)
		}),
	)

	if err != nil {
		t.Errorf("Expected nil, got error: %v", err)
	}

	if retryCount != expectedRetries {
		t.Errorf("Expected 0 retries, but got %d", retryCount)
	}

	defer resp.Body.Close()
	var response ResponseType
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	if response != expectedResponse {
		t.Errorf("Expected response %v, but got %v", expectedResponse, response)
	}
}

// TestRetryWithNoSuccessStatusOnAnyRequest tests the retry mechanism with a server that always returns a 429 status code.
func TestRetryWithNoSuccessStatusOnAnyRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	var backoffTime = 2 * time.Second
	var retryCount uint = 0
	var expectedRetries uint = 3

	sendRequest := func() (*http.Response, error) {
		return http.Get(server.URL + "/notfound")
	}

	startTime := time.Now()

	resp, err := wx.Retry(
		sendRequest,
		wx.WithBackoff(backoffTime),
		wx.WithOnRetry(func(n uint, err error) {
			retryCount = n
			log.Printf("Retrying request after error: %v", err)
		}),
	)

	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)
	expectedMinimumTime := backoffTime * time.Duration(expectedRetries)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	if resp != nil {
		defer resp.Body.Close()
		t.Errorf("Expected nil response, got %v", resp.Body)
	}

	if retryCount != expectedRetries {
		t.Errorf("Expected 3 retries, but got %d", retryCount)
	}

	if elapsedTime < expectedMinimumTime {
		t.Errorf("Expected minimum time of %v, but got %v", expectedMinimumTime, elapsedTime)
	}
}

// TestRetryReturnsWatsonxErrorWithDetails validates that non-200 responses
// are converted into structured WatsonxError using DecodeWatsonxError.
func TestRetryReturnsWatsonxErrorWithDetails(t *testing.T) {
	errorResponse := `{
		"errors": [{
			"code": "invalid_request",
			"message": "Invalid input parameter",
			"more_info": "Check request payload"
		}],
		"trace": "trace-id-123"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorResponse))
	}))
	defer server.Close()

	sendRequest := func() (*http.Response, error) {
		return http.Get(server.URL)
	}

	resp, err := wx.Retry(sendRequest)

	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}

	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	wxErr, ok := err.(*wx.WatsonxError)
	if !ok {
		t.Fatalf("Expected error type *WatsonxError, got %T", err)
	}

	if wxErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, wxErr.StatusCode)
	}

	if len(wxErr.Errors) != 1 {
		t.Fatalf("Expected 1 error detail, got %d", len(wxErr.Errors))
	}

	if wxErr.Errors[0].Code != "invalid_request" {
		t.Errorf("Unexpected error code: %s", wxErr.Errors[0].Code)
	}

	if wxErr.Trace != "trace-id-123" {
		t.Errorf("Expected trace ID, got %s", wxErr.Trace)
	}
}

// TestRetryReturnsWatsonxErrorWithEmptyBody validates handling of
// non-200 responses with an empty response body.
func TestRetryReturnsWatsonxErrorWithEmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sendRequest := func() (*http.Response, error) {
		return http.Get(server.URL)
	}

	resp, err := wx.Retry(sendRequest)

	if resp != nil {
		t.Errorf("Expected nil response, got %v", resp)
	}

	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	wxErr, ok := err.(*wx.WatsonxError)
	if !ok {
		t.Fatalf("Expected error type *WatsonxError, got %T", err)
	}

	if wxErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status %d, got %d", http.StatusInternalServerError, wxErr.StatusCode)
	}

	if len(wxErr.Errors) != 0 {
		t.Errorf("Expected no error details, got %d", len(wxErr.Errors))
	}
}
