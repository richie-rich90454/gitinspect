package mcpserver

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/richie-rich90454/gitinspect/internal/filter"
	internal "github.com/richie-rich90454/gitinspect/internal"
	"github.com/richie-rich90454/gitinspect/internal/repo"
	"github.com/richie-rich90454/gitinspect/internal/token"
)

func ServeStdio(s *server.MCPServer) error {
	return server.ServeStdio(s)
}

func NewServer() *server.MCPServer {
	s := server.NewMCPServer(
		"gitinspect",
		internal.Version,
		server.WithToolCapabilities(false),
	)

	inspectTool := mcp.NewTool("inspect_repo",
		mcp.WithDescription("Inspect a Git repository (local path or remote URL) and return a structured, token-efficient snapshot optimized for LLMs. Returns file contents, dependency lists, and repository stats."),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Local path or remote URL of the Git repository to inspect"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json, text, or yaml (default: json)"),
		),
		mcp.WithNumber("max_tokens",
			mcp.Description("Maximum token budget (default: 6000). Files are sorted by priority and truncated to fit."),
		),
		mcp.WithNumber("max_files",
			mcp.Description("Maximum number of files to include (0 = unlimited)"),
		),
		mcp.WithString("include",
			mcp.Description("Glob pattern to include (e.g. '**/*.go'). Can be comma-separated for multiple patterns."),
		),
		mcp.WithString("exclude",
			mcp.Description("Glob pattern to exclude (e.g. '**/vendor/*'). Can be comma-separated for multiple patterns."),
		),
		mcp.WithBoolean("strip",
			mcp.Description("Strip comments and blank lines from file contents"),
		),
		mcp.WithBoolean("no_cache",
			mcp.Description("Disable caching of remote repo snapshots"),
		),
	)

	s.AddTool(inspectTool, inspectHandler)

	listFilesTool := mcp.NewTool("list_repo_files",
		mcp.WithDescription("List all files in a Git repository with their priority scores and token estimates, without reading file contents. Useful for deciding which files to inspect."),
		mcp.WithString("repo",
			mcp.Required(),
			mcp.Description("Local path or remote URL of the Git repository"),
		),
		mcp.WithString("include",
			mcp.Description("Glob pattern to include (e.g. '**/*.go'). Can be comma-separated."),
		),
		mcp.WithString("exclude",
			mcp.Description("Glob pattern to exclude (e.g. '**/vendor/*'). Can be comma-separated."),
		),
	)

	s.AddTool(listFilesTool, listFilesHandler)

	return s
}

func inspectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	repoArg, err := request.RequireString("repo")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	opts := internal.Options{
		Format:    request.GetString("format", "json"),
		MaxTokens: request.GetInt("max_tokens", 6000),
		MaxFiles:  request.GetInt("max_files", 0),
		Strip:     request.GetBool("strip", false),
		NoCache:   request.GetBool("no_cache", false),
	}

	if include := request.GetString("include", ""); include != "" {
		opts.Include = splitPatterns(include)
	}
	if exclude := request.GetString("exclude", ""); exclude != "" {
		opts.Exclude = splitPatterns(exclude)
	}

	out, err := internal.RunInspect(repoArg, opts)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("inspect failed: %v", err)), nil
	}

	return mcp.NewToolResultText(string(out)), nil
}

func listFilesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	repoArg, err := request.RequireString("repo")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var includes, excludes []string
	if include := request.GetString("include", ""); include != "" {
		includes = splitPatterns(include)
	}
	if exclude := request.GetString("exclude", ""); exclude != "" {
		excludes = splitPatterns(exclude)
	}

	var localPath string
	var tmpDir string

	if _, statErr := os.Stat(repoArg); statErr == nil {
		localPath = repoArg
	} else {
		tmpDir, err = repo.FetchRemote(repoArg)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch remote: %v", err)), nil
		}
		defer repo.Cleanup(tmpDir)
		localPath = tmpDir
	}

	entries, err := repo.ReadLocalRepo(localPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("read repo: %v", err)), nil
	}

	type fileInfo struct {
		path     string
		size     int
		priority int
	}

	var files []fileInfo
	for _, e := range entries {
		if !filter.MatchAny(e.Path, includes, excludes) {
			continue
		}
		files = append(files, fileInfo{
			path:     e.Path,
			size:     e.Size,
			priority: token.PriorityScore(e.Path),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].priority != files[j].priority {
			return files[i].priority > files[j].priority
		}
		return files[i].path < files[j].path
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\nFiles: %d\n\n", repoArg, len(files)))
	for _, f := range files {
		estTokens := (f.size + 3) / 4
		sb.WriteString(fmt.Sprintf("  [priority=%d, ~%d tokens, %d bytes] %s\n", f.priority, estTokens, f.size, f.path))
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func splitPatterns(s string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := trimSpaces(s[start:i])
			if part != "" {
				result = append(result, part)
			}
			start = i + 1
		}
	}
	return result
}

func trimSpaces(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
