package filter

import (
	"path/filepath"
	"strings"
)

// Filter handles include/exclude patterns and .gitignore
type Filter struct {
	includes []string
	excludes []string
}

// New creates a new Filter
func New(includes, excludes []string) *Filter {
	return &Filter{
		includes: includes,
		excludes: excludes,
	}
}

// Match checks if a path should be included
func (f *Filter) Match(path string) bool {
	// Check excludes first
	for _, pattern := range f.excludes {
		if matchGlob(pattern, path) {
			return false
		}
	}

	// If no includes specified, include everything
	if len(f.includes) == 0 {
		return true
	}

	// Check includes
	for _, pattern := range f.includes {
		if matchGlob(pattern, path) {
			return true
		}
	}

	return false
}

// matchGlob is a simple glob matcher (supports * and **)
func matchGlob(pattern, path string) bool {
	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)

	// Split pattern into parts
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	return matchParts(patternParts, pathParts)
}

func matchParts(patternParts, pathParts []string) bool {
	if len(patternParts) == 0 {
		return len(pathParts) == 0
	}

	if patternParts[0] == "**" {
		// ** matches zero or more directories
		if len(patternParts) == 1 {
			return true
		}
		for i := range pathParts {
			if matchParts(patternParts[1:], pathParts[i:]) {
				return true
			}
		}
		return matchParts(patternParts[1:], pathParts)
	}

	if len(pathParts) == 0 {
		return false
	}

	if matchPart(patternParts[0], pathParts[0]) {
		return matchParts(patternParts[1:], pathParts[1:])
	}

	return false
}

func matchPart(pattern, part string) bool {
	// Simple * matching (not full regex)
	// For simplicity, treat * as any sequence of characters
	if pattern == "*" {
		return true
	}
	if strings.Contains(pattern, "*") {
		// Split on * and check if part starts with first part and ends with last part
		parts := strings.Split(pattern, "*")
		if len(parts) == 0 {
			return true
		}
		if !strings.HasPrefix(part, parts[0]) {
			return false
		}
		if len(parts) == 1 {
			return strings.HasPrefix(part, parts[0])
		}
		return strings.HasSuffix(part, parts[len(parts)-1])
	}
	return pattern == part
}
