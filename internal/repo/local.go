package repo

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type FileEntry struct {
	Path    string
	Content []byte
	Size    int
}

func ReadLocalRepo(repoPath string) ([]FileEntry, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return readDir(repoPath, nil)
	}

	w, err := r.Worktree()
	if err != nil {
		return readDir(repoPath, nil)
	}

	var patterns []gitignore.Pattern
	if f, err := w.Filesystem.Open(".gitignore"); err == nil {
		defer f.Close()
		data, err := io.ReadAll(f)
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					patterns = append(patterns, gitignore.ParsePattern(line, nil))
				}
			}
		}
	}
	matcher := gitignore.NewMatcher(patterns)

	return readDir(repoPath, matcher)
}

func readDir(root string, matcher gitignore.Matcher) ([]FileEntry, error) {
	var entries []FileEntry

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		if matcher != nil {
			parts := strings.Split(relPath, "/")
			if matcher.Match(parts, false) {
				return nil
			}
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		entries = append(entries, FileEntry{
			Path:    relPath,
			Content: data,
			Size:    len(data),
		})
		return nil
	})

	return entries, err
}

func StripComments(content string, ext string) string {
	var prefix string
	switch ext {
	case ".go", ".c", ".cpp", ".h", ".hpp", ".js", ".ts", ".java", ".rs":
		prefix = "//"
	case ".py", ".sh", ".bash", ".zsh", ".yaml", ".yml", ".toml", ".rb":
		prefix = "#"
	case ".sql":
		prefix = "--"
	case ".asm":
		prefix = ";"
	default:
		return content
	}

	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, prefix) {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
