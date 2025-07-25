package models

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

// OnRetryFunc is a function type that is called on each retry attempt.
type OnRetryFunc func(attempt uint, err error)

// OnRetryV2Func is a function type that is called on each retry attempt with response.
type OnRetryV2Func func(attempt uint, resp *http.Response, err error)

// Timer interface to abstract time-based operations for retries.
type Timer interface {
	After(time.Duration) <-chan time.Time
}

// RetryIfFunc determines whether a retry should be attempted based on the error.
type RetryIfFunc func(error) bool

// RetryIfFuncV2 determines whether a retry should be attempted based on the response.
type RetryIfFuncV2 func(*http.Response, error) bool

// RetryConfig contains configuration options for the retry mechanism.
type RetryConfig struct {
	retries               uint
	backoff               time.Duration
	maxJitter             time.Duration
	onRetry               OnRetryFunc   // Legacy callback for retries
	onRetryV2             OnRetryV2Func // Callback for retries with response
	retryIf               RetryIfFunc   // Legacy error-based retry function
	retryIfV2             RetryIfFuncV2 // New response-based retry function
	timer                 Timer
	context               context.Context
	returnHTTPStatusAsErr bool // When true, use legacy behavior: convert HTTP status to errors
}

// RetryOption is a function type for modifying RetryConfig options.
type RetryOption func(*RetryConfig)

// timerImpl implements the Timer interface using time.After.
type timerImpl struct{}

func (t timerImpl) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// NewDefaultRetryConfig creates a default RetryConfig with sensible defaults.
func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		retries:   3,
		backoff:   1 * time.Second,
		maxJitter: 1 * time.Second,
		onRetry:   func(n uint, err error) {},                      // no-op onRetry by default
		onRetryV2: func(n uint, resp *http.Response, err error) {}, // no-op onRetry by default
		retryIf:   func(err error) bool { return err != nil },      // retry on any error by default (legacy)
		retryIfV2: func(resp *http.Response, err error) bool {
			return err != nil || (resp != nil && resp.StatusCode >= http.StatusBadRequest)
		}, // retry on any error or 4xx/5xx response by default (new)
		timer:                 &timerImpl{},
		context:               context.Background(),
		returnHTTPStatusAsErr: true, // Legacy behavior: convert HTTP status to errors
	}
}

func NewRetryConfig(options ...RetryOption) *RetryConfig {
	opts := NewDefaultRetryConfig()

	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}
	return &RetryConfig{
		retries:               opts.retries,
		backoff:               opts.backoff,
		maxJitter:             opts.maxJitter,
		onRetry:               opts.onRetry,
		onRetryV2:             opts.onRetryV2,
		retryIf:               opts.retryIf,
		retryIfV2:             opts.retryIfV2,
		timer:                 opts.timer,
		context:               opts.context,
		returnHTTPStatusAsErr: opts.returnHTTPStatusAsErr,
	}
}

// RetryableFuncWithResponse represents a function that returns an HTTP response or an error.
type RetryableFuncWithResponse func() (*http.Response, error)

// Retry retries the provided retryableFunc according to the retry configuration options.
func Retry(retryableFunc RetryableFuncWithResponse, retryConfig *RetryConfig) (*http.Response, error) {
	var lastErr error
	var lastResp *http.Response

	for n := uint(0); n < retryConfig.retries; n++ {
		if err := retryConfig.context.Err(); err != nil {
			return nil, err
		}

		resp, err := retryableFunc()

		// If the response is successful, return it immediately
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		// Store the response and error for potential return
		lastResp = resp
		lastErr = err
		statusAsErr := errors.New(resp.Status)
		errResult := err

		// Determine if we should retry based on error or response
		retryIfV1Err := err
		shouldRetry := false
		if err == nil && resp != nil {
			retryIfV1Err = statusAsErr
			// set errResult and lastErr based on behavior flag
			if retryConfig.returnHTTPStatusAsErr {
				errResult = statusAsErr
				lastErr = statusAsErr
			}
		}

		shouldRetry = retryConfig.retryIf(retryIfV1Err) || retryConfig.retryIfV2(resp, err)
		if !shouldRetry {
			return resp, errResult
		}

		retryConfig.onRetry(n+1, statusAsErr)
		retryConfig.onRetryV2(n+1, resp, err)

		// Apply backoff and jitter
		backoffDuration := retryConfig.backoff
		if retryConfig.maxJitter > 0 {
			jitter := time.Duration(rand.Int63n(int64(retryConfig.maxJitter)))
			backoffDuration += jitter
		}

		select {
		case <-retryConfig.timer.After(backoffDuration):
		case <-retryConfig.context.Done():
			return nil, retryConfig.context.Err()
		}
	}

	if retryConfig.returnHTTPStatusAsErr {
		// Legacy behavior: convert HTTP status to error if no network error
		return nil, lastErr
	}

	// New behavior: only return actual network errors
	return lastResp, lastErr
}

// WithRetries sets the number of retries for the retry configuration.
func WithRetries(retries uint) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.retries = retries
	}
}

// WithBackoff sets the backoff duration between retries.
func WithBackoff(backoff time.Duration) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.backoff = backoff
	}
}

// WithMaxJitter sets the maximum jitter duration to add to the backoff.
func WithMaxJitter(maxJitter time.Duration) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.maxJitter = maxJitter
	}
}

// WithOnRetry sets the callback function to execute on each retry.
func WithOnRetry(onRetry OnRetryFunc) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.onRetry = onRetry
	}
}

// WithOnRetry sets the callback function to execute on each retry.
func WithOnRetryV2(onRetry OnRetryV2Func) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.onRetryV2 = onRetry
	}
}

// WithRetryIf sets the condition to determine whether to retry based on the error.
func WithRetryIf(retryIf RetryIfFunc) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.retryIf = retryIf
	}
}

// WithRetryIfV2 sets the condition to determine whether to retry based on the response.
// This enables the new response-based retry logic and automatically enables the
// new behavior flags.
func WithRetryIfV2(retryIf RetryIfFuncV2) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.retryIf = func(err error) bool { return false } // Disable legacy retryIf
		cfg.onRetry = func(n uint, err error) {}            // Disable legacy onRetry
		cfg.retryIfV2 = retryIf
		cfg.returnHTTPStatusAsErr = false // Enable new behavior by default when using V2
	}
}

// WithReturnHTTPStatusAsErr controls the legacy behavior where HTTP status codes
// are converted to Go errors. When enabled (true, default), HTTP status codes
// are converted to errors (legacy behavior). When disabled (false), only actual
// network/request errors are returned as errors, while HTTP responses (even with
// 4xx/5xx status) return nil error (new correct behavior).
func WithReturnHTTPStatusAsErr(enabled bool) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.returnHTTPStatusAsErr = enabled
	}
}

// Custom wrapper for http.Client that implements the Doer interface.
// - Do
// - DoWithRetry
type HttpClient struct {
	httpClient  *http.Client
	retryConfig *RetryConfig
}

type HttpClientConfig struct {
	retryConfig *RetryConfig
}

func newDefaultHttpClientConfig() *HttpClientConfig {
	return &HttpClientConfig{
		retryConfig: NewDefaultRetryConfig(),
	}
}

// HttpClientConfigOption is a function type for modifying HttpClientConfig options.
type HttpClientConfigOption func(*HttpClientConfig)

func NewHttpClient(options ...HttpClientConfigOption) *HttpClient {

	opts := newDefaultHttpClientConfig()

	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	return &HttpClient{
		httpClient:  &http.Client{},
		retryConfig: opts.retryConfig,
	}
}

func WithRetryConfig(config *RetryConfig) HttpClientConfigOption {
	return func(cfg *HttpClientConfig) {
		cfg.retryConfig = config
	}
}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *HttpClient) DoWithRetry(req *http.Request) (*http.Response, error) {
	return Retry(
		func() (*http.Response, error) {
			return c.httpClient.Do(req)
		},
		c.retryConfig,
	)
}
