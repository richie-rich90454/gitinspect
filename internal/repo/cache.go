package repo

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 1 * time.Hour

type Cache struct {
	Dir string
}

func NewCache() *Cache {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Cache{Dir: ""}
	}
	dir := filepath.Join(homeDir, ".cache", "gitinspect")
	_ = os.MkdirAll(dir, 0755)
	return &Cache{Dir: dir}
}

func GenerateKey(repoURL, commitHash string) string {
	data := repoURL + "#" + commitHash
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

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

func (c *Cache) Set(key string, data []byte) error {
	if c.Dir == "" {
		return nil
	}
	path := filepath.Join(c.Dir, key)
	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
