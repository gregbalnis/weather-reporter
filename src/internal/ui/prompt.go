package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"weather-reporter/src/internal/models"
)

// IsTerminal checks if the file is a terminal.
func IsTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// SelectLocation prompts the user to select a location from a list.
func SelectLocation(locations []models.Location, in io.Reader, out io.Writer, interactive bool) (models.Location, error) {
	if len(locations) == 0 {
		return models.Location{}, fmt.Errorf("no locations to select from")
	}

	// Limit to 10
	displayCount := len(locations)
	if displayCount > 10 {
		displayCount = 10
	}

	if !interactive {
		_, _ = fmt.Fprintln(out, "Multiple locations found:")
		printLocations(out, locations[:displayCount])
		return models.Location{}, fmt.Errorf("multiple locations found, please be more specific")
	}

	_, _ = fmt.Fprintln(out, "Multiple locations found:")
	printLocations(out, locations[:displayCount])

	reader := bufio.NewReader(in)

	for {
		_, _ = fmt.Fprintf(out, "Select location [1-%d]: ", displayCount)
		input, err := reader.ReadString('\n')
		if err != nil {
			return models.Location{}, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > displayCount {
			_, _ = fmt.Fprintf(out, "Invalid selection. Please enter a number between 1 and %d.\n", displayCount)
			continue
		}

		return locations[index-1], nil
	}
}

func printLocations(out io.Writer, locations []models.Location) {
	for i, loc := range locations {
		_, _ = fmt.Fprintf(out, "%d. %s, %s (%s)\n", i+1, loc.Name, loc.Country, loc.Region)
	}
}

// PrintWeather prints the weather information to the output writer.
func PrintWeather(out io.Writer, loc models.Location, w models.WeatherResponse) error {
	if _, err := fmt.Fprintf(out, "Weather for %s, %s (%s)\n", loc.Name, loc.Country, loc.Region); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, "------------------------------------------------"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Temperature:          %s\n", w.QuantityOfTemperature()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Apparent Temperature: %s\n", w.QuantityOfApparentTemperature()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Humidity:             %s\n", w.QuantityOfHumidity()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Precipitation:        %s\n", w.QuantityOfPrecipitation()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Cloud Cover:          %s\n", w.QuantityOfCloudCover()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Pressure:             %s\n", w.QuantityOfPressure()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Wind Speed:           %s\n", w.QuantityOfWindSpeed()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Wind Direction:       %s\n", w.QuantityOfWindDirection()); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Wind Gusts:           %s\n", w.QuantityOfWindGusts()); err != nil {
		return err
	}
	return nil
}
