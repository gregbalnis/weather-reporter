// Package main is the entry point for the weather-reporter CLI application.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"weather-reporter/src/internal/geo"
	"weather-reporter/src/internal/models"
	"weather-reporter/src/internal/ui"
	"weather-reporter/src/internal/weather"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	// Initialize services
	geoClient := geo.NewClient(nil)
	weatherClient := weather.NewClient(nil)

	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, geoClient, weatherClient, defaultInteractiveChecker))
}

type interactiveChecker func(io.Reader) bool

func defaultInteractiveChecker(r io.Reader) bool {
	if f, ok := r.(*os.File); ok {
		return ui.IsTerminal(f)
	}
	return false
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, geoClient models.GeocodingService, weatherClient models.WeatherService, isInteractive interactiveChecker) int {
	fs := flag.NewFlagSet("weather-reporter", flag.ContinueOnError)
	fs.SetOutput(stderr)
	versionFlag := fs.Bool("version", false, "Print version information")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *versionFlag {
		_, _ = fmt.Fprintf(stdout, "weather-reporter version %s\n", Version)
		_, _ = fmt.Fprintf(stdout, "commit: %s\n", Commit)
		_, _ = fmt.Fprintf(stdout, "built at: %s\n", Date)
		return 0
	}

	locationArgs := fs.Args()
	if len(locationArgs) == 0 {
		_, _ = fmt.Fprintln(stdout, "Usage: weather-reporter <location>")
		return 1
	}

	locationName := strings.Join(locationArgs, " ")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Search for location
	locations, err := geoClient.Search(ctx, locationName)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error searching for location: %v\n", err)
		return 1
	}

	if len(locations) == 0 {
		_, _ = fmt.Fprintf(stdout, "Location not found: %s\n", locationName)
		return 0
	}

	var selectedLocation models.Location

	if len(locations) == 1 {
		selectedLocation = locations[0]
	} else {
		interactive := isInteractive(stdin)
		selectedLocation, err = ui.SelectLocation(locations, stdin, stdout, interactive)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error selecting location: %v\n", err)
			return 1
		}
	}

	// 2. Get Weather
	weatherData, err := weatherClient.GetCurrentWeather(ctx, selectedLocation.Latitude, selectedLocation.Longitude)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error fetching weather: %v\n", err)
		return 1
	}

	// 3. Print Weather
	if err := ui.PrintWeather(stdout, selectedLocation, weatherData); err != nil {
		_, _ = fmt.Fprintf(stderr, "Error printing weather: %v\n", err)
		return 1
	}

	return 0
}
