# Data Model: Geocoding SDK Integration

**Phase**: 1 - Core Implementation  
**Date**: January 3, 2026  
**Status**: Reference Implementation

## Entity Definition

### Location

Represents a geographical location that can be used to fetch weather data.

**Purpose**: Store and transport location information from geocoding API results to weather lookup functionality.

**Schema**:
```go
type Location struct {
  ID        int     `json:"id"`        // Unique identifier from geocoding API
  Name      string  `json:"name"`      // City, town, or location name
  Latitude  float64 `json:"latitude"`  // Geographic latitude coordinate
  Longitude float64 `json:"longitude"` // Geographic longitude coordinate
  Country   string  `json:"country"`   // Country name
  Region    string  `json:"admin1"`    // Administrative region (state, province)
}
```

**Constraints**:
- `ID`: Must be positive integer from API (unique per location in API)
- `Name`: Required, non-empty string (typically 2-100 characters)
- `Latitude`: Must be between -90 and 90 degrees
- `Longitude`: Must be between -180 and 180 degrees
- `Country`: Required, non-empty string
- `Region`: Can be empty for regions without administrative subdivisions

**Validation Rules**:
- All fields except Region must be present and non-empty
- Coordinates must be valid geographic coordinates
- ID must be positive

**Usage**:
- Input to `weather.GetCurrentWeather(ctx, location.Latitude, location.Longitude)`
- Display to user via UI components
- Storage/caching if needed

**Relationships**:
- Created by: `GeocodingService.Search()`
- Consumed by: `WeatherService.GetCurrentWeather()`
- Displayed by: UI components

## Interface Definition

### GeocodingService

Public interface that the geocoding implementation must implement. **No changes** - preserved for backward compatibility.

```go
// GeocodingService defines the interface for finding locations
type GeocodingService interface {
  // Search finds locations matching the given name.
  // 
  // Parameters:
  //   ctx: Context for request cancellation and timeouts
  //   name: Location name to search for (e.g., "London", "San Francisco")
  //
  // Returns:
  //   []Location: List of matching locations (up to 10)
  //   error: Non-nil if search fails
  //
  // Errors return user-friendly messages:
  //   - "Unable to search locations. Please try again." for service/network errors
  //   - "Search took too long. Please try again." for timeouts
  Search(ctx context.Context, name string) ([]Location, error)
}
```

**Implementation Requirements**:
- Implement `Search` method exactly as specified
- Handle context cancellation (return error when ctx is canceled)
- Support timeouts via context
- Return user-friendly error messages (no technical details)
- Return up to 10 location matches
- Support English language searches

## SDK Type Mapping

The `open-meteo-geocoding-sdk` provides Location data that must be converted to our internal `Location` model.

### Mapping Function

```go
// mapSDKLocation converts SDK location response to internal Location model
func mapSDKLocation(sdkLoc *sdk.Location) (Location, error) {
  // Validate required fields from SDK response
  if sdkLoc == nil {
    return Location{}, fmt.Errorf("SDK location is nil")
  }
  
  if sdkLoc.Name == "" {
    return Location{}, fmt.Errorf("SDK location missing name")
  }
  
  if sdkLoc.Country == "" {
    return Location{}, fmt.Errorf("SDK location missing country")
  }
  
  // Map SDK fields to internal Location struct
  return Location{
    ID:        sdkLoc.ID,           // Direct mapping
    Name:      sdkLoc.Name,         // Direct mapping
    Latitude:  sdkLoc.Latitude,     // Direct mapping (float64)
    Longitude: sdkLoc.Longitude,    // Direct mapping (float64)
    Country:   sdkLoc.Country,      // Direct mapping
    Region:    sdkLoc.Admin1,       // SDK uses "Admin1" for administrative region
  }, nil
}
```

**Mapping Notes**:
- The SDK's `Admin1` field maps to our `Region` field (both represent administrative subdivisions)
- All numeric fields are `float64` in the SDK
- String fields match our model exactly
- The SDK returns an ID field we can use directly

## Error Handling

### User-Facing Error Messages

All errors from the geocoding service are presented to users with the following messages (no technical details):

| Scenario | User Message | Log Level |
|----------|--------------|-----------|
| Service unavailable | "Unable to search locations. Please try again." | Error |
| Network timeout | "Search took too long. Please try again." | Error |
| Malformed response | "Unable to search locations. Please try again." | Error |
| Invalid input | "Unable to search locations. Please try again." | Error |
| Network unreachable | "Unable to search locations. Please try again." | Error |

### Error Conversion Strategy

```go
// convertSDKError converts SDK error to user-friendly message
func convertSDKError(err error) error {
  if err == nil {
    return nil
  }
  
  // Check SDK error types and return generic user-friendly message
  // Examples:
  // if errors.As(err, &sdk.TimeoutError{}) {
  //   return errors.New("Search took too long. Please try again.")
  // }
  
  // Default: generic message for any error
  return errors.New("Unable to search locations. Please try again.")
}
```

**Rationale**: 
- Users don't need to see technical error details
- Consistent error messages across all failure types
- Security benefit: no information leakage about infrastructure
- Better user experience: clear retry guidance

## Data Flow

### Location Search Flow

```
User Input (location name)
    ↓
GeocodingService.Search(ctx, "London")
    ↓
SDK Client Search Request
    ↓
Open Meteo Geocoding API
    ↓
SDK Client Parse Response
    ↓
mapSDKLocation() - Convert to internal model
    ↓
[]Location (up to 10 results)
    ↓
Display to User / Pass to Weather Service
```

### Error Flow

```
Error in any step
    ↓
convertSDKError() - Convert to user-friendly message
    ↓
Display generic message to user
```

## State Transitions

The `Location` entity doesn't have state transitions - it's a simple data carrier. However, in the overall application flow:

```
[User searches]
    ↓
[Location results displayed]
    ↓
[User selects location]
    ↓
[Weather lookup performed using location coordinates]
    ↓
[Weather displayed]
```

## Validation & Constraints

### Required Fields Validation

When creating a Location from SDK response:
- ID must be > 0
- Name must be non-empty
- Country must be non-empty
- Latitude must be -90 to 90
- Longitude must be -180 to 180
- Region can be empty

### Query Parameter Constraints

When searching:
- Location name: 1-255 characters
- Max results: 10 (hardcoded)
- Language: "en" (hardcoded)

## Integration Points

### With GeocodingService Interface

`Location` is the primary output of the `GeocodingService.Search()` method.

```go
// Search returns []Location
locations, err := geoService.Search(ctx, "London")
```

### With Weather Service

The Location's coordinates are used as input to weather lookups:

```go
weather, err := weatherService.GetCurrentWeather(ctx, location.Latitude, location.Longitude)
```

### With User Interface

Location objects are displayed to the user:

```go
for _, location := range locations {
  fmt.Printf("%s, %s\n", location.Name, location.Country)
}
```

## Type Safety

All Location fields use concrete types (no `interface{}`):
- Numeric coordinates are `float64` (no conversion needed)
- Identifiers are `int` (matches JSON unmarshaling)
- Names are `string` (no special handling needed)

This ensures type safety and simple conversion from SDK types.

## Performance Considerations

- Location struct is small (~60 bytes) - efficient for slices
- No pointers needed - can be copied freely
- JSON marshaling/unmarshaling is standard Go
- No expensive allocations in mapping

## Summary of Changes vs Current Implementation

| Aspect | Current | New |
|--------|---------|-----|
| Source | Custom JSON unmarshaling | SDK provides types |
| Mapping | Direct JSON to struct | SDK response to struct |
| Error handling | Generic error wrapping | SDK error conversion |
| Type validation | None | Validate in mapping function |
| Location fields | Same 6 fields | Same 6 fields (unchanged) |
| Interface | GeocodingService | GeocodingService (unchanged) |

**Key Point**: The Location entity itself does NOT change. Only how it's populated from the SDK changes.
