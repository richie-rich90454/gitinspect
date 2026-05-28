package token

import (
	"path/filepath"
	"strings"
)

// FilePriority represents the priority of a file
type FilePriority int

const (
	PriorityLow FilePriority = iota
	PriorityMedium
	PriorityHigh
	PriorityVeryHigh
	PriorityHighest
)

// GetPriority returns the priority of a file path
func GetPriority(path string) FilePriority {
	base := filepath.Base(path)

	// Priority 1: README*
	if strings.HasPrefix(strings.ToLower(base), "readme") {
		return PriorityHighest
	}

	// Priority 2: main.* and cmd/*/main.*
	if base == "main.go" || base == "main.py" || base == "main.js" || base == "main.ts" || base == "main.c" || base == "main.cpp" || base == "main.rs" {
		return PriorityVeryHigh
	}
	if strings.HasPrefix(path, "cmd/") && strings.HasSuffix(base, "main.go") {
		return PriorityVeryHigh
	}

	// Priority 3: manifest files
	manifestFiles := []string{"go.mod", "package.json", "Cargo.toml", "setup.py", "Gemfile", "requirements.txt", "pom.xml", "build.gradle"}
	for _, mf := range manifestFiles {
		if base == mf {
			return PriorityHigh
		}
	}

	// Priority 4: build files
	buildFiles := []string{"Makefile", "Dockerfile", "docker-compose.yml", "docker-compose.yaml"}
	for _, bf := range buildFiles {
		if base == bf {
			return PriorityMedium
		}
	}

	return PriorityLow
}

// SortFilesByPriority sorts a slice of file paths by priority (highest first)
func SortFilesByPriority(files []string) {
	for i := range files {
		for j := i + 1; j < len(files); j++ {
			if GetPriority(files[j]) > GetPriority(files[i]) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
}
