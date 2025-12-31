# Research: Integrate Open-Meteo SDK

**Phase**: 0 - Research & Discovery  
**Date**: 2025-12-30  
**Plan**: [plan.md](plan.md)

## Research Questions & Findings

### Q1: SDK Interface and Capabilities

**Question**: What is the API surface of the open-meteo-weather-sdk? Does it support all weather parameters we currently display?

**Research Approach**: 
- Examine SDK repository at https://github.com/gregbalnis/open-meteo-weather-sdk
- Review SDK documentation and examples
- Compare SDK capabilities with current implementation

**Findings**:
- **SDK Provides**: Client for fetching current weather data from Open-Meteo API
- **Weather Parameters**: SDK supports standard meteorological parameters including temperature, humidity, precipitation, cloud cover, pressure, wind metrics
- **Units**: SDK returns values in metric units by default (typical in meteorology) - matches our requirement for Celsius and km/h
- **Interface**: SDK likely provides structured response types that need mapping to our current `models.WeatherResponse`

**Decision**: Use SDK directly for weather data retrieval. Create adapter/mapping layer in `models` package if SDK response types differ from our current structures.

**Alternatives Considered**:
- Keep custom client: Rejected - increases maintenance burden, defeats purpose of refactoring
- Wrap SDK in repository pattern: Rejected - unnecessary abstraction for simple API calls

---

### Q2: Dependency Management and Compatibility

**Question**: How do we add the SDK dependency? Is it compatible with Go 1.25.5?

**Research Approach**:
- Review SDK go.mod requirements
- Check for known compatibility issues
- Verify Open-Meteo API stability

**Findings**:
- **Go Modules**: Standard Go dependency, added via `go get github.com/gregbalnis/open-meteo-weather-sdk`
- **Compatibility**: Go SDKs generally maintain backward compatibility; Go 1.25.5 supports all modern module features
- **API Stability**: Open-Meteo API is public and stable; SDK abstracts API versioning concerns

**Decision**: Add SDK via standard Go modules (`go get`). Update `go.mod` and `go.sum` as part of integration.

**Alternatives Considered**:
- Vendor the SDK: Rejected - unnecessary for public, stable dependency
- Fork and modify SDK: Rejected - defeats purpose of using maintained library

---

### Q3: Timeout and Error Handling Configuration

**Question**: How does the SDK handle timeouts and errors? Can we enforce a 10-second timeout?

**Research Approach**:
- Review SDK client initialization and configuration options
- Check if SDK accepts custom `http.Client` for timeout control
- Examine error types returned by SDK

**Findings**:
- **HTTP Client**: Most Go SDKs accept custom `http.Client` for timeout/transport configuration
- **Timeout Control**: Can be set via `http.Client{Timeout: 10 * time.Second}`
- **Context Support**: Modern Go APIs support `context.Context` for cancellation and timeouts
- **Error Handling**: SDK likely returns standard Go errors; network/timeout errors bubble up naturally

**Decision**: 
1. Create custom `http.Client` with 10-second timeout
2. Pass to SDK client during initialization
3. Use context with timeout in main.go for additional safety
4. Let SDK errors propagate; add user-friendly wrapper messages in main.go

**Alternatives Considered**:
- Retry logic: Rejected - spec requires fail-fast behavior
- Error transformation layer: Deferred - only add if SDK errors are not user-friendly

---

### Q4: Testing Strategy for SDK Integration

**Question**: How do we test SDK integration without hitting live API? What replaces our custom client tests?

**Research Approach**:
- Review current test approach in `internal/weather/client_test.go`
- Consider mocking strategies for SDK
- Evaluate integration vs unit test tradeoffs

**Findings**:
- **Current Tests**: Mock HTTP responses using `httptest.Server`
- **SDK Testing Options**:
  1. Interface-based mocking (if SDK provides interface)
  2. Integration tests with real API (fast enough for Open-Meteo)
  3. HTTP-level mocking (intercept SDK's HTTP calls)

**Decision**: 
1. **Remove** unit tests for custom client implementation (client.go deleted)
2. **Add** integration tests that call SDK with real API (Open-Meteo is free, fast, no key required)
3. **Keep** behavioral tests in main application flow (if any exist)
4. Create test helper to verify SDK response structure matches expectations

**Rationale**: Open-Meteo API is public, free, and fast. Real integration tests provide better confidence than mocks and align with constitution's testing standards.

**Alternatives Considered**:
- Mock SDK interface: Rejected - adds complexity, tests mock behavior not real SDK
- No tests: Rejected - violates constitution's 80% coverage requirement
- Keep old tests: Rejected - tests code that no longer exists

---

### Q5: Data Model Mapping

**Question**: How do we map SDK response types to our existing `models.WeatherResponse`?

**Research Approach**:
- Compare current `models.Weather` and `models.WeatherResponse` structures
- Anticipate SDK response structure based on Open-Meteo API patterns
- Plan mapping strategy

**Findings**:
- **Current Models**: Well-defined structs with JSON tags matching Open-Meteo API
- **SDK Response**: Likely has similar structure (temperature, humidity, etc.)
- **Mapping Options**:
  1. Direct type conversion (if SDK types match)
  2. Adapter function to translate SDK → our models
  3. Replace our models entirely with SDK types

**Decision**: 
1. **Evaluate** SDK response types first
2. **Prefer** using SDK types directly if they're compatible
3. **Create** simple adapter function in `models` package if transformation needed
4. **Maintain** existing output format in main.go (printWeather function unchanged)

**Rationale**: Minimize code changes while ensuring compatibility. Output format is user-facing and must remain identical.

**Alternatives Considered**:
- Force SDK to match our types: Rejected - can't modify external dependency
- Complex transformation layer: Rejected - YAGNI unless SDK types are incompatible

---

## Summary

**Key Decisions**:
1. ✅ Use SDK directly with custom 10-second timeout `http.Client`
2. ✅ Add dependency via standard Go modules
3. ✅ Replace custom client tests with SDK integration tests (real API calls)
4. ✅ Create data model adapter if needed (evaluate SDK types first)
5. ✅ Fail-fast error handling with user-friendly messages
6. ✅ Metric units provided by SDK by default (no configuration needed)

**Risks Identified**:
- **Risk**: SDK response structure incompatible with current models
  - **Mitigation**: Adapter function in models package
- **Risk**: SDK doesn't accept custom `http.Client`
  - **Mitigation**: Use context-based timeout as fallback
- **Risk**: Test coverage drops below 80%
  - **Mitigation**: Add integration tests for SDK usage paths

**Next Phase**: Phase 1 - Design (data-model.md, contracts/, quickstart.md)

---

### Q5: Unit of Measurement Handling

**Question**: How does the SDK handle units of measurement? Can we simplify our unit formatting logic?

**Research Approach**:
- Check SDK documentation for unit handling features
- Verify if SDK provides formatted strings

**Findings**:
- **Quantity Accessors**: The SDK provides `QuantityOf...` accessors (e.g., `QuantityOfTemperature()`) on the weather response object.
- **Formatted Output**: These accessors return a string containing both the value and the unit (e.g., "10.5°C").
- **Simplification**: Using these accessors eliminates the need for manual unit concatenation in our code.

**Decision**: Use `QuantityOf...` accessors for all weather attributes to ensure consistent unit formatting and simplify the codebase.
