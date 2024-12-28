package llm

import (
	"github.com/belingud/go-gptcomet/pkg/types"
)

// XAILLM implements the LLM interface for XAI
type XAILLM struct {
	*OpenAILLM
}

// NewXAILLM creates a new XAILLM
func NewXAILLM(config *types.ClientConfig) *XAILLM {
	if config.APIBase == "" {
		config.APIBase = "https://api.x.ai/v1"
	}

	return &XAILLM{
		OpenAILLM: NewOpenAILLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (x *XAILLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "https://api.x.ai/v1",
			PromptMessage: "Enter XAI API base",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "grok-beta",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}
