# Contributing to Weather Reporter

## Development

### Prerequisites
- Go 1.21+
- Make

### Common Tasks
- `make build`: Build the binary
- `make test`: Run tests (includes integration tests requiring network)
- `go test -short ./src/...`: Run unit tests only (skip integration tests)
- `make lint`: Run linter

## Release Process

Releases are automated using GitHub Actions and GoReleaser.

1.  Ensure all changes are merged to `main`.
2.  Create a new tag following semantic versioning (e.g., `v1.0.0`).
    ```bash
    git tag v1.0.0
    git push origin v1.0.0
    ```
3.  The "Release" workflow will automatically run:
    -   Build binaries for Linux, Windows, and macOS.
    -   Create a GitHub Release.
    -   Upload artifacts.
