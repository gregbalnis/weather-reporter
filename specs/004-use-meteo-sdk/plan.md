# Implementation Plan: Integrate Open-Meteo SDK

**Branch**: `004-use-meteo-sdk` | **Date**: 2025-12-31 | **Spec**: [specs/004-use-meteo-sdk/spec.md](specs/004-use-meteo-sdk/spec.md)
**Input**: Feature specification from `/specs/004-use-meteo-sdk/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Refactor the weather reporter application to use the external `open-meteo-weather-sdk` for fetching weather data. This replaces the custom internal weather client, reducing maintenance burden while preserving existing functionality. The SDK provides `QuantityOf...` accessors for formatted output (value + unit), simplifying the display logic.

## Technical Context

**Language/Version**: Go 1.25.5
**Primary Dependencies**: `github.com/gregbalnis/open-meteo-weather-sdk`
**Storage**: N/A
**Testing**: Go standard library `testing` package
**Target Platform**: Linux (Dev Container), Cross-platform CLI
**Project Type**: Single project (CLI)
**Performance Goals**: Weather retrieval must complete within 10 seconds
**Constraints**: Must maintain exact output format and behavior; Fail fast on errors
**Scale/Scope**: Small CLI utility, single external integration

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality**: The integration must follow Effective Go. The SDK usage should be idiomatic.
- **II. Testing Standards**: Existing tests for the custom client will be replaced by tests for the SDK integration. Coverage must be maintained > 80%.
- **III. User Experience Consistency**: The output format must remain identical to the current version. The SDK's `QuantityOf...` accessors will be used to ensure correct unit formatting.
- **IV. Performance Requirements**: A 10-second timeout must be enforced on the SDK client.
- **V. Documentation Standards**: README must be updated to reflect the new dependency.
- **VI. Release & Build Standards**: `go.mod` and `go.sum` will be updated.

## Project Structure

### Documentation (this feature)

```text
specs/004-use-meteo-sdk/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
src/
├── cmd/
│   └── weather-reporter/
│       └── main.go      # Refactored to use SDK
├── internal/
│   ├── weather/         # Custom client removed/replaced with SDK adapter if needed
│   ├── models/          # Updated to align with SDK types
│   └── ui/              # Updated to use SDK response types
```

**Structure Decision**: Option 1: Single project. The `internal/weather` package will be significantly simplified or removed in favor of direct SDK usage or a thin adapter.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A
