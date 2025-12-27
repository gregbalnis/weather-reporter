package ui

import (
	"context"
	"weather-reporter/src/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockGeocodingService is a mock implementation of models.GeocodingService
type MockGeocodingService struct {
	mock.Mock
}

func (m *MockGeocodingService) Search(ctx context.Context, name string) ([]models.Location, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]models.Location), args.Error(1)
}

// MockWeatherService is a mock implementation of models.WeatherService
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetCurrentWeather(ctx context.Context, lat, lon float64) (*models.WeatherResponse, error) {
	args := m.Called(ctx, lat, lon)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WeatherResponse), args.Error(1)
}
