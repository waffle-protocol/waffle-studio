package cli

import (
	"context"
	"time"
)

// Timeout constants for different operations
const (
	DefaultTimeout = 60 * time.Second
	QuickTimeout   = 10 * time.Second
	LongTimeout    = 120 * time.Second
)

// WithTimeout creates a context with the default timeout
func WithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// WithQuickTimeout creates a context with a quick timeout (10s)
func WithQuickTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), QuickTimeout)
}

// WithLongTimeout creates a context with a long timeout (120s)
func WithLongTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), LongTimeout)
}

// WithCustomTimeout creates a context with a custom timeout duration
func WithCustomTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
