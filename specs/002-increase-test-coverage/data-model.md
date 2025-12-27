# Data Model & Mocks: Increase Unit Test Coverage

**Feature**: Increase Unit Test Coverage
**Date**: 2025-12-27

## Mocks

To enable isolated unit testing of the `ui` package, we will implement mocks for the core service interfaces defined in `src/internal/models/interfaces.go`.

### `MockGeocodingService`

Implements `models.GeocodingService`.

```go
type MockGeocodingService struct {
    mock.Mock
}

func (m *MockGeocodingService) Search(ctx context.Context, name string) ([]models.Location, error) {
    args := m.Called(ctx, name)
    return args.Get(0).([]models.Location), args.Error(1)
}
```

### `MockWeatherService`

Implements `models.WeatherService`.

```go
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
```

## Test Data

### Sample Location
```json
{
  "id": 1,
  "name": "London",
  "latitude": 51.50853,
  "longitude": -0.12574,
  "country": "United Kingdom",
  "admin1": "London"
}
```

### Sample Weather Response
```json
{
  "current_weather": {
    "temperature": 15.5,
    "windspeed": 10.2,
    "winddirection": 180,
    "weathercode": 3,
    "time": "2025-12-27T10:00"
  },
  "current_weather_units": {
    "temperature": "°C",
    "windspeed": "km/h",
    "winddirection": "°",
    "weathercode": "wmo code",
    "time": "iso8601"
  }
}
```
