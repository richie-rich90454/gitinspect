package repo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

func ResolveHEAD(url string) (string, error) {
	cmd := exec.Command("git", "ls-remote", url, "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ls-remote failed: %w", err)
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return "", fmt.Errorf("no HEAD found")
	}
	return fields[0], nil
}

func ShallowClone(url, tmpDir string) error {
	cmd := exec.Command("git", "clone", "--depth=1", "--filter=blob:none", url, tmpDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %s: %w", string(out), err)
	}
	return nil
}

func FetchRemote(url string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-")
	if err != nil {
		return "", err
	}

	if err := ShallowClone(url, tmpDir); err != nil {
		_ = os.RemoveAll(tmpDir)

		cloneOpts := &git.CloneOptions{
			URL:               url,
			Depth:             1,
			SingleBranch:     true,
			RecurseSubmodules: git.NoRecurseSubmodules,
		}
		if _, err := git.PlainClone(tmpDir, false, cloneOpts); err != nil {
			_ = os.RemoveAll(tmpDir)
			return "", fmt.Errorf("clone failed (exec and go-git): %w", err)
		}
	}

	return tmpDir, nil
}

func Cleanup(dir string) {
	_ = os.RemoveAll(dir)
}
