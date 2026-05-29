package repo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheSetGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-cache-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	c := &Cache{Dir: tmpDir}
	key := "testkey123"

	data := []byte("hello cache")
	if err := c.Set(key, data); err != nil {
		t.Fatal(err)
	}

	got, ok := c.Get(key)
	if !ok {
		t.Fatal("cache miss after set")
	}
	if string(got) != string(data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestCacheMiss(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-cache-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	c := &Cache{Dir: tmpDir}
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected cache miss for nonexistent key")
	}
}

func TestGenerateKey(t *testing.T) {
	k1 := GenerateKey("https://github.com/foo/bar.git", "abc123")
	k2 := GenerateKey("https://github.com/foo/bar.git", "abc123")
	k3 := GenerateKey("https://github.com/foo/bar.git", "def456")

	if k1 != k2 {
		t.Error("same inputs should produce same key")
	}
	if k1 == k3 {
		t.Error("different inputs should produce different keys")
	}
}

func TestReadLocalRepo(t *testing.T) {
	repoPath := filepath.Join("..", "..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	entries, err := ReadLocalRepo(absPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one file entry")
	}

	found := false
	for _, e := range entries {
		if e.Path == "main.go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected main.go in entries")
	}
}
