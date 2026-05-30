package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/richie-rich90454/gitinspect/internal/output"
)

func TestRunInspectLocal(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "json",
		MaxTokens: 6000,
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	var result output.Result
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	if len(result.Tree) == 0 {
		t.Error("expected non-empty tree")
	}
	if result.Stats.FileCount == 0 {
		t.Error("expected non-zero file count")
	}
	if result.Version == "" {
		t.Error("expected version to be set")
	}
}

func TestRunInspectTextFormat(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "text",
		MaxTokens: 6000,
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRunInspectYAMLFormat(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "yaml",
		MaxTokens: 6000,
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	if len(out) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRunInspectWithStrip(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "json",
		MaxTokens: 6000,
		Strip:     true,
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	var result output.Result
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	if len(result.Tree) == 0 {
		t.Error("expected non-empty tree with strip")
	}
}

func TestRunInspectMaxFiles(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "json",
		MaxTokens: 6000,
		MaxFiles:  1,
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	var result output.Result
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	if result.Stats.FileCount > 1 {
		t.Errorf("expected at most 1 file, got %d", result.Stats.FileCount)
	}
}

func TestRunInspectNonexistentPath(t *testing.T) {
	opts := Options{
		Format:    "json",
		MaxTokens: 6000,
		NoCache:   true,
	}

	_, err := RunInspect("/nonexistent/path/that/does/not/exist", opts)
	if err == nil {
		t.Error("expected error for nonexistent path")
	}
}

func TestRunInspectWithInclude(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "json",
		MaxTokens: 6000,
		Include:   []string{"**/*.go"},
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	var result output.Result
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	for path := range result.Tree {
		if filepath.Ext(path) != ".go" {
			t.Errorf("expected only .go files, got %s", path)
		}
	}
}

func TestRunInspectMaxTokensDefaultsTo6000(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "json",
		MaxTokens: 0,
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	var result output.Result
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	if len(result.Tree) == 0 {
		t.Error("expected non-empty tree with MaxTokens=0 (should default to 6000)")
	}

	opts.MaxTokens = -1
	out, err = RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Tree) == 0 {
		t.Error("expected non-empty tree with MaxTokens=-1 (should default to 6000)")
	}
}

func TestRunInspectExcludeOnly(t *testing.T) {
	repoPath := filepath.Join("..", "testdata", "sample-repo")
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	opts := Options{
		Format:    "json",
		MaxTokens: 6000,
		Exclude:   []string{"**/*.md"},
		NoCache:   true,
	}

	out, err := RunInspect(absPath, opts)
	if err != nil {
		t.Fatal(err)
	}

	var result output.Result
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	for path := range result.Tree {
		if filepath.Ext(path) == ".md" {
			t.Errorf("expected .md files to be excluded, got %s", path)
		}
	}
}

func TestRunInspectEmptyRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitinspect-empty-repo-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	opts := Options{
		Format:    "json",
		MaxTokens: 6000,
		Include:   []string{"**/*.xyz"},
		NoCache:   true,
	}

	out, err := RunInspect(tmpDir, opts)
	if err != nil {
		t.Fatal(err)
	}

	var result output.Result
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	if len(result.Tree) != 0 {
		t.Errorf("expected empty tree for repo with no matching files, got %d files", len(result.Tree))
	}
}
