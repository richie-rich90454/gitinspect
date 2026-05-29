# Changelog

## v0.1.0

Initial release!

- Local repository support with `.gitignore` awareness
- Remote repository support via shallow clone
- Token budget management with file priority sorting
- Include/exclude glob patterns (doublestar)
- Output formats: JSON, text, YAML
- Dependency extraction (go.mod, package.json, Cargo.toml, requirements.txt, Gemfile)
- HTTP server with `POST /inspect` endpoint
- MCP server for AI agent integration
- Disk cache with 1-hour TTL
- Comment stripping
- Multi-platform builds (linux/darwin/windows, amd64/arm64)
- Homebrew, Scoop, Go Install, deb/rpm/apk packages
