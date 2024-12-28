package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/belingud/go-gptcomet/internal/config"
	"github.com/belingud/go-gptcomet/internal/debug"
)

// GetDiff returns the git diff for staged changes
func GetDiff(repoPath string) (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "-U2")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}

	// Filter out index, ---, and +++ lines
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	filteredLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.HasPrefix(line, "index") &&
			!strings.HasPrefix(line, "---") &&
			!strings.HasPrefix(line, "+++") {
			filteredLines = append(filteredLines, line)
		}
	}

	return strings.Join(filteredLines, "\n"), nil
}

// HasStagedChanges checks if there are any staged changes
func HasStagedChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "diff", "--staged", "--quiet")
	cmd.Dir = repoPath

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means there are staged changes
			if exitError.ExitCode() == 1 {
				return true, nil
			}
		}
		return false, fmt.Errorf("failed to check staged changes: %w", err)
	}

	// Exit code 0 means no staged changes
	return false, nil
}

// GetStagedFiles returns a list of staged files
func GetStagedFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 0 || (len(files) == 1 && files[0] == "") {
		return nil, nil
	}

	return files, nil
}

// ShouldIgnoreFile checks if a file should be ignored based on patterns
func ShouldIgnoreFile(file string, ignorePatterns []string) bool {
	for _, pattern := range ignorePatterns {
		matched, err := filepath.Match(pattern, file)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// GetStagedDiffFiltered returns the git diff for staged changes, excluding ignored files
func GetStagedDiffFiltered(repoPath string, cfgManager *config.Manager) (string, error) {
	// First get staged files
	files, err := GetStagedFiles(repoPath)
	debug.Printf("Staged files: %v", files)
	if err != nil {
		return "", err
	}

	// Get ignore patterns from config
	var ignorePatterns []string
	if patterns, ok := cfgManager.Get("file_ignore"); ok {
		if patternList, ok := patterns.([]interface{}); ok {
			for _, p := range patternList {
				if str, ok := p.(string); ok {
					ignorePatterns = append(ignorePatterns, str)
				}
			}
		}
	}

	// Filter files based on ignore patterns
	var filteredFiles []string
	for _, file := range files {
		if !ShouldIgnoreFile(file, ignorePatterns) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	debug.Printf("Filtered files: %v", filteredFiles)

	if len(filteredFiles) == 0 {
		return "", nil
	}

	// Get diff for filtered files
	args := append([]string{"diff", "--staged", "-U2"}, filteredFiles...)
	debug.Printf("Diff args: %v", args)
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}

	return string(output), nil
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}

// GetCommitInfo returns formatted information about the last commit
func GetCommitInfo(repoPath string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--stat", "--pretty=format:Author: %an <%ae>%n%D(%H)%n%n%s%n")
	cmd.Dir = repoPath

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get commit info: %w", err)
	}

	// Get branch name
	branch, err := GetCurrentBranch(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to get branch name: %w", err)
	}

	// Replace the ref info with just the branch name
	output := out.String()
	lines := strings.Split(output, "\n")
	if len(lines) > 1 {
		// Replace the second line (which contains ref info) with just the branch name
		lines[1] = strings.Split(lines[1], "(")[0] + lines[1][strings.LastIndex(lines[1], "("):]
		lines[1] = branch + lines[1][strings.LastIndex(lines[1], "("):]
		output = strings.Join(lines, "\n")
	}

	return output, nil
}

// CreateCommit creates a git commit with the given message
func CreateCommit(repoPath string, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create commit: %s, %w", stderr.String(), err)
	}

	return nil
}
