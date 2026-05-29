# MCP Server

gitinspect includes a built-in **Model Context Protocol (MCP)** server, making it directly invokable by AI coding agents.

## What is MCP?

The [Model Context Protocol](https://modelcontextprotocol.io/) is an open standard that enables AI applications to securely connect to external data sources and tools. It provides a standardized way for LLMs to access and interact with external systems.

## Starting the MCP Server

```bash
gitinspect mcp
```

This starts an MCP server over **stdio** — the standard transport for AI agent integration. The server listens for JSON-RPC messages on stdin and writes responses to stdout.

## Available Tools

### `inspect_repo`

Inspect a Git repository and return a structured snapshot with file contents, dependencies, and stats.

**Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repo` | string | ✅ | — | Local path or remote URL |
| `format` | string | ❌ | `json` | Output format: json, text, yaml |
| `max_tokens` | number | ❌ | `6000` | Maximum token budget |
| `max_files` | number | ❌ | `0` | Max files (0 = unlimited) |
| `include` | string | ❌ | — | Include glob patterns (comma-separated) |
| `exclude` | string | ❌ | — | Exclude glob patterns (comma-separated) |
| `strip` | boolean | ❌ | `false` | Strip comments and blank lines |
| `no_cache` | boolean | ❌ | `false` | Disable caching |

### `list_repo_files`

List all files in a repository with priority scores and token estimates — useful for deciding which files to inspect before reading contents.

**Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `repo` | string | ✅ | — | Local path or remote URL |
| `include` | string | ❌ | — | Include glob patterns (comma-separated) |
| `exclude` | string | ❌ | — | Exclude glob patterns (comma-separated) |

## Two-Step Inspection Pattern

The recommended pattern for AI agents is a two-step approach:

1. **List files** with `list_repo_files` to understand the repository structure
2. **Inspect** with `inspect_repo` using targeted `include`/`exclude` patterns and appropriate `max_tokens`

This avoids wasting the token budget on files the agent doesn't need.
