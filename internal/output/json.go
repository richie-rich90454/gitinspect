package output

import (
	"encoding/json"
	"github.com/your-username/gitinspect/internal/deps"
)

// File represents a file in the snapshot
type File struct {
	Path        string          `json:"path"`
	Content     string          `json:"content"`
	Truncated   bool            `json:"truncated"`
	Tokens      int             `json:"tokens"`
	Dependencies []deps.Dependency `json:"dependencies,omitempty"`
}

// Snapshot represents the repository snapshot
type Snapshot struct {
	RepoURL      string          `json:"repo_url,omitempty"`
	CommitHash   string          `json:"commit_hash,omitempty"`
	Files        []File          `json:"files"`
	TotalTokens  int             `json:"total_tokens"`
	Dependencies []deps.Dependency `json:"dependencies,omitempty"`
}

// ToJSON converts the snapshot to JSON
func ToJSON(snapshot Snapshot) (string, error) {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON parses a JSON snapshot
func FromJSON(data []byte) (*Snapshot, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
