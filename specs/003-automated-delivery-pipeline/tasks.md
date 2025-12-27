# Tasks: Automated Delivery Pipeline

**Feature Branch**: `003-automated-delivery-pipeline`
**Spec**: [specs/003-automated-delivery-pipeline/spec.md](spec.md)
**Plan**: [specs/003-automated-delivery-pipeline/plan.md](plan.md)

## Phase 1: Setup
*Goal: Initialize configuration files and project structure for the pipeline.*

- [ ] T001 Create `.github/workflows/` directory
- [ ] T002 Create `.golangci.yml` with standard linters enabled
- [ ] T003 Create `Makefile` with `test`, `lint`, `build`, `snapshot`, `clean` targets
- [ ] T004 Update `src/cmd/weather-reporter/main.go` to handle version flags and variables

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 Add version variables (`Version`, `Commit`, `Date`) to `src/cmd/weather/main.go`
- [ ] T006 Implement `--version` flag handling in `src/cmd/weather/main.go` per contract

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

*Goal: Establish the core CI/CD infrastructure.*

## Phase 3: User Story 1 - Developer Quality Gate
*Goal: Ensure code quality is enforced on every Pull Request.*

- [ ] T007 [US1] Create `.github/workflows/ci.yml` with Quality job (Lint + Test)
- [ ] T008 [US1] Add testing job (`go test -race`) to `.github/workflows/ci.yml`
- [ ] T009 [US1] Configure caching (`actions/setup-go`) and fail-fast strategy in `.github/workflows/ci.yml`
- [ ] T010 [US1] Add coverage regression check to `.github/workflows/ci.yml`
- [ ] T011 [US1] Verify CI pipeline fails on lint errors (Manual Test)
- [ ] T012 [US1] Verify CI pipeline fails on test failures (Manual Test)
- [ ] T013 [US1] Verify CI pipeline passes on clean code (Manual Test)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

## Phase 4: User Story 2 - Automated Release Publishing
*Goal: Automate cross-platform builds and releases via semantic version tags.*

- [ ] T014 [US2] Create `.goreleaser.yaml` with cross-compilation matrix and ldflags
- [ ] T015 [US2] Configure Syft anchore/sbom-action in Release job
- [ ] T016 [US2] Configure checksums in `.goreleaser.yaml`
- [ ] T017 [US2] Create `.github/workflows/release.yml` triggered by `v*` tags
- [ ] T018 [US2] Configure GoReleaser action and permissions in `.github/workflows/release.yml`
- [ ] T019 [US2] Verify GoReleaser generates cross-platform binaries (Snapshot Test)
- [ ] T020 [US2] Verify SBOM generation (Snapshot Test)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

## Phase 5: User Story 3 - End User Installation
*Goal: Ensure users can verify and install the software.*

- [ ] T021 [US3] Update `README.md` with installation and verification instructions
- [ ] T022 [US3] Verify artifact checksums match generated files (Manual Test)

## Final Phase: Polish
*Goal: Ensure compliance with all standards and final cleanup.*

- [ ] T023 Verify build reproducibility and version metadata embedding
- [ ] T024 Update `README.md` with CI status badges

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P3)**: Can start after User Stories 1 & 2 - Depends on artifacts produced by US2

## Implementation Strategy
We will start by setting up the local tools (Makefile, Linter) to ensure the developer environment is ready. Then we will implement the CI pipeline to enforce quality gates. Finally, we will configure GoReleaser and the Release job to handle artifact publication.
