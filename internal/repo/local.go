package repo

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// FileEntry holds a file's relative path, content, and size.
type FileEntry struct {
	Path    string
	Content []byte
	Size    int
}

var binaryExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true,
	".svg": true, ".webp": true, ".bmp": true, ".tiff": true, ".tif": true,
	".exe": true, ".dll": true, ".so": true, ".dylib": true, ".a": true,
	".o": true, ".obj": true, ".pyc": true, ".pyd": true, ".class": true,
	".jar": true, ".war": true, ".zip": true, ".tar": true, ".gz": true,
	".bz2": true, ".xz": true, ".7z": true, ".rar": true, ".deb": true,
	".rpm": true, ".apk": true, ".dmg": true, ".iso": true, ".woff": true,
	".woff2": true, ".eot": true, ".ttf": true, ".otf": true, ".mp3": true,
	".mp4": true, ".avi": true, ".mov": true, ".wmv": true, ".flac": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true, ".sqlite": true, ".db": true,
}

func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	checkLen := len(data)
	if checkLen > 8000 {
		checkLen = 8000
	}
	for i := 0; i < checkLen; i++ {
		if data[i] == 0 { //nolint:gosec
			return true
		}
	}
	return false
}

// ReadLocalRepo reads all text files from a local Git repository or directory.
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
		data, readErr := io.ReadAll(f)
		closeErr := f.Close()
		if readErr == nil && closeErr == nil {
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

		ext := strings.ToLower(filepath.Ext(relPath))
		if binaryExts[ext] {
			return nil
		}

		data, err := os.ReadFile(path) //nolint:gosec
		if err != nil {
			return nil
		}

		if isBinary(data) {
			return nil
		}

		if !utf8.Valid(data) {
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

// StripComments removes single-line comments and blank lines based on file extension.
func StripComments(content string, ext string) string {
	var prefix string
	switch ext {
	case ".go", ".c", ".cpp", ".h", ".hpp", ".js", ".ts", ".java", ".rs",
		".kt", ".swift", ".tsx", ".jsx":
		prefix = "//"
	case ".py", ".sh", ".bash", ".zsh", ".yaml", ".yml", ".toml", ".rb",
		".r", ".ps1", ".perl", ".pl":
		prefix = "#"
	case ".sql", ".lua", ".hs", ".elm":
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
