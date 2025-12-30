# Implementation Plan: Integrate Open-Meteo SDK

**Branch**: `004-use-meteo-sdk` | **Date**: 2025-12-30 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-use-meteo-sdk/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Replace custom weather API client implementation with the external open-meteo-weather-sdk library to reduce maintenance burden while maintaining 100% functional compatibility. The SDK provides a maintained client for fetching current weather data, eliminating the need for custom HTTP client code. All existing weather data points (temperature, humidity, precipitation, cloud cover, pressure, wind speed/direction/gusts) will continue to be retrieved and displayed identically. The refactoring will use an all-at-once replacement approach with updated SDK-specific tests.

## Technical Context

**Language/Version**: Go 1.25.5  
**Primary Dependencies**: 
- github.com/stretchr/testify v1.11.1 (testing)
- github.com/gregbalnis/open-meteo-weather-sdk (to be added)
**Storage**: N/A (stateless CLI application)  
**Testing**: Go's built-in testing with testify assertions, race detection enabled  
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows)  
**Project Type**: Single project (CLI application)  
**Performance Goals**: Weather data retrieval within 10 seconds timeout  
**Constraints**: 
- 10-second maximum timeout for weather requests
- Fail-fast error handling (no retries)
- Metric units only (Celsius, km/h) - SDK returns metric by default
- Standard output for all messages (errors and data)
**Scale/Scope**: Single-user CLI tool, ~1000 LOC, 4 internal packages

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Initial Check (Pre-Phase 0)

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. Code Quality (Effective Go)** | ✅ PASS | SDK integration maintains Go idioms; existing `gofmt` compliance preserved |
| **II. Testing Standards** | ⚠️ MODIFY | Tests will be updated (not just added); current 80%+ coverage must be maintained |
| **III. User Experience Consistency** | ✅ PASS | Zero user-facing changes; CLI behavior identical before/after |
| **IV. Performance Requirements** | ✅ PASS | 10s timeout enforced; SDK uses similar HTTP patterns to current implementation |
| **V. Documentation Standards** | ⚠️ UPDATE | README.md must be updated if SDK changes dependency installation steps |
| **VI. Release & Build Standards** | ✅ PASS | Standard Go build process maintained; new dependency added to go.mod |

**Overall Gate**: ✅ **PASS** - All constitutional requirements satisfied

---

### Post-Phase 1 Re-Check

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. Code Quality (Effective Go)** | ✅ PASS | Design maintains Go idioms; adapter pattern if needed follows standard practices |
| **II. Testing Standards** | ✅ PASS | Integration tests designed; 80%+ coverage maintained via SDK integration tests |
| **III. User Experience Consistency** | ✅ PASS | Design preserves exact CLI output format and error messages |
| **IV. Performance Requirements** | ✅ PASS | 10s timeout enforced via http.Client and context; SDK respects timeouts |
| **V. Documentation Standards** | ✅ PASS | Quickstart.md documents SDK integration; README updates identified |
| **VI. Release & Build Standards** | ✅ PASS | No build process changes; standard `go get` for dependency |

**Overall Gate**: ✅ **PASS** - Design complies with all constitutional requirements

**Actions Completed**:
1. ✅ Integration test strategy defined (real API calls)
2. ✅ Quickstart.md created with implementation guide
3. ✅ Agent context updated with SDK technology

## Project Structure

### Documentation (this feature)

```text
specs/004-use-meteo-sdk/
├── spec.md              # Feature specification (completed)
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (SDK evaluation)
├── data-model.md        # Phase 1 output (data structure mapping)
├── quickstart.md        # Phase 1 output (integration guide)
├── contracts/           # Phase 1 output (SDK interface contracts)
│   └── sdk-adapter.go   # SDK interface contract
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
weather-reporter/
├── go.mod                              # Add open-meteo-weather-sdk dependency
├── go.sum                              # Updated checksums
├── Makefile                            # Unchanged (standard Go build)
├── README.md                           # Update if SDK installation differs
├── src/
│   ├── cmd/
│   │   └── weather-reporter/
│   │       └── main.go                 # Update to use SDK client instead of internal/weather
│   ├── internal/
│   │   ├── geo/                        # UNCHANGED (location search)
│   │   │   ├── client.go
│   │   │   └── client_test.go
│   │   ├── models/                     # UPDATE (adapt to SDK response types)
│   │   │   ├── interfaces.go
│   │   │   └── models.go               # Map SDK types to existing models
│   │   ├── ui/                         # UNCHANGED (user interaction)
│   │   │   ├── prompt.go
│   │   │   ├── prompt_test.go
│   │   │   └── mocks_test.go
│   │   └── weather/                    # REPLACE or REMOVE
│   │       ├── client.go               # DELETE (replaced by SDK)
│   │       └── client_test.go          # REPLACE with SDK integration tests
│   └── ...
└── tests/                               # Integration tests (if added)
```

**Structure Decision**: Single project structure maintained. The `internal/weather` package will be eliminated entirely, with SDK usage directly in `main.go` or through a thin adapter if needed for testing. Location search (`internal/geo`) and UI (`internal/ui`) packages remain unchanged. Models package updated to handle SDK response mapping.

## Complexity Tracking

**No constitutional violations requiring justification.**

All complexity is within acceptable bounds:
- Single project structure (standard for CLI)
- No new architectural patterns introduced
- Direct SDK usage without unnecessary abstraction layers
- Test updates are maintenance, not added complexity

---

## Phase Summary

### Phase 0: Research ✅ COMPLETE

**Deliverable**: [research.md](research.md)

**Key Findings**:
1. SDK provides all required weather parameters
2. SDK returns metric units by default (matches requirements)
3. Timeout control via custom `http.Client` with 10s timeout
4. Integration tests with real API calls (Open-Meteo is free/fast)
5. Adapter pattern available if SDK types differ from our models

**Risks Mitigated**:
- SDK compatibility verified through research
- Testing strategy defined (integration over mocking)
- Error handling approach clarified (fail-fast)

---

### Phase 1: Design ✅ COMPLETE

**Deliverables**:
- [data-model.md](data-model.md) - Entity definitions and SDK type mapping strategy
- [contracts/sdk-adapter.go](contracts/sdk-adapter.go) - Interface contract for weather client
- [quickstart.md](quickstart.md) - Step-by-step implementation guide

**Key Decisions**:
1. **Integration Approach**: Direct SDK usage with adapter if type mapping needed
2. **Data Model**: Keep existing models, add mapping function if SDK types differ
3. **Testing**: Integration tests using real Open-Meteo API calls
4. **Code Removal**: Delete entire `internal/weather` package (~150 LOC)
5. **Timeout Enforcement**: 10s via `http.Client` passed to SDK

**Design Artifacts**:
- WeatherClient interface contract documented
- SDK adapter pattern specified (Option A vs Option B)
- Test strategy: unit tests → integration tests
- Migration steps outlined in quickstart

---

### Phase 2: Task Breakdown - NEXT STEP

**Command**: `/speckit.tasks` - Creates [tasks.md](tasks.md)

**Expected Output**: Detailed implementation tasks including:
1. Add SDK dependency
2. Examine SDK types
3. Implement adapter (if needed)
4. Update main.go
5. Add integration tests
6. Delete custom client code
7. Update documentation
8. Verify constitution compliance

---

## Implementation Readiness

| Requirement | Status | Notes |
|-------------|--------|-------|
| **Technical Context** | ✅ Complete | All NEEDS CLARIFICATION resolved |
| **Constitution Gates** | ✅ Passing | All principles satisfied |
| **Research** | ✅ Complete | SDK evaluated, decisions made |
| **Design** | ✅ Complete | Data model, contracts, quickstart ready |
| **Test Strategy** | ✅ Defined | Integration tests with real API |
| **Migration Plan** | ✅ Documented | Quickstart provides step-by-step guide |

**Status**: ✅ **READY FOR PHASE 2** (Task breakdown via `/speckit.tasks`)

---

## Quick Reference

- **Spec**: [spec.md](spec.md)
- **Research**: [research.md](research.md)
- **Data Model**: [data-model.md](data-model.md)
- **Contracts**: [contracts/sdk-adapter.go](contracts/sdk-adapter.go)
- **Implementation Guide**: [quickstart.md](quickstart.md)
- **Branch**: `004-use-meteo-sdk`
- **SDK**: https://github.com/gregbalnis/open-meteo-weather-sdk

---

## Notes for Implementation

1. **SDK Note**: User confirmed SDK returns metric units by default - no configuration needed
2. **All-at-once**: Use single changeset for complete replacement (not incremental)
3. **Testing**: Real API calls acceptable (Open-Meteo is free, no key required)
4. **Timeout**: Enforce 10s via `http.Client` + context for safety
5. **Error Output**: Standard output (stdout) for consistency with current implementation

