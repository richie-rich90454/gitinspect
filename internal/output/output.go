package output

import (
	"encoding/json"
	"fmt"
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
	FileCount  int  `json:"file_count" yaml:"file_count"`
	TotalBytes int  `json:"total_bytes" yaml:"total_bytes"`
	Truncated  bool `json:"truncated" yaml:"truncated"`
}

func FormatJSON(r Result) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func FormatText(r Result) ([]byte, error) {
	var sb strings.Builder
	for path, content := range r.Tree {
		sb.WriteString(fmt.Sprintf("File: %s\n---\n%s\n\n", path, content))
	}
	return []byte(sb.String()), nil
}

func FormatYAML(r Result) ([]byte, error) {
	return yaml.Marshal(r)
}
