# Research: Automated Delivery Pipeline

## Decisions

### 1. Release Orchestration
- **Decision**: Use **GoReleaser**.
- **Rationale**: GoReleaser is the industry standard for Go projects. It allows for a declarative configuration (`.goreleaser.yaml`) to manage cross-compilation, archive creation, checksum generation, and publishing to GitHub Releases. It simplifies the build matrix and ensures consistency.
- **Alternatives Considered**:
    - **Custom Shell Scripts**: Hard to maintain, error-prone, difficult to handle cross-compilation and artifact uploading reliably.
    - **GitHub Actions Matrix**: Requires complex logic to aggregate artifacts from multiple jobs before creating a release.

### 2. SBOM Generation
- **Decision**: Use **Syft**.
- **Rationale**: Syft is a robust tool for generating Software Bill of Materials (SBOMs) and integrates natively with GoReleaser. It supports standard formats like SPDX and CycloneDX, which are essential for supply chain security.
- **Implementation Detail**: Use Syft action in the CI pipeline (`anchore/sbom-action@v0`) before GoReleaser runs, as `goreleaser-action` does not include it.
- **Alternatives Considered**:
    - **CycloneDX-Go**: A valid option but requires separate configuration and execution steps compared to Syft's integration.
    - **Go List**: Produces simple text output that is not a standard SBOM format.

### 3. Linting
- **Decision**: Use **golangci-lint**.
- **Rationale**: It is the most popular and comprehensive linter aggregator for Go. It runs fast and supports a wide range of linters (including `staticcheck`, `govet`, etc.) with a single configuration file.
- **Alternatives Considered**:
    - **go vet**: Too basic, misses many common issues.
    - **staticcheck**: Excellent tool, but `golangci-lint` includes it along with many others.

### 4. Versioning Strategy
- **Decision**: **Semantic Versioning** driven by **Conventional Commits**.
- **Rationale**: Allows for automated version calculation and changelog generation. It removes human error from the release process and ensures a standardized history.
- **Alternatives Considered**:
    - **Manual Tagging**: Prone to errors and inconsistencies.

### 5. CI/CD Provider
- **Decision**: **GitHub Actions**.
- **Rationale**: Native integration with the repository, no external setup required, supports OIDC for secure cloud interactions (if needed later), and has a vast marketplace of actions (like `setup-go`, `goreleaser-action`).

## Technical Constraints & Unknowns Resolved
- **Cross-Compilation**: GoReleaser handles `GOOS` and `GOARCH` combinations automatically.
- **Permissions**: We will use the `GITHUB_TOKEN` with specific permissions (`contents: write` only for release jobs).
- **Caching**: We will use `actions/setup-go` which has built-in caching for Go modules.
