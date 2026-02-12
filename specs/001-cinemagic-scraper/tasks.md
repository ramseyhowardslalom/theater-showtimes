# Tasks: Cinemagic Theater Support

**Input**: Design documents from `/specs/001-cinemagic-scraper/`
**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/theater-cinemagic.json](contracts/theater-cinemagic.json)

**Tests**: No test tasks included - tests not explicitly requested in feature specification. Constitution principle I (TDD) allows tests to be written as implementation progresses.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This is a **Web Application** project with:
- Backend: `backend/` (Go language)
- Frontend: `frontend/` (React/TypeScript)
- Documentation: `docs/`
- Specifications: `specs/001-cinemagic-scraper/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 [P] Create Cinemagic Theater entry in backend/data/theaters.json per contracts/theater-cinemagic.json
- [X] T002 [P] Create scraper specification document docs/cinemagic-theater-specs.md
- [X] T003 Create scraper module directory backend/internal/scrapers/cinemagic_theater/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core scraper infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Implement CinemagicScraper struct with Colly collector setup in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T005 [P] Implement title normalization function in backend/internal/scrapers/cinemagic_theater/helpers.go
- [X] T006 [P] Implement film format detection function in backend/internal/scrapers/cinemagic_theater/helpers.go
- [X] T007 Verify ReplaceTheaterShowtimes method exists in backend/internal/storage/storage.go for full refresh strategy
- [X] T008 Register Cinemagic scraper in backend/cmd/scraper/main.go CLI

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - View Cinemagic Showtimes (Priority: P1) 🎯 MVP

**Goal**: Users can browse Cinemagic Theater showtimes alongside existing theaters, seeing movie posters, times, and details enriched with TMDB data

**Independent Test**: Run the scraper for Cinemagic (`go run cmd/scraper/main.go --theater=cinemagic-theater`), start frontend dev server, verify Cinemagic movies appear in movie listings with posters, showtimes, and theater labels

### Implementation for User Story 1

- [X] T009 [US1] Implement movie list extraction from https://tickets.thecinemagictheater.com/now-showing in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T010 [US1] Implement calendar navigation logic to access 3-month showtime window in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T011 [US1] Implement individual movie page scraping for showtime details in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T012 [US1] Implement showtime data extraction (date, time, format) per research.md patterns in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T013 [US1] Integrate TMDB enrichment for scraped movie titles using existing backend/internal/tmdb/client.go
- [X] T014 [US1] Implement TMDB failure handling with limited_info flag per data-model.md in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T015 [US1] Add error handling and logging for scraper failures in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T016 [US1] Implement 3-second rate limiting with Colly LimitRule per research.md in backend/internal/scrapers/cinemagic_theater/scraper.go

**Checkpoint**: At this point, User Story 1 should be fully functional - Cinemagic movies display in UI with TMDB data

---

## Phase 4: User Story 2 - Filter by Cinemagic Theater (Priority: P2)

**Goal**: Users can filter the movie listings to show only movies playing at Cinemagic Theater using the theater dropdown selector

**Independent Test**: Start frontend, open theater dropdown, select "Cinemagic Theater", verify only Cinemagic movies are displayed

### Implementation for User Story 2

**✅ NO TASKS REQUIRED** - This user story is automatically satisfied by the data-driven UI architecture. Once T001 adds Cinemagic to `backend/data/theaters.json` and the API returns theater data, the existing `frontend/src/components/FilterBar/FilterBar.jsx` component automatically includes Cinemagic in the dropdown without code changes.

**Validation**: After Phase 3 completion, verify FilterBar dropdown includes "Cinemagic Theater" option

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Access Cinemagic Ticket Links (Priority: P3)

**Goal**: Users can click a "Tickets" button on movie cards to navigate directly to the Cinemagic event page for purchasing tickets

**Independent Test**: View a movie showing at Cinemagic, click the "Tickets" button, verify it opens the correct Cinemagic event page (https://tickets.thecinemagictheater.com/movie/{movie-slug})

### Implementation for User Story 3

- [X] T017 [US3] Extract event page links during showtime scraping in backend/internal/scrapers/cinemagic_theater/scraper.go
- [X] T018 [US3] Store event page links in showtime.link field per data-model.md in backend/internal/scrapers/cinemagic_theater/scraper.go

**✅ FRONTEND AUTO-COMPLETE** - The existing `frontend/src/components/MovieCard/MovieCard.jsx` component automatically displays "Tickets" buttons when showtime records include a `link` field. No frontend code changes required.

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and final validation

- [X] T019 [P] Update docs/cinemagic-theater-specs.md with final scraper implementation details
- [X] T020 Run end-to-end validation per specs/001-cinemagic-scraper/quickstart.md Steps 6-10
- [X] T021 Performance testing: Verify scraper completes 3-month scrape in under 3 minutes per success criteria SC-005

**Note**: T020 and T021 require live website testing. See [VALIDATION.md](VALIDATION.md) for detailed validation procedures.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (T003 directory creation) - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion (T004-T008)
- **User Story 2 (Phase 4)**: Depends on Phase 1 (T001) only - automatically satisfied by data-driven UI
- **User Story 3 (Phase 5)**: Depends on Foundational phase completion (T004-T008)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Automatically satisfied once theater entry exists (T001) - No implementation needed
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Extends User Story 1 scraping logic but independently testable

### Within Each User Story

**Phase 3 (User Story 1) Task Order**:
1. T009 (movie list) → T010 (calendar nav) → T011 (movie pages) must be sequential
2. T012 (showtime extraction) depends on T011
3. T013 (TMDB enrichment) can run in parallel with T012
4. T014 (TMDB failures) depends on T013
5. T015 (error handling) can be implemented alongside any task
6. T016 (rate limiting) should be added to T004 collector setup

**Phase 5 (User Story 3) Task Order**:
1. T017 and T018 can be combined (extract and store links together)
2. Both integrate into existing scraping logic from Phase 3

### Parallel Opportunities

**Phase 1: Setup Tasks**
```bash
# All Setup tasks can run in parallel (different files):
T001: backend/data/theaters.json
T002: docs/cinemagic-theater-specs.md
T003: Create directory (fast, no conflicts)
```

**Phase 2: Foundational Tasks**
```bash
# These can run in parallel (different files):
T005: backend/internal/scrapers/cinemagic_theater/helpers.go (normalization)
T006: backend/internal/scrapers/cinemagic_theater/helpers.go (format detection)
# Note: Both in same file, but can be implemented as separate functions independently

# Sequential dependencies:
T004 must complete before T008 (register scraper)
T007 is verification, can be done anytime
```

**Phase 3: User Story 1**
- Most tasks are sequential (scraping is inherently sequential)
- T013 (TMDB enrichment) integrates existing code - can be planned while T009-T012 progress
- T015 (error handling) can be added incrementally to any function
- T016 (rate limiting) should be part of T004 setup

**Phase 5: User Story 3**
- T017 and T018 are closely related and should be implemented together

**Phase 6: Polish**
- T019 can be done in parallel with T020 and T021

---

## Parallel Example: Phase 1 (Setup)

```bash
# Launch all Setup tasks together:
Developer A: T001 - Edit backend/data/theaters.json
Developer B: T002 - Create docs/cinemagic-theater-specs.md
Developer C: T003 - Create backend/internal/scrapers/cinemagic_theater/ directory
```

All three tasks operate on different files with no conflicts.

---

## Parallel Example: Phase 2 (Foundational)

```bash
# Parallel setup of helper functions:
Developer A: T004 - Implement CinemagicScraper struct in scraper.go
Developer B: T005 - Implement normalizeTitleForTMDB in helpers.go
Developer C: T006 - Implement extractFilmFormat in helpers.go

# After T004 completes:
Developer A: T008 - Register scraper in cmd/scraper/main.go

# Anytime (verification only):
Developer D: T007 - Verify ReplaceTheaterShowtimes exists
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. **Complete Phase 1: Setup** (T001-T003) → ~30 minutes
2. **Complete Phase 2: Foundational** (T004-T008) → ~2-3 hours (CRITICAL - blocks all stories)
3. **Complete Phase 3: User Story 1** (T009-T016) → ~4-6 hours
4. **STOP and VALIDATE**: 
   - Run scraper: `cd backend && go run cmd/scraper/main.go --theater=cinemagic-theater`
   - Start frontend: `cd frontend && npm run dev`
   - Verify Cinemagic movies appear with posters and showtimes
5. **Deploy/Demo** if ready - this is a complete MVP!

**Total MVP Time**: ~8-10 hours of development

### Incremental Delivery

1. **Complete Setup + Foundational** → Foundation ready (Phases 1-2)
2. **Add User Story 1** → Test independently → Deploy/Demo (MVP! 🎯)
3. **Verify User Story 2** → Already works due to data-driven UI → Deploy/Demo
4. **Add User Story 3** → Test independently → Deploy/Demo
5. **Polish Phase** → Final validation and documentation

Each increment adds value without breaking previous functionality.

### Parallel Team Strategy

With multiple developers:

1. **Team completes Setup together** (30 min)
2. **Team completes Foundational together** (2-3 hours - must be done collaboratively)
3. **Once Foundational is done**:
   - Developer A: User Story 1 (Phase 3) - Core scraping implementation
   - Developer B: User Story 3 (Phase 5) - Ticket links (extends scraping)
   - Developer C: Documentation (Phase 6 - T019)
4. **Integration Point**: Merge and test all stories together
5. **Final Validation**: Phase 6 (T020-T021)

**Note**: User Story 2 requires no development - validate it works once User Story 1 is complete.

---

## Notes

- **[P] tasks**: Different files, no dependencies, can run in parallel
- **[Story] labels**: 
  - [US1] = User Story 1 (View Cinemagic Showtimes)
  - [US2] = User Story 2 (Filter by Theater) - Auto-complete via data-driven UI
  - [US3] = User Story 3 (Ticket Links)
- **No frontend code changes needed**: All UI components are data-driven and automatically adapt to new theater data
- **Tests**: Not included as optional per constitution - can be added during implementation if desired (TDD approach)
- **Data-driven architecture benefit**: User Stories 2 and 3 require minimal backend work; frontend automatically handles new data
- **Rate limiting critical**: 3-second delays prevent IP blocking; scraper will take ~3 minutes for full run
- **TMDB enrichment**: 90% target match rate per success criteria SC-002
- **Independent validation**: Each user story has a clear test procedure in its "Independent Test" section
- **Commit strategy**: Commit after each task or logical group of related tasks
- **Stop at checkpoints**: Validate each user story independently before proceeding to next priority

---

## Success Metrics Mapping

These tasks directly satisfy the success criteria from spec.md:

- **SC-001** (10+ unique movies): Satisfied by T009-T012 (scraping implementation)
- **SC-002** (90% TMDB match): Satisfied by T013 (TMDB enrichment)
- **SC-003** (<1s filter response): Satisfied by T001 (data-driven UI, no code needed)
- **SC-004** (100% correct ticket links): Satisfied by T017-T018 (link extraction)
- **SC-005** (<3 min scrape time): Validated by T021 (performance testing)
- **SC-006** (No frontend changes needed): Validated throughout - no frontend tasks exist!

---

## Ready to Begin

**Next Step**: Start with Phase 1, Task T001 (create theater entry)

**Reference**: See [quickstart.md](quickstart.md) for detailed implementation guidance and code examples

**Branch**: `001-cinemagic-scraper`
