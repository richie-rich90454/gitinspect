# HTTP Server

Run gitinspect as an HTTP server for programmatic access.

## Starting the Server

```bash
gitinspect serve --port 8080
```

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
