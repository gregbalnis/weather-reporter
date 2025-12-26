# Research: Weather by Location

## Unknowns & Clarifications

### 1. Open-Meteo Geocoding API
- **Endpoint**: `https://geocoding-api.open-meteo.com/v1/search`
- **Parameters**:
  - `name`: Location name (required).
  - `count`: Number of results (default 10, max 100). We will use 10.
  - `language`: `en` (default).
  - `format`: `json` (default).
- **Response**: JSON object with `results` array. Each result contains `id`, `name`, `latitude`, `longitude`, `country`, `admin1` (region), etc.
- **Error Handling**: Returns 400 for bad params. Empty `results` if not found.

### 2. Open-Meteo Weather API
- **Endpoint**: `https://api.open-meteo.com/v1/forecast`
- **Parameters**:
  - `latitude`, `longitude`: Coordinates (required).
  - `current`: Comma-separated list of variables.
    - `temperature_2m`
    - `relative_humidity_2m`
    - `apparent_temperature`
    - `precipitation`
    - `cloud_cover`
    - `surface_pressure` (or `pressure_msl`)
    - `wind_speed_10m`
    - `wind_direction_10m`
    - `wind_gusts_10m`
  - `wind_speed_unit`: `kmh` (default).
  - `temperature_unit`: `celsius` (default).
- **Response**: JSON object with `current` (values) and `current_units` (units).

### 3. CLI Argument Parsing
- **Decision**: Use standard `flag` package.
- **Pattern**:
  - `flag.Parse()`
  - `flag.Args()` to get the location name (joined by space if multiple args).
  - If no args, print usage and exit.

### 4. Interactive Prompts
- **Decision**: Use `bufio.NewReader(os.Stdin)` to read user input.
- **Pattern**:
  - Print list of options with index (1-N).
  - Print prompt "Select location [1-N]: ".
  - Read string until newline.
  - Parse integer and validate range.
  - Retry on invalid input.

### 5. Testing Strategy
- **Unit Tests**:
  - Mock HTTP client interface to test service logic without real network calls.
  - Test CLI parsing logic.
- **Integration Tests**:
  - Optional: Real API call (skipped in CI or flagged).

## Library Evaluations

### Open-Meteo Client
- **Candidates**: `HectorMalot/omgo`, Custom `net/http` client.
- **Evaluation**:
  - `omgo`: Good for forecast, but lacks robust Geocoding support. Maintenance is sporadic.
  - Custom: The Open-Meteo API is simple REST/JSON. Writing a custom client allows exact control over fields (reducing payload) and avoids unused dependencies.
- **Decision**: Use **Standard `net/http`**.

### Interactive Prompts
- **Candidates**: `charmbracelet/huh`, `manifoldco/promptui`, `AlecAivazis/survey`, Standard `bufio`.
- **Evaluation**:
  - `huh`: Excellent UX, but heavy dependency tree (Bubble Tea).
  - `survey`: Archived/Unmaintained.
  - `bufio`: Zero dependency. Fits the strict "select by index/number" requirement perfectly.
- **Decision**: Use **Standard `bufio`** to maintain low footprint and strict spec compliance.

### Testing Assertions
- **Candidates**: Standard `testing`, `stretchr/testify`.
- **Evaluation**:
  - `testify`: Industry standard for Go assertions. Reduces boilerplate code significantly compared to `if got != want`.
- **Decision**: Use **`github.com/stretchr/testify`**.

## Decisions
- **HTTP Client**: Use `net/http` with a custom interface for mocking.
- **JSON Parsing**: Use `encoding/json` with struct tags.
- **Output**: Use `fmt.Printf` for formatted output.
- **Testing**: Use `testify` for assertions.
