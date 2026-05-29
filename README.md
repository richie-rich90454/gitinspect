# gitinspect

A CLI tool (and optional HTTP server) that turns any Git repository (local path or remote URL) into a structured, token-efficient snapshot optimized for LLMs and AI agents. Works with any Git host — GitHub, GitLab, Bitbucket, self-hosted — without requiring a full clone.

## Badges

![Go Version](https://img.shields.io/github/go-mod/go-version/your-username/gitinspect)
![License](https://img.shields.io/badge/license-Apache--2.0-blue)
![Build Status](https://github.com/your-username/gitinspect/actions/workflows/ci.yml/badge.svg)

## Quick Install

### Go Install
```bash
go install github.com/your-username/gitinspect/cmd/gitinspect@latest
```

### Homebrew
```bash
brew install your-username/tap/gitinspect
```

### Scoop
```powershell
scoop install gitinspect
```

### Binary Download
Download the latest binary from [Releases](https://github.com/your-username/gitinspect/releases).

## Usage

### Local Repository
```bash
gitinspect /path/to/repo
```

### Remote Repository
```bash
gitinspect https://github.com/user/repo.git
```

### With Options
```bash
gitinspect --format text --max-tokens 10000 --strip /path/to/repo
gitinspect --include "**/*.go" --exclude "**/vendor/*" .
```

### HTTP Server Mode
```bash
gitinspect --server --port 8080
```
Then send a POST request to `/inspect`:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"repo": "https://github.com/user/repo.git", "format": "json"}' \
  http://localhost:8080/inspect
```

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `json` | Output format: `json`, `text`, `yaml` |
| `--max-tokens` | `6000` | Maximum token budget |
| `--max-files` | `0` | Maximum number of files (0 = unlimited) |
| `--include` | (none) | Include glob patterns (repeatable) |
| `--exclude` | (none) | Exclude glob patterns (repeatable) |
| `--strip` | `false` | Strip comments and blank lines |
| `--no-cache` | `false` | Disable cache |
| `--server` | `false` | Start HTTP server |
| `--port` | `8080` | HTTP server port |

## Output Format

### JSON (default)
```json
{
  "tree": { "path/to/file.go": "file content..." },
  "stats": { "file_count": 5, "total_bytes": 1234, "truncated": false },
  "dependencies": ["github.com/foo/bar@v1.0.0"],
  "version": "0.1.0"
}
```

### Text
```
File: path/to/file.go
---
file content...
```

## Comparison

| Feature | gitinspect | gitingest | repomix |
|---------|:----------:|:---------:|:-------:|
| Local repo support | ✅ | ✅ | ✅ |
| Remote repo support | ✅ | ✅ | ✅ |
| Token budget management | ✅ | ✅ | ❌ |
| Dependency extraction | ✅ | ❌ | ❌ |
| HTTP server | ✅ | ❌ | ❌ |
| Output formats | json/text/yaml | text | json |
| Caching | ✅ | ❌ | ❌ |
| .gitignore support | ✅ | ✅ | ✅ |
| File priority sorting | ✅ | ❌ | ❌ |

## Demo

TODO: add demo.gif

## License

Apache License 2.0 — See [LICENSE](LICENSE) for details.
