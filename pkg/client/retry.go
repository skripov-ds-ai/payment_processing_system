package client

import "time"

// RetryConfig determines config for connection retry
type RetryConfig struct {
	MaxAttempts int
	PingDelay   time.Duration
	Delay       time.Duration
}

// NewRetryConfig is constructor for retry config
func NewRetryConfig(maxAttempts int, delay time.Duration) *RetryConfig {
	return &RetryConfig{
		MaxAttempts: maxAttempts,
		Delay:       delay,
	}
}

// Retry for database sql connections
func Retry(fn func() error, retryCfg *RetryConfig) error {
	var err error
	maxAttempts := retryCfg.MaxAttempts
	for maxAttempts > 0 {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(retryCfg.Delay)
		maxAttempts--
	}
	return err
}
