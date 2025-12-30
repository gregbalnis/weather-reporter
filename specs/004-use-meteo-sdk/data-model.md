# Data Model: Integrate Open-Meteo SDK

**Phase**: 1 - Design  
**Date**: 2025-12-30  
**Plan**: [plan.md](plan.md)

## Overview

This document defines how SDK response types map to application data structures. The integration maintains existing model contracts while adapting to SDK-provided types.

## Entities

### WeatherClient (Interface)

**Purpose**: Abstract weather data retrieval to allow testing and SDK encapsulation

**Attributes**:
- None (behavior-only interface)

**Operations**:
- `GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error)`

**Relationships**:
- Implemented by SDK client or adapter
- Used by main application flow

**Lifecycle**:
- Created once during application initialization
- Lives for duration of single application run (CLI process)

**Validation Rules**:
- Context must not be nil
- Latitude: -90 to +90
- Longitude: -180 to +180

**State Transitions**: N/A (stateless)

---

### Location

**Purpose**: Represents a geographic location from geocoding search

**Attributes**:
- `ID`: int - Unique identifier from geocoding API
- `Name`: string - Location name (e.g., "London")
- `Latitude`: float64 - Latitude coordinate
- `Longitude`: float64 - Longitude coordinate
- `Country`: string - Country code (e.g., "GB")
- `Region`: string - Administrative region (e.g., "England")

**Operations**: None (data structure only)

**Relationships**:
- Used as input to WeatherClient
- Selected from geocoding results (geo.Client)

**Lifecycle**: Created from API response, used immediately, discarded after weather fetch

**Validation Rules**:
- Latitude: -90 to +90
- Longitude: -180 to +180
- Name, Country, Region: non-empty strings

**State Transitions**: N/A (immutable)

---

### Weather

**Purpose**: Contains current weather measurements

**Attributes**:
- `Temperature`: float64 - Air temperature in °C
- `Humidity`: int - Relative humidity in %
- `ApparentTemp`: float64 - Feels-like temperature in °C
- `Precipitation`: float64 - Precipitation in mm
- `CloudCover`: int - Cloud coverage in %
- `Pressure`: float64 - Surface pressure in hPa
- `WindSpeed`: float64 - Wind speed in km/h
- `WindDirection`: int - Wind direction in degrees
- `WindGusts`: float64 - Wind gusts in km/h

**Operations**: None (data structure only)

**Relationships**:
- Part of WeatherResponse
- Mapped from SDK response

**Lifecycle**: Created from SDK response, used for display, discarded

**Validation Rules**:
- All numeric fields must be finite (not NaN or Inf)
- Percentages (Humidity, CloudCover): 0-100
- Wind direction: 0-359 degrees

**State Transitions**: N/A (immutable)

---

### WeatherUnits

**Purpose**: Describes units for weather measurements (for display)

**Attributes**:
- `Temperature`: string (e.g., "°C")
- `Humidity`: string (e.g., "%")
- `ApparentTemp`: string (e.g., "°C")
- `Precipitation`: string (e.g., "mm")
- `CloudCover`: string (e.g., "%")
- `Pressure`: string (e.g., "hPa")
- `WindSpeed`: string (e.g., "km/h")
- `WindDirection`: string (e.g., "°")
- `WindGusts`: string (e.g., "km/h")

**Operations**: None (data structure only)

**Relationships**:
- Part of WeatherResponse
- Used by printWeather for formatting

**Lifecycle**: Created from SDK response or constants, used for display

**Validation Rules**: Strings must not be empty

**State Transitions**: N/A (immutable)

---

### WeatherResponse

**Purpose**: Complete weather data response combining measurements and units

**Attributes**:
- `Current`: Weather - Current weather measurements
- `CurrentUnits`: WeatherUnits - Units for measurements

**Operations**: None (data structure only)

**Relationships**:
- Returned by WeatherClient.GetCurrentWeather
- Consumed by main.go printWeather function

**Lifecycle**: Created by SDK/adapter, passed to display, discarded

**Validation Rules**:
- Current must not be nil/zero
- CurrentUnits must not be nil/zero

**State Transitions**: N/A (immutable)

---

## Data Flow

```
User Input (Location Name)
    ↓
geo.Client.Search() → []Location
    ↓
User Selection → Location
    ↓
[SDK Adapter/Client]
    ↓
WeatherClient.GetCurrentWeather(lat, lon)
    ↓
SDK Call → Open-Meteo API
    ↓
SDK Response → WeatherResponse (mapped if needed)
    ↓
printWeather() → stdout
```

## SDK Integration Points

### Current Implementation (To Be Replaced)
```go
// internal/weather/client.go
type Client struct {
    httpClient *http.Client
    baseURL    string
}

func (c *Client) GetCurrentWeather(ctx, lat, lon) (*models.WeatherResponse, error) {
    // Custom HTTP request logic
    // JSON decoding
    // Error handling
}
```

### New Implementation (SDK-Based)
```go
// main.go or internal/weather/adapter.go (if needed)
import sdk "github.com/gregbalnis/open-meteo-weather-sdk"

// Option 1: Direct SDK usage (if types match)
weatherClient := sdk.NewClient(&http.Client{Timeout: 10 * time.Second})
weatherData, err := weatherClient.GetCurrentWeather(ctx, lat, lon)

// Option 2: Adapter pattern (if types differ)
type SDKAdapter struct {
    client *sdk.Client
}

func (a *SDKAdapter) GetCurrentWeather(ctx context.Context, lat, lon float64) (*models.WeatherResponse, error) {
    sdkResp, err := a.client.GetWeather(ctx, lat, lon)
    if err != nil {
        return nil, err
    }
    return mapSDKResponse(sdkResp), nil
}

func mapSDKResponse(sdkResp *sdk.Response) *models.WeatherResponse {
    // Transform SDK types to models types
    return &models.WeatherResponse{
        Current: models.Weather{
            Temperature: sdkResp.Current.Temperature,
            // ... etc
        },
        CurrentUnits: models.WeatherUnits{
            Temperature: sdkResp.Units.Temperature,
            // ... etc
        },
    }
}
```

## Model Changes Required

### Minimal Changes Scenario
If SDK types are compatible:
- ✅ Keep all existing models
- ✅ Remove `internal/weather/client.go`
- ✅ Update `main.go` to use SDK directly
- ✅ No model changes needed

### Adapter Scenario
If SDK types differ:
- ✅ Keep `models.Weather`, `models.WeatherUnits`, `models.WeatherResponse`
- ✅ Add `mapSDKResponse()` function in models package
- ✅ Remove `internal/weather/client.go`
- ✅ Update `main.go` to use adapter

### Interface Extraction (for testing)
```go
// internal/models/interfaces.go
type WeatherClient interface {
    GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error)
}

// Implemented by:
// - SDK client directly (if compatible)
// - SDKAdapter (if transformation needed)
// - MockWeatherClient (for tests)
```

## Assumptions

1. SDK response structure follows Open-Meteo API patterns (JSON with current weather object)
2. SDK supports custom `http.Client` injection for timeout control
3. SDK returns metric units by default (verified by user note)
4. All weather parameters we display are available in SDK response
5. SDK error types are compatible with Go standard errors

## Migration Strategy

1. **Phase 1a**: Add SDK dependency (`go get`)
2. **Phase 1b**: Examine actual SDK types and interfaces
3. **Phase 1c**: Choose integration approach (direct vs adapter)
4. **Phase 1d**: Implement chosen approach
5. **Phase 1e**: Update tests
6. **Phase 1f**: Delete `internal/weather/client.go`

**Rollback Plan**: If SDK integration fails, revert commits and re-evaluate SDK compatibility
