# AI Agent Usage

Use gitinspect through MCP-compatible AI coding agents.

## Prerequisites

1. Install gitinspect (see [Installation](/guide/installation))
2. Configure your AI agent (see [Agent Config](/guide/agent-config))

## Example: Understanding a New Codebase

When you open a project in your AI agent, you can ask:

> "Inspect this repository and give me an overview of the codebase structure"

The agent will call `list_repo_files` to see what's in the repo, then `inspect_repo` to read the key files.

## Example: Focused Code Review

> "Inspect only the Go files in this project, excluding tests and vendor"

The agent will call:

```
inspect_repo(
  repo=".",
  include="**/*.go",
  exclude="**/vendor/*,**/*_test.go",
  max_tokens=8000,
  strip=true
)
```

## Example: Remote Repository Analysis

> "Inspect the repository at https://github.com/go-git/go-git and summarize its structure"

The agent will call:

```
list_repo_files(repo="https://github.com/go-git/go-git")
```

Then based on the file list:

```
inspect_repo(
  repo="https://github.com/go-git/go-git",
  max_tokens=6000,
  strip=true
)
```

## Example: Dependency Audit

> "What dependencies does this project use?"

The agent will call `inspect_repo` and examine the `dependencies` field in the output.

## Two-Step Pattern

The recommended pattern for AI agents:

1. **`list_repo_files`** — Get an overview with file sizes and priority scores
2. **`inspect_repo`** — Read specific files with a targeted token budget

This avoids wasting tokens on files the agent doesn't need.
