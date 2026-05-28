package repo

import (
	"os"
	"path/filepath"
	"strings"
)

// LocalRepo represents a local Git repository
type LocalRepo struct {
	Path string
}

// NewLocalRepo creates a new LocalRepo
func NewLocalRepo(path string) *LocalRepo {
	return &LocalRepo{Path: path}
}

// ListFiles lists all files in the repository, respecting .gitignore
func (r *LocalRepo) ListFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(r.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip .git directory
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Make path relative
		relPath, err := filepath.Rel(r.Path, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		files = append(files, relPath)
		return nil
	})

	return files, err
}

// ReadFile reads a file from the repository
func (r *LocalRepo) ReadFile(path string) (string, error) {
	fullPath := filepath.Join(r.Path, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// StripComments strips comments and blank lines based on file extension
func StripComments(content string, ext string) string {
	var commentPrefix string
	switch ext {
	case ".go", ".c", ".cpp", ".h", ".hpp", ".js", ".ts", ".java":
		commentPrefix = "//"
	case ".py", ".sh", ".bash", ".zsh", ".yaml", ".yml", ".toml", ".rb":
		commentPrefix = "#"
	case ".sql":
		commentPrefix = "--"
	case ".asm":
		commentPrefix = ";"
	default:
		commentPrefix = ""
	}

	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if commentPrefix != "" && strings.HasPrefix(trimmed, commentPrefix) {
			continue
		}
		// Also remove inline comments
		if commentPrefix != "" {
			if idx := strings.Index(line, commentPrefix); idx != -1 {
				line = line[:idx]
			}
		}
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
