# Feature Specification: Cinemagic Theater Support

**Feature Branch**: `001-cinemagic-scraper`  
**Created**: February 11, 2026  
**Status**: Draft  
**Input**: User description: "Support Cinemagic theater listings including UI dropdown integration and scraper implementation"

## Clarifications

### Session 2026-02-11

- Q: TMDB Match Failure Display Strategy - How should unmatched movies appear in the UI? → A: Display unmatched movies but mark them as "Limited Information Available" with a badge
- Q: Scraper Execution Trigger - When/how should the Cinemagic scraper be triggered? → A: On demand.
- Q: Cinemagic Theater Physical Address - What is the street address for theaters.json? → A: address is correct
- Q: Rate Limiting Precision - Should scraper wait 2 or 3 seconds between requests? → A: 3 seconds (conservative)
- Q: Scraper Re-run Data Handling - How should existing Cinemagic data be handled on subsequent scraper runs? → A: Replace all Cinemagic data each run (full refresh)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Cinemagic Showtimes in Listings (Priority: P1)

Users can browse Cinemagic Theater showtimes alongside existing theaters, seeing movie posters, times, and details enriched with TMDB data.

**Why this priority**: This is the core value proposition - expanding the theater coverage to include Cinemagic. Without this, the scraper has no visible user benefit.

**Independent Test**: Run the scraper for Cinemagic, load the frontend, and verify Cinemagic movies appear in the movie listings with posters, showtimes, and theater labels.

**Acceptance Scenarios**:

1. **Given** Cinemagic scraper has run successfully, **When** user visits the homepage, **Then** movies showing at Cinemagic appear in the movie grid with TMDB posters and metadata
2. **Given** movies from multiple theaters are showing, **When** user views a movie card, **Then** Cinemagic appears as a theater label if that movie is showing there
3. **Given** user is viewing movie listings, **When** user clicks on a Cinemagic showtime, **Then** user sees all available showtimes for that movie at Cinemagic

---

### User Story 2 - Filter by Cinemagic Theater (Priority: P2)

Users can filter the movie listings to show only movies playing at Cinemagic Theater using the theater dropdown selector.

**Why this priority**: Filtering enables users to focus on a specific theater when planning their movie experience. This is the second most important feature after basic display.

**Independent Test**: Select "Cinemagic Theater" from the theater dropdown and verify only Cinemagic movies are displayed.

**Acceptance Scenarios**:

1. **Given** user is on the homepage, **When** user opens the theater dropdown, **Then** "Cinemagic Theater" appears as an option alongside other theaters
2. **Given** user selects "Cinemagic Theater" from dropdown, **When** the page updates, **Then** only movies showing at Cinemagic are displayed
3. **Given** user has filtered to Cinemagic, **When** user selects "All Theaters", **Then** movies from all theaters (including Cinemagic) are displayed again

---

### User Story 3 - Access Cinemagic Ticket Links (Priority: P3)

Users can click a "Tickets" button on movie cards to navigate directly to the Cinemagic event page for purchasing tickets.

**Why this priority**: While important for conversion, users can still browse and discover movies without this. It's a convenience feature that comes after core browsing and filtering.

**Independent Test**: View a movie showing at Cinemagic and click the "Tickets" button to verify it opens the correct Cinemagic event page.

**Acceptance Scenarios**:

1. **Given** a movie is showing at Cinemagic, **When** user views the movie card, **Then** a "Tickets" button is visible
2. **Given** user clicks the "Tickets" button, **When** the link opens, **Then** it navigates to the correct Cinemagic event page (https://tickets.thecinemagictheater.com/movie/{movie-slug})
3. **Given** user is on the Cinemagic event page, **When** user views showtimes, **Then** they match the times displayed in the app

---

### Edge Cases

- What happens when Cinemagic website is unreachable during scraping? (Scraper should log error, skip Cinemagic, continue with other theaters)
- What happens when a Cinemagic movie has no TMDB match? (Display movie in main grid with placeholder poster, basic info (title, times, theater), and "Limited Information Available" badge)
- What happens when Cinemagic changes their website structure? (Scraper fails gracefully, logs errors, continues with cached data if available)
- What happens when the same movie shows at both Cinemagic and Clinton Street? (Both theater labels appear on the movie card)
- What happens when navigating through multiple months on Cinemagic calendar? (Scraper should handle calendar navigation for 2+ months of future showtimes)
- What happens when a film format badge is missing from a Cinemagic movie page? (Default to "digital" format)

## Requirements *(mandatory)*

### Functional Requirements

#### Backend Scraper Requirements

- **FR-001**: System MUST scrape showtimes from Cinemagic Theater website (https://tickets.thecinemagictheater.com/now-showing)
- **FR-002**: System MUST extract movie title, date, time, and format for each showing
- **FR-003**: System MUST scrape current month plus next 2 months of showtimes (3 months total)
- **FR-004**: System MUST navigate through Cinemagic's calendar component to access future dates
- **FR-005**: System MUST extract event page links for each movie (e.g., https://tickets.thecinemagictheater.com/movie/sentimental-value)
- **FR-006**: System MUST detect film format from movie page badges (digital, 35mm, 70mm, etc.) or default to "digital"
- **FR-007**: System MUST store Cinemagic theater information in theaters.json with id "cinemagic-theater"
- **FR-008**: System MUST enrich scraped titles with TMDB data (posters, ratings, overview, runtime)
- **FR-009**: System MUST respect rate limiting (3 seconds between requests to Cinemagic domain)
- **FR-010**: System MUST follow the modular scraper architecture pattern (separate module in internal/scrapers/cinemagic_theater/)
- **FR-019**: Scraper MUST be executable on-demand via API endpoint and CLI command (no automatic scheduling)
- **FR-020**: Scraper MUST replace all existing Cinemagic showtimes on each run (full refresh, no incremental updates)

#### Frontend UI Requirements

- **FR-011**: Theater dropdown MUST include "Cinemagic Theater" as a selectable option
- **FR-012**: Movie cards MUST display "Cinemagic Theater" label when movie is showing there
- **FR-013**: Filtering by "Cinemagic Theater" MUST show only Cinemagic movies
- **FR-014**: Movie cards for Cinemagic movies MUST display "Tickets" button linking to event page
- **FR-015**: UI components MUST automatically include Cinemagic without code changes (data-driven from API)
- **FR-018**: Movie cards for movies without TMDB matches MUST display a "Limited Information Available" badge

#### Documentation Requirements

- **FR-016**: System MUST include cinemagic-theater-specs.md documentation file following Clinton Street Theater pattern
- **FR-017**: Documentation MUST specify theater ID, website URLs, scraping URLs, and data extraction patterns

### Key Entities *(include if feature involves data)*

- **Cinemagic Theater**: Theater entity with id "cinemagic-theater", name "Cinemagic Theater", address "2021 SE Hawthorne Blvd, Portland, OR 97214", website URL (https://www.thecinemagictheater.com)
- **Cinemagic Showtime**: Showtime entry linking movie (via tmdb_id), date, time, theater_id ("cinemagic-theater"), format (digital/35mm/etc), and event page link
- **Cinemagic Movie Data**: TMDB-enriched movie metadata for each unique movie showing at Cinemagic (poster, overview, rating, runtime, genres)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Scraper successfully extracts at least 10 unique movies from Cinemagic within 3-month window
- **SC-002**: 90% of Cinemagic movies match to TMDB and display with poster images
- **SC-003**: Users can filter to Cinemagic Theater and see only Cinemagic movies in under 1 second
- **SC-004**: Clicking "Tickets" button navigates to correct Cinemagic event page 100% of the time
- **SC-005**: Scraper completes full 3-month scrape of Cinemagic in under 3 minutes (accounting for 3-second rate limiting)
- **SC-006**: Theater dropdown includes Cinemagic without requiring frontend code changes (proves data-driven architecture)

## Assumptions

- Cinemagic theater website structure remains consistent with current calendar-based navigation
- Cinemagic event pages include film format badges or indicators
- TMDB API will successfully match most mainstream movies from Cinemagic
- Existing FilterBar and MovieCard components are data-driven and require no modifications
- Backend API already returns all theaters and showtimes dynamically
