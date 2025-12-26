package ui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"weather-reporter/src/internal/models"
)

func TestSelectLocation(t *testing.T) {
	locations := []models.Location{
		{ID: 1, Name: "London", Country: "UK", Region: "Greater London"},
		{ID: 2, Name: "London", Country: "Canada", Region: "Ontario"},
	}

	t.Run("Interactive Success", func(t *testing.T) {
		input := "1\n"
		in := strings.NewReader(input)
		var out bytes.Buffer

		loc, err := SelectLocation(locations, in, &out, true)

		assert.NoError(t, err)
		assert.Equal(t, locations[0], loc)
		assert.Contains(t, out.String(), "Multiple locations found:")
		assert.Contains(t, out.String(), "1. London, UK (Greater London)")
	})

	t.Run("Interactive Invalid Input Then Success", func(t *testing.T) {
		input := "invalid\n3\n2\n" // 3 is out of range
		in := strings.NewReader(input)
		var out bytes.Buffer

		loc, err := SelectLocation(locations, in, &out, true)

		assert.NoError(t, err)
		assert.Equal(t, locations[1], loc)
		assert.Contains(t, out.String(), "Invalid selection")
	})

	t.Run("Non-Interactive", func(t *testing.T) {
		in := strings.NewReader("")
		var out bytes.Buffer

		loc, err := SelectLocation(locations, in, &out, false)

		assert.Error(t, err)
		assert.Equal(t, "multiple locations found, please be more specific", err.Error())
		assert.Equal(t, models.Location{}, loc)
		assert.Contains(t, out.String(), "Multiple locations found:")
		assert.Contains(t, out.String(), "1. London, UK (Greater London)")
	})

	t.Run("Limit to 10", func(t *testing.T) {
		manyLocations := make([]models.Location, 12)
		for i := 0; i < 12; i++ {
			manyLocations[i] = models.Location{ID: i, Name: "Loc"}
		}

		in := strings.NewReader("")
		var out bytes.Buffer

		_, err := SelectLocation(manyLocations, in, &out, false)
		assert.Error(t, err)

		output := out.String()
		assert.Contains(t, output, "10. Loc")
		assert.NotContains(t, output, "11. Loc")
	})
}
