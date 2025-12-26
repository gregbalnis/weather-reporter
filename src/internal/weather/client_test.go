package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"weather-reporter/src/internal/models"
)

func TestGetCurrentWeather(t *testing.T) {
	tests := []struct {
		name           string
		lat            float64
		lon            float64
		mockResponse   string
		mockStatusCode int
		expected       *models.WeatherResponse
		expectError    bool
	}{
		{
			name:           "Success",
			lat:            52.52,
			lon:            13.41,
			mockStatusCode: http.StatusOK,
			mockResponse: `{
				"current": {
					"temperature_2m": 20.5,
					"relative_humidity_2m": 60,
					"apparent_temperature": 19.5,
					"precipitation": 0.0,
					"cloud_cover": 20,
					"surface_pressure": 1015.5,
					"wind_speed_10m": 15.0,
					"wind_direction_10m": 180,
					"wind_gusts_10m": 25.0
				},
				"current_units": {
					"temperature_2m": "°C",
					"relative_humidity_2m": "%",
					"apparent_temperature": "°C",
					"precipitation": "mm",
					"cloud_cover": "%",
					"surface_pressure": "hPa",
					"wind_speed_10m": "km/h",
					"wind_direction_10m": "°",
					"wind_gusts_10m": "km/h"
				}
			}`,
			expected: &models.WeatherResponse{
				Current: models.Weather{
					Temperature:   20.5,
					Humidity:      60,
					ApparentTemp:  19.5,
					Precipitation: 0.0,
					CloudCover:    20,
					Pressure:      1015.5,
					WindSpeed:     15.0,
					WindDirection: 180,
					WindGusts:     25.0,
				},
				CurrentUnits: models.WeatherUnits{
					Temperature:   "°C",
					Humidity:      "%",
					ApparentTemp:  "°C",
					Precipitation: "mm",
					CloudCover:    "%",
					Pressure:      "hPa",
					WindSpeed:     "km/h",
					WindDirection: "°",
					WindGusts:     "km/h",
				},
			},
			expectError: false,
		},
		{
			name:           "API Error",
			lat:            0,
			lon:            0,
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   ``,
			expected:       nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/forecast", r.URL.Path)
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := NewClient(server.Client())
			client.baseURL = server.URL

			result, err := client.GetCurrentWeather(context.Background(), tt.lat, tt.lon)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
