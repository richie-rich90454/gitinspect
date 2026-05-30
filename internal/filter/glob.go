// Package filter provides glob-based file matching.
package filter

import "github.com/bmatcuk/doublestar/v4"

// MatchAny returns true if the path matches any include pattern and no exclude pattern.
func MatchAny(path string, includes, excludes []string) bool {
	for _, pattern := range excludes {
		if ok, _ := doublestar.Match(pattern, path); ok {
			return false
		}
	}

	if len(includes) == 0 {
		return true
	}

	for _, pattern := range includes {
		if ok, _ := doublestar.Match(pattern, path); ok {
			return true
		}
	}

	return false
}
