// Package geo provides functionality for searching locations using a geocoding API.
package geo

import (
	"context"
	"errors"
	"net/http"
	"time"

	"weather-reporter/src/internal/models"

	geocoding "github.com/gregbalnis/open-meteo-geocoding-sdk"
)

const defaultBaseURL = "https://geocoding-api.open-meteo.com/v1"

// Client is a geocoding client that implements models.GeocodingService
// using the open-meteo-geocoding-sdk library. It retains baseURL/httpClient
// fields for compatibility with existing tests while delegating requests to
// the SDK client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	sdkClient  *geocoding.Client
}

// NewClient creates a new geocoding client using the open-meteo-geocoding-sdk.
// If httpClient is nil, a default client with a 10s timeout is used.
//
// The client is configured to:
//   - Return up to 10 location results per search
//   - Use English language for location names
//   - Respect context cancellation and timeouts
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	baseURL := defaultBaseURL
	// SDK expects the full endpoint including /search
	sdkBase := baseURL + "/search"

	var opts []geocoding.Option
	opts = append(opts, geocoding.WithHTTPClient(httpClient))
	opts = append(opts, geocoding.WithBaseURL(sdkBase))

	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		sdkClient:  geocoding.NewClient(opts...),
	}
}

// Search searches for locations by name using the SDK.
// It returns up to 10 matching locations.
//
// All errors are converted to user-friendly messages without technical details:
//   - Timeout errors: "Search took too long. Please try again."
//   - All other errors: "Unable to search locations. Please try again."
func (c *Client) Search(ctx context.Context, name string) ([]models.Location, error) {
	// Configure search options
	opts := &geocoding.SearchOptions{
		Count:    10,
		Language: "en",
	}

	// Call SDK
	sdkClient := c.sdkClient
	if c.baseURL != defaultBaseURL {
		// Rebuild SDK client with overridden base URL for tests
		sdkClient = geocoding.NewClient(
			geocoding.WithHTTPClient(c.httpClient),
			geocoding.WithBaseURL(c.baseURL+"/search"),
		)
	}

	sdkLocations, err := sdkClient.Search(ctx, name, opts)
	if err != nil {
		return nil, convertSDKError(err)
	}

	// Map SDK locations to internal model
	locations := make([]models.Location, len(sdkLocations))
	for i, sdkLoc := range sdkLocations {
		locations[i] = mapSDKLocation(sdkLoc)
	}

	return locations, nil
}

// mapSDKLocation converts an SDK location to our internal Location model.
// Note: The SDK does not provide an admin1/region field, so Region will be empty.
func mapSDKLocation(sdkLocation geocoding.Location) models.Location {
	return models.Location{
		ID:        sdkLocation.ID,
		Name:      sdkLocation.Name,
		Latitude:  sdkLocation.Latitude,
		Longitude: sdkLocation.Longitude,
		Country:   sdkLocation.Country,
		Region:    "", // SDK doesn't provide admin1/region field
	}
}

// convertSDKError converts SDK errors to user-friendly error messages.
// All technical details are hidden from users.
func convertSDKError(err error) error {
	if err == nil {
		return nil
	}

	// Context-based errors
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New("Search took too long. Please try again.")
	}
	if errors.Is(err, context.Canceled) {
		return errors.New("Unable to search locations. Please try again.")
	}

	// SDK-defined errors
	if errors.Is(err, geocoding.ErrConcurrencyLimitExceeded) {
		return errors.New("Unable to search locations. Please try again.")
	}
	if errors.Is(err, geocoding.ErrInvalidParameter) {
		return errors.New("Unable to search locations. Please try again.")
	}
	var apiErr *geocoding.APIError
	if errors.As(err, &apiErr) {
		return errors.New("Unable to search locations. Please try again.")
	}

	// Default fallback
	return errors.New("Unable to search locations. Please try again.")
}
