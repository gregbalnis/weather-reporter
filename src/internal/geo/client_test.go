package geo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weather-reporter/src/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResponse   string
		mockStatusCode int
		expected       []models.Location
		expectError    bool
	}{
		{
			name:           "Success",
			query:          "London",
			mockStatusCode: http.StatusOK,
			mockResponse:   `{"results": [{"id": 1, "name": "London", "latitude": 51.5085, "longitude": -0.1257, "country": "United Kingdom", "admin1": "Greater London"}]}`,
			expected: []models.Location{
				{ID: 1, Name: "London", Latitude: 51.5085, Longitude: -0.1257, Country: "United Kingdom", Region: "Greater London"},
			},
			expectError: false,
		},
		{
			name:           "No Results",
			query:          "Nowhere",
			mockStatusCode: http.StatusOK,
			mockResponse:   `{"results": []}`,
			expected:       []models.Location{},
			expectError:    false,
		},
		{
			name:           "API Error",
			query:          "Error",
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   ``,
			expected:       nil,
			expectError:    true,
		},
		{
			name:           "Malformed JSON",
			query:          "BadJSON",
			mockStatusCode: http.StatusOK,
			mockResponse:   `{"results": [{"id": 1, "name": "London", "latitude": "invalid"}]}`, // Invalid latitude type
			expected:       nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/search", r.URL.Path)
				assert.Equal(t, tt.query, r.URL.Query().Get("name"))
				w.WriteHeader(tt.mockStatusCode)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := NewClient(server.Client())
			client.baseURL = server.URL // Override base URL for testing

			results, err := client.Search(context.Background(), tt.query)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(tt.expected) == 0 {
					assert.Empty(t, results)
				} else {
					assert.Equal(t, tt.expected, results)
				}
			}
		})
	}
}

func TestSearch_Timeout(t *testing.T) {
	// Create a server that sleeps longer than the client timeout
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(10 * time.Millisecond)
	}))
	defer server.Close()

	client := NewClient(server.Client())
	client.baseURL = server.URL

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.Search(ctx, "Timeout")
	assert.Error(t, err)
	// The error message might vary slightly depending on where the timeout happens (dial, read, etc.)
	// but it should be a context error or a net error wrapping it.
	// "context deadline exceeded" is standard for ctx timeouts.
	// However, httptest server client might behave slightly differently.
	// Let's check if it's an error at all first (done above).
	// And check for common timeout indicators.
	assert.True(t, 
		assert.Contains(t, err.Error(), "context deadline exceeded") || 
		assert.Contains(t, err.Error(), "Client.Timeout exceeded") ||
		assert.Contains(t, err.Error(), "timeout"),
		"Error should indicate timeout: %v", err,
	)
}

func TestNewClient_Default(t *testing.T) {
	client := NewClient(nil)
	assert.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, defaultBaseURL, client.baseURL)
}
