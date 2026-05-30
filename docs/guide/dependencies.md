# Dependency Extraction

gitinspect automatically extracts dependencies from common manifest files.

## Supported Files

| File | Parser | Output Format |
|------|--------|---------------|
| `go.mod` | Line parser | `github.com/foo/bar@v1.2.3` |
| `package.json` | JSON parser | `react@^18.0.0` |
| `Cargo.toml` | TOML parser | `serde@1.0` |
| `requirements.txt` | Line parser | `flask==2.0` |
| `Gemfile` | Regex parser | `rails` |

## How It Works

Dependencies are extracted from each manifest file as it's processed. They appear in the `dependencies` array of the output:

```json
{
  "dependencies": [
    "github.com/go-git/go-git/v5@v5.12.0",
    "github.com/bmatcuk/doublestar/v4@v4.6.1",
    "gopkg.in/yaml.v3@v3.0.1"
  ]
}
```

## Example

```bash
gitinspect inspect --format json .
```

For a Go project with this `go.mod`:

```
module example.com/app

go 1.25

require (
    github.com/go-git/go-git/v5 v5.12.0
    gopkg.in/yaml.v3 v3.0.1
)
```

The output will include:

```json
{
  "dependencies": [
    "github.com/go-git/go-git/v5@v5.12.0",
    "gopkg.in/yaml.v3@v3.0.1"
  ]
}
```
