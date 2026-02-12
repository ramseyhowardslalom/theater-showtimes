# Implementation Plan: Cinemagic Theater Support

**Branch**: `001-cinemagic-scraper` | **Date**: 2026-02-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-cinemagic-scraper/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Add Cinemagic Theater to the Portland Movie Theater Showtimes application by implementing a modular web scraper that extracts showtimes from https://tickets.thecinemagictheater.com, enriches data with TMDB metadata, and integrates seamlessly into the existing data-driven UI. The scraper will follow the established architecture pattern used by Clinton Street Theater, extracting 3 months of showtimes with film format detection, rate limiting, and graceful error handling.

## Technical Context

**Language/Version**: Go 1.21 (backend), TypeScript 5.3 (frontend), React 18.2  
**Primary Dependencies**: Colly v2 (web scraping), Gin 1.9 (HTTP server), React Router 6.20, Axios 1.6, Vite 5.0 (build tool)  
**Storage**: JSON file-based persistence (backend/data/theaters.json, backend/data/showtimes.json, backend/data/movies.json)  
**Testing**: Go built-in testing package (backend), Vitest 4.0 + React Testing Library (frontend)  
**Target Platform**: Linux/macOS server (backend), Modern web browsers (frontend)  
**Project Type**: Web application (separate backend Go service + frontend React SPA)  
**Performance Goals**: <3 minute scrape time, <500ms API response, <1s filter response, 90% TMDB match rate  
**Constraints**: 3-second rate limiting between requests, on-demand execution only, full data refresh per run  
**Scale/Scope**: Support for 3+ theaters, 100+ unique movies per month, 3-month showtime window

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Applicable Principles from Constitution v1.0.0

**✅ I. Test-Driven Development (NON-NEGOTIABLE)**
- Status: COMPLIANT
- Plan: Write tests first for scraper extraction logic, TMDB enrichment, and data storage
- Rationale: Scraper logic is complex (parsing HTML, calendar navigation, format detection) and benefits from TDD

**✅ II. TypeScript-First Development**
- Status: COMPLIANT
- Plan: Frontend already uses TypeScript 5.3 with strict mode; no new frontend code required (data-driven UI)
- Rationale: This feature adds backend functionality only; existing TypeScript components work without modification

**✅ III. Modular Architecture**
- Status: COMPLIANT
- Plan: Create `backend/internal/scrapers/cinemagic_theater/scraper.go` as independent module
- Rationale: Follows existing pattern from Clinton Street Theater scraper; enables independent testing and maintenance

**✅ IV. Code Quality & Tooling**
- Status: COMPLIANT
- Plan: Go code will pass `go vet` and `golint`; use table-driven tests for multiple scenarios
- Rationale: Standard Go tooling ensures code quality; table-driven tests cover edge cases efficiently

**✅ V. Performance & Caching**
- Status: COMPLIANT
- Plan: 3-second rate limiting enforced, TMDB data cached for 7 days, scraper completes in <3 minutes
- Rationale: Conservative rate limiting protects Cinemagic's servers; TMDB caching reduces API costs

**✅ VI. Graceful Degradation**
- Status: COMPLIANT
- Plan: Scraper failures logged and skipped, TMDB match failures show placeholder with badge, cached data used when available
- Rationale: External dependencies (Cinemagic website, TMDB API) can fail; system remains functional

**✅ VII. Security & Secrets Management**
- Status: COMPLIANT
- Plan: No new secrets required; existing TMDB API key already managed via environment variables
- Rationale: Cinemagic website is publicly accessible (no authentication needed)

### Constitution Compliance Summary

**Overall Status**: ✅ COMPLIANT - No violations

All seven constitution principles are satisfied. This feature fits cleanly within the existing architecture and requires no exceptions or justifications.

## Project Structure

### Documentation (this feature)

```text
specs/001-cinemagic-scraper/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── theater-cinemagic.json
├── checklists/          # Quality validation checklists
│   └── requirements.md
└── spec.md              # Feature specification
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   ├── api/
│   │   └── main.go               # API server (no changes needed)
│   └── scraper/
│       └── main.go               # CLI scraper (register Cinemagic)
├── internal/
│   ├── api/
│   │   ├── handlers.go           # API handlers (no changes needed - data-driven)
│   │   └── router.go             # API routes (no changes needed)
│   ├── models/
│   │   └── models.go             # Data models (Theater, Showtime, Movie - existing)
│   ├── scrapers/
│   │   ├── scraper.go            # Base scraper interface (existing)
│   │   ├── clinton_street_theater/
│   │   │   └── scraper.go        # Clinton Street Theater (reference implementation)
│   │   └── cinemagic_theater/    # 🆕 NEW MODULE
│   │       └── scraper.go        # Cinemagic scraper implementation
│   ├── storage/
│   │   └── storage.go            # JSON storage (existing, may need update for full refresh)
│   └── tmdb/
│       ├── client.go             # TMDB client (existing)
│       └── cache.go              # TMDB cache (existing)
├── configs/
│   └── config.yaml               # Configuration (add Cinemagic theater entry)
├── data/
│   ├── theaters.json             # Theater data (add Cinemagic entry)
│   ├── showtimes.json            # Showtime data (populated by scraper)
│   └── movies.json               # Movie metadata (TMDB-enriched)
└── go.mod                        # Dependencies (no new dependencies)

frontend/
├── src/
│   ├── components/
│   │   ├── FilterBar/
│   │   │   └── FilterBar.jsx     # Theater dropdown (NO CHANGES - data-driven)
│   │   └── MovieCard/
│   │       └── MovieCard.jsx     # Movie display (NO CHANGES - data-driven)
│   ├── pages/
│   │   └── Home.jsx              # Homepage (NO CHANGES - data-driven)
│   └── services/
│       └── api.js                # API client (NO CHANGES)
└── tests/                        # Test files (manual verification tests)

docs/
├── cinemagic-theater-specs.md    # 🆕 NEW - Scraper specifications
└── clinton-street-theater-specs.md  # Reference documentation
```

**Structure Decision**: Web application pattern (Option 2) selected automatically as project has both `backend/` and `frontend/` directories. This feature adds a new scraper module (`backend/internal/scrapers/cinemagic_theater/`) and documentation (`docs/cinemagic-theater-specs.md`). No frontend code changes are required due to data-driven architecture - UI components automatically adapt to new theater data from the API.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A - No constitution violations detected.

---

## Phase 0: Research Complete ✅

See [research.md](research.md) for complete research findings.

**Research Tasks Completed**:
1. Website Structure Analysis
2. Colly Framework Best Practices
3. Film Format Detection Strategy
4. Title Normalization for TMDB
5. Data Refresh Strategy
6. Rate Limiting Configuration
7. TMDB Match Failure Handling

**All NEEDS CLARIFICATION resolved** - Ready for Phase 1 design.

---

## Phase 1: Design & Contracts Complete ✅

**Generated Artifacts**:
- [data-model.md](data-model.md) - 3 entities with validation rules
- [contracts/theater-cinemagic.json](contracts/theater-cinemagic.json) - JSON schema
- [quickstart.md](quickstart.md) - 10-step developer guide
- [.github/agents/copilot-instructions.md](../../.github/agents/copilot-instructions.md) - Agent context updated

**Design Decisions**:
- Theater ID: `cinemagic-theater`
- Address: `2021 SE Hawthorne Blvd, Portland, OR 97214`
- Format detection: Multi-selector with "digital" default
- Title normalization: Strip parenthetical suffixes before TMDB match
- Storage method: `ReplaceTheaterShowtimes()` for full refresh
- TMDB failure handling: Placeholder movie with `limited_info: true` badge

---

## Post-Design Constitution Check ✅

*Re-evaluation after Phase 1 design completion*

### Updated Compliance Status

**✅ I. Test-Driven Development (NON-NEGOTIABLE)**
- Status: COMPLIANT
- Design Impact: `quickstart.md` Step 6-9 include comprehensive testing procedures (unit tests, manual verification, edge cases)
- Evidence: Test structure defined in quickstart guide with specific scenarios

**✅ II. TypeScript-First Development**
- Status: COMPLIANT
- Design Impact: No frontend code changes required (confirmed in design phase)
- Evidence: Existing TypeScript components handle new data without modification

**✅ III. Modular Architecture**
- Status: COMPLIANT
- Design Impact: `data-model.md` confirms clean separation between Theater, Showtime, Movie entities
- Evidence: Single-responsibility principle maintained in `cinemagic_theater/scraper.go` module

**✅ IV. Code Quality & Tooling**
- Status: COMPLIANT
- Design Impact: `quickstart.md` Step 5 includes implementation checklist with Go best practices
- Evidence: Code structure follows existing patterns from Clinton Street Theater reference

**✅ V. Performance & Caching**
- Status: COMPLIANT
- Design Impact: `research.md` confirms 3-second rate limiting with Colly `Limit` rules
- Evidence: Performance targets maintained (<3 min scrape, existing TMDB cache reused)

**✅ VI. Graceful Degradation**
- Status: COMPLIANT
- Design Impact: `data-model.md` adds `limited_info: true` flag for TMDB match failures
- Evidence: System continues functioning with placeholder data + badge when enrichment fails

**✅ VII. Security & Secrets Management**
- Status: COMPLIANT
- Design Impact: No new secrets introduced (confirmed in design)
- Evidence: Uses existing TMDB API key management pattern

### Post-Design Compliance Summary

**Overall Status**: ✅ COMPLIANT - No new violations introduced

Design phase reinforced initial compliance assessment. All seven principles remain satisfied with concrete implementation details now defined in design artifacts.

---

## Next Steps

**Command Complete** - `/speckit.plan` workflow finished.

**Branch**: `001-cinemagic-scraper`  
**Plan Document**: `/private/tmp/theater-showtimes/specs/001-cinemagic-scraper/plan.md`  
**Generated Artifacts**:
- ✅ research.md (Phase 0)
- ✅ data-model.md (Phase 1)
- ✅ contracts/theater-cinemagic.json (Phase 1)
- ✅ quickstart.md (Phase 1)
- ✅ .github/agents/copilot-instructions.md (Agent context)

**Ready for Implementation** - Run `/speckit.tasks` to generate task breakdown.
