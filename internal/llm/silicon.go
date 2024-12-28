package llm

import (
	"github.com/belingud/go-gptcomet/pkg/config"
	"github.com/belingud/go-gptcomet/pkg/types"
)

// SiliconLLM implements the LLM interface for Silicon
type SiliconLLM struct {
	*OpenAILLM
}

// NewSiliconLLM creates a new SiliconLLM
func NewSiliconLLM(config *types.ClientConfig) *SiliconLLM {
	if config.APIBase == "" {
		config.APIBase = "https://api.silicon.ai/v1"
	}
	if config.Model == "" {
		config.Model = "Qwen/Qwen2.5-7B-Instruct"
	}

	return &SiliconLLM{
		OpenAILLM: NewOpenAILLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (s *SiliconLLM) GetRequiredConfig() map[string]config.ConfigRequirement {
	return map[string]config.ConfigRequirement{
		"api_base": {
			DefaultValue:  "https://api.silicon.ai/v1",
			PromptMessage: "Enter Silicon API base",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "Qwen/Qwen2.5-7B-Instruct",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}
