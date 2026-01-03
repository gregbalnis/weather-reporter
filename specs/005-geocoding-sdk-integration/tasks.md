# Tasks: Geocoding SDK Integration

**Branch**: `005-geocoding-sdk-integration`  
**Input**: Design documents from `/specs/005-geocoding-sdk-integration/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/sdk-adapter.go

## Overview

Break down the geocoding SDK integration into actionable tasks organized by implementation phase. All tasks are independent and can be executed in parallel or sequentially based on team capacity.

## Task Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependencies)
- **[Story]**: User story label (US1, US2) - refactoring tasks have no story label
- All file paths are relative to repository root

---

## Phase 1: Research & Design

**Duration**: 1 day  
**Purpose**: Analyze SDK and define implementation strategy  
**Gate**: All unknowns resolved before implementation begins

- [ ] T001 Research open-meteo-geocoding-sdk repository and API documentation
  - **File**: `specs/005-geocoding-sdk-integration/research.md` (reference)
  - **Details**: 
    - Visit SDK repository and understand client API
    - Document client struct name, initialization method, search method signature
    - Identify response types and field names
    - Document error types and handling patterns
    - Verify SDK supports: max 10 results, English language, timeout configuration
  - **Success**: All unknowns answered, notes in research.md

- [ ] T002 [P] Analyze type mappings between SDK and Location model
  - **File**: `specs/005-geocoding-sdk-integration/data-model.md` (reference)
  - **Details**:
    - Compare SDK Location type fields with our Location struct (id, name, latitude, longitude, country, admin1)
    - Identify field name mismatches (e.g., SDK's Admin1 → our Region)
    - Determine type conversions needed
    - Plan mapping function signature
  - **Success**: Type mapping strategy documented and validated

- [ ] T003 [P] Define error handling mapping from SDK errors to user-friendly messages
  - **File**: `specs/005-geocoding-sdk-integration/data-model.md` (reference)
  - **Details**:
    - Identify all SDK error types
    - Define mapping to user-friendly messages per FR-005/006/007
    - Plan error conversion function
  - **Success**: Error mapping strategy complete

- [ ] T004 [P] Review and finalize integration test strategy
  - **File**: `specs/005-geocoding-sdk-integration/quickstart.md` (reference)
  - **Details**:
    - Confirm integration test location: `src/internal/geo/integration_test.go`
    - Finalize test cases: basic search, multiple results, max limit, complete data
    - Confirm can be skipped with `-short` flag
  - **Success**: Integration test plan finalized

**Checkpoint**: All research complete, unknowns resolved, ready for implementation

---

## Phase 2: Core Implementation

**Duration**: 2-3 days  
**Purpose**: Implement SDK adapter and integration test  
**Gate**: All existing tests pass, new integration test validates API contract, zero performance degradation

- [ ] T005 Create integration test file `src/internal/geo/integration_test.go`
  - **File**: `src/internal/geo/integration_test.go` (new file)
  - **Details**:
    - Test function: `TestIntegration_GeocodingAPIContract`
    - Skip short tests: `if testing.Short() { t.Skip(...) }`
    - Sub-tests for: London search, ambiguous queries, result limits, complete data
    - Test function: `TestIntegration_APIContractChange` documenting expected contract
  - **Dependencies**: Task T001-T004 (research complete)
  - **Success**: Integration test file created with all test cases

- [ ] T006 [P] Update `go.mod` to add SDK dependency
  - **File**: `go.mod`
  - **Details**:
    - Run: `go get github.com/gregbalnis/open-meteo-geocoding-sdk@[version]`
    - Pin to specific version (not @latest)
    - Run: `go mod tidy`
    - Verify no version conflicts
  - **Success**: SDK dependency added and pinned in go.mod

- [ ] T007 Create mapping function from SDK Location to internal Location model
  - **File**: `src/internal/geo/client.go`
  - **Details**:
    - Function signature: `mapSDKLocation(sdkLocation interface{}) (models.Location, error)`
    - Map SDK fields to Location fields
    - Handle field name differences (Admin1 → Region)
    - Validate required fields present
    - Handle optional Region field
  - **Dependencies**: Task T002 (type mapping strategy)
  - **Success**: Mapping function implemented and testable

- [ ] T008 Create error conversion function for SDK errors
  - **File**: `src/internal/geo/client.go`
  - **Details**:
    - Function signature: `convertSDKError(err error) error`
    - Map SDK error types to user-friendly messages
    - Ensure no technical details exposed
    - All errors map to "Unable to search locations. Please try again." or "Search took too long. Please try again."
  - **Dependencies**: Task T003 (error handling mapping)
  - **Success**: Error conversion function implemented

- [ ] T009 Implement SDK client wrapper in `src/internal/geo/client.go`
  - **File**: `src/internal/geo/client.go`
  - **Details**:
    - Replace custom HTTP implementation with SDK client wrapper
    - Client struct with sdkClient field
    - `NewClient(httpClient *http.Client) *Client` constructor
    - Import SDK package
    - Initialize SDK client with httpClient if provided
    - Use default timeout (10s) if httpClient nil
  - **Dependencies**: Task T006 (SDK dependency added), T007 (mapping function)
  - **Success**: Client struct and constructor implemented

- [ ] T010 Implement Search method using SDK
  - **File**: `src/internal/geo/client.go`
  - **Details**:
    - Method signature: `(c *Client) Search(ctx context.Context, name string) ([]models.Location, error)`
    - Call SDK search method with context and name parameter
    - Map SDK response to []models.Location using mapSDKLocation()
    - Convert any errors using convertSDKError()
    - Return user-friendly error messages (no technical details)
    - Respect context cancellation and timeouts
  - **Dependencies**: Task T007 (mapping), T008 (error conversion), T009 (client struct)
  - **Success**: Search method fully functional with SDK

- [ ] T011 Run existing unit tests to verify compatibility
  - **File**: `src/internal/geo/client_test.go` (no changes to file)
  - **Details**:
    - Run: `go test ./src/internal/geo -v`
    - All existing tests must pass WITHOUT modification:
      - TestSearch (success case)
      - TestSearch (no results)
      - TestSearch (API error)
      - TestSearch (malformed JSON)
      - TestSearch_Timeout
    - Existing tests use mocks, should still work with new implementation
  - **Dependencies**: Task T010 (Search method implemented)
  - **Success**: All existing unit tests pass

- [ ] T012 [P] Run integration tests to validate API contract
  - **File**: `src/internal/geo/integration_test.go`
  - **Details**:
    - Run: `go test ./src/internal/geo -v -run Integration`
    - Run with network access (not with -short flag)
    - Verify all integration test cases pass
    - Tests validate: London search, multiple results, result limits, complete data
  - **Dependencies**: Task T005 (integration test created), T010 (Search method)
  - **Success**: Integration tests validate API contract stability

- [ ] T013 [P] Build and test the complete application
  - **File**: `bin/weather-reporter` (generated)
  - **Details**:
    - Run: `go build -o bin/weather-reporter ./src/cmd/weather-reporter`
    - Build should succeed with no errors
    - No code changes needed outside geo package
  - **Dependencies**: Task T010 (implementation complete), T006 (dependencies updated)
  - **Success**: Application builds successfully

- [ ] T014 [P] Manual integration test: verify location search works end-to-end
  - **File**: `bin/weather-reporter`
  - **Details**:
    - Test: `./bin/weather-reporter London`
    - Verify: Location resolved, weather displayed
    - Test: `./bin/weather-reporter "San Francisco"`
    - Verify: Works with multi-word location names
    - Test: `./bin/weather-reporter Tokyo`
    - Verify: Works for non-English locations
  - **Dependencies**: Task T013 (application built)
  - **Success**: Application works identically to before refactoring

- [ ] T015 [P] Verify performance: no degradation in response times
  - **File**: `bin/weather-reporter`
  - **Details**:
    - Measure response time: `time ./bin/weather-reporter London`
    - Compare with baseline (should be equal or faster, not slower)
    - Target: <1 second for typical queries
    - Requirement: Zero degradation (no slowdown acceptable)
  - **Dependencies**: Task T013 (application built)
  - **Success**: Performance verified, no degradation

- [ ] T016 [P] Verify error handling: user-friendly messages shown
  - **File**: `bin/weather-reporter`
  - **Details**:
    - Test offline (disable network): Should show user-friendly message
    - Test timeout (very short timeout): Should show user-friendly message
    - Verify: No technical error details exposed
    - Verify: Messages match requirement: "Unable to search locations. Please try again."
  - **Dependencies**: Task T010 (error handling implemented)
  - **Success**: Error messages are user-friendly

**Checkpoint**: SDK integration complete, all tests passing, no regressions

---

## Phase 3: Quality & Cleanup

**Duration**: 1 day  
**Purpose**: Final validation and code quality  
**Gate**: All quality checks pass, ready for production

- [ ] T017 [P] Run full test suite including all existing tests
  - **File**: `src/internal/geo/`
  - **Details**:
    - Run: `go test ./src/internal/geo -v`
    - Verify all unit tests pass
    - Verify all integration tests pass (if network available, or skip with -short)
    - Ensure 100% of existing tests still pass
  - **Success**: Full test suite passes

- [ ] T018 [P] Code review: verify implementation matches adapter contract
  - **File**: `specs/005-geocoding-sdk-integration/contracts/sdk-adapter.go`
  - **Details**:
    - Compare implementation against contract
    - Verify all functions documented
    - Verify error handling matches specification
    - Verify interface implementation is complete
  - **Success**: Implementation matches contract

- [ ] T019 [P] Run linting and formatting checks
  - **File**: `src/internal/geo/client.go`
  - **Details**:
    - Run: `go fmt ./src/internal/geo/`
    - Run: `go vet ./src/internal/geo/`
    - Fix any issues
  - **Success**: Code passes linting

- [ ] T020 Run final build verification
  - **File**: `bin/weather-reporter`
  - **Details**:
    - Run: `go build -o bin/weather-reporter ./src/cmd/weather-reporter`
    - Verify: Build succeeds
    - Verify: No warnings or errors
  - **Success**: Clean build

- [ ] T021 [P] Update documentation if needed
  - **Files**: `README.md`, `CONTRIBUTING.md`, `Makefile` (if applicable)
  - **Details**:
    - Check if any documentation mentions geocoding implementation
    - Update if references to custom HTTP client
    - Add note about SDK dependency if appropriate
    - Document any new testing approach (integration tests)
  - **Success**: Documentation is current

- [ ] T022 Create commit message documenting the refactoring
  - **File**: Feature branch changes
  - **Details**:
    - Summarize: Replaced custom HTTP client with open-meteo-geocoding-sdk
    - List: Files changed, tests added, benefits
    - Reference: Issue/spec if applicable
  - **Success**: Clear commit message for history

**Checkpoint**: All quality checks pass, ready for merge

---

## Phase 4: Documentation & Finalization

**Duration**: 0.5 day  
**Purpose**: Update design documents with implementation findings  
**Gate**: Documentation reflects actual implementation

- [ ] T023 [P] Update research.md with implementation findings
  - **File**: `specs/005-geocoding-sdk-integration/research.md`
  - **Details**:
    - Add actual SDK type names discovered
    - Add actual SDK method signatures used
    - Document any surprises or gotchas
    - Confirm type mapping worked as expected
    - Note any version compatibility issues found
  - **Success**: Research document reflects actual implementation

- [ ] T024 [P] Update data-model.md with actual mappings
  - **File**: `specs/005-geocoding-sdk-integration/data-model.md`
  - **Details**:
    - Update with actual SDK field names
    - Update with actual SDK error types
    - Document any differences from planned mapping
    - Verify Location struct unchanged
  - **Success**: Data model document is accurate

- [ ] T025 [P] Update quickstart.md with actual commands
  - **File**: `specs/005-geocoding-sdk-integration/quickstart.md`
  - **Details**:
    - Verify all steps work as documented
    - Update any version numbers
    - Confirm test commands match actual tests
    - Update any SDK package names if different
  - **Success**: Quickstart is accurate and tested

- [ ] T026 Mark plan.md as complete
  - **File**: `specs/005-geocoding-sdk-integration/plan.md`
  - **Details**:
    - Update status from "Draft" to "Complete"
    - Add completion date
    - Document any deviations from plan
    - Mark all phases as complete
  - **Success**: Plan marked complete with notes

**Checkpoint**: Documentation complete and accurate

---

## Phase 5: Delivery & Code Review

**Duration**: Varies  
**Purpose**: Prepare for merge and deployment  
**Gate**: Code review approved, ready for production

- [ ] T027 Commit changes to feature branch
  - **Details**:
    - All changes committed with clear messages
    - Branch ready for pull request
    - No uncommitted changes
  - **Success**: Branch ready for review

- [ ] T028 Create pull request with:
  - **Details**:
    - Title: "refactor: Replace custom geocoding client with SDK"
    - Description: Reference issue, list changes, note benefits
    - Link to spec: specs/005-geocoding-sdk-integration/spec.md
    - Checklist: All tests pass, no regressions, documented
  - **Success**: PR created for review

- [ ] T029 Address code review feedback
  - **Details**:
    - Respond to comments
    - Make requested changes
    - Re-run tests after changes
    - Push updates to PR
  - **Success**: Feedback addressed, PR approved

- [ ] T030 Merge to main branch
  - **Details**:
    - All checks pass
    - Approvals received
    - Squash or rebase as per project conventions
    - Delete feature branch after merge
  - **Success**: Changes merged to main

- [ ] T031 [P] Verify deployment/build
  - **Details**:
    - Confirm CI/CD pipeline passes
    - Verify application builds in CI environment
    - Confirm all tests run and pass in CI
  - **Success**: CI/CD pipeline green

- [ ] T032 [P] Smoke test on staging if applicable
  - **Details**:
    - Run application on staging environment
    - Verify location search still works
    - Verify error handling works
    - Verify performance acceptable
  - **Success**: Staging environment works

**Checkpoint**: Changes successfully merged and verified

---

## Dependencies & Execution Strategy

### Phase Dependencies

```
Phase 1 (Research)
    ↓
Phase 2 (Implementation) [Can parallelize T005-T016]
    ↓
Phase 3 (Quality) [Can parallelize T017-T019, T021]
    ↓
Phase 4 (Documentation) [Can parallelize T023-T026]
    ↓
Phase 5 (Code Review & Merge) [Sequential]
```

### Task Dependencies Within Phases

**Phase 1**: All tasks can run in parallel
- T001, T002, T003, T004 are all independent research activities

**Phase 2**: Some parallelization possible
- T005 (test file) can start immediately
- T006 (go.mod update) can start immediately
- T007-T008 (functions) depend on T005 research
- T009-T010 (implementation) depend on T006, T007, T008
- T011-T012 (testing) depend on T010 (implementation)
- T013-T016 (manual testing) can run in parallel, depend on T013 build

**Phase 3**: Most can parallelize
- T017 (tests) depends on T010 (implementation)
- T018-T019, T021 can run in parallel
- T020 (build) can run independently

### Parallel Opportunities

```
Parallel within Phase 2 (after research done):
├─ T005 (Integration test) ──┐
├─ T006 (go.mod) ────────────┼─→ T007-T008 (mapping/error)
└─ (Research outputs) ────┘        ↓
                              T009-T010 (Implementation)
                                    ↓
                    ┌─ T011 (Unit tests)
                    ├─ T012 (Integration tests) ──┐
                    ├─ T013 (Build)                │
                    ├─ T014 (Manual test)          ├─ Parallel
                    ├─ T015 (Performance)          │
                    └─ T016 (Error handling) ──┘

Parallel within Phase 3:
├─ T017 (Run tests) ──┐
├─ T018 (Code review) │
├─ T019 (Linting)     ├─ Parallel
├─ T021 (Docs)        │
└─ T020 (Build) ──┘
```

### Minimum Timeline (Parallel)

- **Phase 1**: 1 day (all research tasks in parallel)
- **Phase 2**: 2-3 days (research → mapping/errors → implementation → tests)
- **Phase 3**: 0.5 days (quality tasks in parallel)
- **Phase 4**: 0.5 days (doc updates in parallel)
- **Phase 5**: 1-2 days (depends on review feedback)

**Total**: 5-8 days end-to-end (varies by review feedback)

### Recommended Team Allocation

**Single Developer**: Execute phases sequentially (5-8 days)
**Two Developers**: 
- Developer 1: T001, T007, T009, T010, T020
- Developer 2: T002-T004, T005-T006, T011-T016
- Both: T017-T032 in parallel where possible

---

## Success Criteria & Verification

### Phase 1 Success
- [ ] All research questions answered
- [ ] Type mapping strategy defined and validated
- [ ] Error handling strategy documented
- [ ] Integration test approach finalized

### Phase 2 Success
- [ ] All existing unit tests pass without modification
- [ ] New integration tests pass and validate API contract
- [ ] Zero performance degradation verified
- [ ] Error messages are user-friendly
- [ ] Application builds successfully
- [ ] Manual testing confirms functionality

### Phase 3 Success
- [ ] Full test suite passes
- [ ] Code review approved
- [ ] Linting passes
- [ ] Build succeeds cleanly
- [ ] Documentation updated if needed

### Phase 4 Success
- [ ] Research.md updated with findings
- [ ] Data-model.md reflects actual implementation
- [ ] Quickstart.md verified accurate
- [ ] Plan.md marked complete

### Phase 5 Success
- [ ] Pull request approved
- [ ] Changes merged to main
- [ ] CI/CD pipeline passes
- [ ] Ready for production deployment

---

## Rollback Plan

If critical issues discovered after implementation:

1. Keep feature branch available for reference
2. All previous behavior preserved by existing tests
3. Can revert to main if needed (pre-SDK version still available)
4. Integration test will alert to API contract changes

---

## Notes

- All file paths are relative to repository root (`/workspaces/weather-reporter/`)
- Existing tests in `src/internal/geo/client_test.go` should NOT be modified
- New integration test in `src/internal/geo/integration_test.go` documents API contract
- SDK dependency version should be pinned (not @latest)
- All error messages must be user-friendly per specification
- Zero performance degradation is required - not negotiable
