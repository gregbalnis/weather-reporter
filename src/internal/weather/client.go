// Package weather provides functionality for fetching weather data.
package weather

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"weather-reporter/src/internal/models"

	meteosdk "github.com/gregbalnis/open-meteo-weather-sdk"
)

// Client is a client for the weather API.
type Client struct {
	sdkClient *meteosdk.Client
}

// NewClient creates a new weather client.
// If httpClient is nil, a default client with a 10s timeout is used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return &Client{
		sdkClient: meteosdk.NewClient(meteosdk.WithHTTPClient(httpClient)),
	}
}

// GetCurrentWeather fetches the current weather for the given coordinates.
func (c *Client) GetCurrentWeather(ctx context.Context, lat, lon float64) (models.WeatherResponse, error) {
	resp, err := c.sdkClient.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("received nil response from SDK")
	}
	return &weatherResponseAdapter{resp}, nil
}

type weatherResponseAdapter struct {
	*meteosdk.CurrentWeather
}

// QuantityOfTemperature returns the temperature.
func (w *weatherResponseAdapter) QuantityOfTemperature() string {
	return w.CurrentWeather.QuantityOfTemperature()
}

// QuantityOfHumidity returns the humidity.
func (w *weatherResponseAdapter) QuantityOfHumidity() string {
	return w.QuantityOfRelativeHumidity()
}

// QuantityOfApparentTemperature returns the apparent temperature.
func (w *weatherResponseAdapter) QuantityOfApparentTemperature() string {
	return w.CurrentWeather.QuantityOfApparentTemperature()
}

// QuantityOfPrecipitation returns the precipitation.
func (w *weatherResponseAdapter) QuantityOfPrecipitation() string {
	return w.CurrentWeather.QuantityOfPrecipitation()
}

// QuantityOfCloudCover returns the cloud cover.
func (w *weatherResponseAdapter) QuantityOfCloudCover() string {
	return w.CurrentWeather.QuantityOfCloudCover()
}

// QuantityOfPressure returns the pressure.
func (w *weatherResponseAdapter) QuantityOfPressure() string {
	return w.QuantityOfSurfacePressure()
}

// QuantityOfWindSpeed returns the wind speed.
func (w *weatherResponseAdapter) QuantityOfWindSpeed() string {
	return w.CurrentWeather.QuantityOfWindSpeed()
}

// QuantityOfWindDirection returns the wind direction.
func (w *weatherResponseAdapter) QuantityOfWindDirection() string {
	return w.CurrentWeather.QuantityOfWindDirection()
}

// QuantityOfWindGusts returns the wind gusts.
func (w *weatherResponseAdapter) QuantityOfWindGusts() string {
	return w.CurrentWeather.QuantityOfWindGusts()
}
