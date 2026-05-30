# Output Formats

gitinspect supports three output formats, controlled by the `--format` flag.

## JSON (default)

```bash
gitinspect inspect --format json .
```

```json
{
  "tree": {
    "README.md": "# my-project\n...",
    "main.go": "package main\n...",
    "go.mod": "module github.com/user/project\n..."
  },
  "stats": {
    "file_count": 3,
    "total_bytes": 456,
    "total_tokens": 114,
    "truncated": false
  },
  "dependencies": [
    "github.com/go-git/go-git/v5@v5.12.0"
  ],
  "version": "0.1.0"
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `tree` | `map[string]string` | File paths to file contents |
| `stats.file_count` | `int` | Number of files included |
| `stats.total_bytes` | `int` | Total bytes of all file contents |
| `stats.total_tokens` | `int` | Estimated total tokens across all files |
| `stats.truncated` | `bool` | Whether any content was truncated |
| `dependencies` | `[]string` | Extracted dependency strings |
| `version` | `string` | gitinspect version |

## Text

```bash
gitinspect inspect --format text .
```

```
File: README.md
---
# my-project
...

File: main.go
---
package main
...
```

Useful for piping into LLM prompts directly.

## YAML

```bash
gitinspect inspect --format yaml .
```

```yaml
tree:
  README.md: |
    # my-project
  main.go: |
    package main
stats:
  file_count: 2
  total_bytes: 123
  total_tokens: 31
  truncated: false
dependencies:
  - github.com/go-git/go-git/v5@v5.12.0
version: 0.1.0
```
