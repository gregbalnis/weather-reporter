# Quickstart: Geocoding SDK Integration

**Phase**: 1 - Core Implementation  
**Date**: January 3, 2026  
**Status**: Implementation Guide

## Overview

This guide walks through implementing and verifying the geocoding SDK integration.

## Prerequisites

- Go 1.25.5 installed
- Access to internet (for integration tests)
- Current working directory: `/workspaces/weather-reporter`

## Step 1: Understand the SDK

### Research the SDK Repository

1. Visit: https://github.com/gregbalnis/open-meteo-geocoding-sdk
2. Review the README for:
   - Client initialization
   - Search method signature
   - Response types
   - Error handling

### Key Questions to Answer

- [ ] What is the exact client struct name?
- [ ] How do you initialize the client?
- [ ] What method name is used for searching?
- [ ] What are the parameter types?
- [ ] What is the return type structure?
- [ ] What error types does it define?

## Step 2: Create the Adapter

### File: `src/internal/geo/client.go`

Replace the current custom HTTP implementation with an SDK adapter:

```go
package geo

import (
  "context"
  "fmt"
  
  sdk "github.com/gregbalnis/open-meteo-geocoding-sdk"
  "weather-reporter/src/internal/models"
)

// Client wraps the open-meteo-geocoding-sdk to implement GeocodingService
type Client struct {
  sdkClient *sdk.Client
}

// NewClient creates a new geocoding client using the SDK
func NewClient(httpClient *http.Client) *Client {
  // Initialize SDK client
  // Note: Adapt to actual SDK API
  sdkClient := sdk.NewClient(httpClient)
  
  return &Client{
    sdkClient: sdkClient,
  }
}

// Search implements GeocodingService.Search using the SDK
func (c *Client) Search(ctx context.Context, name string) ([]models.Location, error) {
  // Call SDK with user-friendly error handling
  // Return []models.Location mapped from SDK response
  // Errors must be user-friendly (no technical details)
}

// mapSDKLocation converts SDK location to internal Location model
func mapSDKLocation(sdkLoc *sdk.Location) models.Location {
  return models.Location{
    ID:        sdkLoc.ID,
    Name:      sdkLoc.Name,
    Latitude:  sdkLoc.Latitude,
    Longitude: sdkLoc.Longitude,
    Country:   sdkLoc.Country,
    Region:    sdkLoc.Admin1,
  }
}

// convertError converts SDK errors to user-friendly messages
func convertError(err error) error {
  if err == nil {
    return nil
  }
  // Convert any SDK error to user-friendly message
  return errors.New("Unable to search locations. Please try again.")
}
```

### Implementation Checklist

- [ ] Import SDK package
- [ ] Create Client struct wrapping SDK client
- [ ] Implement NewClient() constructor
- [ ] Implement Search() method
- [ ] Add mapSDKLocation() conversion function
- [ ] Add convertError() error mapping function
- [ ] Preserve context handling (cancellation, timeouts)
- [ ] Handle all error scenarios with user-friendly messages

## Step 3: Add Integration Test

### File: `src/internal/geo/integration_test.go` (new file)

Create a short integration test that validates the API contract:

```go
package geo

import (
  "context"
  "net/http"
  "testing"
  "time"

  "github.com/stretchr/testify/assert"
  "weather-reporter/src/internal/models"
)

// TestIntegration_GeocodingAPIContract verifies the Open Meteo Geocoding API
// response structure hasn't changed in ways that would break our integration.
//
// This test makes REAL API calls and should be skipped in offline environments:
//   go test -short  (skips integration tests)
func TestIntegration_GeocodingAPIContract(t *testing.T) {
  if testing.Short() {
    t.Skip("Skipping integration test in short mode")
  }

  // Create a client with a real HTTP client
  client := NewClient(&http.Client{
    Timeout: 10 * time.Second,
  })

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  // Test 1: Search for a well-known location
  t.Run("SearchForLondon", func(t *testing.T) {
    results, err := client.Search(ctx, "London")
    
    assert.NoError(t, err, "Search should succeed for London")
    assert.Greater(t, len(results), 0, "Should return at least one result")
    
    // Verify first result has all required fields
    london := results[0]
    assert.Greater(t, london.ID, 0, "ID must be positive")
    assert.NotEmpty(t, london.Name, "Name must not be empty")
    assert.NotEmpty(t, london.Country, "Country must not be empty")
    
    // Verify coordinates are valid
    assert.GreaterOrEqual(t, london.Latitude, -90.0, "Latitude must be >= -90")
    assert.LessOrEqual(t, london.Latitude, 90.0, "Latitude must be <= 90")
    assert.GreaterOrEqual(t, london.Longitude, -180.0, "Longitude must be >= -180")
    assert.LessOrEqual(t, london.Longitude, 180.0, "Longitude must be <= 180")
  })

  // Test 2: Ambiguous query returns multiple results
  t.Run("AmbiguousQueryReturnsMultiple", func(t *testing.T) {
    results, err := client.Search(ctx, "Springfield")
    
    assert.NoError(t, err, "Search should succeed for Springfield")
    assert.Greater(t, len(results), 1, "Should return multiple Springfield results")
  })

  // Test 3: Results don't exceed 10
  t.Run("ResultsRespectMaxLimit", func(t *testing.T) {
    results, err := client.Search(ctx, "New")
    
    assert.NoError(t, err, "Search should succeed")
    assert.LessOrEqual(t, len(results), 10, "Should not exceed 10 results")
  })

  // Test 4: All results have complete Location data
  t.Run("AllResultsHaveCompleteData", func(t *testing.T) {
    results, err := client.Search(ctx, "Paris")
    
    assert.NoError(t, err, "Search should succeed")
    assert.Greater(t, len(results), 0, "Should return results")
    
    for i, loc := range results {
      assert.Greater(t, loc.ID, 0, "Result %d: ID must be positive", i)
      assert.NotEmpty(t, loc.Name, "Result %d: Name must not be empty", i)
      assert.NotEmpty(t, loc.Country, "Result %d: Country must not be empty", i)
      assert.NotZero(t, loc.Latitude, "Result %d: Latitude must be set", i)
      assert.NotZero(t, loc.Longitude, "Result %d: Longitude must be set", i)
    }
  })
}

// TestIntegration_APIContractChange detects if the Open Meteo API has changed
// in breaking ways. This test documents the expected API response structure.
func TestIntegration_APIContractChange(t *testing.T) {
  if testing.Short() {
    t.Skip("Skipping integration test in short mode")
  }

  client := NewClient(nil) // Use default HTTP client

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  // Search for a well-known, stable location
  results, err := client.Search(ctx, "London")

  if err != nil {
    t.Fatalf("API call failed: %v\nThis might indicate API breaking changes", err)
  }

  if len(results) == 0 {
    t.Fatalf("No results for London search\nThis indicates API breaking changes")
  }

  // Verify the expected Location structure is present
  london := results[0]
  
  // Expected: London, United Kingdom
  if london.Name != "London" && !assert.Contains(t, london.Name, "London") {
    t.Logf("WARNING: Expected 'London' in result name, got: %s", london.Name)
    t.Logf("API response structure may have changed")
  }

  if london.Country != "United Kingdom" {
    t.Logf("WARNING: Expected 'United Kingdom', got: %s", london.Country)
    t.Logf("API response structure may have changed")
  }

  t.Logf("✓ API contract verified for location search")
  t.Logf("  - Location: %s, %s", london.Name, london.Country)
  t.Logf("  - Coordinates: (%.2f, %.2f)", london.Latitude, london.Longitude)
}
```

### Test Verification

Run the integration tests:

```bash
# Run all tests including integration tests
go test ./src/internal/geo -v

# Run only unit tests (skip integration tests)
go test ./src/internal/geo -v -short

# Run only integration tests
go test ./src/internal/geo -v -run Integration
```

Expected output:
```
TestIntegration_GeocodingAPIContract/SearchForLondon PASS
TestIntegration_GeocodingAPIContract/AmbiguousQueryReturnsMultiple PASS
TestIntegration_GeocodingAPIContract/ResultsRespectMaxLimit PASS
TestIntegration_GeocodingAPIContract/AllResultsHaveCompleteData PASS
TestIntegration_APIContractChange PASS
```

## Step 4: Update Dependencies

### File: `go.mod`

Add the SDK dependency:

```bash
go get github.com/gregbalnis/open-meteo-geocoding-sdk@v0.1.0
```

This will update `go.mod` to include:
```
require github.com/gregbalnis/open-meteo-geocoding-sdk v0.1.0
```

Verify:
```bash
go mod tidy
```

## Step 5: Verify Existing Tests Pass

### Run Unit Tests

Ensure all existing tests pass without modification:

```bash
go test ./src/internal/geo -v
```

Expected: All tests from `client_test.go` pass

### Expected Test Cases

From existing `client_test.go`:
- ✅ Success (returns location)
- ✅ No Results (returns empty slice)
- ✅ API Error (returns error)
- ✅ Malformed JSON (returns error)
- ✅ Timeout (returns error)

**Important**: These tests should NOT be modified. The adapter should make them pass with the new SDK implementation.

## Step 6: Build and Manual Test

### Build the Application

```bash
go build -o bin/weather-reporter ./src/cmd/weather-reporter
```

### Test Location Search

```bash
./bin/weather-reporter London
```

Expected output:
```
Location: London, United Kingdom
Temperature: ...
```

### Test Multiple Queries

```bash
./bin/weather-reporter "San Francisco"
./bin/weather-reporter "Sydney"
./bin/weather-reporter "Tokyo"
```

All should work identically to before the refactoring.

### Test Error Handling

Test that errors are user-friendly:

1. **Offline test** (disable network access):
   ```bash
   # Output should show: "Unable to search locations. Please try again."
   ```

2. **Timeout test** (set very short timeout):
   - Output should show: "Unable to search locations. Please try again."

## Step 7: Performance Verification

### Compare Response Times

Measure that there's no performance degradation:

```bash
# Before refactoring
time ./bin/weather-reporter London

# After refactoring  
time ./bin/weather-reporter London
```

Both should take approximately the same time (typically <1 second for well-known locations).

## Troubleshooting

### Import Error: SDK not found

**Problem**: `cannot find package github.com/gregbalnis/open-meteo-geocoding-sdk`

**Solution**:
```bash
go get github.com/gregbalnis/open-meteo-geocoding-sdk@latest
go mod tidy
```

### Integration Test Fails: Network Timeout

**Problem**: Integration tests fail with timeout errors

**Solution**: 
- Verify internet connectivity
- Use `-short` flag to skip integration tests if offline
- Increase timeout in test if network is slow

### Existing Tests Fail After Refactoring

**Problem**: `client_test.go` tests fail with new implementation

**Solution**:
- Verify SDK API matches expected behavior
- Check that type mapping is correct
- Verify error handling returns correct format
- DO NOT modify existing tests - fix implementation instead

## Success Criteria Checklist

- [ ] SDK dependency added to go.mod
- [ ] `client.go` refactored to use SDK
- [ ] Integration test created and passing
- [ ] All existing unit tests pass without modification
- [ ] Application builds successfully
- [ ] Manual testing works for various locations
- [ ] Error messages are user-friendly
- [ ] Performance is equal or better
- [ ] No regressions in functionality

## Next Steps

Once verification is complete:

1. Commit changes to feature branch
2. Create pull request for code review
3. Verify CI/CD pipeline passes
4. Merge to main branch
5. Deploy to production
