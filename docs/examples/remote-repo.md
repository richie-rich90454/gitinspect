# Remote Repository

Inspect a remote Git repository without cloning the entire history.

## Basic Usage

```bash
gitinspect inspect https://github.com/user/repo.git
```

## How It Works

1. gitinspect resolves the HEAD commit via `git ls-remote`
2. It checks the cache for an existing snapshot
3. If not cached, it performs a shallow clone (`--depth 1 --filter=blob:none`)
4. The repository is analyzed and the result is cached
5. The temporary clone is cleaned up

## Any Git Host

Works with any Git remote — GitHub, GitLab, Bitbucket, self-hosted:

```bash
# GitHub
gitinspect inspect https://github.com/user/repo.git

# GitLab
gitinspect inspect https://gitlab.com/user/repo.git

# Bitbucket
gitinspect inspect https://bitbucket.org/user/repo.git

# Self-hosted
gitinspect inspect https://git.company.com/team/project.git
```

## With Options

```bash
gitinspect inspect \
  --format text \
  --max-tokens 4000 \
  --strip \
  https://github.com/user/repo.git
```

## Disable Cache

Force a fresh clone:

```bash
gitinspect inspect --no-cache https://github.com/user/repo.git
```

## Performance Tips

- **Use caching** (default) — repeated inspections of the same repo are instant
- **Set a token budget** — `--max-tokens 4000` is usually enough for an overview
- **Use `--strip`** — removing comments saves significant tokens
- **Use `--include`** — only inspect the files you need
