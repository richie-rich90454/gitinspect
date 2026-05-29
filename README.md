# gitinspect

A CLI tool (and optional HTTP server) that turns any Git repository (local path or remote URL) into a structured, token-efficient snapshot optimized for LLMs and AI agents. Works with any Git host — GitHub, GitLab, Bitbucket, self-hosted — without requiring a full clone.

## Badges

![Go Version](https://img.shields.io/github/go-mod/go-version/your-username/gitinspect)
![License](https://img.shields.io/badge/license-Apache--2.0-blue)
![Build Status](https://github.com/your-username/gitinspect/actions/workflows/ci.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/your-username/gitinspect)

## Installation

### Quick Install (macOS / Linux)

The fastest way to get started — one command:

```bash
curl -fsSL https://raw.githubusercontent.com/your-username/gitinspect/main/install.sh | bash
```

Or with `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/your-username/gitinspect/main/install.sh | bash
```

> 💡 Set `INSTALL_DIR` to change the install location (default: `/usr/local/bin`):
> ```bash
> curl -fsSL ... | INSTALL_DIR=~/.local/bin bash
> ```

### Homebrew (macOS / Linux)

```bash
brew tap your-username/tap
brew install gitinspect
```

### Scoop (Windows)

```bash
scoop bucket add gitinspect https://github.com/your-username/scoop-bucket
scoop install gitinspect
```

### Go Install

Requires Go 1.22+:

```bash
go install github.com/your-username/gitinspect/cmd/gitinspect@latest
```

### Docker

```bash
docker run --rm -v /path/to/repo:/repo your-username/gitinspect /repo
```

### Binary Download

Download the latest binary for your platform from the [Releases page](https://github.com/your-username/gitinspect/releases):

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

Packages are published with each release:

```bash
# Debian / Ubuntu
sudo dpkg -i gitinspect_*_linux_amd64.deb

# RHEL / Fedora
sudo rpm -i gitinspect_*_linux_amd64.rpm

# Alpine
sudo apk add gitinspect_*_linux_amd64.apk
```

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
