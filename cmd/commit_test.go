package cmd

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/belingud/go-gptcomet/internal/git"
	"github.com/belingud/go-gptcomet/internal/testutils"
	"github.com/belingud/go-gptcomet/pkg/config"
	"github.com/belingud/go-gptcomet/pkg/types"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommitCmd(t *testing.T) {
	// Create a new commit command
	cmd := NewCommitCmd()
	require.NotNil(t, cmd)

	// Test flags
	flags := map[string]bool{
		"config":  false,
		"dry-run": false,
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if _, ok := flags[flag.Name]; ok {
			flags[flag.Name] = true
		}
	})

	for name, found := range flags {
		assert.True(t, found, "flag %q not found", name)
	}
}

func TestCommitCmd(t *testing.T) {
	// Create a temporary git repository
	repoPath, cleanup := testutils.TestGitRepo(t)
	defer cleanup()

	// Create a dummy file and stage it
	testFileContent := "This is a test file."
	err := os.WriteFile(repoPath+"/test.txt", []byte(testFileContent), 0644)
	require.NoError(t, err)
	err = testutils.RunGitCommand(t, repoPath, "add", "test.txt")
	require.NoError(t, err)

	// Create a dummy config file
	configContent := `
provider: openai
openai:
  api_key: "test_api_key"
  model: "gpt-4"
  api_base: "https://api.openai.com/v1"
`
	configPath, cleanupConfig := testutils.TestConfig(t, configContent)
	defer cleanupConfig()

	// Create a new commit command with dry-run flag
	cmd := NewCommitCmd()
	cmd.SetArgs([]string{"--config", configPath, "--dry-run"})

	// Redirect stdout to capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute the command
	err = cmd.Execute()
	require.NoError(t, err)

	// Check if the commit message is printed
	output := buf.String()
	assert.Contains(t, output, "Generated commit message:")

	// Ensure no commit was created
	_, err = git.GetLastCommitHash(repoPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fatal: your current branch 'master' does not have any commits yet")
}

func TestCommitCmd_NoStagedChanges(t *testing.T) {
	// Create a temporary git repository
	repoPath, cleanup := testutils.TestGitRepo(t)
	defer cleanup()

	// Create a new commit command
	cmd := NewCommitCmd()
	cmd.SetArgs([]string{"--config", repoPath + "/config.yaml"})

	// Redirect stdout to capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute the command
	err := cmd.Execute()

	// Check for the expected error message
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no staged changes found")
}

func TestCommitCmd_NoStagedChangesAfterFiltering(t *testing.T) {
	// Create a temporary git repository
	repoPath, cleanup := testutils.TestGitRepo(t)
	defer cleanup()

	// Create a dummy ignored file
	err := os.WriteFile(repoPath+"/README.md", []byte("This is a README file."), 0644)
	require.NoError(t, err)

	// Stage the ignored file
	err = testutils.RunGitCommand(t, repoPath, "add", "README.md")
	require.NoError(t, err)

	// Create a config file that ignores README.md
	configContent := `
file_ignore:
  - README.md
provider: openai
openai:
  api_key: "test_api_key"
`
	configPath, cleanupConfig := testutils.TestConfig(t, configContent)
	defer cleanupConfig()

	// Create a new commit command
	cmd := NewCommitCmd()
	cmd.SetArgs([]string{"--config", configPath})

	// Redirect stdout to capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute the command
	err = cmd.Execute()

	// Check for the expected error message
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no staged changes found after filtering")
}

// Mock LLM implementation for testing
type mockLLM struct {
	name                  string
	generateCommitMessage func(diff string, prompt string) (string, error)
	translateMessage      func(prompt string, message string, lang string) (string, error)
	makeRequest           func(ctx context.Context, client *http.Client, message string, history []types.Message) (string, error)
}

func (m *mockLLM) GetRequiredConfig() map[string]config.ConfigRequirement {
	return map[string]config.ConfigRequirement{}
}

func (m *mockLLM) BuildHeaders() map[string]string {
	return map[string]string{}
}

func (m *mockLLM) BuildURL() string {
	return "https://mock.api"
}

func (m *mockLLM) FormatMessages(message string, history []types.Message) (interface{}, error) {
	return message, nil
}

func (m *mockLLM) ParseResponse(response []byte) (string, error) {
	return string(response), nil
}

func (m *mockLLM) GetUsage(data []byte) (string, error) {
	return "", nil
}

func (m *mockLLM) MakeRequest(ctx context.Context, client *http.Client, message string, history []types.Message) (string, error) {
	if m.makeRequest != nil {
		return m.makeRequest(ctx, client, message, history)
	}
	return "mock response", nil
}

func (m *mockLLM) GenerateCommitMessage(diff string, prompt string) (string, error) {
	if m.generateCommitMessage != nil {
		return m.generateCommitMessage(diff, prompt)
	}
	return "Test commit message", nil
}

func (m *mockLLM) TranslateMessage(prompt string, message string, lang string) (string, error) {
	if m.translateMessage != nil {
		return m.translateMessage(prompt, message, lang)
	}
	return message, nil
}

func (m *mockLLM) Name() string {
	return m.name
}
