package ui

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"weather-reporter/src/internal/models"

	"github.com/stretchr/testify/assert"
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

func TestIsTerminal(t *testing.T) {
	// Create a temp file
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }() // clean up
	defer func() { _ = tmpfile.Close() }()

	// It's a file, not a terminal
	assert.False(t, IsTerminal(tmpfile))
}

type errorReader struct{}

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestSelectLocation_ReadError(t *testing.T) {
	locations := []models.Location{
		{ID: 1, Name: "London"},
	}
	var out bytes.Buffer

	_, err := SelectLocation(locations, &errorReader{}, &out, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read input")
}

type mockWeatherResponse struct{}

func (m mockWeatherResponse) QuantityOfTemperature() string         { return "20°C" }
func (m mockWeatherResponse) QuantityOfHumidity() string            { return "50%" }
func (m mockWeatherResponse) QuantityOfApparentTemperature() string { return "18°C" }
func (m mockWeatherResponse) QuantityOfPrecipitation() string       { return "0mm" }
func (m mockWeatherResponse) QuantityOfCloudCover() string          { return "10%" }
func (m mockWeatherResponse) QuantityOfPressure() string            { return "1013hPa" }
func (m mockWeatherResponse) QuantityOfWindSpeed() string           { return "10km/h" }
func (m mockWeatherResponse) QuantityOfWindDirection() string       { return "N" }
func (m mockWeatherResponse) QuantityOfWindGusts() string           { return "15km/h" }

func TestPrintWeather(t *testing.T) {
loc := models.Location{
Name:    "Test City",
Country: "Test Country",
Region:  "Test Region",
}
w := mockWeatherResponse{}
var out bytes.Buffer

err := PrintWeather(&out, loc, w)
assert.NoError(t, err)

output := out.String()
assert.Contains(t, output, "Weather for Test City, Test Country (Test Region)")
assert.Contains(t, output, "Temperature:          20°C")
assert.Contains(t, output, "Apparent Temperature: 18°C")
assert.Contains(t, output, "Humidity:             50%")
assert.Contains(t, output, "Precipitation:        0mm")
assert.Contains(t, output, "Cloud Cover:          10%")
assert.Contains(t, output, "Pressure:             1013hPa")
assert.Contains(t, output, "Wind Speed:           10km/h")
assert.Contains(t, output, "Wind Direction:       N")
assert.Contains(t, output, "Wind Gusts:           15km/h")
}

type errorWriter struct{}

func (e errorWriter) Write(p []byte) (n int, err error) {
return 0, fmt.Errorf("write error")
}

func TestPrintWeather_Error(t *testing.T) {
loc := models.Location{Name: "Test"}
w := mockWeatherResponse{}
out := errorWriter{}

err := PrintWeather(out, loc, w)
assert.Error(t, err)
assert.Equal(t, "write error", err.Error())
}
