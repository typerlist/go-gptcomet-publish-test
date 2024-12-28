package provider

// GetAnswerPath returns the configured answer path based on the provider
func GetAnswerPath(provider string) string {
	switch provider {
	case "openai", "azure-openai", "xai", "groq":
		return "choices.0.message.content"
	case "claude":
		return "content.0.text"
	case "palm", "gemini", "vertex":
		return "candidates.0.content.parts.0.text"
	case "cohere":
		return "generations.0.text"
	case "mistral":
		return "choices.0.message.content"
	case "ollama":
		return "response"
	case "tongyi":
		return "output.text"
	case "chatglm":
		return "data.choices.0.message.content"
	default:
		return "choices.0.message.content"
	}
}

// GetCompletionPath returns the configured completion path based on the provider
func GetCompletionPath(provider string) string {
	switch provider {
	case "openai", "xai", "groq":
		return "/chat/completions"
	case "azure-openai":
		return "/openai/deployments/{model}/chat/completions"
	case "claude":
		return "/v1/messages"
	case "palm":
		return "/v1beta3/models/{model}:generateText"
	case "gemini", "vertex":
		return "/v1beta/models/{model}:generateContent"
	case "cohere":
		return "/v1/generate"
	case "mistral":
		return "/v1/chat/completions"
	case "ollama":
		return "/api/generate"
	case "tongyi":
		return "/compatible-mode/v1/chat/completions"
	case "chatglm":
		return "/api/paas/v4/chat/completions"
	default:
		return "/chat/completions"
	}
}

// GetDefaultModel returns the default model based on the provider
func GetDefaultModel(provider string) string {
	switch provider {
	case "openai":
		return "gpt-4"
	case "azure-openai":
		return "gpt-4"
	case "xai":
		return "grok-beta"
	case "claude":
		return "claude-3-sonnet"
	case "palm":
		return "text-bison-001"
	case "gemini":
		return "gemini-pro"
	case "vertex":
		return "chat-bison"
	case "cohere":
		return "command"
	case "mistral":
		return "mistral-medium"
	case "ollama":
		return "llama2"
	case "tongyi":
		return "qwen-turbo"
	case "chatglm":
		return "glm-4-flash"
	case "groq":
		return "llama3-8b-8192"
	default:
		return "gpt-4"
	}
}

// GetDefaultAPIBase returns the default API base URL based on the provider
func GetDefaultAPIBase(provider string) string {
	switch provider {
	case "openai":
		return "https://api.openai.com/v1"
	case "azure-openai":
		return "https://{resource}.openai.azure.com"
	case "xai":
		return "https://api.x.ai/v1"
	case "claude":
		return "https://api.anthropic.com/v1"
	case "palm":
		return "https://generativelanguage.googleapis.com"
	case "gemini":
		return "https://generativelanguage.googleapis.com/v1beta"
	case "vertex":
		return "https://{region}-aiplatform.googleapis.com"
	case "cohere":
		return "https://api.cohere.ai"
	case "mistral":
		return "https://api.mistral.ai"
	case "ollama":
		return "http://localhost:11434"
	case "tongyi":
		return "https://dashscope.aliyuncs.com"
	case "chatglm":
		return "https://open.bigmodel.cn"
	case "groq":
		return "https://api.groq.com/openai/v1"
	default:
		return "https://api.openai.com/v1"
	}
}

// GetProviders returns a list of supported providers
func GetProviders() []string {
	return []string{
		"openai",
		"azure-openai",
		"xai",
		"claude",
		"palm",
		"gemini",
		"vertex",
		"cohere",
		"mistral",
		"ollama",
		"tongyi",
		"chatglm",
		"groq",
	}
}

// IsValidProvider checks if a provider is supported
func IsValidProvider(provider string) bool {
	switch provider {
	case "openai", "azure-openai", "xai", "claude", "palm", "gemini", "vertex",
		"cohere", "mistral", "ollama", "tongyi", "chatglm", "groq":
		return true
	default:
		return false
	}
}
