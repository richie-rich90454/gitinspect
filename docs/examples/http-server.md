# HTTP Server

Run gitinspect as an HTTP server for programmatic access.

## Starting the Server

```bash
gitinspect serve --port 8080
```

### Authentication

Use `--api-key` to require authentication on all requests:

```bash
gitinspect serve --port 8080 --api-key my-secret-key
```

When an API key is set, all requests must include an `Authorization: Bearer <key>` header:

```bash
curl -X POST http://localhost:8080/inspect \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer my-secret-key" \
  -d '{"repo": "https://github.com/user/repo.git"}'
```

Requests without a valid key receive a `401 Unauthorized` response.

### Concurrency

Use `--max-concurrent` to control how many remote fetches run in parallel (default: 4):

```bash
gitinspect serve --port 8080 --max-concurrent 8
```

Increase this value for high-throughput scenarios with many remote repositories.

## API Endpoint

### `POST /inspect`

**Request:**

```bash
curl -X POST http://localhost:8080/inspect \
  -H "Content-Type: application/json" \
  -d '{
    "repo": "https://github.com/user/repo.git",
    "format": "json",
    "max_tokens": 6000,
    "strip": true
  }'
```

**Request Body Fields:**

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `repo` | string | (required) | Repository path or URL |
| `format` | string | `json` | Output format |
| `max_tokens` | number | `6000` | Token budget |
| `max_files` | number | `0` | Max files |
| `include` | string[] | `[]` | Include patterns |
| `exclude` | string[] | `[]` | Exclude patterns |
| `strip` | boolean | `false` | Strip comments |
| `no_cache` | boolean | `false` | Disable cache |

**Response:**

- `200` with JSON/text/YAML body on success
- `400` for invalid request
- `500` for inspection errors

## Example: Inspect Local Repo

```bash
curl -X POST http://localhost:8080/inspect \
  -H "Content-Type: application/json" \
  -d '{"repo": "/path/to/local/repo"}'
```

## Example: With Filters

```bash
curl -X POST http://localhost:8080/inspect \
  -H "Content-Type: application/json" \
  -d '{
    "repo": "https://github.com/user/repo.git",
    "format": "text",
    "max_tokens": 4000,
    "include": ["**/*.go"],
    "exclude": ["**/vendor/*"],
    "strip": true
  }'
```

## Integration with LLMs

The HTTP server is useful for building custom LLM pipelines:

```python
import requests

response = requests.post("http://localhost:8080/inspect", json={
    "repo": "https://github.com/user/repo.git",
    "format": "json",
    "max_tokens": 4000,
    "strip": True,
})

snapshot = response.json()
# Feed snapshot["tree"] into your LLM
```
