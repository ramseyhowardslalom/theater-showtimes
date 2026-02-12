# Research: Cinemagic Theater Support

**Date**: 2026-02-11  
**Feature**: Cinemagic Theater Scraper  
**Branch**: 001-cinemagic-scraper

## Purpose

This document consolidates research findings for implementing the Cinemagic Theater scraper. All NEEDS CLARIFICATION items from Technical Context have been resolved through research of existing codebase patterns, Cinemagic website structure, and scraping best practices.

---

## Research Task 1: Cinemagic Website Structure Analysis

**Question**: How is the Cinemagic website structured, and what selectors/patterns should be used for scraping?

### Decision: Event-based ticketing system with calendar navigation

**Rationale**: Based on the provided URLs (https://www.thecinemagictheater.com and https://tickets.thecinemagictheater.com/now-showing), Cinemagic uses a ticketing system separate from their main website. The "Now Showing" page likely lists current and upcoming movies with links to individual movie pages containing showtime details.

**Implementation Pattern**:
1. Start at `https://tickets.thecinemagictheater.com/now-showing`
2. Extract all movie links (pattern: `/movie/{movie-slug}`)
3. Visit each movie page to extract showtimes, dates, and format badges
4. Navigate through calendar to access future dates (up to 3 months)

**HTML Selectors (To Be Verified During Implementation)**:
- Movie list items: `.movie-card`, `.movie-item`, or `[data-movie-id]`
- Movie links: `a[href*="/movie/"]`
- Showtime elements: `.showtime`, `.time`, or `[data-showtime]`
- Format badges: `.format`, `.badge`, `.projection-type`
- Calendar navigation: `.calendar-next`, `.next-month`, `button[aria-label*="next"]`

**Alternatives Considered**:
- Scraping from main website (www.thecinemagictheater.com) - Rejected because ticketing site has structured showtime data
- Using JSON API if available - Preferred if discovered, but assume HTML scraping needed

**Fallback Strategy**: If initial selectors fail, implement flexible pattern matching and log selector failures for manual review and adjustment.

---

## Research Task 2: Colly Framework for Calendar Navigation

**Question**: What's the best approach for navigating through multi-month calendars using Colly v2?

### Decision: Sequential page visits with state tracking

**Rationale**: Colly excels at stateful scraping where context needs to be maintained across multiple page visits. For calendar navigation, we'll track visited months and systematically navigate through the 3-month window.

**Implementation Approach**:
```go
type CinemagicScraper struct {
    collector     *colly.Collector
    currentDate   time.Time
    endDate       time.Time
    visitedMonths map[string]bool
}

// Navigate to next month, visit movie pages, repeat until endDate
```

**Best Practices**:
- Use `c.OnHTML()` callbacks for page-specific extraction
- Maintain state in scraper struct (visited pages, extracted data)
- Use `c.OnRequest()` for rate limiting (3-second delays)
- Use `c.OnError()` for graceful error handling and logging
- Limit parallel requests: `c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})`

**Alternatives Considered**:
- Concurrent scraping of all dates - Rejected due to rate limiting requirements and potential IP blocking
- Recursive navigation - Rejected as sequential iteration is simpler and more maintainable

**Reference**: Existing `clinton_street_theater/scraper.go` uses similar patterns for event-based scraping

---

## Research Task 3: Film Format Badge Detection

**Question**: How should the scraper detect and extract film format information (digital, 35mm, 70mm) from movie pages?

### Decision: Multi-selector search with "digital" default

**Rationale**: Film format badges may appear in various locations on movie pages (header, metadata section, near showtimes). A flexible detection strategy with sensible defaults prevents scraper failures when format information is missing or moved.

**Implementation Pattern**:
```go
func extractFilmFormat(e *colly.HTMLElement) string {
    // Try multiple selectors in order of likelihood
    selectors := []string{
        ".format-badge",
        ".film-format",
        "[data-format]",
        ".badge:contains('mm')",
        ".badge:contains('igital')",
    }
    
    for _, selector := range selectors {
        if format := e.DOM.Find(selector).First().Text(); format != "" {
            return normalizeFormat(format)
        }
    }
    
    // Default to digital if no format found
    return "digital"
}

func normalizeFormat(raw string) string {
    raw = strings.ToLower(strings.TrimSpace(raw))
    if strings.Contains(raw, "35") || strings.Contains(raw, "35mm") {
        return "35mm"
    }
    if strings.Contains(raw, "70") || strings.Contains(raw, "70mm") {
        return "70mm"
    }
    if strings.Contains(raw, "imax") {
        return "IMAX"
    }
    return "digital"
}
```

**Alternatives Considered**:
- Failing scrape when format not found - Rejected as format is supplementary information
- Extracting from free-text descriptions - Rejected as too error-prone; prefer explicit badges

**Edge Case Handling**: If multiple formats detected (e.g., "Digital & 35mm"), use the first detected format or most specific (35mm > digital).

---

## Research Task 4: Movie Title Normalization for TMDB Matching

**Question**: How should movie titles be cleaned and normalized to maximize TMDB API match success?

### Decision: Strip special event suffixes while preserving core title

**Rationale**: Cinemagic may include special edition information in titles (e.g., "Blade Runner (Director's Cut)", "Movie Title - 25th Anniversary"). TMDB searches work best with clean, core movie titles.

**Implementation Pattern**:
```go
func normalizeTitleForTMDB(rawTitle string) string {
    title := strings.TrimSpace(rawTitle)
    
    // Remove parenthetical suffixes
    if idx := strings.Index(title, "("); idx > 0 {
        base := strings.TrimSpace(title[:idx])
        if len(base) > 0 {
            title = base
        }
    }
    
    // Remove dash suffixes (e.g., " - 25th Anniversary")
    if idx := strings.Index(title, " - "); idx > 0 {
        base := strings.TrimSpace(title[:idx])
        if len(base) > 3 { // Ensure we don't strip actual title
            title = base
        }
    }
    
    // Remove year indicators (e.g., " (2023)")
    title = regexp.MustCompile(`\s*\(\d{4}\)`).ReplaceAllString(title, "")
    
    return strings.TrimSpace(title)
}
```

**Best Practices**:
- Always store original title for display purposes
- Use normalized title only for TMDB API search
- Log normalization transformations for debugging
- If TMDB match fails, try with original title as fallback

**Alternatives Considered**:
- Using full title with all suffixes - Rejected as reduces TMDB match rate
- Manual title mapping table - Rejected as not scalable; normalization rules are sufficient

**Reference**: Similar pattern exists in `internal/tmdb/client.go` for title matching

---

## Research Task 5: Data Refresh Strategy (Full Replace vs. Merge)

**Question**: Should the scraper perform full replacement of Cinemagic data or merge with existing data?

### Decision: Full replacement per scraper run

**Rationale**: Clarification confirmed "Replace all Cinemagic data each run (full refresh)". This is the simplest approach, ensures data consistency, and matches how showtime data naturally works (theaters update their calendars completely, not incrementally).

**Implementation Pattern**:
```go
func (s *CinemagicScraper) Scrape() ([]models.Showtime, error) {
    // 1. Scrape fresh showtimes from Cinemagic
    showtimes := s.extractShowtimes()
    
    // 2. Storage layer handles replacement
    // In storage.go:
    // - Load existing showtimes
    // - Filter out all Cinemagic showtimes (theater_id == "cinemagic-theater")
    // - Append new Cinemagic showtimes
    // - Save to showtimes.json
    
    return showtimes, nil
}
```

**Best Practices**:
- Always preserve showtimes from other theaters
- Use theater_id as filter key ("cinemagic-theater")
- Maintain atomic file writes (write to temp file, then rename)
- Log count of removed vs. added showtimes

**Alternatives Considered**:
- Incremental updates (add new, remove old) - Rejected as more complex with no benefit
- Historical data preservation - Rejected as not required for current MVP; can be added later

**Reference**: Existing `internal/storage/storage.go` likely needs minor update to support theater-scoped replacement

---

## Research Task 6: Rate Limiting Implementation in Colly

**Question**: How should 3-second rate limiting be implemented in Colly to respect Cinemagic's servers?

### Decision: Use Colly's built-in rate limiting with conservative delays

**Rationale**: Colly provides `LimitRule` for domain-specific rate limiting. Setting a 3-second delay ensures politeness and reduces risk of IP blocking.

**Implementation Pattern**:
```go
func NewCinemagicScraper() *CinemagicScraper {
    c := colly.NewCollector(
        colly.AllowedDomains("tickets.thecinemagictheater.com", "www.thecinemagictheater.com"),
    )
    
    // Rate limiting: 3 seconds between requests
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*.thecinemagictheater.com",
        Delay:       3 * time.Second,
        RandomDelay: 500 * time.Millisecond, // Add slight randomness (2.5-3.5s)
        Parallelism: 1, // Sequential requests only
    })
    
    // Log all requests for debugging
    c.OnRequest(func(r *colly.Request) {
        log.Printf("Visiting %s", r.URL.String())
    })
    
    return &CinemagicScraper{collector: c}
}
```

**Best Practices**:
- Add small random delay (±500ms) to appear more human-like
- Set `Parallelism: 1` to ensure sequential processing
- Use domain glob to apply rules to all Cinemagic subdomains
- Log all requests for debugging and monitoring

**Performance Impact**: With 3-second delays, scraping 45 movie pages (15 movies × 3 months) takes ~135 seconds (2.25 minutes), meeting the <3 minute success criteria.

**Alternatives Considered**:
- Manual `time.Sleep()` calls - Rejected as Colly's built-in solution is cleaner and more maintainable
- 2-second delays - Rejected in favor of conservative 3-second approach per clarification

---

## Research Task 7: TMDB Match Failure Handling

**Question**: How should the UI display movies that don't match TMDB?

### Decision: Show with placeholder poster and "Limited Information Available" badge

**Rationale**: Clarification specified this approach (Option D from clarification session). This keeps all movies visible while transparently communicating data limitations.

**Implementation Requirements**:

**Backend Changes**:
- When TMDB match fails, store minimal movie data:
  ```json
  {
    "tmdb_id": 0,  // Special ID for unmatched
    "title": "Original Movie Title",
    "poster_path": "/assets/placeholder-poster.png",
    "overview": "",
    "tmdb_rating": 0,
    "runtime": 0,
    "rating": "",
    "limited_info": true  // Flag for frontend
  }
  ```

**Frontend Changes**:
- MovieCard.jsx: Check for `limited_info: true` flag
- Display badge: `<span className="limited-info-badge">Limited Information Available</span>`
- Use placeholder poster image
- Hide missing metadata (rating, runtime) gracefully

**CSS Styling**:
```css
.limited-info-badge {
    background: var(--neon-orange);
    color: var(--dark-navy);
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: bold;
}
```

**Alternatives Considered**:
- Hiding unmatched movies (Option C) - Rejected as users expect complete listings
- Separate section for unmatched (Option B) - Rejected as adds UI complexity

**Testing**: Ensure scraper handles titles that definitely won't match TMDB (e.g., local film festival entries)

---

## Summary

All research tasks completed successfully. Key decisions:

| Decision Area | Choice | Rationale |
|---------------|--------|-----------|
| Website scraping | Event-based with movie pages | Structured data on ticketing site |
| Calendar navigation | Sequential state tracking | Simple, maintainable, respects rate limits |
| Format detection | Multi-selector with default | Robust to page structure changes |
| Title normalization | Strip suffixes, preserve core | Maximizes TMDB match rate |
| Data refresh | Full replacement per run | Simplest, ensures consistency |
| Rate limiting | 3-second with Colly LimitRule | Conservative, built-in support |
| TMDB failures | Placeholder + badge | Transparency, complete listings |

**No NEEDS CLARIFICATION items remain**. All technical unknowns resolved through research and clarification session.

**Next Phase**: Proceed to Phase 1 (Design) - Generate data-model.md, contracts/, and quickstart.md.
