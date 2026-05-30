package repo

import (
	"fmt"
	"os"
)

// ResolveRepo determines whether arg is a local path or remote URL and returns the local path,
// temp directory (if remote), whether it's remote, the commit hash (if remote), and any error.
// If remote, the caller MUST call repo.Cleanup on the returned tmpDir when done.
func ResolveRepo(arg string) (localPath, tmpDir string, isRemote bool, commitHash string, err error) {
	info, statErr := os.Stat(arg)
	if statErr == nil {
		if info.IsDir() {
			return arg, "", false, "", nil
		}
		return "", "", false, "", fmt.Errorf("%q is not a directory", arg)
	}
	if !os.IsNotExist(statErr) {
		return "", "", false, "", fmt.Errorf("cannot access %q: %w", arg, statErr)
	}

	if err := ValidateRemoteURL(arg); err != nil {
		return "", "", false, "", err
	}

	head, headErr := ResolveHEAD(arg)
	if headErr == nil {
		commitHash = head
	}

	tmp, fetchErr := FetchRemote(arg)
	if fetchErr != nil {
		return "", "", false, "", fmt.Errorf("fetch remote: %w", fetchErr)
	}

	return tmp, tmp, true, commitHash, nil
}
