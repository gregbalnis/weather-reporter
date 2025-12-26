package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"weather-reporter/src/internal/models"
)

const defaultBaseURL = "https://api.open-meteo.com/v1"

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    defaultBaseURL,
	}
}

func (c *Client) GetCurrentWeather(ctx context.Context, lat, lon float64) (*models.WeatherResponse, error) {
	u, err := url.Parse(c.baseURL + "/forecast")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("latitude", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("longitude", strconv.FormatFloat(lon, 'f', -1, 64))

	currentVars := []string{
		"temperature_2m",
		"relative_humidity_2m",
		"apparent_temperature",
		"precipitation",
		"cloud_cover",
		"surface_pressure",
		"wind_speed_10m",
		"wind_direction_10m",
		"wind_gusts_10m",
	}
	q.Set("current", strings.Join(currentVars, ","))
	q.Set("wind_speed_unit", "kmh")
	q.Set("temperature_unit", "celsius")

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var weatherResp models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &weatherResp, nil
}
