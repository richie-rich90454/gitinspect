# CLI Reference

## Commands

```
gitinspect inspect [flags] <repo-path-or-url>   Inspect a repository
gitinspect serve [flags]                          Start HTTP server
gitinspect mcp                                     Start MCP server (for AI agents)
```

## `gitinspect inspect`

Inspect a Git repository and output a structured snapshot.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `json` | Output format: `json`, `text`, `yaml` |
| `--max-tokens` | `6000` | Maximum token budget |
| `--max-files` | `0` | Maximum number of files (0 = unlimited) |
| `--include` | (none) | Include glob patterns (repeatable) |
| `--exclude` | (none) | Exclude glob patterns (repeatable) |
| `--strip` | `false` | Strip comments and blank lines |
| `--no-cache` | `false` | Disable cache |
| `--quiet` | `false` | Suppress progress output to stderr |

### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Runtime error (repo not found, clone failed, etc.) |
| `2` | Usage error (missing arguments, invalid flags) |

### Examples

```bash
# Basic inspection
gitinspect inspect .

# Remote repository
gitinspect inspect https://github.com/user/repo.git

# Text output with comment stripping
gitinspect inspect --format text --strip .

# Only Go files, exclude vendor
gitinspect inspect --include "**/*.go" --exclude "**/vendor/*" .

# Limit to 4000 tokens
gitinspect inspect --max-tokens 4000 /path/to/repo

# Quiet mode (for scripting)
gitinspect inspect --quiet --format json . > snapshot.json
```

## `gitinspect serve`

Start an HTTP server with a `POST /inspect` endpoint.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | HTTP server port |

### Request Body

```json
{
  "repo": "https://github.com/user/repo.git",
  "format": "json",
  "max_tokens": 6000,
  "max_files": 0,
  "include": [],
  "exclude": [],
  "strip": false,
  "no_cache": false
}
```

### Example

```bash
# Start server
gitinspect serve --port 8080

# Call the endpoint
curl -X POST -H "Content-Type: application/json" \
  -d '{"repo": "https://github.com/user/repo.git"}' \
  http://localhost:8080/inspect
```

## `gitinspect mcp`

Start an MCP server over stdio for AI agent integration. See [MCP Server](/guide/mcp) for details.

```bash
gitinspect mcp
```
