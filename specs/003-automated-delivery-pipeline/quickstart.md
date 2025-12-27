# Quickstart: Automated Delivery Pipeline

## Prerequisites
- **Go**: 1.25+
- **Make**: Standard version
- **GoReleaser**: Install via `go install github.com/goreleaser/goreleaser/v2@latest` (optional, for local snapshots)
- **golangci-lint**: Install via `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## Local Development

### Running Tests
Run all unit tests with race detection:
```bash
make test
```

### Linting Code
Run the linter locally:
```bash
make lint
```

### Building Locally
Build the binary for your current OS/Arch:
```bash
make build
```
The binary will be placed in `bin/weather-reporter`.

### Checking Version
```bash
./bin/weather-reporter --version
```

## Release Process

### Creating a Local Snapshot
To test the release process locally without publishing:
```bash
make snapshot
```
Artifacts will be generated in the `dist/` directory.

### Triggering a Release
1. Ensure all changes are committed and pushed to `main`.
2. Create a semantic version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. The GitHub Actions pipeline will automatically:
   - Run linters and tests.
   - Build cross-platform binaries.
   - Generate SBOMs and checksums.
   - Publish a new Release on GitHub.
