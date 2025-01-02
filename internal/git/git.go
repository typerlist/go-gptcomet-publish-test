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

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

// GetDiff retrieves the staged git diff for the specified repository path.
// It runs the "git diff --staged -U2" command and filters out lines that start with "index", "---", and "+++".
//
// Parameters:
//   - repoPath: The file path to the git repository.
//
// Returns:
//   - A string containing the filtered diff output.
//   - An error if the command fails or if the specified path is not a git repository.
func GetDiff(repoPath string) (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "-U2")
	cmd.Dir = repoPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 128 {
				return "", fmt.Errorf("not a git repository (or any of the parent directories): %w\nGit output: %s", err, stderr.String())
			} else {
				return "", fmt.Errorf("git diff command failed with exit code %d: %w\nGit output: %s", exitError.ExitCode(), err, stderr.String())
			}
		}
		return "", fmt.Errorf("failed to get diff: %w\nGit output: %s", err, stderr.String())
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

// HasStagedChanges checks if there are any staged changes in the git repository at the given path.
// It runs "git diff --staged --quiet" command and interprets the exit code to determine if there
// are staged changes.
//
// Parameters:
//   - repoPath: The file system path to the git repository
//
// Returns:
//   - bool: true if there are staged changes, false otherwise
//   - error: nil if the command executed successfully, error otherwise
//
// The function returns true if the git diff command exits with code 1 (staged changes present),
// false if it exits with code 0 (no staged changes), and an error for any other exit code or
// if the command fails to execute.
func HasStagedChanges(repoPath string) (bool, error) {
	cmd := exec.Command("git", "diff", "--staged", "--quiet")
	cmd.Dir = repoPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means there are staged changes
			if exitError.ExitCode() == 1 {
				return true, nil
			} else {
				return false, fmt.Errorf("git diff command failed with exit code %d: %w\nGit output: %s", exitError.ExitCode(), err, stderr.String())
			}
		}
		return false, fmt.Errorf("failed to check staged changes: %w\nGit output: %s", err, stderr.String())
	}

	// Exit code 0 means no staged changes
	return false, nil
}

// GetStagedFiles returns a list of files that are currently staged for commit in the git repository
// at the specified path. It executes the 'git diff --staged --name-only' command to get the list
// of staged files.
//
// Parameters:
//   - repoPath: The file system path to the git repository
//
// Returns:
//   - []string: A slice containing the paths of all staged files, or nil if no files are staged
//   - error: An error if the git command fails or if there are issues accessing the repository
//
// The function will return (nil, nil) if there are no staged files in the repository.
// If the git command fails, it returns a detailed error message including the exit code.
func GetStagedFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	cmd.Dir = repoPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git diff command failed with exit code %d: %w\nGit output: %s", exitError.ExitCode(), err, stderr.String())
		}
		return nil, fmt.Errorf("failed to get staged files: %w\nGit output: %s", err, stderr.String())
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

// GetStagedDiffFiltered returns the git diff for staged changes, excluding files that match the patterns
// specified in the config manager under the "file_ignore" key.
//
// Parameters:
//   - repoPath: The file system path to the git repository
//   - cfgManager: The config manager to use for retrieving ignore patterns
//
// Returns:
//   - string: The filtered diff output
//   - error: An error if the git command fails or if there are issues accessing the repository
//
// The function will return an empty string if there are no staged files in the repository.
// If the git command fails, it returns a detailed error message including the exit code.
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

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git diff command failed with exit code %d: %w\nGit output: %s", exitError.ExitCode(), err, stderr.String())
		}
		return "", fmt.Errorf("failed to get diff: %w\nGit output: %s", err, stderr.String())
	}

	return string(output), nil
}

// GetCurrentBranch returns the name of the current branch in the git repository
// at the specified path.
//
// Parameters:
//   - repoPath: The file system path to the git repository
//
// Returns:
//   - string: The name of the current branch
//   - error: An error if the git command fails or if there are issues accessing the repository
func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git rev-parse command failed with exit code %d: %w\nGit output: %s", exitError.ExitCode(), err, stderr.String())
		}
		return "", fmt.Errorf("failed to get current branch: %w\nGit output: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

// GetCommitInfo returns formatted information about the commit
// If commitHash is empty, returns info about the last commit
//
// Parameters:
//   - repoPath: The file system path to the git repository
//   - commitHash: The hash of the commit to get info for (or empty for the last commit)
//
// Returns:
//   - string: The formatted commit info
//   - error: An error if the git command fails or if there are issues accessing the repository
func GetCommitInfo(repoPath string, commitHash string) (string, error) {
	if commitHash == "" {
		// Get last commit hash
		hash, err := GetLastCommitHash(repoPath)
		if err != nil {
			return "", err
		}
		commitHash = hash
	}

	cmd := exec.Command("git", "log", "-1", "--stat", commitHash, "--pretty=format:Author: %an <%ae>%n%D(%H)%n%n%s%n")
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git log command failed with exit code %d: %w\nGit output: %s", exitError.ExitCode(), err, stderr.String())
		}
		return "", fmt.Errorf("failed to get commit info: %w\nGit output: %s", err, stderr.String())
	}

	// Get branch name
	branch, err := GetCurrentBranch(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to get branch name: %w", err)
	}

	// Replace the ref info with just the branch name and add colors
	output := stdout.String()
	lines := strings.Split(output, "\n")
	if len(lines) > 1 {
		// Replace the second line (which contains ref info) with just the branch name
		lines[1] = strings.Split(lines[1], "(")[0] + lines[1][strings.LastIndex(lines[1], "("):]
		lines[1] = branch + lines[1][strings.LastIndex(lines[1], "("):]

		// Add colors to the stats
		for i := 4; i < len(lines); i++ {
			line := lines[i]
			if strings.Contains(line, "|") {
				parts := strings.Split(line, "|")
				if len(parts) == 2 {
					stats := strings.TrimSpace(parts[1])
					coloredStats := strings.ReplaceAll(stats, "+", colorGreen+"+")
					coloredStats = strings.ReplaceAll(coloredStats, "-", colorReset+colorRed+"-")
					lines[i] = parts[0] + "| " + coloredStats + colorReset
				}
			}
		}
		output = strings.Join(lines, "\n")
	}

	return output, nil
}

// GetLastCommitHash returns the hash of the last commit
// Parameters:
//   - repoPath: The file system path to the git repository
//
// Returns:
//   - string: The hash of the last commit
//   - error: An error if the git command fails or if there are issues accessing the repository
func GetLastCommitHash(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git rev-parse command failed with exit code %d: %w\nGit output: %s", exitError.ExitCode(), err, stderr.String())
		}
		return "", fmt.Errorf("failed to get last commit hash: %w\nGit output: %s", err, stderr.String())
	}
	return strings.TrimSpace(string(output)), nil
}

// CreateCommit creates a git commit with the given message
//
// Parameters:
//   - repoPath: The file system path to the git repository
//   - message: The commit message
//
// Returns:
//   - error: An error if the git command fails or if there are issues accessing the repository
func CreateCommit(repoPath string, message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("git commit command failed with exit code %d: %w", exitError.ExitCode(), err)
		}
		return fmt.Errorf("failed to create commit: %s, %w", stderr.String(), err)
	}

	return nil
}
