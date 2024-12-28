package llm

import (
	"context"
	"net/http"

	"github.com/belingud/go-gptcomet/pkg/types"
)

// OllamaLLM implements the LLM interface for Ollama
type OllamaLLM struct {
	*BaseLLM
}

// NewOllamaLLM creates a new OllamaLLM
func NewOllamaLLM(config *types.ClientConfig) *OllamaLLM {
	if config.APIBase == "" {
		config.APIBase = "http://localhost:11434/api"
	}
	if config.CompletionPath == "" {
		config.CompletionPath = "chat"
	}
	if config.AnswerPath == "" {
		config.AnswerPath = "message.content"
	}
	if config.Model == "" {
		config.Model = "llama2"
	}

	return &OllamaLLM{
		BaseLLM: NewBaseLLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (o *OllamaLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "http://localhost:11434/api",
			PromptMessage: "Enter Ollama API base",
		},
		"model": {
			DefaultValue:  "llama2",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}

// FormatMessages formats messages for Ollama API
func (o *OllamaLLM) FormatMessages(message string, history []types.Message) (interface{}, error) {
	messages := make([]map[string]string, 0, len(history)+1)
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": message,
	})

	payload := map[string]interface{}{
		"model":    o.Config.Model,
		"messages": messages,
		"options":  map[string]interface{}{},
	}

	if o.Config.Temperature != 0 {
		payload["options"].(map[string]interface{})["temperature"] = o.Config.Temperature
	}
	if o.Config.TopP != 0 {
		payload["options"].(map[string]interface{})["top_p"] = o.Config.TopP
	}

	return payload, nil
}

// BuildHeaders builds request headers
func (o *OllamaLLM) BuildHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	for k, v := range o.Config.ExtraHeaders {
		headers[k] = v
	}
	return headers
}

// GetUsage returns usage information for the provider
func (o *OllamaLLM) GetUsage(data []byte) (string, error) {
	// Ollama doesn't provide token usage information
	return "", nil
}

// MakeRequest makes a request to the API
func (o *OllamaLLM) MakeRequest(ctx context.Context, client *http.Client, message string, history []types.Message) (string, error) {
	return o.BaseLLM.MakeRequest(ctx, client, o, message, history)
}
