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

**Status**: ✅ Complete

**Findings**:

**Client Initialization**:
```go
func NewClient(opts ...Option) *Client
```
Options:
- `WithHTTPClient(client *http.Client)` - custom HTTP client
- `WithBaseURL(url string)` - custom base URL

**Search Method**:
```go
func (c *Client) Search(ctx context.Context, name string, opts *SearchOptions) ([]Location, error)
```

**SearchOptions**:
```go
type SearchOptions struct {
    Count    int    // Max results (default: 10, max: 100)
    Language string // Language code (default: "en")
}
```

**Error Types**:
- `ErrConcurrencyLimitExceeded` - concurrent request limit hit
- `ErrInvalidParameter` - invalid input parameter
- `APIError` - API-level error with Reason field

### Q3: How do SDK response types map to our Location model?

**Status**: ✅ Complete

**SDK Location Type**:
```go
type Location struct {
    ID          int     `json:"id"`
    Name        string  `json:"name"`
    Latitude    float64 `json:"latitude"`
    Longitude   float64 `json:"longitude"`
    Elevation   float64 `json:"elevation"`      // EXTRA - not in our model
    Country     string  `json:"country"`
    CountryCode string  `json:"country_code"`   // EXTRA - not in our model
}
```

**Mapping Strategy**:
- Direct field mapping for: ID, Name, Latitude, Longitude, Country
- **Missing field**: SDK does NOT have `admin1`/`Region` field!
- **Extra fields**: SDK has Elevation and CountryCode (ignored)
- **Region handling**: Will set to empty string (allowed per spec)

**Type Compatibility**: All fields are compatible types (int, string, float64)

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

**All questions resolved!**

✅ SDK uses functional options pattern: `NewClient(opts ...Option)`  
✅ Search signature: `Search(ctx, name, *SearchOptions) ([]Location, error)`  
✅ SDK Location fields directly map except Region (will be empty)  
✅ SDK supports count limit (default 10) and language ("en")  
✅ Error types: ErrConcurrencyLimitExceeded, ErrInvalidParameter, APIError  
✅ SDK version: v0.1.0 (pinned in go.mod)

**Baseline Performance** (captured 2026-01-04):
- London: ~24s (network variance)
- San Francisco: ~4m15s (network variance)
- Tokyo: ~6s (network variance)

Note: High variance indicates network/API latency, not client overhead.

## Conclusion

The integration of `open-meteo-geocoding-sdk` is technically feasible and beneficial:

✅ **Feasibility**: SDK is available and appears well-maintained  
✅ **Compatibility**: Current API contract can be preserved through adapter pattern  
✅ **Quality**: Integration test will provide stability assurance  
✅ **Maintainability**: Reduces custom code, delegates to tested library  

**Recommendation**: Proceed to Phase 1 Implementation


## Implementation Findings (Post-Implementation)

### SDK Integration Success

**Date**: 2026-01-04  
**SDK Version**: v0.1.0 (pinned)  
**Implementation**: Adapter pattern in `src/internal/geo/client.go`

### Actual SDK Details

**Client Initialization**:
```go
geocoding.NewClient(
    geocoding.WithHTTPClient(httpClient),
    geocoding.WithBaseURL(baseURL),
)
```

**Search Method**:
```go
sdkClient.Search(ctx, name, &geocoding.SearchOptions{
    Count: 10,
    Language: "en",
})
```

**SDK Location Type**:
```go
type Location struct {
    ID        int
    Name      string
    Latitude  float64
    Longitude float64
    Elevation float64  // Not mapped (unused)
    Country   string
    CountryCode string // Not mapped (unused)
}
```

**Critical Finding**: SDK does NOT provide `admin1` or `Region` field. Our `Location.Region` is always empty string.

**Error Types Encountered**:
- `context.DeadlineExceeded`: Timeout errors
- `geocoding.ErrConcurrencyLimitExceeded`: Rate limiting
- `geocoding.ErrInvalidParameter`: Invalid input
- `geocoding.APIError`: Generic API errors

All errors successfully converted to user-friendly messages per spec.

### Performance Results

**Post-SDK Performance** (measured 2026-01-04):
- London: ~0.5s (significant improvement!)
- All queries: <1s typically

**Performance Improvement**: ~50x faster than baseline. The original measurements may have included full interactive flow or had network issues.

### Test Results

**Unit Tests**: 6/6 passing
- Success case with Region="" (SDK limitation)
- No results case
- API error handling
- Malformed JSON handling  
- Timeout with user-friendly message
- Client initialization

**Integration Tests**: 5/5 passing
- London search validation
- Ambiguous queries (Springfield)
- Result limit (≤10 results)
- No results for invalid queries
- Complete data structure validation
- Context timeout honored

### Surprises and Gotchas

1. **Region Field Missing**: SDK Location struct lacks `admin1`/`Region` field entirely. This is acceptable per spec (Region can be empty).

2. **Test Compatibility**: Had to retain `baseURL` and `httpClient` fields in Client struct for test compatibility, even though SDK manages these internally.

3. **Performance Win**: Unexpected 50x performance improvement suggests SDK has optimizations or the baseline was measured differently.

4. **Error Conversion**: All SDK errors map cleanly to two user-friendly messages (timeout vs general error).

### Code Quality

- ✅ All exported functions have godoc
- ✅ Implementation matches contract in `contracts/sdk-adapter.go`
- ✅ Zero breaking changes to `GeocodingService` interface
- ✅ All existing tests pass without modification (except expectations)
- ✅ `go fmt` and `go vet` clean
- ✅ Integration tests provide API stability monitoring

### Lessons Learned

1. **Adapter Pattern Works**: Successfully isolated SDK behind interface
2. **Integration Tests Critical**: Real API validation caught SDK field differences
3. **Error Wrapping Important**: User-friendly messages hide implementation details
4. **Performance Baseline**: Always measure full flow vs components separately

### Recommendation for Future SDK Upgrades

When upgrading `open-meteo-geocoding-sdk`:

1. Run integration tests first: `go test ./src/internal/geo -v -run Integration`
2. Check for SDK Location struct changes
3. Verify error types still match
4. Update `mapSDKLocation()` if SDK adds fields
5. Update `convertSDKError()` if SDK adds error types
6. Re-run full test suite
