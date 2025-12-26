package geo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"weather-reporter/src/internal/models"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/search", r.URL.Path)
				assert.Equal(t, tt.query, r.URL.Query().Get("name"))
				w.WriteHeader(tt.mockStatusCode)
				w.Write([]byte(tt.mockResponse))
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
