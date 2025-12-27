// Package geo provides functionality for searching locations using a geocoding API.
package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"weather-reporter/src/internal/models"
)

const defaultBaseURL = "https://geocoding-api.open-meteo.com/v1"

// Client is a client for the geocoding API.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new geocoding client.
// If httpClient is nil, a default client with a 10s timeout is used.
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

type searchResponse struct {
	Results []models.Location `json:"results"`
}

// Search searches for locations by name.
func (c *Client) Search(ctx context.Context, name string) ([]models.Location, error) {
	u, err := url.Parse(c.baseURL + "/search")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("name", name)
	q.Set("count", "10")
	q.Set("language", "en")
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var searchResp searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return searchResp.Results, nil
}
