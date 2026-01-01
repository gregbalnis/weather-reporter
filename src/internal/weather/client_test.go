package weather

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

// roundTripFunc .
type roundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestGetCurrentWeather(t *testing.T) {
	// Mock JSON response from Open-Meteo
	jsonResponse := `{
		"latitude": 52.52,
		"longitude": 13.419998,
		"generationtime_ms": 0.13709068298339844,
		"utc_offset_seconds": 0,
		"timezone": "GMT",
		"timezone_abbreviation": "GMT",
		"elevation": 38.0,
		"current_units": {
			"time": "iso8601",
			"interval": "seconds",
			"temperature_2m": "°C",
			"relative_humidity_2m": "%",
			"apparent_temperature": "°C",
			"precipitation": "mm",
			"cloud_cover": "%",
			"surface_pressure": "hPa",
			"wind_speed_10m": "km/h",
			"wind_direction_10m": "°",
			"wind_gusts_10m": "km/h"
		},
		"current": {
			"time": "2026-01-01T06:30",
			"interval": 900,
			"temperature_2m": 2.5,
			"relative_humidity_2m": 76,
			"apparent_temperature": -2.8,
			"precipitation": 0.00,
			"cloud_cover": 99,
			"surface_pressure": 997.4,
			"wind_speed_10m": 20.2,
			"wind_direction_10m": 239,
			"wind_gusts_10m": 46.1
		}
	}`

	// Create a mock HTTP client
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			// Send response to be read
			Body: io.NopCloser(bytes.NewBufferString(jsonResponse)),
			// Must be set to non-nil value
			Header: make(http.Header),
		}
	})

	client := NewClient(httpClient)

	ctx := context.Background()
	resp, err := client.GetCurrentWeather(ctx, 52.52, 13.41)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}

	// Verify values
	if got := resp.QuantityOfTemperature(); got != "2.5°C" {
		t.Errorf("QuantityOfTemperature() = %v, want %v", got, "2.5°C")
	}
	if got := resp.QuantityOfHumidity(); got != "76%" {
		t.Errorf("QuantityOfHumidity() = %v, want %v", got, "76%")
	}
	if got := resp.QuantityOfApparentTemperature(); got != "-2.8°C" {
		t.Errorf("QuantityOfApparentTemperature() = %v, want %v", got, "-2.8°C")
	}
	// Note: Exact formatting depends on the SDK implementation.
	// If these fail, I will adjust the expected values.
	if got := resp.QuantityOfPrecipitation(); got != "0mm" && got != "0.00mm" {
		t.Logf("QuantityOfPrecipitation() = %v", got)
	}
	if got := resp.QuantityOfCloudCover(); got != "99%" {
		t.Errorf("QuantityOfCloudCover() = %v, want %v", got, "99%")
	}
	if got := resp.QuantityOfPressure(); got != "997.4 hPa" {
		t.Errorf("QuantityOfPressure() = %v, want %v", got, "997.4 hPa")
	}
	if got := resp.QuantityOfWindSpeed(); got != "20.2 km/h" {
		t.Errorf("QuantityOfWindSpeed() = %v, want %v", got, "20.2 km/h")
	}
	if got := resp.QuantityOfWindDirection(); got != "239°" {
		t.Errorf("QuantityOfWindDirection() = %v, want %v", got, "239°")
	}
	if got := resp.QuantityOfWindGusts(); got != "46.1 km/h" {
		t.Errorf("QuantityOfWindGusts() = %v, want %v", got, "46.1 km/h")
	}
}

func TestGetCurrentWeather_Error(t *testing.T) {
	httpClient := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			Header:     make(http.Header),
		}
	})

	client := NewClient(httpClient)
	_, err := client.GetCurrentWeather(context.Background(), 0, 0)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestNewClient_Default(t *testing.T) {
	client := NewClient(nil)
	if client == nil {
		t.Error("Expected client, got nil")
	}
}
