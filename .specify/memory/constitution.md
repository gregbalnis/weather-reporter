<!--
Sync Impact Report:
- Version change: 1.1.0 -> 1.2.0
- List of modified principles: Documentation Standards (added README requirement)
- Added sections: V. Documentation Standards
- Removed sections: None
- Templates requiring updates: ✅ tasks-template.md
- Follow-up TODOs: None
-->
# weather-reporter Constitution

## Core Principles

### I. Code Quality (Effective Go)
We strictly adhere to the idioms and best practices outlined in [Effective Go](https://go.dev/doc/effective_go). Code MUST be formatted with `gofmt`. Naming conventions, error handling, and concurrency patterns MUST follow Go community standards. Clarity and simplicity are preferred over cleverness. Public APIs MUST be documented.

### II. Testing Standards
Testing is mandatory. All packages MUST have unit tests (`_test.go`) colocated with source code. Test coverage MUST meet a minimum of 80% for unit tests. Integration tests are REQUIRED for external interactions (APIs, file systems). Tests MUST be deterministic and fast.

### III. User Experience Consistency
The CLI and output MUST be consistent and predictable. Use standard flags and arguments. Output SHOULD be human-readable by default, with options for machine-readable formats (e.g., JSON) where appropriate. Error messages MUST be actionable and clear to the end-user, distinguishing between user errors and system failures.

### IV. Performance Requirements
The application MUST be efficient with resources (CPU, Memory). Avoid unnecessary allocations in hot paths. Network operations MUST have timeouts. Performance critical paths SHOULD be benchmarked. Latency for user-facing operations SHOULD be minimized.

### V. Documentation Standards
A `README.md` file is REQUIRED at the root of the repository. It MUST document the current state of the project, including what it does, why it is useful, how to get started, and how to contribute. This file MUST be updated after implementation is complete and before a Pull Request is created to reflect any changes in functionality, usage, or configuration.

## Technical Stack

**Language**: Go (Latest Stable)
**Dependency Management**: Go Modules
**Linter**: `golangci-lint` (Standard configuration)
**Build Tool**: Standard `go build`

## Development Workflow

**Branching**: Feature branches off `main`.
**Commits**: Follow Conventional Commits specification.
**Review**: Pull Request required for all changes. Code review MUST verify compliance with Core Principles.
**CI/CD**: Automated tests and linters MUST pass before merging.

## Governance

This Constitution supersedes all other practices. Amendments require documentation, approval, and a migration plan.

**Rules**:
1.  All PRs/reviews MUST verify compliance with "Effective Go" and this Constitution.
2.  Complexity MUST be justified.
3.  New dependencies MUST be vetted for license and maintenance status.
4.  Versioning follows Semantic Versioning (SemVer).

**Version**: 1.2.0 | **Ratified**: 2025-12-25 | **Last Amended**: 2025-12-27
