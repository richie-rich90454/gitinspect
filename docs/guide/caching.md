# Caching

gitinspect caches snapshots of remote repositories to avoid repeated clones.

## How It Works

- Cache directory: `~/.cache/gitinspect/`
- Cache key: `sha256(repo_url + "#" + commit_hash)`
- TTL: **1 hour** — entries older than 1 hour are automatically invalidated

## Cache Flow

1. When inspecting a remote repo, gitinspect resolves the HEAD commit hash
2. It checks the cache for an entry keyed by `sha256(url + "#" + hash)`
3. If a valid cache entry exists (modtime < 1 hour), it's returned immediately
4. If not, the repo is cloned and the result is cached

## Disabling the Cache

```bash
# Disable for this run
gitinspect inspect --no-cache https://github.com/user/repo.git
```

## Clearing the Cache

```bash
rm -rf ~/.cache/gitinspect/
```
