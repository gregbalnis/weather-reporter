# Quickstart: Weather Reporter

## Prerequisites
- Go 1.21+ installed
- Internet connection (for Open-Meteo API)

## Build
```bash
go build -o weather-reporter ./cmd/weather
```

## Usage

### Basic Usage
Get weather for a specific location:
```bash
./weather-reporter "New York"
```

### Handling Multiple Matches
If multiple locations match, follow the interactive prompt:
```text
$ ./weather-reporter London
Multiple locations found:
1. London, United Kingdom (Greater London)
2. London, Canada (Ontario)
3. London, United States (Kentucky)
...
Select location [1-10]: 1
```

### Non-Interactive Mode
For scripts, the tool will list matches and exit if ambiguous:
```bash
$ ./weather-reporter London | cat
Multiple locations found:
1. London, United Kingdom (Greater London)
...
```

## Development

### Running Tests
```bash
go test ./...
```
