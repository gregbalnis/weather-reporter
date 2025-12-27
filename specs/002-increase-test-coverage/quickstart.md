# Quickstart: Running Tests

**Feature**: Increase Unit Test Coverage
**Date**: 2025-12-27

## Prerequisites

*   Go 1.25+ installed
*   Repository cloned

## Running Tests

To run all unit tests in the project:

```bash
go test ./...
```

## Checking Coverage

To check code coverage percentages:

```bash
go test -cover ./src/internal/...
```

To generate a detailed HTML coverage report:

```bash
go test -coverprofile=coverage.out ./src/internal/...
go tool cover -html=coverage.out
```

## Verification

Ensure that all packages report **>= 80.0%** coverage.

```bash
# Example output
# ok      weather-reporter/src/internal/geo       0.015s  coverage: 85.0% of statements
# ok      weather-reporter/src/internal/ui        0.007s  coverage: 82.3% of statements
# ok      weather-reporter/src/internal/weather   0.014s  coverage: 88.5% of statements
```
