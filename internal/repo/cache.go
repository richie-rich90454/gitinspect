// Package repo provides local and remote repository reading with caching.
package repo

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 1 * time.Hour

// Cache provides a disk-based cache keyed by SHA-256 hash.
type Cache struct {
	Dir string
}

// NewCache creates a Cache using ~/.cache/gitinspect/.
func NewCache() *Cache {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Cache{Dir: ""}
	}
	dir := filepath.Join(homeDir, ".cache", "gitinspect")
	_ = os.MkdirAll(dir, 0755)
	return &Cache{Dir: dir}
}

// GenerateKey returns a SHA-256 key from a repo URL and commit hash.
func GenerateKey(repoURL, commitHash string) string {
	data := repoURL + "\x00" + commitHash
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// Get retrieves cached data by key. Returns false on miss or expiry.
func (c *Cache) Get(key string) ([]byte, bool) {
	if c.Dir == "" {
		return nil, false
	}
	path := filepath.Join(c.Dir, key)

	info, err := os.Stat(path)
	if err != nil {
		return nil, false
	}
	if time.Since(info.ModTime()) > cacheTTL {
		_ = os.Remove(path)
		return nil, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return data, true
}

// Set writes data to the cache under the given key.
func (c *Cache) Set(key string, data []byte) error {
	if c.Dir == "" {
		return nil
	}
	path := filepath.Join(c.Dir, key)
	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		_ = os.Remove(path)
		if err2 := os.WriteFile(path, data, 0644); err2 != nil {
			return err2
		}
	}
	return nil
}
