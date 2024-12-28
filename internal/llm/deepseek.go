package llm

import (
	"github.com/belingud/go-gptcomet/pkg/types"
)

// DeepSeekLLM implements the LLM interface for DeepSeek
type DeepSeekLLM struct {
	*OpenAILLM
}

// NewDeepSeekLLM creates a new DeepSeekLLM
func NewDeepSeekLLM(config *types.ClientConfig) *DeepSeekLLM {
	if config.APIBase == "" {
		config.APIBase = "https://api.deepseek.com/v1"
	}
	if config.Model == "" {
		config.Model = "deepseek-chat"
	}

	return &DeepSeekLLM{
		OpenAILLM: NewOpenAILLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (d *DeepSeekLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "https://api.deepseek.com/v1",
			PromptMessage: "Enter DeepSeek API base",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "deepseek-chat",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}
