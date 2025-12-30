package llm

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

// AnthropicProvider implements Provider using Anthropic's API directly
type AnthropicProvider struct {
	llm     *anthropic.LLM
	modelID string
}

// AnthropicConfig holds configuration for Anthropic provider
type AnthropicConfig struct {
	APIKey  string // Anthropic API key (required)
	ModelID string // Model ID, defaults to claude-3-haiku-20240307
}

// NewAnthropicProvider creates a new Anthropic API provider
func NewAnthropicProvider(cfg AnthropicConfig) (*AnthropicProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("Anthropic API key is required")
	}
	if cfg.ModelID == "" {
		cfg.ModelID = "claude-haiku-4-5-20251001"
	}

	llm, err := anthropic.New(
		anthropic.WithToken(cfg.APIKey),
		anthropic.WithModel(cfg.ModelID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Anthropic LLM: %w", err)
	}

	return &AnthropicProvider{
		llm:     llm,
		modelID: cfg.ModelID,
	}, nil
}

// GenerateText implements Provider
func (p *AnthropicProvider) GenerateText(ctx context.Context, prompt string) (string, error) {
	response, err := llms.GenerateFromSinglePrompt(ctx, p.llm, prompt,
		llms.WithMaxTokens(1024),
		llms.WithTemperature(0.3),
	)
	if err != nil {
		return "", fmt.Errorf("anthropic generation failed: %w", err)
	}
	return response, nil
}

// Name implements Provider
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// ModelID returns the model being used
func (p *AnthropicProvider) ModelID() string {
	return p.modelID
}
