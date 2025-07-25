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
		wx.NewRetryConfig(
			wx.WithOnRetry(func(n uint, err error) {
				retryCount = n
				log.Printf("Retrying request after error: %v", err)
			}),
		),
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

// TestLegacyRetryWithNoSuccessStatusOnAnyRequest tests the retry mechanism with a server that always returns a 429 status code.
func TestLegacyRetryWithNoSuccessStatusOnAnyRequest(t *testing.T) {
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
		wx.NewRetryConfig(
			wx.WithBackoff(backoffTime),
			wx.WithOnRetry(func(n uint, err error) {
				retryCount = n
				log.Printf("Retrying request after error: %v", err)
			}),
		),
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

// TestRetryWithNoSuccessStatusOnAnyRequest tests the retry mechanism with a server that always returns a 429 status code.
func TestRetryWithNoSuccessStatusOnAnyRequest(t *testing.T) {
	expectedStatusCode := http.StatusTooManyRequests

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatusCode)
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
		wx.NewRetryConfig(
			wx.WithReturnHTTPStatusAsErr(false), // Use new behavior: only return actual network errors
			wx.WithBackoff(backoffTime),
			wx.WithOnRetryV2(func(n uint, resp *http.Response, err error) {
				retryCount = n
				if err != nil {
					t.Errorf("In OnRetry, expected nil, got error: %v", err)
				}

				if resp == nil {
					t.Errorf("In OnRetry, expected non-nil response, got nil")
				}

				if resp != nil && resp.StatusCode != expectedStatusCode {
					t.Errorf("Expected status code %d, got %d", expectedStatusCode, resp.StatusCode)
				}

				log.Printf("Retrying request after response with status code: %d", resp.StatusCode)
			}),
		),
	)

	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)
	expectedMinimumTime := backoffTime * time.Duration(expectedRetries)

	if err != nil {
		t.Errorf("Expected nil, got error: %v", err)
	}

	if resp == nil {
		t.Errorf("Expected non-nil response, got nil")
	}

	if resp != nil && resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code %d, got %d", expectedStatusCode, resp.StatusCode)
	}

	if retryCount != expectedRetries {
		t.Errorf("Expected 3 retries, but got %d", retryCount)
	}

	if elapsedTime < expectedMinimumTime {
		t.Errorf("Expected minimum time of %v, but got %v", expectedMinimumTime, elapsedTime)
	}
}
