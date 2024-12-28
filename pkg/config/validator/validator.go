package validator

import (
	"fmt"
	"strings"

	"github.com/belingud/go-gptcomet/pkg/errors"
	"github.com/belingud/go-gptcomet/pkg/types"
)

// ConfigValidator validates configuration
type ConfigValidator struct {
	RequiredFields map[string][]string
}

// NewConfigValidator creates a new ConfigValidator
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		RequiredFields: map[string][]string{
			"vertex": {"api_base", "project_id", "location"},
			"ollama": {"api_base"},
			"openai": {"api_base", "api_key"},
			"anthropic": {"api_base", "api_key"},
			"gemini": {"api_base", "api_key"},
		},
	}
}

// ValidateConfig validates the configuration
func (v *ConfigValidator) ValidateConfig(config *types.ClientConfig) error {
	if config == nil {
		return errors.NewError(errors.ErrValidation, "config is nil", nil)
	}

	if config.Provider == "" {
		return errors.NewError(errors.ErrValidation, "provider is required", nil)
	}

	requiredFields, ok := v.RequiredFields[config.Provider]
	if !ok {
		return errors.NewError(errors.ErrValidation, fmt.Sprintf("unsupported provider: %s", config.Provider), nil)
	}

	var missingFields []string
	for _, field := range requiredFields {
		if !v.hasField(config, field) {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return errors.NewError(
			errors.ErrValidation,
			fmt.Sprintf("missing required fields: %s", strings.Join(missingFields, ", ")),
			nil,
		).WithContext("provider", config.Provider)
	}

	return nil
}

// hasField checks if a field exists and has a non-empty value
func (v *ConfigValidator) hasField(config *types.ClientConfig, field string) bool {
	switch field {
	case "api_base":
		return config.APIBase != ""
	case "api_key":
		return config.APIKey != ""
	case "project_id":
		return config.ProjectID != ""
	case "location":
		return config.Location != ""
	default:
		return false
	}
}
