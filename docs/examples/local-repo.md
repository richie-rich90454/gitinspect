# Local Repository

Inspect a repository on your local filesystem.

## Basic Usage

```bash
gitinspect inspect /path/to/repo
```

## Current Directory

```bash
gitinspect inspect .
```

## With Comment Stripping

Remove comments and blank lines to save tokens:

```bash
gitinspect inspect --strip .
```

## Filter by File Type

Only include Go files:

```bash
gitinspect inspect --include "**/*.go" .
```

## Exclude Directories

Skip vendor and test directories:

```bash
gitinspect inspect --exclude "**/vendor/*" --exclude "**/*_test.go" .
```

## Limit Token Budget

Get a quick overview with a small budget:

```bash
gitinspect inspect --max-tokens 2000 --format text .
```

## Combine Options

Full example combining multiple options:

```bash
gitinspect inspect \
  --format text \
  --max-tokens 8000 \
  --max-files 20 \
  --include "**/*.go" \
  --exclude "**/vendor/*" \
  --exclude "**/*_test.go" \
  --strip \
  .
```

## .gitignore Support

gitinspect automatically respects `.gitignore` files when reading local repositories. Files matching patterns in `.gitignore` are excluded from the snapshot.

## Output to File

```bash
gitinspect inspect --quiet --format json . > snapshot.json
```
