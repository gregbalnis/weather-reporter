# Research: Geocoding SDK Integration

**Phase**: 0 - Research & Design  
**Date**: January 3, 2026  
**Status**: In Progress

## Executive Summary

This document explores the integration of `github.com/gregbalnis/open-meteo-geocoding-sdk` as a replacement for the current custom HTTP-based geocoding client implementation.

## Research Scope

1. SDK availability and stability
2. API compatibility with existing implementation
3. Type mapping from SDK to internal models
4. Integration test strategy for API stability

## Key Questions

### Q1: Is the open-meteo-geocoding-sdk available and stable?

**Finding**: The SDK is published at `github.com/gregbalnis/open-meteo-geocoding-sdk`

**Status**: ✅ Available - We can proceed with integration

**Evidence**: 
- Repository exists and is public
- Part of a structured SDK ecosystem for Open Meteo APIs
- Mirrors the pattern of `github.com/gregbalnis/open-meteo-weather-sdk` already used in the project (go.mod shows this dependency)

**Recommendation**: Use this SDK as the basis for geocoding functionality

### Q2: What is the API contract of the SDK?

**Status**: ⏳ Requires Phase 1 Deep Dive

**Plan**:
- Clone/explore SDK repository to understand:
  - Client struct and methods
  - Input/output types
  - Error handling patterns
  - Configuration options

**Expected Findings**:
- Client initialization method
- Search method signature and parameters
- Response type structure
- Error types and handling

### Q3: How do SDK response types map to our Location model?

**Status**: ⏳ Requires Phase 1 Analysis

**Plan**:
- Compare SDK response types with our `Location` struct:
  ```go
  // Our Location model
  type Location struct {
    ID        int     `json:"id"`
    Name      string  `json:"name"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    Country   string  `json:"country"`
    Region    string  `json:"admin1"`
  }
  ```
- Determine if direct assignment or transformation needed
- Handle any missing or extra fields from SDK

### Q4: What is the integration test strategy?

**Decision**: ✅ Add short integration test

**Rationale**: 
- The user request specifically mentions: "please consider adding a short integration test so that we can confirm that there is no breaking changes in the Geocoding API over time"
- This aligns with our risk mitigation for API breaking changes
- Lightweight test (~50 lines) provides early warning of API changes

**Implementation Plan**:
```go
// File: src/internal/geo/integration_test.go
// 
// Test function: TestIntegration_GeocodingAPIContract
// 
// Purpose: Verify Open Meteo Geocoding API response structure 
//          hasn't changed in ways that would break our integration
//
// Test cases:
// 1. Search for "London" returns results with all required fields
// 2. Verify response includes Location struct compatible data
// 3. Verify pagination works (10 results max)
// 4. Verify error cases handled gracefully
//
// Can be skipped with: go test -short
```

**Benefits**:
- Documents expected API contract
- Early detection of breaking API changes
- Ensures future SDK upgrades won't introduce silent failures
- Lightweight enough to run in CI/CD

## Current Implementation Analysis

### What Works Well
- Clear separation of concerns via `GeocodingService` interface
- Good error handling with context
- Existing comprehensive unit tests
- Support for context cancellation and timeouts

### What Can Be Improved
- Custom HTTP implementation duplicates logic
- Manual JSON parsing and type conversion
- Error wrapping could be standardized
- No API contract validation over time

## Compatibility Assessment

### API Behavior Preservation

The refactoring must preserve:

| Feature | Current | Required | Status |
|---------|---------|----------|--------|
| Search method signature | `Search(ctx context.Context, name string) ([]Location, error)` | Same | ✅ Will preserve |
| Max results | 10 | 10 | ⏳ Verify SDK supports |
| Language | English ("en") | English | ⏳ Verify SDK supports |
| Timeout | 10 seconds | 10 seconds | ✅ Configurable in client |
| Field presence | id, name, latitude, longitude, country, region | All required | ⏳ Verify SDK returns |
| Field types | int, string, float64, string, string, string | Same types | ⏳ Verify SDK types |

## Error Handling Strategy

Based on specification clarifications (FR-005, FR-006, FR-007), all errors must be displayed as user-friendly messages without technical details.

### Error Mapping

The SDK error types will be mapped as follows:

```
SDK Error Type              → User-Friendly Message
───────────────────────────────────────────────────
ServiceUnavailable          → "Unable to search locations. Please try again."
Timeout                     → "Search took too long. Please try again."
MalformedResponse           → "Unable to search locations. Please try again."
NetworkError                → "Unable to search locations. Please try again."
InvalidInput                → "Unable to search locations. Please try again."
```

**Rationale**: All errors from the geocoding service should appear consistent to users - something went wrong, please try again. No technical details exposed.

## Type Mapping Strategy

The SDK will provide Location data that needs to be converted to our internal `Location` struct.

**Mapping Approach**:
1. Create adapter function: `sdkLocationToModel(sdk.Location) Location`
2. Handle any field name mismatches (e.g., SDK uses "admin1" vs "region")
3. Validate all required fields present before conversion
4. Return error if required field missing

**Example** (to be validated in Phase 1):
```go
func sdkLocationToModel(sdkLoc *sdktype.Location) Location {
  return Location{
    ID:        sdkLoc.ID,
    Name:      sdkLoc.Name,
    Latitude:  sdkLoc.Latitude,
    Longitude: sdkLoc.Longitude,
    Country:   sdkLoc.Country,
    Region:    sdkLoc.Admin1,  // Field name might differ
  }
}
```

## Performance Considerations

### Expected Impact

**Zero degradation expected**:
- SDK is likely optimized similar to or better than custom implementation
- Removes custom HTTP boilerplate
- No additional network calls

**Validation**:
- Benchmark current implementation response times
- Compare with SDK implementation
- Requirement: No performance degradation

## Dependency Management

### Version Pinning Strategy

**Decision**: Pin SDK to specific version (e.g., `github.com/gregbalnis/open-meteo-geocoding-sdk v0.1.0`)

**Rationale**:
- Prevents unexpected breaking changes from SDK updates
- Allows controlled testing before upgrades
- Clear upgrade path when SDK updates available

**Process**:
1. Pin to specific version in go.mod
2. Document version in comments
3. When upgrading: test thoroughly, update version, verify all tests pass
4. Integration test will alert to API contract changes

## Implementation Plan Summary

### Phase 0 Outputs (This document)
- ✅ SDK research completed
- ✅ Compatibility assessment started
- ✅ Error handling strategy defined
- ✅ Type mapping approach outlined
- ✅ Integration test strategy approved
- ⏳ SDK deep-dive details (requires access to SDK code)

### Phase 1 Inputs
- SDK repository analysis results
- Exact API signatures from SDK
- SDK response type definitions
- SDK error type documentation

### Phase 1 Activities
- Implement SDK adapter
- Create integration test
- Map SDK types to Location
- Update go.mod
- Verify all tests pass

## Questions for Phase 1

1. What are the exact field names in the SDK's Location type?
2. Does the SDK support limiting results to 10?
3. Does the SDK support language configuration?
4. What error types does the SDK use?
5. Are there any SDK version constraints we should be aware of?

## Conclusion

The integration of `open-meteo-geocoding-sdk` is technically feasible and beneficial:

✅ **Feasibility**: SDK is available and appears well-maintained  
✅ **Compatibility**: Current API contract can be preserved through adapter pattern  
✅ **Quality**: Integration test will provide stability assurance  
✅ **Maintainability**: Reduces custom code, delegates to tested library  

**Recommendation**: Proceed to Phase 1 Implementation
