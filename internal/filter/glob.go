// Package filter provides glob-based file matching.
package filter

import "github.com/bmatcuk/doublestar/v4"

// MatchAny returns true if the path matches any include pattern and no exclude pattern.
func MatchAny(path string, includes, excludes []string) (bool, error) {
	for _, pattern := range excludes {
		ok, err := doublestar.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}

	if len(includes) == 0 {
		return true, nil
	}

	for _, pattern := range includes {
		ok, err := doublestar.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}
