package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Result struct {
	Tree         map[string]string `json:"tree" yaml:"tree"`
	Stats        Stats             `json:"stats" yaml:"stats"`
	Dependencies []string          `json:"dependencies" yaml:"dependencies"`
	Version      string            `json:"version" yaml:"version"`
}

type Stats struct {
	FileCount   int  `json:"file_count" yaml:"file_count"`
	TotalBytes  int  `json:"total_bytes" yaml:"total_bytes"`
	TotalTokens int  `json:"total_tokens" yaml:"total_tokens"`
	Truncated   bool `json:"truncated" yaml:"truncated"`
}

func FormatJSON(r Result) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func FormatText(r Result) ([]byte, error) {
	paths := make([]string, 0, len(r.Tree))
	for path := range r.Tree {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var sb strings.Builder
	for _, path := range paths {
		sb.WriteString(fmt.Sprintf("File: %s\n---\n%s\n\n", path, r.Tree[path]))
	}
	return []byte(sb.String()), nil
}

func FormatYAML(r Result) ([]byte, error) {
	return yaml.Marshal(r)
}
