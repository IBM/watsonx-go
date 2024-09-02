package test

import (
	"encoding/json"
	wx "github.com/IBM/watsonx-go/pkg/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	defer resp.Body.Close()

	if err != nil {
		t.Errorf("Expected nil, got error: %v", err)
	}

	if retryCount != expectedRetries {
		t.Errorf("Expected 0 retries, but got %d", retryCount)
	}

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
