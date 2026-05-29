# Agent Configuration

Configure gitinspect as an MCP server in your AI coding tool.

## Claude Desktop / Claude Code

Add to `claude_desktop_config.json` or `.claude/settings.json`:

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

## Cursor

Add to `.cursor/mcp.json` in your project root:

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

## Windsurf

Add to `.windsurf/mcp.json`:

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

## VS Code (GitHub Copilot)

Add to `.vscode/mcp.json`:

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

## Custom Binary Path

If gitinspect is not in your PATH, use the full path:

```json
{
  "mcpServers": {
    "gitinspect": {
      "command": "/usr/local/bin/gitinspect",
      "args": ["mcp"]
    }
  }
}
```

On Windows:

```json
{
  "mcpServers": {
    "gitinspect": {
      "command": "C:\\Users\\you\\scoop\\shims\\gitinspect.exe",
      "args": ["mcp"]
    }
  }
}
```

## Environment Variables

You can pass environment variables to the MCP server:

```json
{
  "mcpServers": {
    "gitinspect": {
      "command": "gitinspect",
      "args": ["mcp"],
      "env": {
        "HOME": "/Users/you"
      }
    }
  }
}
```
