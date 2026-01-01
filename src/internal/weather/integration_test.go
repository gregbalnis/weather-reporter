package weather_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"weather-reporter/src/internal/weather"
)

func TestClient_GetCurrentWeather_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a client with a real HTTP client
	client := weather.NewClient(http.DefaultClient)

	// Use a known location (London)
	lat := 51.5074
	lon := -0.1278

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		t.Fatalf("Failed to get weather: %v", err)
	}

	if resp == nil {
		t.Fatal("Response is nil")
	}

	// Verify that we get formatted strings back
	if resp.QuantityOfTemperature() == "" {
		t.Error("QuantityOfTemperature is empty")
	}
	if resp.QuantityOfWindSpeed() == "" {
		t.Error("QuantityOfWindSpeed is empty")
	}
}

func TestClient_GetCurrentWeather_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a client with a very short timeout
	httpClient := &http.Client{
		Timeout: 1 * time.Millisecond,
	}
	client := weather.NewClient(httpClient)

	lat := 51.5074
	lon := -0.1278

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.GetCurrentWeather(ctx, lat, lon)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
