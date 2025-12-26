# Data Model: Weather by Location

## Entities

### Location
Represents a geographical location found via search.

| Field | Type | Description | Source |
|-------|------|-------------|--------|
| ID | int | Unique identifier | Geocoding API (`id`) |
| Name | string | Name of the location | Geocoding API (`name`) |
| Latitude | float64 | Latitude coordinate | Geocoding API (`latitude`) |
| Longitude | float64 | Longitude coordinate | Geocoding API (`longitude`) |
| Country | string | Country name | Geocoding API (`country`) |
| Region | string | Administrative region (e.g., State) | Geocoding API (`admin1`) |

### Weather
Represents the current weather conditions for a location.

| Field | Type | Description | Source |
|-------|------|-------------|--------|
| Temperature | float64 | Air temperature (2m) | Weather API (`temperature_2m`) |
| Humidity | int | Relative humidity (2m) | Weather API (`relative_humidity_2m`) |
| ApparentTemp | float64 | Feels-like temperature | Weather API (`apparent_temperature`) |
| Precipitation | float64 | Precipitation amount | Weather API (`precipitation`) |
| CloudCover | int | Cloud cover percentage | Weather API (`cloud_cover`) |
| Pressure | float64 | Surface pressure | Weather API (`surface_pressure`) |
| WindSpeed | float64 | Wind speed (10m) | Weather API (`wind_speed_10m`) |
| WindDirection | int | Wind direction (degrees) | Weather API (`wind_direction_10m`) |
| WindGusts | float64 | Wind gusts (10m) | Weather API (`wind_gusts_10m`) |

## Value Objects

### WeatherUnits
Stores the units for the weather data (e.g., "°C", "%", "km/h").

| Field | Type | Description |
|-------|------|-------------|
| Temperature | string | Unit for temperature |
| ... | ... | Units for other fields |
