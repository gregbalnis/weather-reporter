// Package contracts defines the interface contract for weather data retrieval
// This file serves as documentation for SDK integration requirements
package contracts

import "context"

// WeatherClient defines the contract that any weather data source must implement.
type WeatherClient interface {
	// GetCurrentWeather retrieves current weather data for the specified coordinates.
	// Returns a response object that provides formatted quantity accessors.
	GetCurrentWeather(ctx context.Context, lat, lon float64) (WeatherResponse, error)
}

// WeatherResponse defines the interface for the weather data response.
// It provides accessors that return formatted strings (value + unit).
type WeatherResponse interface {
	QuantityOfTemperature() string         // e.g., "10.5°C"
	QuantityOfHumidity() string            // e.g., "85%"
	QuantityOfApparentTemperature() string // e.g., "8.0°C"
	QuantityOfPrecipitation() string       // e.g., "0.0 mm"
	QuantityOfCloudCover() string          // e.g., "20%"
	QuantityOfPressure() string            // e.g., "1015 hPa"
	QuantityOfWindSpeed() string           // e.g., "15.0 km/h"
	QuantityOfWindDirection() string       // e.g., "180°"
	QuantityOfWindGusts() string           // e.g., "25.0 km/h"
}

// SDK Integration Requirements:
//
// 1. The SDK implementation MUST:
//    - Accept custom http.Client with 10-second timeout
//    - Return data in metric units
//    - Support context-based cancellation
//    - Provide QuantityOf... accessors for all 9 weather parameters
//
// 2. Error Handling Contract:
//    - Fail fast (no automatic retries)
//    - Propagate timeout errors clearly
//
// 3. Testing Requirements:
//    - Mock implementations MUST implement the WeatherResponse interface
//    - Integration tests MUST verify all accessors return non-empty strings
