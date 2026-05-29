# Token Budget

gitinspect estimates tokens using a simple heuristic: `tokens = ceil(len(content) / 4)`. This gives a reasonable approximation for both code and natural language.

## How the Budget Works

1. Files are sorted by [priority](/guide/priority)
2. Each file's content is estimated for tokens
3. Files are added in priority order
4. When adding the next file would exceed `--max-tokens`, it is truncated
5. If a single file exceeds the remaining budget, it is included but truncated

## Truncation

When a file's content exceeds the remaining token budget, gitinspect truncates it:

- **First 1000 characters** are kept
- **Last 500 characters** are kept
- A `... [truncated] ...` marker is inserted between them
- `stats.truncated` is set to `true` in the output

## Configuring the Budget

```bash
# Default: 6000 tokens
gitinspect inspect .

# Larger budget for bigger repos
gitinspect inspect --max-tokens 20000 .

# Small budget for quick overviews
gitinspect inspect --max-tokens 2000 .
```

## Estimating Your Needs

| Use Case | Recommended `--max-tokens` |
|----------|---------------------------|
| Quick overview | 2000 |
| Standard analysis | 6000 (default) |
| Deep analysis | 20000 |
| Full repo dump | 100000+ |
