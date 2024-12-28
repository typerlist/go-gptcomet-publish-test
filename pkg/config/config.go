package config

import (
	"github.com/belingud/go-gptcomet/pkg/errors"
	"github.com/belingud/go-gptcomet/pkg/types"
)

// Manager handles configuration management
type Manager struct {
	defaults map[string]*types.ClientConfig
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		defaults: map[string]*types.ClientConfig{
			"vertex": {
				APIBase:   "https://us-central1-aiplatform.googleapis.com/v1",
				Model:     "gemini-1.5-flash",
				Location:  "us-central1",
				MaxTokens: 1024,
			},
			"openai": {
				APIBase:   "https://api.openai.com/v1",
				Model:     "gpt-4",
				MaxTokens: 1024,
			},
			"anthropic": {
				APIBase:   "https://api.anthropic.com/v1",
				Model:     "claude-3-opus",
				MaxTokens: 1024,
			},
			"gemini": {
				APIBase:   "https://generativelanguage.googleapis.com/v1beta/models",
				Model:     "gemini-pro",
				MaxTokens: 1024,
			},
			"ollama": {
				APIBase:   "http://localhost:11434/api",
				Model:     "llama2",
				MaxTokens: 1024,
			},
			"chatglm": {
				APIBase:   "https://open.bigmodel.cn/api/paas/v4",
				Model:     "glm-4-flash",
				MaxTokens: 1024,
			},
			"cohere": {
				APIBase:   "https://api.cohere.com/v1",
				Model:     "command-r-plus",
				MaxTokens: 1024,
			},
			"deepseek": {
				APIBase:   "https://api.deepseek.com/beta",
				Model:     "deepseek-chat",
				MaxTokens: 1024,
			},
			"sambanova": {
				APIBase:   "https://api.sambanova.com/v1",
				Model:     "sambanova-1",
				MaxTokens: 1024,
			},
			"tongyi": {
				APIBase:   "https://api.tongyi.ai/v1",
				Model:     "tongyi-1",
				MaxTokens: 1024,
			},
		},
	}
}

// ApplyDefaults applies default values to the configuration
func (m *Manager) ApplyDefaults(provider string, cfg *types.ClientConfig) error {
	defaults, ok := m.defaults[provider]
	if !ok {
		return errors.NewError(
			errors.ErrConfig,
			"no defaults found for provider",
			nil,
		).WithContext("provider", provider)
	}

	// Only apply defaults for empty values
	if cfg.APIBase == "" {
		cfg.APIBase = defaults.APIBase
	}
	if cfg.Model == "" {
		cfg.Model = defaults.Model
	}
	if cfg.Location == "" && defaults.Location != "" {
		cfg.Location = defaults.Location
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = defaults.MaxTokens
	}
	if cfg.Temperature == 0 {
		cfg.Temperature = 0.7 // 常用默认值
	}
	if cfg.TopP == 0 {
		cfg.TopP = 1.0 // 常用默认值
	}

	return nil
}

// ValidateConfig validates the configuration
func (m *Manager) ValidateConfig(provider string, cfg *types.ClientConfig) error {
	// 基本验证
	if cfg.APIBase == "" {
		return errors.NewError(
			errors.ErrConfig,
			"api_base is required",
			nil,
		)
	}

	if cfg.Model == "" {
		return errors.NewError(
			errors.ErrConfig,
			"model is required",
			nil,
		)
	}

	// 提供商特定验证
	switch provider {
	case "vertex":
		if cfg.ProjectID == "" {
			return errors.NewError(
				errors.ErrConfig,
				"project_id is required for Vertex AI",
				nil,
			)
		}
		if cfg.Location == "" {
			return errors.NewError(
				errors.ErrConfig,
				"location is required for Vertex AI",
				nil,
			)
		}
	case "anthropic":
		if cfg.AnthropicVersion == "" {
			return errors.NewError(
				errors.ErrConfig,
				"anthropic_version is required for Anthropic",
				nil,
			)
		}
	}

	return nil
}
