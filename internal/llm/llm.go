package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/belingud/go-gptcomet/pkg/types"
	"github.com/tidwall/gjson"
)

// LLM is the interface that all LLM providers must implement
type LLM interface {
	// GetRequiredConfig returns provider-specific configuration requirements
	GetRequiredConfig() map[string]ConfigRequirement

	// FormatMessages formats messages for the provider's API
	FormatMessages(message string, history []types.Message) (interface{}, error)

	// BuildURL builds the API URL
	BuildURL() string

	// BuildHeaders builds request headers
	BuildHeaders() map[string]string

	// ParseResponse parses the response from the API
	ParseResponse(response []byte) (string, error)

	// GetUsage returns usage information for the provider
	GetUsage(data []byte) (string, error)

	// MakeRequest makes a request to the API
	MakeRequest(ctx context.Context, client *http.Client, message string, history []types.Message) (string, error)
}

// ConfigRequirement represents a configuration requirement
type ConfigRequirement struct {
	DefaultValue  string
	PromptMessage string
}

// BaseLLM provides common functionality for all LLM providers
type BaseLLM struct {
	Config *types.ClientConfig
}

// NewBaseLLM creates a new BaseLLM
func NewBaseLLM(config *types.ClientConfig) *BaseLLM {
	return &BaseLLM{
		Config: config,
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (b *BaseLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "",
			PromptMessage: "Enter API Base URL",
		},
		"model": {
			DefaultValue:  "",
			PromptMessage: "Enter model name",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}

// FormatMessages provides a default implementation for formatting messages (OpenAI format)
func (b *BaseLLM) FormatMessages(message string, history []types.Message) (interface{}, error) {
	messages := make([]types.Message, 0, len(history)+1)
	if history != nil {
		messages = append(messages, history...)
	}
	messages = append(messages, types.Message{
		Role:    "user",
		Content: message,
	})

	payload := map[string]interface{}{
		"model":      b.Config.Model,
		"messages":   messages,
		"max_tokens": b.Config.MaxTokens,
	}
	if b.Config.Temperature != 0 {
		payload["temperature"] = b.Config.Temperature
	}
	if b.Config.TopP != 0 {
		payload["top_p"] = b.Config.TopP
	}
	if b.Config.FrequencyPenalty != 0 {
		payload["frequency_penalty"] = b.Config.FrequencyPenalty
	}
	if b.Config.PresencePenalty != 0 {
		payload["presence_penalty"] = b.Config.PresencePenalty
	}

	return payload, nil
}

// BuildHeaders provides a default implementation for building headers
func (b *BaseLLM) BuildHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	if b.Config.APIKey != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", b.Config.APIKey)
	}
	for k, v := range b.Config.ExtraHeaders {
		headers[k] = v
	}
	return headers
}

// BuildURL provides a default implementation for building URL
func (b *BaseLLM) BuildURL() string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(b.Config.APIBase, "/"), strings.TrimPrefix(b.Config.CompletionPath, "/"))
}

// ParseResponse provides a default implementation for parsing response
func (b *BaseLLM) ParseResponse(response []byte) (string, error) {
	result := gjson.GetBytes(response, b.Config.AnswerPath)
	if !result.Exists() {
		return "", fmt.Errorf("failed to find answer in response: %s", string(response))
	}
	return result.String(), nil
}

// GetUsage provides a default implementation for getting usage information
func (b *BaseLLM) GetUsage(data []byte) (string, error) {
	usage := gjson.GetBytes(data, "usage")
	if !usage.Exists() {
		return "", nil
	}

	var promptTokens, completionTokens, totalTokens int64

	// Try different field names used by different providers
	promptTokens = usage.Get("prompt_tokens").Int()
	completionTokens = usage.Get("completion_tokens").Int()
	totalTokens = usage.Get("total_tokens").Int()

	return fmt.Sprintf(
		"Token usage> prompt: %d, completion: %d, total: %d",
		promptTokens,
		completionTokens,
		totalTokens,
	), nil
}

// MakeRequest provides a default implementation for making requests
func (b *BaseLLM) MakeRequest(ctx context.Context, client *http.Client, provider LLM, message string, history []types.Message) (string, error) {
	url := provider.BuildURL()
	headers := provider.BuildHeaders()
	payload, err := provider.FormatMessages(message, history)
	if err != nil {
		return "", fmt.Errorf("failed to format messages: %w", err)
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	usage, err := provider.GetUsage(respBody)
	if err != nil {
		return "", fmt.Errorf("failed to get usage: %w", err)
	}
	if usage != "" {
		fmt.Printf("%s\n", usage)
	}

	return provider.ParseResponse(respBody)
}
