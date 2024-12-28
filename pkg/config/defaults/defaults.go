package defaults

// ProviderDefaults contains default configurations for each provider
var ProviderDefaults = map[string]map[string]interface{}{
	"vertex": {
		"api_base":   "https://us-central1-aiplatform.googleapis.com/v1",
		"model":      "gemini-1.5-flash",
		"location":   "us-central1",
		"max_tokens": 1024,
	},
	"ollama": {
		"api_base":   "http://localhost:11434/api",
		"model":      "llama2",
		"max_tokens": 1024,
	},
	"openai": {
		"api_base":   "https://api.openai.com/v1",
		"model":      "gpt-4-turbo",
		"max_tokens": 1024,
	},
	"anthropic": {
		"api_base":   "https://api.anthropic.com/v1",
		"model":      "claude-3-opus",
		"max_tokens": 1024,
	},
	"gemini": {
		"api_base":   "https://generativelanguage.googleapis.com/v1",
		"model":      "gemini-1.5-pro",
		"max_tokens": 1024,
	},
}
