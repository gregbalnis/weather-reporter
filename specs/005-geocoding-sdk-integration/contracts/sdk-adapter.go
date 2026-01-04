// Package geo provides functionality for searching locations using a geocoding API.
package geo

import (
	"context"
	"errors"
	"net/http"

	"github.com/gregbalnis/open-meteo-geocoding-sdk"
	"weather-reporter/src/internal/models"
)

// Client implements the models.GeocodingService interface using the
// open-meteo-geocoding-sdk library instead of custom HTTP implementation.
//
// This adapter wraps the SDK client and provides conversion from SDK
// response types to our internal Location model, with user-friendly
// error handling.
type Client struct {
	// sdkClient holds the underlying SDK client
	sdkClient *geocoding.Client
}

// NewClient creates a new geocoding client using the open-meteo-geocoding-sdk.
//
// If httpClient is nil, the SDK's default HTTP client with standard timeouts is used.
// The adapter is configured to support:
//   - Location searches with max 10 results
//   - English language responses
//   - 10-second request timeout (SDK default)
//
// Returns a GeocodingService interface implementation.
func NewClient(httpClient *http.Client) models.GeocodingService {
	var opts []geocoding.Option
	if httpClient != nil {
		opts = append(opts, geocoding.WithHTTPClient(httpClient))
	}
	
	return &Client{
		sdkClient: geocoding.NewClient(opts...),
	}
}

// Search finds locations matching the given name using the SDK.
//
// Implementation contract:
//   - Uses SDK's location search functionality
//   - Returns up to 10 matching locations
//   - Searches in English language
//   - Respects context cancellation and timeouts
//   - Converts SDK response types to internal Location model
//   - Returns user-friendly error messages (no technical details)
//
// All errors are converted to messages that provide helpful guidance to users
// without exposing implementation details:
//   - Service/Network errors: "Unable to search locations. Please try again."
//   - Timeout errors: "Search took too long. Please try again."
//   - Any other errors: "Unable to search locations. Please try again."
//
// Parameters:
//
//	ctx: Context for cancellation and timeouts
//	name: Location name to search for (e.g., "London", "New York")
//
// Returns:
//
//	[]Location: Slice of locations matching the search (may be empty if no matches)
//	error: Non-nil if search fails (always with user-friendly message)
func (c *Client) Search(ctx context.Context, name string) ([]models.Location, error) {
	// Call SDK with default options (10 results, English language)
	opts := &geocoding.SearchOptions{
		Count:    10,
		Language: "en",
	}
	
	sdkLocations, err := c.sdkClient.Search(ctx, name, opts)
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

// mapSDKLocation converts a location from the SDK response format to
// our internal Location model.
//
// Expected mapping:
//
//	SDK.ID        → Location.ID (int)
//	SDK.Name      → Location.Name (string)
//	SDK.Latitude  → Location.Latitude (float64)
//	SDK.Longitude → Location.Longitude (float64)
//	SDK.Country   → Location.Country (string)
//	(no SDK field) → Location.Region (empty string - SDK lacks admin1 field)
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
//
// This function ensures that all errors returned to users are helpful
// without exposing technical implementation details or infrastructure info.
//
// Error mapping strategy:
//
//	context.DeadlineExceeded → "Search took too long. Please try again."
//	context.Canceled          → "Unable to search locations. Please try again."
//	geocoding.ErrConcurrencyLimitExceeded → "Unable to search locations. Please try again."
//	geocoding.ErrInvalidParameter → "Unable to search locations. Please try again."
//	geocoding.APIError        → "Unable to search locations. Please try again."
//	Any other error           → "Unable to search locations. Please try again."
//
// Returns nil if input error is nil.
func convertSDKError(err error) error {
	if err == nil {
		return nil
	}
	
	// Check for context timeout
	if errors.Is(err, context.DeadlineExceeded) {
		return errors.New("Search took too long. Please try again.")
	}
	
	// All other errors map to generic message
	return errors.New("Unable to search locations. Please try again.")
}
