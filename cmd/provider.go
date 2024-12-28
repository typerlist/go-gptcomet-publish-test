package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/belingud/go-gptcomet/internal/config"
	"github.com/belingud/go-gptcomet/internal/debug"
	"github.com/belingud/go-gptcomet/internal/provider"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Provider represents an LLM provider configuration
type Provider struct {
	Name            string
	DefaultAPIBase  string
	DefaultModel    string
	DefaultMaxToken int
	AnswerPath      string
}

var providers = []Provider{
	{
		Name:            "openai",
		DefaultAPIBase:  provider.GetDefaultAPIBase("openai"),
		DefaultModel:    provider.GetDefaultModel("openai"),
		DefaultMaxToken: 1024,
		AnswerPath:      provider.GetAnswerPath("openai"),
	},
	{
		Name:            "xai",
		DefaultAPIBase:  provider.GetDefaultAPIBase("xai"),
		DefaultModel:    provider.GetDefaultModel("xai"),
		DefaultMaxToken: 1024,
		AnswerPath:      provider.GetAnswerPath("xai"),
	},
	{
		Name:            "claude",
		DefaultAPIBase:  provider.GetDefaultAPIBase("claude"),
		DefaultModel:    provider.GetDefaultModel("claude"),
		DefaultMaxToken: 1024,
		AnswerPath:      provider.GetAnswerPath("claude"),
	},
	{
		Name:            "manual",
		DefaultAPIBase:  "",
		DefaultModel:    "",
		DefaultMaxToken: 1024,
		AnswerPath:      "",
	},
}

func readMaskedInput(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println() // add newline
	return string(bytePassword), nil
}

// NewProviderCmd creates a new provider command
func NewProviderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "newprovider",
		Short: "Add a new API provider interactively",
		RunE: func(cmd *cobra.Command, args []string) error {
			debug.Println("Starting provider configuration")

			// Create a reader for user input
			reader := bufio.NewReader(os.Stdin)

			// Display available providers
			fmt.Println("Available providers:")
			for i, p := range providers {
				fmt.Printf("%d. %s\n", i+1, p.Name)
			}

			// Get provider choice
			var selectedProvider Provider
			for {
				fmt.Print("Select a provider (enter number): ")
				choice, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read provider choice: %w", err)
				}
				choice = strings.TrimSpace(choice)
				index, err := strconv.Atoi(choice)
				if err != nil || index < 1 || index > len(providers) {
					fmt.Println("Invalid choice. Please try again.")
					continue
				}
				selectedProvider = providers[index-1]
				break
			}

			var provider, apiBase, model string
			var maxTokens int
			var answerPath string

			if selectedProvider.Name == "manual" {
				// Get provider name for manual input
				fmt.Print("Enter provider name: ")
				provider, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read provider name: %w", err)
				}
				provider = strings.TrimSpace(provider)
				if provider == "" {
					return fmt.Errorf("provider name cannot be empty")
				}

				// Get API base
				fmt.Print("Enter API base URL: ")
				apiBase, err = reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read API base: %w", err)
				}
				apiBase = strings.TrimSpace(apiBase)
				if apiBase == "" {
					return fmt.Errorf("API base cannot be empty")
				}

				// Get model
				fmt.Print("Enter model name: ")
				model, err = reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read model: %w", err)
				}
				model = strings.TrimSpace(model)
				if model == "" {
					return fmt.Errorf("model name cannot be empty")
				}

				// Get answer path
				fmt.Print("Enter answer path (e.g., choices.0.message.content): ")
				answerPath, err = reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read answer path: %w", err)
				}
				answerPath = strings.TrimSpace(answerPath)
				if answerPath == "" {
					return fmt.Errorf("answer path cannot be empty")
				}

				maxTokens = 1024 // Default value for manual input
			} else {
				provider = selectedProvider.Name
				apiBase = selectedProvider.DefaultAPIBase
				model = selectedProvider.DefaultModel
				maxTokens = selectedProvider.DefaultMaxToken
				answerPath = selectedProvider.AnswerPath

				// Allow customization of pre-configured values
				fmt.Printf("Enter API base URL [%s]: ", apiBase)
				customAPIBase, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read API base: %w", err)
				}
				customAPIBase = strings.TrimSpace(customAPIBase)
				if customAPIBase != "" {
					apiBase = customAPIBase
				}

				fmt.Printf("Enter model name [%s]: ", model)
				customModel, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read model: %w", err)
				}
				customModel = strings.TrimSpace(customModel)
				if customModel != "" {
					model = customModel
				}
			}

			// Get API key (with masked input)
			apiKey, err := readMaskedInput("Enter API key: ")
			if err != nil {
				return fmt.Errorf("failed to read API key: %w", err)
			}
			if apiKey == "" {
				return fmt.Errorf("API key cannot be empty")
			}

			// Get model max tokens
			fmt.Printf("Enter model max tokens [%d]: ", maxTokens)
			maxTokensStr, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read max tokens: %w", err)
			}
			maxTokensStr = strings.TrimSpace(maxTokensStr)
			if maxTokensStr != "" {
				maxTokens, err = strconv.Atoi(maxTokensStr)
				if err != nil {
					return fmt.Errorf("invalid max tokens value: %w", err)
				}
			}

			// Create config manager
			cfgManager, err := config.New()
			if err != nil {
				return err
			}

			// Check if provider already exists
			if _, exists := cfgManager.Get(provider); exists {
				fmt.Printf("Provider '%s' already exists. Do you want to overwrite it? [y/N]: ", provider)
				answer, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read answer: %w", err)
				}
				answer = strings.ToLower(strings.TrimSpace(answer))
				if answer != "y" && answer != "yes" {
					fmt.Println("Operation cancelled")
					return nil
				}
			}

			// Set the provider configuration
			providerConfig := map[string]interface{}{
				"api_key":          apiKey,
				"api_base":         apiBase,
				"model":            model,
				"model_max_tokens": maxTokens,
				"answer_path":      answerPath,
			}
			if err := cfgManager.Set(provider, providerConfig); err != nil {
				return fmt.Errorf("failed to set provider config: %w", err)
			}

			fmt.Printf("\nProvider configuration saved:\n")
			fmt.Printf("  Provider: %s\n", provider)
			fmt.Printf("  API Base: %s\n", apiBase)
			fmt.Printf("  Model: %s\n", model)
			fmt.Printf("  Max Tokens: %d\n", maxTokens)
			fmt.Printf("  Answer Path: %s\n", answerPath)
			return nil
		},
	}

	return cmd
}
