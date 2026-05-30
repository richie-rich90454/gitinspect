# Feature Overview

A detailed look at what gitinspect can do and how each feature works.

## Core Features

### Local Repository Support

Walk the worktree, read file contents, and respect `.gitignore`. Binary files and non-UTF-8 content are automatically skipped.

```bash
gitinspect inspect /path/to/local/repo
```

### Remote Repository Support

Shallow clone (`--depth 1`) into a temp directory, analyze, and clean up. Snapshots are cached by commit hash so repeated inspections are instant.

```bash
gitinspect inspect https://github.com/user/repo.git
```

### Token Budget Management

Estimate tokens as `len(content)/4`. Sort files by priority. Stop when the budget is exceeded. If a single file exceeds the budget, include it but truncate the content.

```bash
gitinspect inspect --max-tokens 4000 .
```

### File Priority Sorting

Files are sorted by importance before inclusion:

| Priority | Files |
|----------|-------|
| 4 (highest) | `README*` |
| 3 | `main.*`, `cmd/*/main.*` |
| 2 | `go.mod`, `package.json`, `Cargo.toml`, `setup.py`, `Gemfile`, `requirements.txt` |
| 1 | `Makefile`, `Dockerfile`, `docker-compose.yml` |
| 0 | All other files |

### Include/Exclude Globs

Filter files using glob patterns. Supports `**` for recursive matching.

```bash
gitinspect inspect --include "**/*.go" --exclude "**/vendor/*" .
```

### Output Formats

Three output formats for different use cases:

- **JSON** — Structured data for programmatic consumption
- **Text** — Human-readable, pipe-friendly output
- **YAML** — Configuration-friendly format

```bash
gitinspect inspect --format json .
gitinspect inspect --format text .
gitinspect inspect --format yaml .
```

### Dependency Extraction

Automatically parse dependency manifests:

| File | Format | Example Output |
|------|--------|----------------|
| `go.mod` | Go modules | `github.com/go-git/go-git/v5@v5.12.0` |
| `package.json` | Node.js | `react@^18.0.0` |
| `Cargo.toml` | Rust | `serde@1.0` |
| `requirements.txt` | Python | `flask==2.0` |
| `Gemfile` | Ruby | `rails` |

### Comment Stripping

Remove single-line comments and blank lines based on file extension:

- `//` for Go, C, C++, Java, JavaScript, TypeScript, Rust, Kotlin, Swift
- `#` for Python, Shell, YAML, TOML, Ruby, R, PowerShell
- `--` for SQL, Lua, Haskell
- `;` for Assembly

```bash
gitinspect inspect --strip .
```

### Binary File Detection

Automatically skips binary files (images, executables, archives, fonts, databases, etc.) and non-UTF-8 content. No garbled output in your snapshots.

### Caching

Remote repo snapshots are cached in `~/.cache/gitinspect/` keyed by `sha256(repo_url + commit_hash)` with a 1-hour TTL. Disable with `--no-cache`.

```bash
gitinspect inspect https://github.com/user/repo.git          # First run: clones
gitinspect inspect https://github.com/user/repo.git          # Second run: cached
gitinspect inspect --no-cache https://github.com/user/repo.git  # Force fresh
```

### HTTP Server

Run as an HTTP server with a `POST /inspect` endpoint for programmatic access.

```bash
gitinspect serve --port 8080
```

### MCP Server (AI Agents)

Built-in MCP server over stdio for AI agent integration. Two tools:

- **`inspect_repo`** — Full inspection with file contents, dependencies, and stats
- **`list_repo_files`** — Lightweight file listing with priority scores and token estimates

```bash
gitinspect mcp
```

See [MCP Integration](/guide/mcp) for agent configuration.
