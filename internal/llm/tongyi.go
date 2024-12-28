package llm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/belingud/go-gptcomet/pkg/types"
	"github.com/tidwall/gjson"
)

// TongyiLLM implements the LLM interface for Tongyi (DashScope)
type TongyiLLM struct {
	*BaseLLM
}

// NewTongyiLLM creates a new TongyiLLM
func NewTongyiLLM(config *types.ClientConfig) *TongyiLLM {
	if config.APIBase == "" {
		config.APIBase = "https://dashscope.aliyuncs.com/api/v1"
	}
	if config.CompletionPath == "" {
		config.CompletionPath = "services/aigc/text-generation/generation"
	}
	if config.AnswerPath == "" {
		config.AnswerPath = "output.text"
	}
	if config.Model == "" {
		config.Model = "qwen-turbo"
	}

	return &TongyiLLM{
		BaseLLM: NewBaseLLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (t *TongyiLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "https://dashscope.aliyuncs.com/api/v1",
			PromptMessage: "Enter Tongyi API base",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "qwen-turbo",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}

// FormatMessages formats messages for Tongyi API
func (t *TongyiLLM) FormatMessages(message string, history []types.Message) (interface{}, error) {
	messages := make([]map[string]string, 0, len(history)+1)
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": message,
	})

	payload := map[string]interface{}{
		"model": t.Config.Model,
		"input": map[string]interface{}{
			"messages": messages,
		},
		"parameters": map[string]interface{}{
			"max_tokens": t.Config.MaxTokens,
		},
	}

	if t.Config.Temperature != 0 {
		payload["parameters"].(map[string]interface{})["temperature"] = t.Config.Temperature
	}
	if t.Config.TopP != 0 {
		payload["parameters"].(map[string]interface{})["top_p"] = t.Config.TopP
	}

	return payload, nil
}

// BuildHeaders builds request headers
func (t *TongyiLLM) BuildHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", t.Config.APIKey),
	}
	for k, v := range t.Config.ExtraHeaders {
		headers[k] = v
	}
	return headers
}

// GetUsage returns usage information for the provider
func (t *TongyiLLM) GetUsage(data []byte) (string, error) {
	usage := gjson.GetBytes(data, "usage")
	if !usage.Exists() {
		return "", nil
	}

	return fmt.Sprintf(
		"Token usage> input: %d, output: %d, total: %d",
		usage.Get("input_tokens").Int(),
		usage.Get("output_tokens").Int(),
		usage.Get("total_tokens").Int(),
	), nil
}

// MakeRequest makes a request to the API
func (t *TongyiLLM) MakeRequest(ctx context.Context, client *http.Client, message string, history []types.Message) (string, error) {
	return t.BaseLLM.MakeRequest(ctx, client, t, message, history)
}
