# Research: Increase Unit Test Coverage

**Feature**: Increase Unit Test Coverage
**Date**: 2025-12-27

## 1. Current Coverage Analysis

Ran `go test -cover ./src/internal/...` to establish baseline.

| Package | Coverage | Status | Notes |
| :--- | :--- | :--- | :--- |
| `src/internal/geo` | 80.0% | ✅ Pass | Right on the boundary. Should increase buffer. |
| `src/internal/weather` | 81.5% | ✅ Pass | Good, but verify edge cases. |
| `src/internal/ui` | 79.3% | ❌ Fail | Needs immediate improvement. |
| `src/internal/models` | N/A | N/A | Structs/Interfaces only. |

**Conclusion**: The primary focus is `src/internal/ui`. `geo` and `weather` meet the metric but should be hardened to ensure they stay above 80% as code evolves.

## 2. Testability & Mocking Strategy

### `geo` and `weather` Packages
*   **Design**: Both clients accept `*http.Client` in their constructor.
*   **Strategy**: Use `httptest.NewServer` or a custom `http.RoundTripper` to mock API responses.
*   **Edge Cases**:
    *   Network timeouts (mock `http.Client` timeout).
    *   Non-200 status codes.
    *   Malformed JSON bodies.

### `ui` Package
*   **Design**: Likely consumes `models.GeocodingService` and `models.WeatherService`.
*   **Strategy**: Generate mocks for these interfaces using `testify/mock`.
*   **Edge Cases**:
    *   User input validation (empty, long strings).
    *   Service errors (simulated via mocks).
    *   Reader/Writer mocking for stdin/stdout.

## 3. Tooling
*   **Library**: `github.com/stretchr/testify` is already a dependency.
*   **Sub-packages**:
    *   `assert`: For readable assertions.
    *   `mock`: For generating mocks.
    *   `suite`: Optional, but standard `testing.T` is preferred for simplicity in this project.

## 4. Decisions
*   **Decision**: Focus on `ui` package first to cross the 80% threshold.
*   **Decision**: Add "buffer" tests to `geo` and `weather` to handle edge cases (timeouts, bad JSON) which will likely increase coverage further and improve robustness.
*   **Decision**: Use `httptest` for client tests and `testify/mock` for service consumers.
