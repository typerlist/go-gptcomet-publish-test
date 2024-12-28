package llm

import (
	"github.com/belingud/go-gptcomet/pkg/types"
)

// ChatGLMLLM implements the LLM interface for ChatGLM
type ChatGLMLLM struct {
	*OpenAILLM
}

// NewChatGLMLLM creates a new ChatGLMLLM
func NewChatGLMLLM(config *types.ClientConfig) *ChatGLMLLM {
	if config.APIBase == "" {
		config.APIBase = "https://open.bigmodel.cn/api/paas/v4"
	}
	if config.Model == "" {
		config.Model = "chatglm_turbo"
	}

	return &ChatGLMLLM{
		OpenAILLM: NewOpenAILLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (c *ChatGLMLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "https://open.bigmodel.cn/api/paas/v4",
			PromptMessage: "Enter ChatGLM API base",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "chatglm_turbo",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}
