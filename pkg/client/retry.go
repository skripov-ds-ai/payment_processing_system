package client

import "time"

type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

func NewRetryConfig(maxAttempts int, delay time.Duration) *RetryConfig {
	return &RetryConfig{
		MaxAttempts: maxAttempts,
		Delay:       delay,
	}
}
