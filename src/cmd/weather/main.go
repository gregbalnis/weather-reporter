package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"weather-reporter/src/internal/geo"
	"weather-reporter/src/internal/models"
	"weather-reporter/src/internal/ui"
	"weather-reporter/src/internal/weather"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: weather-reporter <location>")
		os.Exit(1)
	}

	locationName := strings.Join(args, " ")

	// Initialize services
	geoClient := geo.NewClient(nil)
	weatherClient := weather.NewClient(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Search for location
	locations, err := geoClient.Search(ctx, locationName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching for location: %v\n", err)
		os.Exit(1)
	}

	if len(locations) == 0 {
		fmt.Printf("Location not found: %s\n", locationName)
		os.Exit(0)
	}

	var selectedLocation models.Location

	if len(locations) == 1 {
		selectedLocation = locations[0]
	} else {
		isInteractive := ui.IsTerminal(os.Stdin)
		selectedLocation, err = ui.SelectLocation(locations, os.Stdin, os.Stdout, isInteractive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting location: %v\n", err)
			os.Exit(1)
		}
	}

	// 2. Get Weather
	weatherData, err := weatherClient.GetCurrentWeather(ctx, selectedLocation.Latitude, selectedLocation.Longitude)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching weather: %v\n", err)
		os.Exit(1)
	}

	// 3. Print Weather
	printWeather(selectedLocation, weatherData)
}

func printWeather(loc models.Location, w *models.WeatherResponse) {
	fmt.Printf("Weather for %s, %s (%s)\n", loc.Name, loc.Country, loc.Region)
	fmt.Println("------------------------------------------------")
	fmt.Printf("Temperature:          %.1f %s\n", w.Current.Temperature, w.CurrentUnits.Temperature)
	fmt.Printf("Apparent Temperature: %.1f %s\n", w.Current.ApparentTemp, w.CurrentUnits.ApparentTemp)
	fmt.Printf("Humidity:             %d %s\n", w.Current.Humidity, w.CurrentUnits.Humidity)
	fmt.Printf("Precipitation:        %.1f %s\n", w.Current.Precipitation, w.CurrentUnits.Precipitation)
	fmt.Printf("Cloud Cover:          %d %s\n", w.Current.CloudCover, w.CurrentUnits.CloudCover)
	fmt.Printf("Pressure:             %.1f %s\n", w.Current.Pressure, w.CurrentUnits.Pressure)
	fmt.Printf("Wind Speed:           %.1f %s\n", w.Current.WindSpeed, w.CurrentUnits.WindSpeed)
	fmt.Printf("Wind Direction:       %d %s\n", w.Current.WindDirection, w.CurrentUnits.WindDirection)
	fmt.Printf("Wind Gusts:           %.1f %s\n", w.Current.WindGusts, w.CurrentUnits.WindGusts)
}
