package client

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"github.com/belingud/go-gptcomet/internal/debug"
	"github.com/belingud/go-gptcomet/internal/llm"
	"github.com/belingud/go-gptcomet/pkg/types"
)

// Client represents an LLM client
type Client struct {
	config *types.ClientConfig
	llm    llm.LLM
}

// New creates a new client with the given config
func New(config *types.ClientConfig) *Client {
	var provider llm.LLM
	switch config.Provider {
	case "openai":
		provider = llm.NewOpenAILLM(config)
	case "claude":
		provider = llm.NewClaudeLLM(config)
	case "gemini":
		provider = llm.NewGeminiLLM(config)
	case "mistral":
		provider = llm.NewMistralLLM(config)
	case "xai":
		provider = llm.NewXAILLM(config)
	case "cohere":
		provider = llm.NewCohereLLM(config)
	case "tongyi":
		provider = llm.NewTongyiLLM(config)
	case "deepseek":
		provider = llm.NewDeepSeekLLM(config)
	case "chatglm":
		provider = llm.NewChatGLMLLM(config)
	case "azure":
		provider = llm.NewAzureLLM(config)
	case "vertex":
		provider = llm.NewVertexLLM(config)
	case "kimi":
		provider = llm.NewKimiLLM(config)
	case "ollama":
		provider = llm.NewOllamaLLM(config)
	case "silicon":
		provider = llm.NewSiliconLLM(config)
	case "sambanova":
		provider = llm.NewSambanovaLLM(config)
	default:
		// Default to OpenAI if provider is not specified
		provider = llm.NewOpenAILLM(config)
	}

	return &Client{
		config: config,
		llm:    provider,
	}
}

// RawChat sends a chat completion request and returns the raw JSON response
func (c *Client) RawChat(messages []types.Message) (string, error) {
	return c.sendRawRequest(&types.CompletionRequest{
		Model:    c.config.Model,
		Messages: messages,
	})
}

// Chat sends a chat message to the LLM provider
func (c *Client) Chat(ctx context.Context, message string, history []types.Message) (*types.CompletionResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	content, err := c.llm.MakeRequest(ctx, client, message, history)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return &types.CompletionResponse{
		Content: content,
		Raw:     make(map[string]interface{}),
	}, nil
}

// sendRawRequest sends a completion request to the LLM provider and returns the raw JSON response
func (c *Client) sendRawRequest(req *types.CompletionRequest) (string, error) {
	// Create a transport with proxy if configured
	transport := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: false}, // 默认验证证书
	}

	if c.config.Proxy != "" {
		debug.Printf("Using proxy: %s", c.config.Proxy)
		proxyURL, err := url.Parse(c.config.Proxy)
		if err != nil {
			return "", fmt.Errorf("failed to parse proxy URL: %w", err)
		}

		// 根据代理类型设置不同的配置
		switch proxyURL.Scheme {
		case "http", "https":
			transport.Proxy = http.ProxyURL(proxyURL)
			debug.Printf("Using HTTP/HTTPS proxy: %s", proxyURL.String())
		case "socks5":
			auth := &proxy.Auth{}
			if proxyURL.User != nil {
				auth.User = proxyURL.User.Username()
				if password, ok := proxyURL.User.Password(); ok {
					auth.Password = password
				}
			}
			dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
			if err != nil {
				return "", fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
			debug.Printf("Using SOCKS5 proxy: %s", proxyURL.String())
		default:
			return "", fmt.Errorf("unsupported proxy scheme: %s", proxyURL.Scheme)
		}

		// 如果代理有认证信息，添加 Proxy-Authorization 头
		if proxyURL.User != nil {
			auth := proxyURL.User.String()
			basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			headers := c.llm.BuildHeaders()
			headers["Proxy-Authorization"] = "Basic " + basicAuth
			debug.Printf("Added proxy authentication")
		}
	}

	// Create a client with the configured transport and timeout
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(c.config.Timeout) * time.Second,
	}

	// Make the request using the LLM provider
	return c.llm.MakeRequest(context.Background(), client, req.Messages[len(req.Messages)-1].Content, req.Messages[:len(req.Messages)-1])
}

// getDefaultAnswerPath returns the default answer path for the given provider
func getDefaultAnswerPath(provider string) string {
	switch provider {
	case "openai", "tongyi", "deepseek", "chatglm", "azure", "kimi", "silicon", "sambanova":
		return "choices.0.message.content"
	case "claude":
		return "content.0.text"
	case "gemini":
		return "candidates.0.content.parts.0.text"
	case "mistral":
		return "choices.0.message.content"
	case "xai":
		return "choices.0.message.content"
	case "cohere":
		return "text"
	case "vertex":
		return "predictions.0.candidates.0"
	case "ollama":
		return "message.content"
	default:
		return "choices.0.message.content"
	}
}

// getAnswerPath returns the configured answer path or the default value
func (c *Client) getAnswerPath() string {
	if c.config.AnswerPath != "" {
		return c.config.AnswerPath
	}
	return getDefaultAnswerPath(c.config.Provider)
}

// TranslateMessage translates the given message to the specified language
func (c *Client) TranslateMessage(prompt string, message string, lang string) (string, error) {
	// Format the prompt
	formattedPrompt := fmt.Sprintf(prompt, message, lang)

	// Send the request
	resp, err := c.Chat(context.Background(), formattedPrompt, nil)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.Content), nil
}

// GenerateCommitMessage generates a commit message for the given diff
func (c *Client) GenerateCommitMessage(diff string, prompt string) (string, error) {
	formattedPrompt := fmt.Sprintf(prompt, diff)

	// Send the request
	resp, err := c.Chat(context.Background(), formattedPrompt, nil)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.Content), nil
}

// GenerateCodeExplanation generates an explanation for the given code in the specified language
func (c *Client) GenerateCodeExplanation(message, lang string) (string, error) {
	const prompt = "Explain the following %s code:\n\n%s"
	formattedPrompt := fmt.Sprintf(prompt, lang, message)

	// Send the request
	resp, err := c.Chat(context.Background(), formattedPrompt, nil)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.Content), nil
}

// Stream sends a chat message to the LLM provider and streams the response
func (c *Client) Stream(ctx context.Context, message string, history []types.Message) (*types.CompletionResponse, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	content, err := c.llm.MakeRequest(ctx, client, message, history)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return &types.CompletionResponse{
		Content: content,
		Raw:     make(map[string]interface{}),
	}, nil
}

// getClient returns an HTTP client configured with proxy settings if specified
func (c *Client) getClient() (*http.Client, error) {
	// Create a transport with proxy if configured
	transport := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: false}, // 默认验证证书
	}

	if c.config.Proxy != "" {
		fmt.Printf("Using proxy: %s\n", c.config.Proxy)
		proxyURL, err := url.Parse(c.config.Proxy)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}

		// 根据代理类型设置不同的配置
		switch proxyURL.Scheme {
		case "http", "https":
			transport.Proxy = http.ProxyURL(proxyURL)
		case "socks5":
			auth := &proxy.Auth{}
			if proxyURL.User != nil {
				auth.User = proxyURL.User.Username()
				if password, ok := proxyURL.User.Password(); ok {
					auth.Password = password
				}
			}
			dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
			if err != nil {
				return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
			debug.Printf("Using SOCKS5 proxy: %s", proxyURL.String())
		default:
			return nil, fmt.Errorf("unsupported proxy scheme: %s", proxyURL.Scheme)
		}

		// 如果代理有认证信息，添加 Proxy-Authorization 头
		if proxyURL.User != nil {
			auth := proxyURL.User.String()
			basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			headers := c.llm.BuildHeaders()
			headers["Proxy-Authorization"] = "Basic " + basicAuth
			debug.Printf("Added proxy authentication")
		}
	}

	// Create a client with the configured transport and timeout
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(c.config.Timeout) * time.Second,
	}

	return client, nil
}
