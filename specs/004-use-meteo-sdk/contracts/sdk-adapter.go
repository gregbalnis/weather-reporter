// Package contracts defines the interface contract for weather data retrieval
// This file serves as documentation for SDK integration requirements
package contracts

import "context"

// WeatherClient defines the contract that any weather data source must implement.
// This interface abstracts the weather data provider, allowing for:
// - SDK-based implementation (production)
// - Mock implementation (testing)
// - Alternative implementations (future)
type WeatherClient interface {
	// GetCurrentWeather retrieves current weather data for the specified coordinates.
	//
	// Parameters:
	//   - ctx: Context for cancellation and timeout control (MUST not be nil)
	//   - lat: Latitude in decimal degrees (-90 to +90)
	//   - lon: Longitude in decimal degrees (-180 to +180)
	//
	// Returns:
	//   - WeatherResponse: Current weather data with measurements and units
	//   - error: Network errors, timeout errors, invalid coordinates, API failures
	//
	// Timeout Behavior:
	//   - MUST respect context deadline (10 seconds max from spec)
	//   - MUST return error on timeout (no retries)
	//
	// Error Handling:
	//   - Context cancellation: return context.Canceled
	//   - Timeout: return context.DeadlineExceeded or wrapped timeout error
	//   - Network errors: return error describing failure
	//   - Invalid coordinates: return validation error
	//   - API errors: return error with API message
	//
	// Example Usage:
	//   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//   defer cancel()
	//   resp, err := client.GetCurrentWeather(ctx, 64.1466, -21.9426) // Reykjavik
	//   if err != nil {
	//       return fmt.Errorf("failed to fetch weather: %w", err)
	//   }
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error)
}

// WeatherResponse represents the complete weather data response.
// All fields use metric units (Celsius, km/h, hPa, mm).
type WeatherResponse struct {
	Current      Weather      // Current weather measurements
	CurrentUnits WeatherUnits // Units for display formatting
}

// Weather contains all current weather measurements in metric units.
type Weather struct {
	Temperature   float64 // Air temperature (°C)
	Humidity      int     // Relative humidity (%)
	ApparentTemp  float64 // Feels-like temperature (°C)
	Precipitation float64 // Precipitation amount (mm)
	CloudCover    int     // Cloud coverage (%)
	Pressure      float64 // Surface pressure (hPa)
	WindSpeed     float64 // Wind speed (km/h)
	WindDirection int     // Wind direction (degrees, 0-359)
	WindGusts     float64 // Wind gusts (km/h)
}

// WeatherUnits describes the units for each weather measurement (for display).
type WeatherUnits struct {
	Temperature   string // e.g., "°C"
	Humidity      string // e.g., "%"
	ApparentTemp  string // e.g., "°C"
	Precipitation string // e.g., "mm"
	CloudCover    string // e.g., "%"
	Pressure      string // e.g., "hPa"
	WindSpeed     string // e.g., "km/h"
	WindDirection string // e.g., "°"
	WindGusts     string // e.g., "km/h"
}

// SDK Integration Requirements:
//
// 1. The SDK implementation MUST:
//    - Accept custom http.Client with 10-second timeout
//    - Return data in metric units (verified: SDK uses metric by default)
//    - Support context-based cancellation
//    - Provide all 9 weather parameters listed in Weather struct
//
// 2. Error Handling Contract:
//    - Fail fast (no automatic retries)
//    - Propagate timeout errors clearly
//    - Wrap API errors with context
//
// 3. Testing Requirements:
//    - Mock implementations MUST honor context timeout
//    - Integration tests MUST verify all weather fields populated
//    - Error cases MUST be tested (timeout, network failure, invalid coords)
//
// 4. Type Mapping:
//    - If SDK types match this contract: use directly
//    - If SDK types differ: create adapter function to transform
//    - Goal: minimize transformation complexity
//
// Example SDK Adapter (if needed):
//
//   type SDKAdapter struct {
//       client *openmeteo.Client
//   }
//
//   func (a *SDKAdapter) GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error) {
//       sdkResp, err := a.client.FetchWeather(ctx, lat, lon)
//       if err != nil {
//           return nil, fmt.Errorf("SDK error: %w", err)
//       }
//       return &WeatherResponse{
//           Current: Weather{
//               Temperature: sdkResp.Temperature,
//               // ... map remaining fields
//           },
//           CurrentUnits: WeatherUnits{
//               Temperature: "°C",
//               // ... set units
//           },
//       }, nil
//   }
