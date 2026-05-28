package repo

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheDirName = "gitinspect"
	cacheTTL     = 1 * time.Hour
)

// CacheEntry represents a cached snapshot
type CacheEntry struct {
	Key       string    `json:"key"`
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
}

// Cache handles caching snapshots
type Cache struct {
	dir string
}

// NewCache creates a new Cache
func NewCache() (*Cache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cacheDir := filepath.Join(homeDir, ".cache", cacheDirName)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}
	return &Cache{dir: cacheDir}, nil
}

// GenerateKey generates a cache key from repo URL and commit hash
func GenerateKey(repoURL, commitHash string) string {
	data := repoURL + "#" + commitHash
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// Get retrieves a cached entry
func (c *Cache) Get(key string) (*CacheEntry, error) {
	path := filepath.Join(c.dir, key+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}

	// Check TTL
	if time.Since(entry.Timestamp) > cacheTTL {
		_ = os.Remove(path)
		return nil, fmt.Errorf("cache expired")
	}

	return &entry, nil
}

// Set stores a cache entry
func (c *Cache) Set(key string, data []byte) error {
	entry := CacheEntry{
		Key:       key,
		Timestamp: time.Now(),
		Data:      data,
	}
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	path := filepath.Join(c.dir, key+".json")
	return os.WriteFile(path, jsonData, 0644)
}
