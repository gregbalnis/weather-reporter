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
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error)
}
