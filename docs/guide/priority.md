# File Priority

gitinspect sorts files by priority before including them in the snapshot. This ensures the most important files are always included within the token budget.

## Priority Levels

| Priority | Score | Files |
|----------|-------|-------|
| **Highest** | 4 | `README*` — any file starting with "README" |
| **Very High** | 3 | `main.*`, `cmd/*/main.*` — entry point files |
| **High** | 2 | `go.mod`, `package.json`, `Cargo.toml`, `setup.py`, `Gemfile`, `requirements.txt` |
| **Medium** | 1 | `Makefile`, `Dockerfile`, `docker-compose.yml` |
| **Low** | 0 | All other files |

## How It Works

1. All files are collected from the repository
2. Each file is assigned a priority score (0–4)
3. Files are sorted highest-first
4. Files are added to the snapshot in priority order until the token budget is exhausted

## Example

Given a repository with these files:

```
src/utils.go       (priority 0)
Makefile           (priority 1)
go.mod             (priority 2)
main.go            (priority 3)
README.md          (priority 4)
```

With `--max-tokens 2000`, gitinspect will include them in this order:
1. `README.md` (priority 4)
2. `main.go` (priority 3)
3. `go.mod` (priority 2)
4. `Makefile` (priority 1)
5. `src/utils.go` (priority 0) — if tokens remain
