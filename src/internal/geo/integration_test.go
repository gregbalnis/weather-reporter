// Package geo provides functionality for searching locations using a geocoding API.
package geo

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"weather-reporter/src/internal/models"
)

// TestIntegration_GeocodingAPIContract validates that the Open-Meteo Geocoding API
// maintains its contract over time. This test makes real network calls and should be
// skipped in short test runs or environments without network access.
//
// Purpose: Detect if the upstream API changes its response format or behavior,
// which would break our adapter implementation.
func TestIntegration_GeocodingAPIContract(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a real client with reasonable timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := NewClient(httpClient)

	t.Run("London_Search_Returns_Valid_Data", func(t *testing.T) {
		ctx := context.Background()
		locations, err := client.Search(ctx, "London")

		assert.NoError(t, err, "London search should succeed")
		assert.NotEmpty(t, locations, "Should return at least one result")

		// Validate first result has expected structure
		london := locations[0]
		assert.NotZero(t, london.ID, "ID should be populated")
		assert.NotEmpty(t, london.Name, "Name should be populated")
		assert.NotZero(t, london.Latitude, "Latitude should be populated")
		assert.NotZero(t, london.Longitude, "Longitude should be populated")
		assert.NotEmpty(t, london.Country, "Country should be populated")
		// Region may be empty (SDK doesn't provide admin1 field)

		// Sanity check: London should be in UK/GB
		assert.Contains(t, []string{"United Kingdom", "UK", "England"}, london.Country, "London should be in United Kingdom")

		// Sanity check: London coordinates should be roughly correct
		assert.InDelta(t, 51.5, london.Latitude, 1.0, "London latitude should be ~51.5")
		assert.InDelta(t, -0.1, london.Longitude, 1.0, "London longitude should be ~-0.1")
	})

	t.Run("Ambiguous_Query_Returns_Multiple_Results", func(t *testing.T) {
		ctx := context.Background()
		locations, err := client.Search(ctx, "Springfield")

		assert.NoError(t, err, "Springfield search should succeed")
		assert.NotEmpty(t, locations, "Should return multiple results for ambiguous query")
		assert.GreaterOrEqual(t, len(locations), 2, "Springfield exists in multiple countries/states")

		// Each result should have complete data
		for i, loc := range locations {
			assert.NotZero(t, loc.ID, "Result %d: ID should be populated", i)
			assert.NotEmpty(t, loc.Name, "Result %d: Name should be populated", i)
			assert.NotZero(t, loc.Latitude, "Result %d: Latitude should be populated", i)
			assert.NotZero(t, loc.Longitude, "Result %d: Longitude should be populated", i)
			assert.NotEmpty(t, loc.Country, "Result %d: Country should be populated", i)
		}
	})

	t.Run("Result_Limit_Respected", func(t *testing.T) {
		ctx := context.Background()
		locations, err := client.Search(ctx, "London")

		assert.NoError(t, err, "Search should succeed")
		assert.LessOrEqual(t, len(locations), 10, "Should return at most 10 results (SDK default limit)")
	})

	t.Run("No_Results_Returns_Empty_Slice", func(t *testing.T) {
		ctx := context.Background()
		locations, err := client.Search(ctx, "ZZZNonexistentCityXYZ123")

		assert.NoError(t, err, "No results should not be an error")
		assert.Empty(t, locations, "Should return empty slice for no results")
	})

	t.Run("Complete_Data_Structure", func(t *testing.T) {
		ctx := context.Background()
		locations, err := client.Search(ctx, "Tokyo")

		assert.NoError(t, err, "Tokyo search should succeed")
		assert.NotEmpty(t, locations, "Should return at least one result")

		tokyo := locations[0]

		// Verify all Location struct fields are accessible
		// This ensures the SDK → internal model mapping is complete
		var testLoc models.Location = tokyo
		_ = testLoc.ID
		_ = testLoc.Name
		_ = testLoc.Latitude
		_ = testLoc.Longitude
		_ = testLoc.Country
		_ = testLoc.Region

		// Validate Tokyo specifics
		assert.Equal(t, "Tokyo", tokyo.Name, "First result should be Tokyo")
		assert.Contains(t, []string{"Japan", "JP"}, tokyo.Country, "Tokyo should be in Japan")
		assert.InDelta(t, 35.7, tokyo.Latitude, 1.0, "Tokyo latitude should be ~35.7")
		assert.InDelta(t, 139.7, tokyo.Longitude, 1.0, "Tokyo longitude should be ~139.7")
	})
}

// TestIntegration_APIContractChange documents the expected API contract behavior.
// If this test fails, it indicates the upstream API has changed its response format
// or behavior, which requires adapter updates.
//
// Expected Contract:
//   - GET /v1/search?name={query}&count=10&language=en&format=json
//   - Response: { "results": [ { "id": int, "name": string, "latitude": float, "longitude": float,
//     "country": string } ] }
//   - Empty results: { "results": [] } or { "results": null }
//   - Errors: HTTP 4xx/5xx status codes
//   - Timeout: Respects context deadlines
//
// If this test fails:
// 1. Check if API response structure changed
// 2. Update mapSDKLocation() function in client.go
// 3. Update convertSDKError() if error handling changed
// 4. Update this documentation with new contract
func TestIntegration_APIContractChange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := NewClient(httpClient)

	t.Run("Context_Timeout_Honored", func(t *testing.T) {
		// Create a context that times out immediately
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait for context to definitely timeout
		time.Sleep(10 * time.Millisecond)

		_, err := client.Search(ctx, "London")

		assert.Error(t, err, "Should return error when context times out")
		assert.Contains(t, err.Error(), "Search took too long", "Should return user-friendly timeout message")
	})

	t.Run("Network_Error_Handling", func(t *testing.T) {
		// Create client with invalid base URL to force network error
		// Note: This test may be fragile depending on SDK error handling
		// If SDK doesn't expose a way to override base URL, this test can be removed
		t.Skip("Network error simulation requires baseURL override mechanism")
	})
}
