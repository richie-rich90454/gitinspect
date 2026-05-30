package internal

import (
	"encoding/json"
)

type cacheData struct {
	Tree         map[string]string `json:"tree"`
	Paths        []string          `json:"paths"`
	TotalBytes   int               `json:"total_bytes"`
	TotalTokens  int               `json:"total_tokens"`
	Truncated    bool              `json:"truncated"`
	Dependencies []string          `json:"dependencies"`
}

func marshalCache(r *inspectResult) ([]byte, error) {
	return json.Marshal(cacheData{
		Tree:         r.tree,
		Paths:        r.paths,
		TotalBytes:   r.totalBytes,
		TotalTokens:  r.totalTokens,
		Truncated:    r.truncated,
		Dependencies: r.dependencies,
	})
}

func unmarshalCache(data []byte, r *inspectResult) error {
	var cd cacheData
	if err := json.Unmarshal(data, &cd); err != nil {
		return err
	}
	r.tree = cd.Tree
	r.paths = cd.Paths
	r.totalBytes = cd.TotalBytes
	r.totalTokens = cd.TotalTokens
	r.truncated = cd.Truncated
	r.dependencies = cd.Dependencies
	return nil
}
