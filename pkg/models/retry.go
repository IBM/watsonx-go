package models

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// OnRetryFunc is a function type that is called on each retry attempt.
type OnRetryFunc func(attempt uint, err error)

// Timer interface to abstract time-based operations for retries.
type Timer interface {
	After(time.Duration) <-chan time.Time
}

// RetryIfFunc determines whether a retry should be attempted based on the error.
type RetryIfFunc func(error) bool

// RetryConfig contains configuration options for the retry mechanism.
type RetryConfig struct {
	retries   uint
	backoff   time.Duration
	maxJitter time.Duration
	onRetry   OnRetryFunc
	retryIf   RetryIfFunc
	timer     Timer
	context   context.Context
}

// RetryOption is a function type for modifying RetryConfig options.
type RetryOption func(*RetryConfig)

// timerImpl implements the Timer interface using time.After.
type timerImpl struct{}

func (t timerImpl) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// newDefaultRetryConfig creates a default RetryConfig with sensible defaults.
func newDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		retries:   3,
		backoff:   1 * time.Second,
		maxJitter: 0 * time.Second,                            // no jitter by default
		onRetry:   func(n uint, err error) {},                 // no-op onRetry by default
		retryIf:   func(err error) bool { return err != nil }, // retry on any error by default
		timer:     &timerImpl{},
		context:   context.Background(),
	}
}

// RetryableFuncWithResponse represents a function that returns an HTTP response or an error.
type RetryableFuncWithResponse func() (*http.Response, error)

// Retry retries the provided retryableFunc according to the retry configuration options.
func Retry(retryableFunc RetryableFuncWithResponse, options ...RetryOption) ([]byte, error) {
	opts := newDefaultRetryConfig()

	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	var lastErr error
	for n := uint(0); n < opts.retries; n++ {
		if err := opts.context.Err(); err != nil {
			return nil, err
		}

		resp, err := retryableFunc()
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return body, nil
		}

		if err == nil && resp != nil {
			err = errors.New(resp.Status)
		}

		if !opts.retryIf(err) {
			return nil, err
		}

		lastErr = err
		opts.onRetry(n+1, err)

		backoffDuration := opts.backoff
		if opts.maxJitter > 0 {
			jitter := time.Duration(rand.Int63n(int64(opts.maxJitter)))
			backoffDuration += jitter
		}

		select {
		case <-opts.timer.After(backoffDuration):
		case <-opts.context.Done():
			return nil, opts.context.Err()
		}
	}

	return nil, lastErr
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

// WithRetryIf sets the condition to determine whether to retry based on the error.
func WithRetryIf(retryIf RetryIfFunc) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.retryIf = retryIf
	}
}

// WithContext sets the context for the retry configuration.
func WithContext(ctx context.Context) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.context = ctx
	}
}
