package config

// ConfigRequirement represents a configuration requirement
type ConfigRequirement struct {
	DefaultValue  string
	PromptMessage string
}

// ProviderConfig represents provider-specific configuration
type ProviderConfig struct {
	APIBase    string
	Model      string
	MaxTokens  int
	Location   string
	ProjectID  string
	APIKey     string
	APIVersion string
}
