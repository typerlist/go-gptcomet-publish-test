package llm

import (
	"fmt"
	"strings"

	"github.com/belingud/go-gptcomet/pkg/types"
)

// AzureLLM implements the LLM interface for Azure OpenAI
type AzureLLM struct {
	*OpenAILLM
}

// NewAzureLLM creates a new AzureLLM
func NewAzureLLM(config *types.ClientConfig) *AzureLLM {
	if config.Model == "" {
		config.Model = "gpt-35-turbo"
	}
	// Azure requires a specific API version
	if !strings.Contains(config.APIBase, "api-version=") {
		config.APIBase = fmt.Sprintf("%s?api-version=2023-12-01-preview", config.APIBase)
	}

	return &AzureLLM{
		OpenAILLM: NewOpenAILLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (a *AzureLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "",
			PromptMessage: "Enter Azure OpenAI endpoint",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "gpt-35-turbo",
			PromptMessage: "Enter deployment name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}

// BuildHeaders builds request headers for Azure OpenAI
func (a *AzureLLM) BuildHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
		"api-key":     a.Config.APIKey,
	}
	for k, v := range a.Config.ExtraHeaders {
		headers[k] = v
	}
	return headers
}
