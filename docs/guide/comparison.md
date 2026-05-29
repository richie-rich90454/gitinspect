# Comparison

How gitinspect compares to similar tools.

## vs gitingest & repomix

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
| Comment stripping | ✅ | ❌ | ❌ |
| Install via Homebrew | ✅ | ❌ | ✅ |
| Install via Scoop | ✅ | ❌ | ❌ |
| One-line install script | ✅ | ❌ | ❌ |
| Linux packages (deb/rpm) | ✅ | ❌ | ❌ |

## Key Differentiators

### Token Budget Management

gitinspect is the only tool that lets you set a **token budget** and automatically prioritizes files to fit within it. This is critical for LLM context windows.

### MCP Server

gitinspect is the only tool with a built-in **MCP server**, making it directly invokable by AI coding agents like Claude, Cursor, and Windsurf.

### Dependency Extraction

gitinspect automatically extracts dependencies from `go.mod`, `package.json`, `Cargo.toml`, `requirements.txt`, and `Gemfile` — giving your LLM awareness of the project's dependency graph.

### File Priority

Files are sorted by importance: README → entry points → manifests → build files → others. The most relevant code is always included first.
