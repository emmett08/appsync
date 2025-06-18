## Usage

### Incremental rollout

| Stage                  | Command example                                             |
|------------------------|-------------------------------------------------------------|
| OTT Reference App only | `appsync sync --team core-platform --app ott-reference-app` |
| All DPE repos          | `appsync sync --team core-platform`                         |
| Full catalogue         | `appsync sync` (omit `--team` and `--app`)                  |

Invoke the same binary with different filters in your CI pipeline to promote gradually.

---

## Building and running

```bash
# Build the CLI
go build -o appsync ./cmd/appsync

# Acquire a GitHub token via the GitHub CLI
export GITHUB_TOKEN="$(gh auth token)"

# Run sync against the OTT Reference App repository
./appsync sync \
  --root /path/to/dpe/tenants/.applications \
  --owner NBCUDTC \
  --repo ott-reference-app \
  --team core-platform \
  --app ott-reference-app \
  --mode pr \
  --token "$GITHUB_TOKEN"
```

---

## Testing

```bash
# Run BDD feature tests
go test ./features -v
```
