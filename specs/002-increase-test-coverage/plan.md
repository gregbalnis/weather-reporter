# Implementation Plan: Increase Unit Test Coverage

**Branch**: `002-increase-test-coverage` | **Date**: 2025-12-27 | **Spec**: [specs/002-increase-test-coverage/spec.md](../spec.md)
**Input**: Feature specification from `/specs/002-increase-test-coverage/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Increase unit test coverage to a minimum of 80% across all internal packages (`geo`, `weather`, `ui`) to comply with the Constitution and ensure code quality. This involves identifying coverage gaps, implementing new tests using `testify`, and ensuring all tests are deterministic and fast.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `github.com/stretchr/testify`
**Storage**: N/A
**Testing**: `go test`, `github.com/stretchr/testify`
**Target Platform**: Linux (CLI)
**Project Type**: Single project (CLI)
**Performance Goals**: Individual unit tests < 100ms
**Constraints**: Minimum 80% code coverage per package
**Scale/Scope**: ~3 core packages to test

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality**: New tests must follow Effective Go idioms.
- **II. Testing Standards**: This feature directly implements the 80% coverage requirement.
- **III. User Experience Consistency**: N/A (internal quality).
- **IV. Performance Requirements**: Tests must be fast (sub-second suite).

## Project Structure

### Documentation (this feature)

```text
specs/002-increase-test-coverage/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
src/
├── cmd/
│   └── weather/         # Main entry point (excluded from strict 80% if just wiring)
└── internal/
    ├── geo/             # Location logic (Target: >80% coverage)
    ├── models/          # Interfaces and data structures
    ├── ui/              # User interaction (Target: >80% coverage)
    └── weather/         # Weather logic (Target: >80% coverage)
```

**Structure Decision**: Standard Go project layout with `cmd` and `internal`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A - No violations.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
