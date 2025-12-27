# Data Model: Release Artifacts & Versioning

## Release Artifacts

The pipeline will produce the following artifacts for each release:

| Artifact Type | Naming Convention | Description |
| :--- | :--- | :--- |
| **Binary (Linux)** | `weather-reporter_linux_amd64.tar.gz` | Compressed binary for Linux (x86_64) |
| **Binary (Linux ARM)** | `weather-reporter_linux_arm64.tar.gz` | Compressed binary for Linux (ARM64) |
| **Binary (macOS)** | `weather-reporter_darwin_amd64.tar.gz` | Compressed binary for macOS (Intel) |
| **Binary (macOS ARM)** | `weather-reporter_darwin_arm64.tar.gz` | Compressed binary for macOS (Apple Silicon) |
| **Binary (Windows)** | `weather-reporter_windows_amd64.zip` | Compressed binary for Windows (x86_64) |
| **Checksums** | `checksums.txt` | SHA256 checksums for all artifacts |
| **SBOM** | `sbom.spdx.json` | Software Bill of Materials in SPDX format |

## Version Metadata

The application will embed the following metadata at build time using `-ldflags`:

| Field | Variable in Code | Source | Example |
| :--- | :--- | :--- | :--- |
| **Version** | `main.version` | Git Tag | `v1.2.3` |
| **Commit** | `main.commit` | Git Commit Hash | `a1b2c3d` |
| **Date** | `main.date` | Build Timestamp | `2025-12-27T10:00:00Z` |

## Configuration Models

### GoReleaser Config (`.goreleaser.yaml`)
- **Builds**: Defines targets (OS/Arch), binary name, and ldflags.
- **Archives**: Defines compression format (tar.gz/zip) and file structure.
- **Checksum**: Defines hashing algorithm (sha256).
- **SBOMs**: Defines generation tool (syft) and format.
- **Release**: Defines GitHub release settings (draft, prerelease).

### CI Workflow Model
- **Triggers**:
    - `push` to `main` (CI checks)
    - `pull_request` (CI checks)
    - `push` tags `v*` (Release)
- **Jobs**:
    - `lint`: Runs `golangci-lint`
    - `test`: Runs `go test -race`
    - `release`: Runs `goreleaser` (only on tags)
