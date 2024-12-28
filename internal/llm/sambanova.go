package llm

import (
	"github.com/belingud/go-gptcomet/pkg/types"
)

// SambanovaLLM implements the LLM interface for SambaNova
type SambanovaLLM struct {
	*OpenAILLM
}

// NewSambanovaLLM creates a new SambanovaLLM
func NewSambanovaLLM(config *types.ClientConfig) *SambanovaLLM {
	if config.APIBase == "" {
		config.APIBase = "https://api.sambanova.ai/v1"
	}
	if config.Model == "" {
		config.Model = "sambanova-gpt"
	}

	return &SambanovaLLM{
		OpenAILLM: NewOpenAILLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (s *SambanovaLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "https://api.sambanova.ai/v1",
			PromptMessage: "Enter SambaNova API base",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "sambanova-gpt",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}
