package git

import "github.com/belingud/go-gptcomet/internal/config"

// VCSType represents the type of version control system
type VCSType string

const (
	Git VCSType = "git"
	SVN VCSType = "svn"
)

// VCS defines the interface for version control operations
type VCS interface {
	GetDiff(repoPath string) (string, error)
	HasStagedChanges(repoPath string) (bool, error)
	GetStagedFiles(repoPath string) ([]string, error)
	GetStagedDiffFiltered(repoPath string, cfgManager *config.Manager) (string, error)
	GetCurrentBranch(repoPath string) (string, error)
	GetCommitInfo(repoPath, commitHash string) (string, error)
	GetLastCommitHash(repoPath string) (string, error)
	CreateCommit(repoPath, message string) error
}

// NewVCS creates a new VCS instance based on the type
func NewVCS(vcsType VCSType) (VCS, error) {
	switch vcsType {
	case Git:
		return &GitVCS{}, nil
	case SVN:
		return &SVNVCS{}, nil
	default:
		return &GitVCS{}, nil
	}
}
