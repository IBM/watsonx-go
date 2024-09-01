package models

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

type OnRetryFunc func(attempt uint, err error)

// Timer to track time for a retry
type Timer interface {
	After(time.Duration) <-chan time.Time
}

type RetryIfFunc func(error) bool

type RetryConfig struct {
	retries   uint
	backoff   time.Duration
	maxJitter time.Duration
	onRetry   OnRetryFunc
	retryIf   RetryIfFunc
	timer     Timer
	context   context.Context
}

type RetryOption func(*RetryConfig)

type timerImpl struct{}

func (t timerImpl) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func newDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		retries:   3,
		backoff:   1 * time.Second,
		maxJitter: 0 * time.Second, // no jitter by default
		onRetry:   func(n uint, err error) {},
		retryIf:   func(err error) bool { return err != nil }, // default to retry on any error
		timer:     &timerImpl{},
		context:   context.Background(),
	}
}

type RetryableFunc func() error

type RetryableFuncWithResponse func() (*http.Response, error)

// Retry retries the provided retryableFunc according to the retry configuration options.
func Retry(retryableFunc RetryableFuncWithResponse, options ...RetryOption) error {
	opts := newDefaultRetryConfig()

	for _, opt := range options {
		if opt != nil {
			opt(opts)
		}
	}

	var lastErr error
	for n := uint(0); n < opts.retries; n++ {
		if err := opts.context.Err(); err != nil {
			return err
		}

		resp, err := retryableFunc()
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		if err == nil {
			err = errors.New(resp.Status)
		}

		if !opts.retryIf(err) {
			return err
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
			return opts.context.Err()
		}
	}

	return lastErr
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
