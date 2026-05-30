# gitinspect

A CLI tool (and optional HTTP/MCP server) that turns any Git repository (local path or remote URL) into a structured, token-efficient snapshot optimized for LLMs and AI agents. Works with any Git host — GitHub, GitLab, Bitbucket, self-hosted — without requiring a full clone.

## Badges

![Go Version](https://img.shields.io/github/go-mod/go-version/richie-rich90454/gitinspect)
![License](https://img.shields.io/badge/license-Apache--2.0-blue)
![Build Status](https://github.com/richie-rich90454/gitinspect/actions/workflows/ci.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/richie-rich90454/gitinspect)

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

Download the latest binary for your platform from the [Releases page](https://github.com/richie-rich90454/gitinspect/releases):

| Platform | Architecture | File |
|----------|-------------|------|
| Linux | amd64 | `gitinspect_X.Y.Z_linux_amd64.tar.gz` |
| Linux | arm64 | `gitinspect_X.Y.Z_linux_arm64.tar.gz` |
| macOS | amd64 | `gitinspect_X.Y.Z_darwin_amd64.tar.gz` |
| macOS | arm64 | `gitinspect_X.Y.Z_darwin_arm64.tar.gz` |
| Windows | amd64 | `gitinspect_X.Y.Z_windows_amd64.zip` |
| Windows | arm64 | `gitinspect_X.Y.Z_windows_arm64.zip` |

Extract and place the binary in your `PATH`:

```bash
# Linux / macOS
tar xzf gitinspect_*_linux_amd64.tar.gz
chmod +x gitinspect
sudo mv gitinspect /usr/local/bin/

# Windows — extract the zip and add gitinspect.exe to your PATH
```

### Debian / RPM / APK

```bash
# Debian / Ubuntu
sudo dpkg -i gitinspect_*_linux_amd64.deb

# RHEL / Fedora
sudo rpm -i gitinspect_*_linux_amd64.rpm

# Alpine
sudo apk add gitinspect_*_linux_amd64.apk
```

## Usage

### Commands

```
gitinspect inspect [flags] <repo-path-or-url>   Inspect a repository
gitinspect serve [flags]                          Start HTTP server
gitinspect mcp                                     Start MCP server (for AI agents)
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
Then send a POST request to `/inspect`:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"repo": "https://github.com/user/repo.git", "format": "json"}' \
  http://localhost:8080/inspect
```

## AI Agent Integration (MCP)

gitinspect includes a built-in **Model Context Protocol (MCP)** server, making it directly invokable by AI coding agents like Claude, Cursor, Windsurf, and others.

### Starting the MCP Server

```bash
gitinspect mcp
```

This starts an MCP server over stdio — the standard transport for AI agent integration.

### Configuring in AI Agents

#### Claude Desktop / Claude Code

Add to your `claude_desktop_config.json` or `.claude/settings.json`:

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

#### Cursor

Add to your `.cursor/mcp.json`:

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

#### Windsurf

Add to your `.windsurf/mcp.json`:

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

#### VS Code (GitHub Copilot)

Add to your `.vscode/mcp.json`:

```json
{
  "servers": {
    "gitinspect": {
      "command": "gitinspect",
      "args": ["mcp"]
    }
  }
}
```

### MCP Tools Available

| Tool | Description |
|------|-------------|
| `inspect_repo` | Inspect a Git repo and return a structured snapshot with file contents, dependencies, and stats. Supports `format`, `max_tokens`, `max_files`, `include`, `exclude`, `strip`, `no_cache` parameters. |
| `list_repo_files` | List all files in a repo with priority scores and token estimates — useful for deciding which files to inspect before reading contents. |

### Example Agent Interactions

An AI agent can use gitinspect like this:

```
Agent: I'll inspect the repository structure first.
→ Calls: list_repo_files(repo="/path/to/project")
← Gets: File list with priority scores and token estimates

Agent: Now let me read the key files.
→ Calls: inspect_repo(repo="/path/to/project", max_tokens=4000, strip=true)
← Gets: Structured JSON with file contents, dependencies, and stats
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
| `--quiet` | `false` | Suppress progress output (for scripting) |
| `--port` | `8080` | HTTP server port (serve command) |

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
| MCP server (AI agents) | ✅ | ❌ | ❌ |
| Output formats | json/text/yaml | text | json |
| Caching | ✅ | ❌ | ❌ |
| .gitignore support | ✅ | ✅ | ✅ |
| File priority sorting | ✅ | ❌ | ❌ |
| Install via Homebrew | ✅ | ❌ | ✅ |
| Install via Scoop | ✅ | ❌ | ❌ |
| One-line install script | ✅ | ❌ | ❌ |

## Demo

TODO: add demo.gif

## Contributing

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes (`git commit -am 'Add my feature'`)
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a Pull Request

## License

Apache License 2.0 — See [LICENSE](LICENSE) for details.
