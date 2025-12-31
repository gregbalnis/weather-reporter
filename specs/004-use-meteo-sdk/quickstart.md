# Quickstart: Open-Meteo SDK Integration

## Overview

This feature integrates the `open-meteo-weather-sdk` to fetch weather data.

## Prerequisites

- Go 1.25.5 or later
- Internet connection (to reach Open-Meteo API)

## Running the Application

Build and run the application as usual:

```bash
go build -o weather-reporter ./cmd/weather-reporter
./weather-reporter "London"
```

## Testing

Run the tests to verify the SDK integration:

```bash
go test ./...
```

## Key Changes

- The internal weather client has been replaced by the SDK.
- Output formatting now relies on the SDK's `QuantityOf...` methods.
