# Tasks: Weather by Location

**Feature Branch**: `001-weather-by-location`
**Spec**: [specs/001-weather-by-location/spec.md](specs/001-weather-by-location/spec.md)

## Phase 1: Setup
- [ ] T001 Initialize Go module `weather-reporter` in `go.mod`
- [ ] T002 Create project directory structure (`src/cmd/weather`, `src/internal/geo`, `src/internal/weather`, `src/internal/ui`, `src/internal/models`)
- [ ] T003 Install `github.com/stretchr/testify` dependency

## Phase 2: Foundational
- [ ] T004 [P] Define domain models (Location, Weather, WeatherResponse) in `src/internal/models/models.go`
- [ ] T005 [P] Define service interfaces (GeocodingService, WeatherService) in `src/internal/models/interfaces.go`

## Phase 3: User Story 1 - Get Weather for Unique Location (P1)
**Goal**: User can get weather for a specific location when the name is unique.
**Independent Test**: Run `./weather-reporter "Reykjavik"` and see weather output.

- [ ] T006 [P] [US1] Implement `GeocodingService` client in `src/internal/geo/client.go`
- [ ] T007 [P] [US1] Implement `WeatherService` client in `src/internal/weather/client.go`
- [ ] T008 [P] [US1] Add unit tests for `GeocodingService` in `src/internal/geo/client_test.go`
- [ ] T009 [P] [US1] Add unit tests for `WeatherService` in `src/internal/weather/client_test.go`
- [ ] T010 [US1] Implement main entry point with argument parsing in `src/cmd/weather/main.go`
- [ ] T011 [US1] Implement "Unique Location" workflow (Search -> Get -> Print) in `src/cmd/weather/main.go`

## Phase 4: User Story 2 - Disambiguate Location (P1)
**Goal**: User can select from a list of locations when search is ambiguous.
**Independent Test**: Run `./weather-reporter "London"`, see list, select 1, see weather.

- [ ] T012 [P] [US2] Implement `SelectLocation` prompt function in `src/internal/ui/prompt.go`
- [ ] T013 [P] [US2] Implement non-interactive mode detection in `src/internal/ui/prompt.go`
- [ ] T014 [US2] Add unit tests for UI logic in `src/internal/ui/prompt_test.go`
- [ ] T015 [US2] Update `src/cmd/weather/main.go` to handle multiple matches using `ui.SelectLocation`

## Phase 5: User Story 3 - Handle Unknown Location (P2)
**Goal**: User is informed when no location matches their search.
**Independent Test**: Run `./weather-reporter "InvalidLocationName"`, see "not found" message.

- [ ] T016 [US3] Update `src/cmd/weather/main.go` to handle 0 results from search
- [ ] T017 [US3] Verify error message output in `src/cmd/weather/main.go`

## Phase 6: Polish & Cross-Cutting Concerns
- [ ] T018 Verify all tests pass with `go test ./...`
- [ ] T019 Ensure code formatting with `go fmt ./...`

## Dependencies
1. **Setup** (T001-T003) must be done first.
2. **Foundational** (T004-T005) blocks all User Stories.
3. **US1** (T006-T011) provides the core loop.
4. **US2** (T012-T015) and **US3** (T016-T017) extend US1 and can be done in parallel or sequence, but US2 is higher priority.

## Implementation Strategy
- **MVP**: Complete Phase 1, 2, and 3. This gives a working tool for unique locations.
- **Incremental**: Add US2 for better usability, then US3 for error handling.
