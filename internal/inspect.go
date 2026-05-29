package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-username/gitinspect/internal/deps"
	"github.com/your-username/gitinspect/internal/filter"
	"github.com/your-username/gitinspect/internal/output"
	"github.com/your-username/gitinspect/internal/repo"
	"github.com/your-username/gitinspect/internal/token"
)

const Version = "0.1.0"

type Options struct {
	Format    string
	MaxTokens int
	MaxFiles  int
	Include   []string
	Exclude   []string
	Strip     bool
	NoCache   bool
}

func RunInspect(repoArg string, opts Options) ([]byte, error) {
	var localPath string
	var tmpDir string
	var isRemote bool
	var commitHash string

	if _, err := os.Stat(repoArg); err == nil {
		localPath = repoArg
	} else {
		isRemote = true
		head, err := repo.ResolveHEAD(repoArg)
		if err == nil {
			commitHash = head
		}

		cache := repo.NewCache()
		if !opts.NoCache && commitHash != "" {
			key := repo.GenerateKey(repoArg, commitHash)
			if data, ok := cache.Get(key); ok {
				return data, nil
			}
		}

		tmpDir, err = repo.FetchRemote(repoArg)
		if err != nil {
			return nil, fmt.Errorf("fetch remote: %w", err)
		}
		defer repo.Cleanup(tmpDir)
		localPath = tmpDir
	}

	entries, err := repo.ReadLocalRepo(localPath)
	if err != nil {
		return nil, fmt.Errorf("read repo: %w", err)
	}

	var paths []string
	for _, e := range entries {
		if filter.MatchAny(e.Path, opts.Include, opts.Exclude) {
			paths = append(paths, e.Path)
		}
	}

	token.SortByPriority(paths)

	if opts.MaxFiles > 0 && len(paths) > opts.MaxFiles {
		paths = paths[:opts.MaxFiles]
	}

	tree := make(map[string]string)
	totalBytes := 0
	totalTokens := 0
	truncated := false
	var allDeps []string

	for _, p := range paths {
		var content string
		for _, e := range entries {
			if e.Path == p {
				content = string(e.Content)
				break
			}
		}
		if content == "" {
			continue
		}

		if opts.Strip {
			ext := filepath.Ext(p)
			content = repo.StripComments(content, ext)
		}

		tokens := token.Estimate(content)
		if totalTokens+tokens > opts.MaxTokens {
			content = token.Truncate(content)
			truncated = true
		}

		fileDeps := deps.Extract(p, content)
		allDeps = append(allDeps, fileDeps...)

		tree[p] = content
		totalBytes += len(content)
		totalTokens += token.Estimate(content)

		if totalTokens >= opts.MaxTokens {
			truncated = true
			break
		}
	}

	result := output.Result{
		Tree: tree,
		Stats: output.Stats{
			FileCount:  len(tree),
			TotalBytes: totalBytes,
			Truncated:  truncated,
		},
		Dependencies: allDeps,
		Version:      Version,
	}

	var out []byte
	switch strings.ToLower(opts.Format) {
	case "text":
		out, err = output.FormatText(result)
	case "yaml":
		out, err = output.FormatYAML(result)
	default:
		out, err = output.FormatJSON(result)
	}
	if err != nil {
		return nil, err
	}

	if isRemote && !opts.NoCache && commitHash != "" {
		cache := repo.NewCache()
		key := repo.GenerateKey(repoArg, commitHash)
		_ = cache.Set(key, out)
	}

	return out, nil
}

func RunInspectRaw(repoArg string, opts Options) (*output.Result, error) {
	var localPath string
	var tmpDir string
	var isRemote bool
	var commitHash string

	if _, err := os.Stat(repoArg); err == nil {
		localPath = repoArg
	} else {
		isRemote = true
		head, err := repo.ResolveHEAD(repoArg)
		if err == nil {
			commitHash = head
		}
		tmpDir, err = repo.FetchRemote(repoArg)
		if err != nil {
			return nil, fmt.Errorf("fetch remote: %w", err)
		}
		defer repo.Cleanup(tmpDir)
		localPath = tmpDir
	}

	entries, err := repo.ReadLocalRepo(localPath)
	if err != nil {
		return nil, fmt.Errorf("read repo: %w", err)
	}

	var paths []string
	for _, e := range entries {
		if filter.MatchAny(e.Path, opts.Include, opts.Exclude) {
			paths = append(paths, e.Path)
		}
	}

	token.SortByPriority(paths)

	if opts.MaxFiles > 0 && len(paths) > opts.MaxFiles {
		paths = paths[:opts.MaxFiles]
	}

	tree := make(map[string]string)
	totalBytes := 0
	totalTokens := 0
	truncated := false
	var allDeps []string

	for _, p := range paths {
		var content string
		for _, e := range entries {
			if e.Path == p {
				content = string(e.Content)
				break
			}
		}
		if content == "" {
			continue
		}

		if opts.Strip {
			ext := filepath.Ext(p)
			content = repo.StripComments(content, ext)
		}

		tokens := token.Estimate(content)
		if totalTokens+tokens > opts.MaxTokens {
			content = token.Truncate(content)
			truncated = true
		}

		fileDeps := deps.Extract(p, content)
		allDeps = append(allDeps, fileDeps...)

		tree[p] = content
		totalBytes += len(content)
		totalTokens += token.Estimate(content)

		if totalTokens >= opts.MaxTokens {
			truncated = true
			break
		}
	}

	result := &output.Result{
		Tree: tree,
		Stats: output.Stats{
			FileCount:  len(tree),
			TotalBytes: totalBytes,
			Truncated:  truncated,
		},
		Dependencies: allDeps,
		Version:      Version,
	}

	if isRemote && !opts.NoCache && commitHash != "" {
		cache := repo.NewCache()
		key := repo.GenerateKey(repoArg, commitHash)
		data, _ := json.Marshal(result)
		_ = cache.Set(key, data)
	}

	return result, nil
}
