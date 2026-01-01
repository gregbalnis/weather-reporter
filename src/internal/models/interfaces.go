// Package models defines the data structures and interfaces for the application.
package models

import "context"

// GeocodingService defines the interface for finding locations.
type GeocodingService interface {
	// Search finds locations matching the given name.
	Search(ctx context.Context, name string) ([]Location, error)
}

// WeatherService defines the interface for fetching weather.
type WeatherService interface {
	// GetCurrentWeather returns the current weather for the given coordinates.
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
