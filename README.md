# envlens

> Tool to diff and audit environment variable changes across deployments or .env files.

## Installation

```bash
go install github.com/yourname/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envlens.git && cd envlens && go build ./...
```

## Usage

Compare two `.env` files to see what changed:

```bash
envlens diff .env.staging .env.production
```

Audit a deployment by comparing a local file against exported environment variables:

```bash
envlens diff .env.local <(printenv)
```

**Example output:**

```
+ NEW_FEATURE_FLAG=true
- DEPRECATED_API_KEY=abc123
~ DATABASE_URL  [changed]
  API_TIMEOUT   [unchanged]
```

### Flags

| Flag | Description |
|------|-------------|
| `--redact` | Mask secret values in output |
| `--format json` | Output results as JSON |
| `--only-changed` | Show only differing variables |

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)