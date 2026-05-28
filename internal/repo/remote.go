package repo

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// RemoteRepo represents a remote Git repository
type RemoteRepo struct {
	URL    string
	Commit string
}

// NewRemoteRepo creates a new RemoteRepo
func NewRemoteRepo(url string, commit string) *RemoteRepo {
	return &RemoteRepo{URL: url, Commit: commit}
}

// Clone clones the remote repository to a temporary directory
func (r *RemoteRepo) Clone() (string, error) {
	tempDir, err := os.MkdirTemp("", "gitinspect-")
	if err != nil {
		return "", err
	}

	cloneOpts := &git.CloneOptions{
		URL:               r.URL,
		Depth:             1,
		SingleBranch:     true,
		RecurseSubmodules: git.NoRecurseSubmodules,
	}

	repo, err := git.PlainClone(tempDir, false, cloneOpts)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return "", err
	}

	// If a specific commit is requested, check it out
	if r.Commit != "" {
		w, err := repo.Worktree()
		if err != nil {
			_ = os.RemoveAll(tempDir)
			return "", err
		}
		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(r.Commit),
		})
		if err != nil {
			_ = os.RemoveAll(tempDir)
			return "", err
		}
	}

	return tempDir, nil
}

// Cleanup cleans up the temporary directory
func (r *RemoteRepo) Cleanup(tempDir string) error {
	return os.RemoveAll(tempDir)
}
