package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/richie-rich90454/gitinspect/internal/deps"
	"github.com/richie-rich90454/gitinspect/internal/filter"
	"github.com/richie-rich90454/gitinspect/internal/output"
	"github.com/richie-rich90454/gitinspect/internal/repo"
	"github.com/richie-rich90454/gitinspect/internal/token"
)

const Version = "0.1.0"

func PriorityScore(path string) int {
	return token.PriorityScore(path)
}

type Options struct {
	Format    string
	MaxTokens int
	MaxFiles  int
	Include   []string
	Exclude   []string
	Strip     bool
	NoCache   bool
}

type inspectResult struct {
	tree         map[string]string
	paths        []string
	totalBytes   int
	totalTokens  int
	truncated    bool
	dependencies []string
}

func runInspectCore(repoArg string, opts Options) (*inspectResult, error) {
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
				var cached inspectResult
				if unmarshalCache(data, &cached) == nil {
					return &cached, nil
				}
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

	entryMap := make(map[string]string, len(entries))
	for _, e := range entries {
		entryMap[e.Path] = string(e.Content)
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
		content, ok := entryMap[p]
		if !ok || content == "" {
			continue
		}

		if opts.Strip {
			ext := filepath.Ext(p)
			content = repo.StripComments(content, ext)
		}

		tokens := token.Estimate(content)
		if totalTokens+tokens > opts.MaxTokens {
			truncatedContent := token.Truncate(content)
			if truncatedContent != content {
				content = truncatedContent
				truncated = true
			} else {
				if totalTokens > 0 {
					truncated = true
					break
				}
			}
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

	sort.Strings(allDeps)

	result := &inspectResult{
		tree:         tree,
		paths:        paths,
		totalBytes:   totalBytes,
		totalTokens:  totalTokens,
		truncated:    truncated,
		dependencies: allDeps,
	}

	if isRemote && !opts.NoCache && commitHash != "" {
		cache := repo.NewCache()
		key := repo.GenerateKey(repoArg, commitHash)
		if data, err := marshalCache(result); err == nil {
			_ = cache.Set(key, data)
		}
	}

	return result, nil
}

func RunInspect(repoArg string, opts Options) ([]byte, error) {
	res, err := runInspectCore(repoArg, opts)
	if err != nil {
		return nil, err
	}

	result := output.Result{
		Tree: res.tree,
		Stats: output.Stats{
			FileCount:   len(res.tree),
			TotalBytes:  res.totalBytes,
			TotalTokens: res.totalTokens,
			Truncated:   res.truncated,
		},
		Dependencies: res.dependencies,
		Version:      Version,
	}

	switch strings.ToLower(opts.Format) {
	case "text":
		return output.FormatText(result)
	case "yaml":
		return output.FormatYAML(result)
	default:
		return output.FormatJSON(result)
	}
}

func RunInspectRaw(repoArg string, opts Options) (*output.Result, error) {
	res, err := runInspectCore(repoArg, opts)
	if err != nil {
		return nil, err
	}

	return &output.Result{
		Tree: res.tree,
		Stats: output.Stats{
			FileCount:   len(res.tree),
			TotalBytes:  res.totalBytes,
			TotalTokens: res.totalTokens,
			Truncated:   res.truncated,
		},
		Dependencies: res.dependencies,
		Version:      Version,
	}, nil
}
