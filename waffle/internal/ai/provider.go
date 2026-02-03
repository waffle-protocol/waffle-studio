package ai

import "context"

// Provider defines the interface for AI code processing providers
type Provider interface {
	// ProcessCode takes a prompt and code content, returns the modified code
	ProcessCode(ctx context.Context, prompt string, content []byte) ([]byte, error)

	// Close releases any resources held by the provider
	Close() error
}
