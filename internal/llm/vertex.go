package llm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/belingud/go-gptcomet/pkg/types"
	"github.com/tidwall/gjson"
)

// VertexLLM implements the LLM interface for Google Cloud Vertex AI
type VertexLLM struct {
	*BaseLLM
}

// NewVertexLLM creates a new VertexLLM
func NewVertexLLM(config *types.ClientConfig) *VertexLLM {
	if config.APIBase == "" {
		config.APIBase = "https://us-central1-aiplatform.googleapis.com/v1"
	}
	if config.Model == "" {
		config.Model = "gemini-pro"
	}
	if config.CompletionPath == "" {
		config.CompletionPath = "models/%s:predict"
	}
	if config.AnswerPath == "" {
		config.AnswerPath = "predictions.0.candidates.0"
	}

	return &VertexLLM{
		BaseLLM: NewBaseLLM(config),
	}
}

// GetRequiredConfig returns provider-specific configuration requirements
func (v *VertexLLM) GetRequiredConfig() map[string]ConfigRequirement {
	return map[string]ConfigRequirement{
		"api_base": {
			DefaultValue:  "https://us-central1-aiplatform.googleapis.com/v1",
			PromptMessage: "Enter Vertex AI API base",
		},
		"api_key": {
			DefaultValue:  "",
			PromptMessage: "Enter API key",
		},
		"model": {
			DefaultValue:  "gemini-pro",
			PromptMessage: "Enter model name",
		},
		"max_tokens": {
			DefaultValue:  "1024",
			PromptMessage: "Enter max tokens",
		},
	}
}

// BuildURL builds the API URL for Vertex AI
func (v *VertexLLM) BuildURL() string {
	return fmt.Sprintf("%s/%s", v.Config.APIBase, fmt.Sprintf(v.Config.CompletionPath, v.Config.Model))
}

// FormatMessages formats messages for Vertex AI
func (v *VertexLLM) FormatMessages(message string, history []types.Message) (interface{}, error) {
	context := ""
	examples := []map[string]string{}

	if history != nil {
		for i := 0; i < len(history); i += 2 {
			if i+1 < len(history) {
				examples = append(examples, map[string]string{
					"input":  history[i].Content,
					"output": history[i+1].Content,
				})
			}
		}
	}

	payload := map[string]interface{}{
		"instances": []map[string]interface{}{
			{
				"context":  context,
				"examples": examples,
				"messages": []map[string]string{
					{
						"author":  "user",
						"content": message,
					},
				},
			},
		},
		"parameters": map[string]interface{}{
			"maxOutputTokens": v.Config.MaxTokens,
			"temperature":     v.Config.Temperature,
			"topP":            v.Config.TopP,
		},
	}

	return payload, nil
}

// GetUsage returns usage information for the provider
func (v *VertexLLM) GetUsage(data []byte) (string, error) {
	usage := gjson.GetBytes(data, "metadata.tokenMetadata")
	if !usage.Exists() {
		return "", nil
	}

	return fmt.Sprintf(
		"Token usage> input: %d, output: %d, total: %d",
		usage.Get("inputTokenCount").Int(),
		usage.Get("outputTokenCount").Int(),
		usage.Get("totalTokenCount").Int(),
	), nil
}

// MakeRequest makes a request to the API
func (v *VertexLLM) MakeRequest(ctx context.Context, client *http.Client, message string, history []types.Message) (string, error) {
	return v.BaseLLM.MakeRequest(ctx, client, v, message, history)
}
