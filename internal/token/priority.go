package token

import (
	"path/filepath"
	"strings"
)

// PriorityScore returns a priority score (0-4) for a file based on its path.
func PriorityScore(path string) int {
	base := filepath.Base(path)

	if strings.HasPrefix(strings.ToLower(base), "readme") {
		return 4
	}

	if strings.HasPrefix(path, "cmd/") && strings.HasPrefix(base, "main.") {
		return 3
	}
	if strings.HasPrefix(base, "main.") {
		return 3
	}

	manifests := []string{"go.mod", "package.json", "Cargo.toml", "setup.py", "Gemfile", "requirements.txt"}
	for _, m := range manifests {
		if base == m {
			return 2
		}
	}

	builds := []string{"Makefile", "Dockerfile", "docker-compose.yml", "docker-compose.yaml"}
	for _, b := range builds {
		if base == b {
			return 1
		}
	}

	return 0
}

type scoredFile struct {
	path  string
	score int
}

// SortByPriority sorts files by descending priority score.
func SortByPriority(files []string) {
	scored := make([]scoredFile, len(files))
	for i, f := range files {
		scored[i] = scoredFile{path: f, score: PriorityScore(f)}
	}

	for i := 1; i < len(scored); i++ {
		for j := i; j > 0 && scored[j].score > scored[j-1].score; j-- {
			scored[j], scored[j-1] = scored[j-1], scored[j]
		}
	}

	for i, s := range scored {
		files[i] = s.path
	}
}
