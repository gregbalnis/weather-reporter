# Data Model: Open-Meteo SDK Integration

**Feature**: Integrate Open-Meteo SDK
**Source**: [specs/004-use-meteo-sdk/spec.md](specs/004-use-meteo-sdk/spec.md)

## Entities

### Weather Response (SDK)

The SDK returns a weather response object (likely `CurrentWeather` or similar) which provides accessors for weather data.

| Field/Accessor | Type | Description | Example |
| :--- | :--- | :--- | :--- |
| `QuantityOfTemperature()` | `string` | Current air temperature with unit | "10.5°C" |
| `QuantityOfHumidity()` | `string` | Relative humidity with unit | "85%" |
| `QuantityOfApparentTemperature()` | `string` | Feels-like temperature with unit | "8.0°C" |
| `QuantityOfPrecipitation()` | `string` | Precipitation amount with unit | "0.0 mm" |
| `QuantityOfCloudCover()` | `string` | Cloud cover percentage with unit | "20%" |
| `QuantityOfPressure()` | `string` | Atmospheric pressure with unit | "1015 hPa" |
| `QuantityOfWindSpeed()` | `string` | Wind speed with unit | "15.0 km/h" |
| `QuantityOfWindDirection()` | `string` | Wind direction with unit | "180°" |
| `QuantityOfWindGusts()` | `string` | Wind gusts with unit | "25.0 km/h" |

### Application Data Flow

1.  **Input**: User provides location (e.g., "London").
2.  **Geocoding**: Application resolves location to coordinates (Lat/Lon) using existing Geo client (or SDK if supported, but spec implies replacing weather client).
3.  **Weather Retrieval**: Application calls SDK `GetCurrentWeather(lat, lon)`.
4.  **Output**: Application uses `QuantityOf...` accessors to format the output for the UI.

## Schema Definitions

No database schema changes. The data model is purely in-memory structs provided by the SDK.

## Validation Rules

- **Coordinates**: Latitude must be between -90 and 90. Longitude must be between -180 and 180. (Handled by SDK or existing validation).
- **Timeout**: Request must complete within 10 seconds.
