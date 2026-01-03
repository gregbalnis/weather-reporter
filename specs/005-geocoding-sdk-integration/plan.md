# Implementation Plan: Geocoding SDK Integration

**Branch**: `005-geocoding-sdk-integration` | **Date**: January 3, 2026 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/005-geocoding-sdk-integration/spec.md`

## Summary

Replace the custom Open Meteo Geocoding API client implementation with the official `open-meteo-geocoding-sdk` library. This refactoring reduces maintenance burden by delegating geocoding logic to an established external library while preserving all existing functionality and the public `GeocodingService` interface. The refactoring will be verified by both existing unit tests and a new integration test that validates API response structures remain stable over time.

## Technical Context

**Language/Version**: Go 1.25.5  
**Primary Dependencies**: `github.com/gregbalnis/open-meteo-geocoding-sdk` (new), `github.com/stretchr/testify` (existing)  
**Storage**: N/A  
**Testing**: Go `testing` package with `httptest`, `testify/assert`  
**Target Platform**: CLI application (Linux/macOS/Windows)  
**Project Type**: Single project (CLI tool)  
**Performance Goals**: Zero performance degradation - response times must remain at current levels (typically <1s for local searches)  
**Constraints**: 10-second timeout for requests, max 10 results per query, English language support  
**Scale/Scope**: Single CLI tool with integrated geocoding service used by weather lookup feature

## Current Architecture

### System Overview

The weather-reporter CLI tool uses a layered architecture:

```
CLI Entry Point (main.go)
    ↓
Application Logic (run function)
    ├→ GeocodingService interface (models/interfaces.go)
    │   └→ Custom HTTP Client (geo/client.go) - TARGET FOR REPLACEMENT
    │
    └→ WeatherService interface (models/interfaces.go)
        └→ Custom HTTP Client (weather/client.go)
```

### Current Geocoding Implementation

**File**: `src/internal/geo/client.go`
- Custom HTTP client directly calling Open Meteo Geocoding API
- Parses JSON responses into `models.Location` structs
- Configuration: 10-second timeout, 10 max results, English language
- Error handling: Returns formatted error messages

**Tests**: `src/internal/geo/client_test.go`
- Unit tests using `httptest.Server` for mocking
- Coverage: Success, no results, API errors (500), malformed JSON, timeouts

**Location Model**: `src/internal/models/models.go`
```go
type Location struct {
  ID        int     `json:"id"`
  Name      string  `json:"name"`
  Latitude  float64 `json:"latitude"`
  Longitude float64 `json:"longitude"`
  Country   string  `json:"country"`
  Region    string  `json:"admin1"`
}
```

**Interface**: `src/internal/models/interfaces.go`
```go
type GeocodingService interface {
  Search(ctx context.Context, name string) ([]Location, error)
}
```

### Dependencies & Integration Points

- **Depends On**: None (HTTP calls to external API)
- **Used By**: 
  - `main.go` - Instantiates and passes to run function
  - `run()` function - Uses for location lookup
  - UI components - Display location results
- **External**: Open Meteo Geocoding API (`https://geocoding-api.open-meteo.com/v1/search`)

## Constitution Check

### Code Quality Standards

✅ **Reduces duplication** - Eliminates custom HTTP client code  
✅ **Maintains interface isolation** - Public interface unchanged, implementation details hidden  
✅ **Improves maintainability** - Delegates to tested external library  
✅ **Preserves testability** - All existing tests continue to work  

### Testing & Coverage

✅ **No regression in test coverage** - All existing tests pass without modification  
✅ **Integration test for API stability** - New test validates API contract stability  
✅ **Performance verification** - Benchmarks ensure no degradation  

### Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|-----------|
| SDK API breaking changes in future versions | High | Medium | Pin SDK version in go.mod, new integration test detects changes |
| SDK returns different JSON structure | Medium | Low | Adapter maps SDK types to Location model, comprehensive unit tests validate |
| SDK error handling differs from custom client | Medium | Low | Wrapper functions convert SDK errors to user-friendly messages |
| Performance overhead from SDK abstraction | Medium | Low | No-degradation requirement enforced by tests, benchmark verification |
| Dependency version conflicts | Low | Low | Explicit go.mod dependency management, resolve conflicts before merge |

## Implementation Phases

### Phase 1: Research & Design

**Duration**: 1 day  
**Status**: To Start

**Goals**:
- Understand SDK API, types, and response structures
- Map SDK types to existing Location model
- Define error handling strategy
- Plan integration test for API stability detection

**Outputs**:
- `research.md` - SDK analysis, compatibility assessment, migration strategy
- `data-model.md` - Location entity definition and type mapping
- `contracts/sdk-adapter.go` - Interface definition for SDK wrapper

**Success Gate**: 
- [ ] SDK API fully understood and documented
- [ ] Type mapping strategy defined and validated
- [ ] Integration test plan defined
- [ ] All implementation unknowns resolved

### Phase 2: Core Implementation

**Duration**: 2-3 days  
**Status**: To Start

**Goals**:
- Implement SDK adapter wrapping geocoding SDK client
- Add Location model mapping from SDK types
- Update go.mod with SDK dependency (pinned version)
- Create integration test validating API contract
- Verify all existing unit tests pass with new implementation
- Ensure error handling matches specification

**Outputs**:
- Modified `src/internal/geo/client.go` - Updated to use SDK instead of custom HTTP
- New `src/internal/geo/integration_test.go` - Integration test for API stability
- Updated `go.mod` - SDK dependency added with pinned version
- All existing tests passing

**Test Coverage**:
- Unit tests: All existing tests pass (Search, No Results, API Error, Malformed JSON, Timeout)
- Integration tests: New tests validating API response structure hasn't changed
- Manual verification: Application continues to work for location searches

**Success Gate**:
- [ ] All existing tests pass without modification
- [ ] New integration test successfully validates API contract
- [ ] Zero performance degradation measured
- [ ] Error messages are user-friendly (no technical details)
- [ ] Code review approved

### Phase 3: Quality Assurance & Deployment

**Duration**: 1 day  
**Status**: To Start

**Goals**:
- Final validation and testing
- Documentation updates
- Deployment verification
- Monitor for API changes during rollout

**Outputs**:
- Updated documentation (README if SDK mentioned)
- Deployment verification checklist
- Release notes explaining the refactoring benefits

**Success Gate**:
- [ ] All quality checks pass
- [ ] No regressions in functionality or performance
- [ ] Documentation updated
- [ ] Ready for production deployment

## Project Structure

### Documentation (this feature)

```
specs/005-geocoding-sdk-integration/
├── spec.md                      # Specification (completed)
├── plan.md                      # This file (in progress)
├── research.md                  # Phase 0 output (pending)
├── data-model.md                # Phase 1 output (pending)
├── quickstart.md                # Phase 1 output (pending)
├── contracts/
│   └── sdk-adapter.go           # Interface definition (pending)
└── checklists/
    └── requirements.md          # Specification checklist
```

### Source Code Changes

**Primary Files**:
```
src/internal/geo/
├── client.go                    # MODIFY: Replace HTTP client with SDK client
├── client_test.go               # PRESERVE: All existing unit tests pass
└── integration_test.go          # NEW: Integration test for API stability

src/internal/models/
├── models.go                    # PRESERVE: Location struct unchanged
└── interfaces.go                # PRESERVE: GeocodingService interface unchanged

go.mod                           # UPDATE: Add SDK dependency (pinned version)
```

**No Changes Required**:
- `src/cmd/weather-reporter/main.go` - Uses interface, no implementation detail changes
- `src/internal/weather/` - No dependencies on geocoding implementation
- `src/internal/ui/` - No dependencies on geocoding implementation

## Implementation Strategy

### Adapter Pattern Approach

Instead of directly replacing the custom HTTP client, we use an adapter pattern:

1. **Wrapper Function**: Create a wrapper around the SDK client
2. **Type Mapping**: Map SDK response types to existing `Location` struct
3. **Error Handling**: Convert SDK errors to user-friendly messages
4. **Interface Preservation**: Maintain the existing `GeocodingService` interface

This approach:
- Isolates SDK changes to a single location (adapter)
- Makes it easy to switch implementations in the future
- Preserves all existing public interfaces
- Simplifies testing

### Error Handling Strategy

Map SDK errors to user-friendly messages as per FR-005/006/007:

```
SDK Error Type          → User Message (no technical details)
ServiceUnavailable      → "Unable to search locations. Please try again."
Timeout                 → "Search took too long. Please try again."
MalformedResponse       → "Unable to search locations. Please try again."
NetworkError            → "Unable to search locations. Please try again."
```

### Integration Test Strategy

**File**: `src/internal/geo/integration_test.go` (new)

Tests against the real Open Meteo Geocoding API (required; may be skipped when running with `-short` if network unavailable):

```go
func TestIntegration_LocationSearchContract(t *testing.T) {
  // Verify API response structure matches expected contract
  // Test cases:
  // 1. Search for "London" returns valid Location objects
  // 2. Verify all required fields present (id, name, latitude, longitude, country, admin1)
  // 3. Verify field types correct (id=int, coordinates=float64, names=string)
  // 4. Verify multiple results for ambiguous queries
  // 5. Verify max 10 results returned
}
```

Benefits:
- Documents expected API contract for future reference
- Detects API breaking changes early
- Can be run periodically to verify API stability
- Can be skipped in CI with `-short` flag if network unavailable

## Data Model

### Location Entity

The `Location` struct remains unchanged:

```go
type Location struct {
  ID        int     `json:"id"`           // Unique location identifier
  Name      string  `json:"name"`         // City or location name
  Latitude  float64 `json:"latitude"`     // Geographic latitude
  Longitude float64 `json:"longitude"`    // Geographic longitude
  Country   string  `json:"country"`      // Country name
  Region    string  `json:"admin1"`       // State/province/administrative region
}
```

### SDK Type Mapping

The SDK will provide response data that needs to be mapped to the Location struct. Mapping strategy will be documented in `data-model.md` after Phase 0 research.

## Contracts & Interfaces

### GeocodingService Interface (No Changes)

```go
// GeocodingService defines the interface for finding locations
type GeocodingService interface {
  // Search finds locations matching the given name
  Search(ctx context.Context, name string) ([]Location, error)
}
```

The public interface remains unchanged. Only the implementation changes from custom HTTP client to SDK client.

### SDK Adapter Contract

To be defined in Phase 0 (`contracts/sdk-adapter.go`), but will include:
- Wrapper struct implementing `GeocodingService`
- Constructor function `NewClient(httpClient *http.Client) GeocodingService`
- Internal mapping function from SDK Location types to our Location struct
- Error conversion from SDK errors to user-friendly messages

## Quickstart Verification

Once implementation is complete, verify the feature with:

**1. Build**
```bash
cd /workspaces/weather-reporter
go build -o bin/weather-reporter ./src/cmd/weather-reporter
```

**2. Run Unit Tests**
```bash
go test ./src/internal/geo -v
# Expected: All tests pass (existing + new integration test)
```

**3. Test Application**
```bash
./bin/weather-reporter San Francisco
# Expected: Successful location resolution and weather display
```

**4. Verify Error Handling**
- Offline test: Disable network → Should show "Unable to search locations. Please try again."
- Timeout test: Use very short timeout → Should show "Unable to search locations. Please try again."

**5. Performance Verification**
```bash
# Compare response times before/after (should be equal or faster)
time ./bin/weather-reporter London
```

## Next Steps

1. **Phase 1 (Research & Design)**
   - Research open-meteo-geocoding-sdk API and types
   - Create data-model.md with mapping strategy
   - Define contracts for SDK adapter
   - Create research.md documenting findings

2. **Phase 2 (Implementation)**
   - Implement SDK adapter in geo/client.go
   - Add integration test
   - Update go.mod
   - Run and verify all tests

3. **Phase 3 (Quality & Deployment)**
   - Final verification
   - Documentation updates
   - Prepare for deployment
