package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"github.com/your-username/gitinspect/internal/deps"
	"github.com/your-username/gitinspect/internal/filter"
	"github.com/your-username/gitinspect/internal/output"
	"github.com/your-username/gitinspect/internal/repo"
	"github.com/your-username/gitinspect/internal/token"
)

// InspectRequest represents the request body for /inspect
type InspectRequest struct {
	RepoURL    string   `json:"repo_url"`
	CommitHash string   `json:"commit_hash,omitempty"`
	Format     string   `json:"format,omitempty"` // json, text, yaml
	MaxTokens  int      `json:"max_tokens,omitempty"`
	MaxFiles   int      `json:"max_files,omitempty"`
	Includes   []string `json:"includes,omitempty"`
	Excludes   []string `json:"excludes,omitempty"`
	Strip      bool     `json:"strip,omitempty"`
	NoCache    bool     `json:"no_cache,omitempty"`
}

// Serve starts the HTTP server
func Serve(port int) error {
	http.HandleFunc("/inspect", inspectHandler)
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func inspectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InspectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set defaults
	if req.Format == "" {
		req.Format = "json"
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 6000
	}

	// Process the repository
	snapshot, err := processRepo(req.RepoURL, req.CommitHash, req.MaxTokens, req.MaxFiles, req.Includes, req.Excludes, req.Strip, req.NoCache)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format the output
	var result string
	switch req.Format {
	case "json":
		result, err = output.ToJSON(*snapshot)
	case "text":
		result, err = output.ToText(*snapshot)
	case "yaml":
		result, err = output.ToYAML(*snapshot)
	default:
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type
	switch req.Format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
	case "yaml":
		w.Header().Set("Content-Type", "text/yaml")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}

func processRepo(repoURL, commitHash string, maxTokens, maxFiles int, includes, excludes []string, strip, noCache bool) (*output.Snapshot, error) {
	var localPath string
	var tempDir string
	var err error
	var isRemote bool

	// Check if it's a local path or remote URL
	if _, err := repo.NewLocalRepo(repoURL).ListFiles(); err == nil {
		localPath = repoURL
	} else {
		isRemote = true
		remoteRepo := repo.NewRemoteRepo(repoURL, commitHash)
		tempDir, err = remoteRepo.Clone()
		if err != nil {
			return nil, err
		}
		defer func() { _ = remoteRepo.Cleanup(tempDir) }()
		localPath = tempDir
	}

	localRepo := repo.NewLocalRepo(localPath)

	// List files
	files, err := localRepo.ListFiles()
	if err != nil {
		return nil, err
	}

	// Filter files
	fil := filter.New(includes, excludes)
	var filteredFiles []string
	for _, f := range files {
		if fil.Match(f) {
			filteredFiles = append(filteredFiles, f)
		}
	}

	// Sort files by priority
	token.SortFilesByPriority(filteredFiles)

	// Limit number of files if needed
	if maxFiles > 0 && len(filteredFiles) > maxFiles {
		filteredFiles = filteredFiles[:maxFiles]
	}

	// Read files and build snapshot
	snapshot := output.Snapshot{
		RepoURL:    repoURL,
		CommitHash: commitHash,
		Files:      []output.File{},
	}

	totalTokens := 0
	var allDeps []deps.Dependency

	for _, path := range filteredFiles {
		content, err := localRepo.ReadFile(path)
		if err != nil {
			continue // Skip unreadable files
		}

		// Strip comments if needed
		if strip {
			ext := filepath.Ext(path)
			content = repo.StripComments(content, ext)
		}

		// Estimate tokens
		tokens := token.Estimate(content)

		// Truncate if needed
		truncatedContent, truncated := content, false
		if totalTokens+tokens > maxTokens {
			truncatedContent, truncated = token.Truncate(content, maxTokens-totalTokens)
			tokens = token.Estimate(truncatedContent)
		}

		// Extract dependencies
		fileDeps, _ := deps.Extract(path, content)
		allDeps = append(allDeps, fileDeps...)

		// Add file to snapshot
		snapshot.Files = append(snapshot.Files, output.File{
			Path:        path,
			Content:     truncatedContent,
			Truncated:   truncated,
			Tokens:      tokens,
			Dependencies: fileDeps,
		})

		totalTokens += tokens

		// Stop if we hit the token limit
		if totalTokens >= maxTokens {
			break
		}
	}

	snapshot.TotalTokens = totalTokens
	snapshot.Dependencies = allDeps

	// Cache the snapshot if remote and not no-cache
	if isRemote && !noCache {
		cache, err := repo.NewCache()
		if err == nil {
			key := repo.GenerateKey(repoURL, commitHash)
			jsonData, _ := json.Marshal(snapshot)
			_ = cache.Set(key, jsonData)
		}
	}

	return &snapshot, nil
}
