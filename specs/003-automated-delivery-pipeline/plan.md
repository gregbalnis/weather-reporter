# Implementation Plan: Automated Delivery Pipeline

**Branch**: `003-automated-delivery-pipeline` | **Date**: 2025-12-27 | **Spec**: [specs/003-automated-delivery-pipeline/spec.md](spec.md)
**Input**: Feature specification from `/specs/003-automated-delivery-pipeline/spec.md`

## Summary

Implement a robust, automated delivery pipeline using GitHub Actions and GoReleaser. This includes enforcing quality gates (linting, testing) on Pull Requests, automating semantic versioned releases with cross-compiled binaries, generating SBOMs and checksums, and ensuring reproducible builds via containerized environments and declarative configuration.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `goreleaser` (build/release), `golangci-lint` (linting), `syft` (SBOM)
**Storage**: N/A
**Testing**: `go test` (standard library)
**Target Platform**: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64)
**Project Type**: CLI Tool
**Performance Goals**: CI feedback < 5 mins, Release < 10 mins
**Constraints**: Public repository, GitHub Actions free tier limits
**Scale/Scope**: Single repository, multiple build targets

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **I. Code Quality**: Linter (`golangci-lint`) will be enforced.
- [x] **II. Testing Standards**: Unit tests with race detection will be enforced (80% coverage).
- [x] **III. UX Consistency**: Version command will be standardized.
- [x] **IV. Performance**: Caching will be used to optimize build times.
- [x] **V. Documentation**: README will be updated.
- [x] **VI. Release & Build Standards**: This feature directly implements this section (Reproducible, Secure, Automated, SemVer, Cross-compiled, Integrity checks).

## Project Structure

### Documentation (this feature)

```text
specs/003-automated-delivery-pipeline/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A for this feature)
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
.github/
└── workflows/
    └── ci.yml           # CI/CD Workflow definition
.goreleaser.yaml         # GoReleaser configuration
.golangci.yml            # Linter configuration
Makefile                 # Local task runner
src/
└── cmd/
    └── weather-reporter/
        └── main.go      # Updated with version logic
```

**Structure Decision**: Standard Go project layout with root-level configuration files for build tools.

## Complexity Tracking

N/A - Complexity is justified by the requirement for a secure and automated supply chain.

## Implementation Steps

### Step 1: Linter Configuration & Fixes
- **Goal**: Establish a baseline for code quality.
- **Action**: Create `.golangci.yml` with standard linters enabled. Run it locally and fix any existing issues.
- **Verification**: `golangci-lint run` passes locally.

### Step 2: Versioning Implementation
- **Goal**: Allow the application to report its version.
- **Action**: Update `src/cmd/weather-reporter/main.go` to define `version`, `commit`, `date` variables and handle the `--version` flag.
- **Verification**: `go run -ldflags "-X main.version=dev" src/cmd/weather-reporter/main.go --version` prints the version.

### Step 3: Makefile Creation
- **Goal**: Standardize local development tasks.
- **Action**: Create a `Makefile` with targets: `test` (with coverage), `lint` (with golangci-lint), `build`, `snapshot` (local release), `clean`.
- **Verification**: `make test`, `make lint`, `make build` all execute successfully.

### Step 4: GoReleaser Configuration
- **Goal**: Define the release process declaratively.
- **Action**: Create `.goreleaser.yaml`. Configure:
    - Builds for Linux, macOS, Windows (amd64/arm64).
    - Binary renaming.
    - Archives (tar.gz/zip).
    - Checksums.
    - SBOM generation via Syft.
    - Changelog generation.
- **Verification**: `goreleaser release --snapshot --clean` generates artifacts in `dist/`.

### Step 5: GitHub Actions Workflow
- **Goal**: Automate the pipeline.
- **Action**: Create `.github/workflows/ci.yml`.
    - **Triggers**: Push to main, PRs, Tags.
    - **Job 1: Quality**:
        - Checkout (`actions/checkout@v6`)
        - Setup Go (`actions/setup-go@v6`, cache enabled)
        - Lint (`golangci/golangci-lint-action@v9`)
        - Test (`go test -race` + coverage check). Fail if coverage < 80%.
    - **Job 2: Release**:
        - Needs Quality. Runs only on tags.
        - Checkout (`actions/checkout@v6`, fetch-depth: 0 for changelog)
        - Setup Go (`actions/setup-go@v6`)
        - Create SBOM (`anchore/sbom-action@v0`)
        - Run GoReleaser (`goreleaser/goreleaser-action@v6`)
    - **Permissions**: `contents: write` (for release), `read` (default).
    - **Concurrency**: Cancel in-progress group.
- **Verification**: Push to branch triggers Quality job. Tag push triggers Release job.

### Step 6: Documentation Update
- **Goal**: Document the new build/release process.
- **Action**: Update `README.md` with badges (CI status), installation instructions (from release), and development commands (Makefile).
- **Verification**: README is accurate and helpful.
