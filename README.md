# taco-vulndb

Vulnerability database builder and distributor for the TACO ecosystem.

Fetches vulnerability data from multiple sources, merges with precedence rules,
and distributes via OCI registries, HTTP, or file export.

This tool imports its core vulndb logic from `github.com/jabbo/taco/pkg/vulndb`
to avoid code duplication. A `replace` directive in `go.mod` points to the
local `../taco` directory for development.

## Sources

- NVD (National Vulnerability Database)
- OSV (Open Source Vulnerabilities)
- GHSA (GitHub Security Advisories)
- Alpine SecDB
- Debian Security Tracker
- Ubuntu CVE Tracker
- Red Hat Security Data
- Amazon Linux Security Advisories (ALAS)
- CISA Known Exploited Vulnerabilities (KEV)

## Usage

```sh
# Update from all sources
taco-vulndb update

# Push to OCI registry
taco-vulndb push ghcr.io/myorg/vulndb:latest

# Serve over HTTP
taco-vulndb serve --addr :8080

# Export for distribution
taco-vulndb export --output vulndb.json.gz
```

## Docker

Build from the parent directory (so the `taco` module is available):

```sh
docker build -f taco-vulndb/Dockerfile -t taco-vulndb .
docker build -f taco-vulndb/Dockerfile.vulndb -t taco-vulndb-image .
```
