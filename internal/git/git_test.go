package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/belingud/go-gptcomet/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()
	dir, cleanup := testutils.TestDir(t)

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	err := cmd.Run()
	require.NoError(t, err, "Failed to initialize git repository")

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to configure git email")

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to configure git username")

	// Create and add a test file
	testFile := filepath.Join(dir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err, "Failed to create test file")

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to stage test file")

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to create initial commit")

	return dir, cleanup
}

func TestGetDiff(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create a change
	testFile := filepath.Join(dir, "test.txt")
	err := os.WriteFile(testFile, []byte("modified content"), 0644)
	require.NoError(t, err)

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	// Get diff
	diff, err := GetDiff(dir)
	require.NoError(t, err)
	assert.Contains(t, diff, "modified content")
}

func TestHasStagedChanges(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Initially no staged changes
	hasChanges, err := HasStagedChanges(dir)
	require.NoError(t, err)
	assert.False(t, hasChanges)

	// Create and stage a change
	testFile := filepath.Join(dir, "test.txt")
	err = os.WriteFile(testFile, []byte("modified content"), 0644)
	require.NoError(t, err)

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	// Now should have staged changes
	hasChanges, err = HasStagedChanges(dir)
	require.NoError(t, err)
	assert.True(t, hasChanges)
}

func TestGetStagedFiles(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create and stage multiple files
	files := []string{"file1.txt", "file2.txt"}
	for _, f := range files {
		path := filepath.Join(dir, f)
		err := os.WriteFile(path, []byte("test content"), 0644)
		require.NoError(t, err)

		cmd := exec.Command("git", "add", f)
		cmd.Dir = dir
		err = cmd.Run()
		require.NoError(t, err)
	}

	// Get staged files
	stagedFiles, err := GetStagedFiles(dir)
	require.NoError(t, err)
	assert.Equal(t, len(files), len(stagedFiles))
	for _, f := range files {
		assert.Contains(t, stagedFiles, f)
	}
}

func TestCreateCommit(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create and stage a change
	testFile := filepath.Join(dir, "test.txt")
	err := os.WriteFile(testFile, []byte("modified content"), 0644)
	require.NoError(t, err)

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	// Create commit
	message := "test commit"
	err = CreateCommit(dir, message)
	require.NoError(t, err)

	// Verify commit was created
	cmd = exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = dir
	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(output), message)
}

func TestGetCommitInfo(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create and commit a change
	testFile := filepath.Join(dir, "test.txt")
	err := os.WriteFile(testFile, []byte("modified content"), 0644)
	require.NoError(t, err)

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	message := "test commit"
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	// Get commit info
	hash, err := GetLastCommitHash(dir)
	require.NoError(t, err)

	info, err := GetCommitInfo(dir, hash)
	require.NoError(t, err)
	assert.NotEmpty(t, info)
	assert.Contains(t, info, message)
}

func TestGetCurrentBranch(t *testing.T) {
	dir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Get current branch
	branch, err := GetCurrentBranch(dir)
	require.NoError(t, err)
	assert.Equal(t, "master", branch)

	// Create and checkout new branch
	newBranch := "test-branch"
	cmd := exec.Command("git", "checkout", "-b", newBranch)
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err)

	branch, err = GetCurrentBranch(dir)
	require.NoError(t, err)
	assert.Equal(t, newBranch, branch)
}
