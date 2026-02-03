package ai

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const systemPrompt = `You are a code modification assistant. When given code and a modification request:
1. Analyze the code carefully
2. Apply the requested changes
3. Return ONLY the modified code, nothing else
4. Do not include explanations, markdown formatting, or code blocks
5. Preserve the original formatting and style as much as possible`

// GeminiProvider implements the Provider interface using Google's Gemini AI
type GeminiProvider struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiProvider creates a new Gemini AI provider
func NewGeminiProvider(ctx context.Context, apiKey string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	return &GeminiProvider{
		client: client,
		model:  model,
	}, nil
}

// ProcessCode sends the code and prompt to Gemini and returns the modified code
func (p *GeminiProvider) ProcessCode(ctx context.Context, prompt string, content []byte) ([]byte, error) {
	// Build the user prompt
	userPrompt := fmt.Sprintf("Code:\n```\n%s\n```\n\nRequest: %s", string(content), prompt)

	// Generate response
	resp, err := p.model.GenerateContent(ctx, genai.Text(userPrompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	result, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from Gemini")
	}

	return []byte(result), nil
}

// Close releases resources held by the Gemini provider
func (p *GeminiProvider) Close() error {
	return p.client.Close()
}
