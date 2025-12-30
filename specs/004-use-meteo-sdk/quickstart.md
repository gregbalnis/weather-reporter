# Quickstart: Integrate Open-Meteo SDK

**Phase**: 1 - Design  
**Date**: 2025-12-30  
**Audience**: Developers implementing this feature

## Overview

This guide provides step-by-step instructions for integrating the open-meteo-weather-sdk to replace the custom weather client implementation.

## Prerequisites

- Go 1.25.5 installed
- Git access to the weather-reporter repository
- Branch `004-use-meteo-sdk` checked out
- Familiarity with Go modules and testing

## Implementation Steps

### Step 1: Add SDK Dependency

Add the SDK to project dependencies:

```bash
cd /home/gbalnis/projects/weather-reporter
go get github.com/gregbalnis/open-meteo-weather-sdk
go mod tidy
```

**Verification**:
```bash
grep "open-meteo-weather-sdk" go.mod
# Should show: github.com/gregbalnis/open-meteo-weather-sdk v<version>
```

---

### Step 2: Examine SDK Interface

Before modifying code, understand the SDK's API:

```bash
# View SDK godoc locally
go doc github.com/gregbalnis/open-meteo-weather-sdk
```

**Key Questions to Answer**:
1. What is the client initialization function? (e.g., `NewClient()`)
2. Does it accept custom `http.Client`?
3. What is the weather fetch method signature?
4. What types does it return?
5. How are errors reported?

**Document findings** in research.md if they differ from assumptions.

---

### Step 3: Create SDK Adapter (If Needed)

**Option A**: SDK types match our models → Skip to Step 4

**Option B**: SDK types differ → Create adapter in `internal/models/`:

```go
// internal/models/sdk_adapter.go
package models

import (
    "context"
    "fmt"
    sdk "github.com/gregbalnis/open-meteo-weather-sdk"
)

// SDKWeatherClient adapts the SDK to our WeatherClient interface
type SDKWeatherClient struct {
    client *sdk.Client
}

// NewSDKWeatherClient creates a weather client using the SDK
func NewSDKWeatherClient(httpClient *http.Client) *SDKWeatherClient {
    return &SDKWeatherClient{
        client: sdk.NewClient(httpClient),
    }
}

// GetCurrentWeather implements the WeatherClient interface
func (s *SDKWeatherClient) GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error) {
    resp, err := s.client.FetchCurrentWeather(ctx, lat, lon)
    if err != nil {
        return nil, fmt.Errorf("SDK weather fetch failed: %w", err)
    }
    
    // Map SDK response to our models
    return &WeatherResponse{
        Current: Weather{
            Temperature:   resp.Current.Temperature,
            Humidity:      resp.Current.Humidity,
            ApparentTemp:  resp.Current.ApparentTemperature,
            Precipitation: resp.Current.Precipitation,
            CloudCover:    resp.Current.CloudCover,
            Pressure:      resp.Current.Pressure,
            WindSpeed:     resp.Current.WindSpeed,
            WindDirection: resp.Current.WindDirection,
            WindGusts:     resp.Current.WindGusts,
        },
        CurrentUnits: WeatherUnits{
            Temperature:   resp.Units.Temperature,
            Humidity:      resp.Units.Humidity,
            ApparentTemp:  resp.Units.ApparentTemperature,
            Precipitation: resp.Units.Precipitation,
            CloudCover:    resp.Units.CloudCover,
            Pressure:      resp.Units.Pressure,
            WindSpeed:     resp.Units.WindSpeed,
            WindDirection: resp.Units.WindDirection,
            WindGusts:     resp.Units.WindGusts,
        },
    }, nil
}
```

---

### Step 4: Update main.go

Replace custom weather client with SDK:

**Before** (using internal/weather):
```go
import "weather-reporter/src/internal/weather"

weatherClient := weather.NewClient(nil)
```

**After** (using SDK directly or via adapter):
```go
import (
    "net/http"
    "time"
    "weather-reporter/src/internal/models"
)

// Create HTTP client with 10-second timeout
httpClient := &http.Client{
    Timeout: 10 * time.Second,
}

// Use SDK via adapter
weatherClient := models.NewSDKWeatherClient(httpClient)
```

**No other changes needed** - the interface remains the same:
```go
weatherData, err := weatherClient.GetCurrentWeather(ctx, selectedLocation.Latitude, selectedLocation.Longitude)
```

---

### Step 5: Update or Remove Tests

#### 5a: Remove Old Client Tests

```bash
# Delete custom client tests
rm src/internal/weather/client_test.go

# Delete custom client implementation
rm src/internal/weather/client.go

# If weather package is now empty, remove it
rmdir src/internal/weather
```

#### 5b: Add SDK Integration Tests

Create `internal/models/sdk_adapter_test.go`:

```go
package models_test

import (
    "context"
    "net/http"
    "testing"
    "time"
    
    "weather-reporter/src/internal/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestSDKWeatherClient_GetCurrentWeather(t *testing.T) {
    // Integration test with real API (Open-Meteo is free and fast)
    httpClient := &http.Client{Timeout: 10 * time.Second}
    client := models.NewSDKWeatherClient(httpClient)
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Reykjavik coordinates
    lat, lon := 64.1466, -21.9426
    
    resp, err := client.GetCurrentWeather(ctx, lat, lon)
    require.NoError(t, err, "Should fetch weather successfully")
    require.NotNil(t, resp, "Response should not be nil")
    
    // Verify all fields are populated
    assert.NotZero(t, resp.Current.Temperature, "Temperature should be set")
    assert.NotZero(t, resp.Current.Humidity, "Humidity should be set")
    assert.NotEmpty(t, resp.CurrentUnits.Temperature, "Temperature unit should be set")
}

func TestSDKWeatherClient_Timeout(t *testing.T) {
    httpClient := &http.Client{Timeout: 10 * time.Second}
    client := models.NewSDKWeatherClient(httpClient)
    
    // Create context with very short timeout to force timeout error
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
    defer cancel()
    
    _, err := client.GetCurrentWeather(ctx, 64.1466, -21.9426)
    require.Error(t, err, "Should return timeout error")
}

func TestSDKWeatherClient_InvalidCoordinates(t *testing.T) {
    httpClient := &http.Client{Timeout: 10 * time.Second}
    client := models.NewSDKWeatherClient(httpClient)
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Invalid latitude
    _, err := client.GetCurrentWeather(ctx, 999, 0)
    assert.Error(t, err, "Should return error for invalid coordinates")
}
```

---

### Step 6: Run Tests

```bash
# Run all tests with coverage
make test

# Or directly:
go test -race -coverprofile=coverage.out ./src/...
go tool cover -func=coverage.out
```

**Success Criteria**:
- ✅ All tests pass
- ✅ Coverage remains ≥80%
- ✅ No race conditions detected

---

### Step 7: Manual Testing

Test the full application flow:

```bash
# Build application
make build

# Test with unique location
./bin/weather-reporter Reykjavik

# Test with ambiguous location (should prompt)
./bin/weather-reporter London

# Test with non-existent location
./bin/weather-reporter AtlantisUnderSea

# Test timeout behavior (if possible to simulate)
```

**Verify**:
- Output format matches previous implementation exactly
- Error messages are user-friendly
- Location selection works identically
- Timeout occurs within 10 seconds on slow networks

---

### Step 8: Update Documentation

If SDK installation requires any special steps, update README.md:

```bash
# Check if SDK has special requirements
cat go.mod | grep open-meteo-weather-sdk
```

**Add to README.md only if needed**:
```markdown
## Dependencies

This project uses the open-meteo-weather-sdk for weather data retrieval:
- Repository: https://github.com/gregbalnis/open-meteo-weather-sdk
- License: [Check SDK repository]
```

---

### Step 9: Verify Constitution Compliance

Run linter and ensure standards are met:

```bash
make lint

# Check specific items:
# ✅ Code formatted with gofmt
# ✅ Test coverage ≥80%
# ✅ Error messages are clear
# ✅ No new complexity introduced
```

---

### Step 10: Commit and Push

```bash
git add .
git commit -m "feat: integrate open-meteo-weather-sdk to replace custom weather client

- Add github.com/gregbalnis/open-meteo-weather-sdk dependency
- Remove internal/weather custom client implementation
- Add SDK adapter in internal/models for type mapping
- Update main.go to use SDK with 10s timeout
- Replace unit tests with SDK integration tests
- Maintain 100% functional compatibility with existing behavior

Reduces maintenance burden by ~150 LOC while preserving all functionality.
Closes #[issue-number]"

git push origin 004-use-meteo-sdk
```

---

## Troubleshooting

### Issue: SDK types don't match our models

**Solution**: Use adapter pattern (Step 3 Option B) to transform types.

### Issue: SDK doesn't accept custom http.Client

**Solution**: Use context-based timeout as primary control:
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
```

### Issue: Test coverage drops below 80%

**Solution**: 
1. Ensure integration tests cover all code paths
2. Add table-driven tests for error cases
3. Mock only external dependencies, not SDK itself

### Issue: SDK not found or build fails

**Solution**:
```bash
go clean -modcache
go mod download
go mod tidy
```

---

## Success Checklist

Before creating PR:

- [ ] SDK dependency added to go.mod
- [ ] internal/weather/client.go deleted
- [ ] internal/weather/client_test.go deleted
- [ ] SDK adapter created (if needed)
- [ ] main.go updated to use SDK
- [ ] Integration tests added and passing
- [ ] Test coverage ≥80%
- [ ] Linter passes (make lint)
- [ ] Manual testing confirms identical behavior
- [ ] README.md updated (if needed)
- [ ] All commits follow Conventional Commits format

---

## Expected Outcome

- **Code Reduction**: ~150 lines removed (custom client implementation)
- **Test Changes**: Unit tests → Integration tests (same coverage)
- **User Impact**: Zero (identical behavior)
- **Maintenance**: Reduced (delegate to SDK maintainers)

## Next Steps

After integration complete:
1. Create Pull Request
2. Request code review
3. Verify CI/CD passes
4. Merge to main
5. Monitor for any issues
