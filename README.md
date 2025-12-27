# Weather Reporter

![CI](https://github.com/owner/weather-reporter/actions/workflows/ci.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/owner/weather-reporter)

Weather Reporter is a command-line interface (CLI) tool written in Go that fetches current weather information for any location using the [Open-Meteo API](https://open-meteo.com/).

## Features

- **Location Search**: Search for cities, towns, or villages by name.
- **Interactive Selection**: Disambiguates between locations with the same name (e.g., "London, UK" vs "London, Canada") via an interactive prompt.
- **Detailed Weather Data**: Displays temperature, humidity, wind speed, precipitation, and more.
- **Metric Units**: All data is presented in metric units (Celsius, km/h, mm).
- **Script Friendly**: Detects non-interactive environments and exits gracefully with information.

## Prerequisites

- [Go](https://go.dev/dl/) 1.21 or higher.
- An active internet connection.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/weather-reporter.git
   cd weather-reporter
   ```

2. Build the binary:
   ```bash
   go build -o weather-reporter ./src/cmd/weather-reporter
   ```

## Usage

### Basic Usage

Get the weather for a specific location:

```bash
./weather-reporter "New York"
```

**Output:**
```text
Weather for New York, United States (New York)
------------------------------------------------
Temperature:          15.2 °C
Apparent Temperature: 14.0 °C
Humidity:             60 %
...
```

### Handling Multiple Matches

If multiple locations match your query, the tool will ask you to select the correct one:

```bash
$ ./weather-reporter London
Multiple locations found:
1. London, United Kingdom (England)
2. London, Canada (Ontario)
3. London, United States (Ohio)
...
Select location [1-10]: 1
```

### Non-Interactive Mode (Scripts)

If you run the tool in a non-interactive environment (e.g., piped to another command), it will list the matches and exit with an error to prevent hanging:

```bash
$ ./weather-reporter London | cat
Multiple locations found:
1. London, United Kingdom (England)
...
Error selecting location: multiple locations found, please be more specific
```

## Development

### Running Tests

To run the unit tests:

```bash
go test ./...
```

### Project Structure

- `src/cmd/weather-reporter`: Main entry point.
- `src/internal/geo`: Geocoding service client.
- `src/internal/weather`: Weather service client.
- `src/internal/ui`: User interaction logic.
- `src/internal/models`: Shared data models.

## License

This project is licensed under the terms of the [LICENSE](LICENSE) file.
