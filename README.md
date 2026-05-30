# gitinspect

A CLI tool (and optional HTTP/MCP server) that turns any Git repository — local path or remote URL — into a structured, token-efficient snapshot optimized for LLMs and AI agents. Works with any Git host without requiring a full clone.

## Badges

![Go Version](https://img.shields.io/github/go-mod/go-version/richie-rich90454/gitinspect)
![License](https://img.shields.io/badge/license-Apache--2.0-blue)
![Build Status](https://github.com/richie-rich90454/gitinspect/actions/workflows/ci.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/richie-rich90454/gitinspect)

## Features

### Token Budget Management

Estimate tokens as `len(content)/4`. Sort files by priority. Stop when the budget is exceeded. Never waste context window space on irrelevant files.

### AI Agent Ready (MCP)

Built-in MCP server over stdio. Works with Claude, Cursor, Windsurf, VS Code Copilot, and any MCP-compatible agent. Two tools: `inspect_repo` for full snapshots, `list_repo_files` for lightweight file listings.

### Any Git Host

GitHub, GitLab, Bitbucket, self-hosted — works with any Git repository. Remote repos use shallow clone (`--depth 1`) and are cached by commit hash.

### Smart File Priority

Files are sorted by importance: README first, then entry points (`main.*`), then manifests (`go.mod`, `package.json`, etc.), then build files (`Makefile`, `Dockerfile`), then everything else.

### Dependency Extraction

Automatically parse `go.mod`, `package.json`, `Cargo.toml`, `requirements.txt`, and `Gemfile` for dependency awareness.

### Binary File Detection

Automatically skips binary files (images, executables, archives, fonts, etc.) and non-UTF-8 content. No garbled output.

### Comment Stripping

Remove single-line comments and blank lines based on file extension. Supports `//` (Go, C, JS, Java, Rust, etc.), `#` (Python, Shell, YAML, Ruby, etc.), `--` (SQL, Lua, Haskell), and `;` (Assembly).

### Multiple Output Formats

JSON (default), plain text, or YAML. Perfect for piping into LLMs, saving to files, or programmatic consumption.

## Installation

### Quick Install (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/richie-rich90454/gitinspect/main/install.sh | bash
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/richie-rich90454/gitinspect/main/install.sh | bash
```

> 💡 Set `INSTALL_DIR` to change the install location (default: `/usr/local/bin`):
> ```bash
> curl -fsSL ... | INSTALL_DIR=~/.local/bin bash
> ```

### Homebrew (macOS / Linux)

```bash
brew tap richie-rich90454/tap
brew install gitinspect
```

### Scoop (Windows)

```bash
scoop bucket add gitinspect https://github.com/richie-rich90454/scoop-bucket
scoop install gitinspect
```

### Go Install

```bash
go install github.com/richie-rich90454/gitinspect/cmd/gitinspect@latest
```

### Binary Download

Download the latest binary for your platform from the [Releases page](https://github.com/richie-rich90454/gitinspect/releases).

### Debian / RPM / APK

```bash
sudo dpkg -i gitinspect_*_linux_amd64.deb
sudo rpm -i gitinspect_*_linux_amd64.rpm
sudo apk add gitinspect_*_linux_amd64.apk
```

## Usage

### Commands

```
gitinspect inspect [flags] <repo-path-or-url>   Inspect a repository
gitinspect serve [flags]                          Start HTTP server
gitinspect mcp                                     Start MCP server (for AI agents)
gitinspect version                                 Print version
```

### Local Repository

```bash
gitinspect inspect /path/to/repo
```

### Remote Repository

```bash
gitinspect inspect https://github.com/user/repo.git
```

### With Options

```bash
gitinspect inspect --format text --max-tokens 10000 --strip /path/to/repo
gitinspect inspect --include "**/*.go" --exclude "**/vendor/*" .
```

### HTTP Server Mode

```bash
gitinspect serve --port 8080
```

Then send a POST request:

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"repo": "https://github.com/user/repo.git", "format": "json"}' \
  http://localhost:8080/inspect
```

## AI Agent Integration (MCP)

gitinspect includes a built-in **Model Context Protocol (MCP)** server, making it directly invokable by AI coding agents.

### Starting the MCP Server

```bash
gitinspect mcp
```

This starts an MCP server over stdio — the standard transport for AI agent integration.

### Configuring in AI Agents

Add to your agent's MCP config file:

```json
{
  "mcpServers": {
    "gitinspect": {
      "command": "gitinspect",
      "args": ["mcp"]
    }
  }
}
```

Works with Claude Desktop, Cursor, Windsurf, VS Code Copilot, and any MCP-compatible agent.

### MCP Tools

**`inspect_repo`** — Inspect a Git repo and return a structured snapshot with file contents, dependencies, and stats. Supports `format`, `max_tokens`, `max_files`, `include`, `exclude`, `strip`, `no_cache` parameters.

**`list_repo_files`** — List all files with priority scores and token estimates, without reading file contents. Useful for deciding which files to inspect before reading them.

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
| `--quiet` | `false` | Suppress progress output (for scripting) |
| `--port` | `8080` | HTTP server port (serve command) |

## Output Examples

### JSON (default)

```bash
gitinspect inspect --format json .
```

```json
{
  "tree": {
    "README.md": "# my-project\nA cool project.\n",
    "main.go": "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
    "go.mod": "module github.com/user/project\n\ngo 1.25\n"
  },
  "stats": {
    "file_count": 3,
    "total_bytes": 120,
    "total_tokens": 30,
    "truncated": false
  },
  "dependencies": [
    "github.com/go-git/go-git/v5@v5.12.0"
  ],
  "version": "0.1.0"
}
```

### Text

```bash
gitinspect inspect --format text .
```

```
File: README.md
---
# my-project
A cool project.

File: go.mod
---
module github.com/user/project

go 1.25

File: main.go
---
package main

func main() {
        println("hello")
}
```

### YAML

```bash
gitinspect inspect --format yaml .
```

```yaml
tree:
  README.md: |
    # my-project
  go.mod: |
    module github.com/user/project
  main.go: |
    package main
stats:
  file_count: 3
  total_bytes: 120
  total_tokens: 30
  truncated: false
dependencies:
  - github.com/go-git/go-git/v5@v5.12.0
version: 0.1.0
```

## File Priority

Files are sorted by priority before inclusion. Higher-priority files are included first within the token budget:

| Priority | Files |
|----------|-------|
| 4 (highest) | `README*` |
| 3 | `main.*`, `cmd/*/main.*` |
| 2 | `go.mod`, `package.json`, `Cargo.toml`, `setup.py`, `Gemfile`, `requirements.txt` |
| 1 | `Makefile`, `Dockerfile`, `docker-compose.yml` |
| 0 | All other files |

## Token Budget

Tokens are estimated as `ceil(len(content) / 4)`. When the budget is exceeded:

1. Files are added in priority order
2. If a file would exceed the remaining budget, it is truncated (first 1000 + last 500 chars with a truncation marker)
3. If a single file alone exceeds the entire budget, it is included but truncated
4. `stats.truncated` is set to `true` in the output

## Dependency Extraction

Supported manifest files:

| File | Format | Example Output |
|------|--------|----------------|
| `go.mod` | Go modules | `github.com/go-git/go-git/v5@v5.12.0` |
| `package.json` | Node.js | `react@^18.0.0` |
| `Cargo.toml` | Rust | `serde@1.0` |
| `requirements.txt` | Python | `flask==2.0` |
| `Gemfile` | Ruby | `rails` |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

Apache License 2.0 — See [LICENSE](LICENSE) for details.
