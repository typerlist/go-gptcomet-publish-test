package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/belingud/go-gptcomet/internal/config"
)

// SVNVCS implements the VCS interface for SVN
type SVNVCS struct{}

func (s *SVNVCS) GetDiff(repoPath string) (string, error) {
	cmd := exec.Command("svn", "diff")
	return s.runCommand(cmd, repoPath)
}

func (s *SVNVCS) HasStagedChanges(repoPath string) (bool, error) {
	cmd := exec.Command("svn", "status")
	output, err := s.runCommand(cmd, repoPath)
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(output)) > 0, nil
}

func (s *SVNVCS) GetStagedFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("svn", "status")
	output, err := s.runCommand(cmd, repoPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(output, "\n") {
		if len(line) > 7 && line[0] != '?' {
			files = append(files, strings.TrimSpace(line[7:]))
		}
	}
	return files, nil
}

func (s *SVNVCS) GetStagedDiffFiltered(repoPath string, cfgManager *config.Manager) (string, error) {
	files, err := s.GetStagedFiles(repoPath)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	cmd := exec.Command("svn", append([]string{"diff"}, files...)...)
	return s.runCommand(cmd, repoPath)
}

func (s *SVNVCS) GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("svn", "info", "--show-item", "url")
	return s.runCommand(cmd, repoPath)
}

func (s *SVNVCS) GetCommitInfo(repoPath, commitHash string) (string, error) {
	args := []string{"log", "-r", "HEAD", "-v"}
	if commitHash != "" {
		args = []string{"log", "-r", commitHash, "-v"}
	}
	cmd := exec.Command("svn", args...)
	output, err := s.runCommand(cmd, repoPath)
	if err != nil {
		return "", err
	}

	branch, err := s.GetCurrentBranch(repoPath)
	if err != nil {
		return "", err
	}
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

func (s *SVNVCS) GetLastCommitHash(repoPath string) (string, error) {
	cmd := exec.Command("svn", "info", "--show-item", "revision")
	return s.runCommand(cmd, repoPath)
}

func (s *SVNVCS) CreateCommit(repoPath, message string) error {
	cmd := exec.Command("svn", "commit", "-m", message)
	_, err := s.runCommand(cmd, repoPath)
	return err
}

// runCommand 执行命令并返回输出
func (s *SVNVCS) runCommand(cmd *exec.Cmd, repoPath string) (string, error) {
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
