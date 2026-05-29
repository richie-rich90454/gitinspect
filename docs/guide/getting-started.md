# Getting Started

## What is gitinspect?

gitinspect is a CLI tool that turns any Git repository — local or remote — into a structured, token-efficient snapshot optimized for LLMs and AI agents. It works with any Git host without requiring a full clone.

## Quick Install

::: code-group

```bash [macOS / Linux (curl)]
curl -fsSL https://raw.githubusercontent.com/your-username/gitinspect/main/install.sh | bash
```

```bash [macOS / Linux (wget)]
wget -qO- https://raw.githubusercontent.com/your-username/gitinspect/main/install.sh | bash
```

```bash [Homebrew]
brew tap your-username/tap
brew install gitinspect
```

```powershell [Windows (Scoop)]
scoop bucket add gitinspect https://github.com/your-username/scoop-bucket
scoop install gitinspect
```

```bash [Go Install]
go install github.com/your-username/gitinspect/cmd/gitinspect@latest
```

:::

## Your First Inspection

```bash
# Inspect the current directory
gitinspect inspect .

# Inspect a remote repository
gitinspect inspect https://github.com/user/repo.git

# Get plain text output
gitinspect inspect --format text --strip .

# Start as an MCP server for AI agents
gitinspect mcp
```

## Next Steps

- [Installation](/guide/installation) — detailed install options
- [CLI Reference](/guide/cli) — all flags and commands
- [MCP Integration](/guide/mcp) — use with AI coding agents
- [Examples](/examples/local-repo) — real-world usage patterns
