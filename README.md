# taco-vulndb

A CLI tool that builds, updates, and distributes the TACO vulnerability database. It aggregates CVE and security advisory data from multiple upstream sources, merges them with precedence rules, and makes the resulting database available through OCI registries, an HTTP server, or file export.

## Features

- **Multi-source aggregation** -- fetches vulnerability data from nine sources with configurable source selection
- **Incremental and full updates** -- first run performs a full historical fetch; subsequent runs fetch only the last 7 days unless `--full` is specified
- **OCI registry distribution** -- push and pull the database as OCI artifacts to any container registry (e.g., GHCR)
- **Built-in HTTP server** -- serve the database over HTTP with JSON, gzip, metadata, and health-check endpoints
- **File export** -- export the database as a gzip-compressed file for offline or custom distribution
- **Download and import** -- download a pre-built database from a URL or load one from a local file
- **Status reporting** -- inspect the local cache including entry counts, staleness, and per-source breakdowns
- **Automated daily builds** -- GitHub Actions workflow updates the database on a daily cron schedule and publishes to GHCR

## Vulnerability Sources

| Source | Description |
|--------|-------------|
| NVD | National Vulnerability Database |
| OSV | Open Source Vulnerabilities |
| GHSA | GitHub Security Advisories |
| Alpine SecDB | Alpine Linux security database |
| Debian | Debian Security Tracker |
| Ubuntu | Ubuntu CVE Tracker |
| Red Hat | Red Hat Security Data |
| ALAS | Amazon Linux Security Advisories |
| CISA KEV | CISA Known Exploited Vulnerabilities |

## Prerequisites

- Go 1.25 or later
- Docker (optional, for container builds)
- API keys (optional, increases rate limits):
  - `TACO_NVD_API_KEY` -- NVD API key
  - `GITHUB_TOKEN` -- GitHub token for GHSA

## Installation

### From source

```sh
git clone https://github.com/tacosec/taco-vulndb.git
cd taco-vulndb
make build
```

The binary is written to `bin/taco-vulndb`.

### Docker

```sh
# Build the CLI image
docker build -t taco-vulndb:latest .

# Build a minimal image containing only the pre-built database
docker build -f Dockerfile.vulndb -t taco-vulndb-image:latest .
```

## Usage

### Update the database

```sh
# Fetch from all sources (incremental by default)
taco-vulndb update

# Force a full historical fetch
taco-vulndb update --full

# Fetch from specific sources only
taco-vulndb update --sources nvd,osv,ghsa
```

### Check database status

```sh
taco-vulndb status
```

### Serve over HTTP

```sh
# Start on default port 8080
taco-vulndb serve

# Use a custom address
taco-vulndb serve --addr :9090
```

Endpoints:

| Path | Description |
|------|-------------|
| `GET /vulndb.json` | Database file (JSON) |
| `GET /vulndb.json.gz` | Database file (gzip-compressed) |
| `GET /meta.json` | Database metadata |
| `GET /health` | Health check |

### Distribute via OCI registry

```sh
# Push to a registry
taco-vulndb push ghcr.io/myorg/taco-vulndb:latest

# Pull from a registry
taco-vulndb pull ghcr.io/myorg/taco-vulndb:latest
```

### Export and import

```sh
# Export as gzip
taco-vulndb export --output vulndb.json.gz

# Build a standalone database file from NVD
taco-vulndb build --output ./vulndb.json --days 120

# Download a pre-built database
taco-vulndb download --url https://example.com/vulndb.json.gz

# Load a local file into the cache
taco-vulndb load --file /path/to/vulndb.json
```

### Global flags

| Flag | Description |
|------|-------------|
| `--debug` | Enable debug logging |
| `--quiet` | Suppress non-essential output |

## Tech Stack

- **Language:** Go
- **CLI framework:** [Cobra](https://github.com/spf13/cobra)
- **Core library:** [taco-lib](https://github.com/Joeri-Abbo/taco-lib) (shared vulndb logic)
- **OCI support:** [go-containerregistry](https://github.com/google/go-containerregistry)
- **CI/CD:** GitHub Actions (CI, daily DB update, Docker publish)
- **Container registry:** GitHub Container Registry (GHCR)

## Development

```sh
# Run tests
make test

# Run linter
make lint

# Run go vet
make vet

# Clean build artifacts
make clean
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes and add tests
4. Ensure `make test` and `make lint` pass
5. Commit your changes and open a pull request

## License

See the repository for license details.
