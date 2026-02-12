# Quickstart Guide: Cinemagic Theater Support

**Feature**: Cinemagic Theater Scraper  
**Branch**: 001-cinemagic-scraper  
**Date**: 2026-02-11

## Overview

This guide walks developers through implementing, testing, and deploying the Cinemagic Theater scraper. Follow steps in order for a smooth development experience.

---

## Prerequisites

Before starting implementation:

- [x] Feature specification approved ([spec.md](spec.md))
- [x] Technical research completed ([research.md](research.md))
- [x] Data model designed ([data-model.md](data-model.md))
- [ ] Development environment set up (Go 1.21+, Node.js 18+)
- [ ] Existing codebase familiarized (Clinton Street Theater scraper as reference)

**Required Tools**:
- Go 1.21 or later
- Git
- Code editor (VS Code recommended)
- Web browser for manual testing

**Optional Tools**:
- Postman or curl for API testing
- Chrome DevTools for inspecting Cinemagic website HTML

---

## Step 1: Create Theater Entry

**File**: `backend/data/theaters.json`

**Action**: Add Cinemagic Theater to the theaters array.

```json
{
  "id": "cinemagic-theater",
  "name": "Cinemagic Theater",
  "address": "2021 SE Hawthorne Blvd, Portland, OR 97214",
  "website": "https://www.thecinemagictheater.com"
}
```

**Validation**: 
```bash
# Verify theaters.json is valid JSON
cat backend/data/theaters.json | jq .
```

**Expected Result**: JSON parses successfully,Cinemagic Theater appears in array.

---

## Step 2: Create Scraper Module

**File**: `backend/internal/scrapers/cinemagic_theater/scraper.go`

**Action**: Implement scraper following the research decisions.

### Basic Structure

```go
package cinemagic_theater

import (
    "fmt"
    "log"
    "time"
    "theater-showtimes/internal/models"
    "github.com/gocolly/colly/v2"
)

type CinemagicScraper struct {
    collector *colly.Collector
}

func NewCinemagicScraper() *CinemagicScraper {
    c := colly.NewCollector(
        colly.AllowedDomains("tickets.thecinemagictheater.com", "www.thecinemagictheater.com"),
    )
    
    // Rate limiting: 3 seconds between requests
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*.thecinemagictheater.com",
        Delay:       3 * time.Second,
        RandomDelay: 500 * time.Millisecond,
        Parallelism: 1,
    })
    
    return &CinemagicScraper{collector: c}
}

func (s *CinemagicScraper) Scrape() ([]models.Showtime, error) {
    showtimes := []models.Showtime{}
    
    // TODO: Implement scraping logic from research.md
    // 1. Visit https://tickets.thecinemagictheater.com/now-showing
    // 2. Extract movie links
    // 3. Visit each movie page
    // 4. Extract showtimes, dates, formats
    // 5. Normalize titles for TMDB matching
    
    return showtimes, nil
}
```

**Reference**: See `backend/internal/scrapers/clinton_street_theater/scraper.go` for similar implementation pattern.

**Testing During Development**:
```bash
# Run scraper in isolation
cd backend
go run cmd/scraper/main.go --theater=cinemagic-theater
```

---

## Step 3: Update Storage Layer

**File**: `backend/internal/storage/storage.go`

**Action**: Add or verify `ReplaceTheaterShowtimes()` method exists (see data-model.md for implementation).

**Test**:
```go
// In storage_test.go
func TestReplaceTheaterShowtimes(t *testing.T) {
    storage := NewJSONStorage("test_data/")
    
    // Create test showtimes for Cinemagic
    cinemagicShowtimes := []models.Showtime{
        {TheaterID: "cinemagic-theater", /* ... */},
    }
    
    err := storage.ReplaceTheaterShowtimes("cinemagic-theater", cinemagicShowtimes)
    assert.Nil(t, err)
    
    // Verify only Cinemagic showtimes changed, others preserved
}
```

---

## Step 4: Register Scraper in CLI

**File**: `backend/cmd/scraper/main.go`

**Action**: Add Cinemagic to available scrapers.

```go
import (
    // ... existing imports
    "theater-showtimes/internal/scrapers/cinemagic_theater"
)

func main() {
    // ... existing code
    
    scrapers := map[string]Scraper{
        "clinton-street-theater": clinton_street_theater.NewClintonStreetScraper(),
        "cinemagic-theater":       cinemagic_theater.NewCinemagicScraper(), // NEW
    }
    
    // ... rest of main logic
}
```

**Test**:
```bash
cd backend
go run cmd/scraper/main.go --theater=cinemagic-theater --dry-run
```

**Expected**: No errors, scraper executes (even if no data returned yet).

---

## Step 5: Implement Scraping Logic

**Guided Implementation Checklist**:

### 5.1 Extract Movie Links from Now Showing
- [ ] Visit `https://tickets.thecinemagictheater.com/now-showing`
- [ ] Use Colly `OnHTML()` to find movie links
- [ ] Extract href attributes (pattern: `/movie/{slug}`)
- [ ] Store movie URLs in slice

### 5.2 Visit Each Movie Page
- [ ] Iterate through movie URLs
- [ ] Use `collector.Visit(url)` to fetch movie page
- [ ] Implement 3-second delay between visits (Colly LimitRule handles this)

### 5.3 Extract Showtime Data
- [ ] Parse movie title from page
- [ ] Detect film format badge (use multi-selector approach from research.md)
- [ ] Extract showtime list (dates and times)
- [ ] Parse event page link (current movie page URL)

### 5.4 Navigate Calendar (if applicable)
- [ ] If calendar component exists, find "next month" button
- [ ] Click/visit next month page
- [ ] Repeat extraction for 3 total months
- [ ] Track visited months to avoid duplicates

### 5.5 Build Showtime Objects
- [ ] Create `models.Showtime` for each screening
- [ ] Set `theater_id` to "cinemagic-theater"
- [ ] Set `format` (detected or default to "digital")
- [ ] Set `link` to movie page URL
- [ ] Normalize title for TMDB matching (see research.md)

### 5.6 Error Handling
- [ ] Use `c.OnError()` to log HTTP errors
- [ ] Continue scraping if individual movie fails
- [ ] Return partial results on errors

---

## Step 6: Test TMDB Integration

**File**: `backend/internal/tmdb/client.go` (existing)

**Action**: Test that normalized Cinemagic titles match TMDB correctly.

```bash
# Manual test with sample title
go run -m testing backend/internal/tmdb/client_test.go
```

**Scenarios to Test**:
1. Standard title: "The Matrix" → Should match
2. Director's cut: "Blade Runner (The Final Cut)" → Normalize to "Blade Runner"
3. Anniversary: "Jaws - 45th Anniversary" → Normalize to "Jaws"
4. Unmatched indie film: "Local Film Festival Entry" → Create limited_info movie

**Expected**: 90%+ match rate for mainstream movies.

---

## Step 7: Add Limited Info Badge to Frontend

**File**: `frontend/src/components/MovieCard/MovieCard.jsx`

**Action**: Display "Limited Information Available" badge for unmatched movies.

```jsx
{movie.limited_info && (
    <span className="limited-info-badge">
        Limited Information Available
    </span>
)}
```

**CSS**: Add to `frontend/src/components/MovieCard/MovieCard.css`

```css
.limited-info-badge {
    background: var(--neon-orange);
    color: var(--dark-navy);
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.75rem;
    font-weight: bold;
    margin-top: 8px;
    display: inline-block;
}
```

**Test**: Create a test movie with `limited_info: true` in movies.json and verify badge appears.

---

## Step 8: Run End-to-End Test

### 8.1 Scrape Cinemagic Data

```bash
cd backend
go run cmd/scraper/main.go --theater=cinemagic-theater
```

**Expected Output**:
```
Scraping Cinemagic Theater...
Visiting https://tickets.thecinemagictheater.com/now-showing
Found 12 movies
Visiting movie page: /movie/example-movie-1
Extracted 5 showtimes
...
Scrape complete: 63 showtimes extracted
Saved to backend/data/showtimes.json
```

**Validation Checks**:
- [ ] At least 10 unique movies extracted (SC-001)
- [ ] Scrape completed in under 3 minutes (SC-005)
- [ ] No HTTP errors or panics
- [ ] showtimes.json updated with Cinemagic entries

### 8.2 Start Backend API

```bash
cd backend
go run cmd/api/main.go
```

**Test Endpoints**:
```bash
# Get all theaters (should include Cinemagic)
curl http://localhost:8080/api/theaters | jq '.[] | select(.id == "cinemagic-theater")'

# Get all showtimes (should include Cinemagic)
curl http://localhost:8080/api/showtimes | jq '[.[] | select(.theater_id == "cinemagic-theater") | .title] | unique'

# Get movies (should include Cinemagic movies)
curl http://localhost:8080/api/movies | jq 'length'
```

### 8.3 Start Frontend

```bash
cd frontend
npm run dev
```

**Manual UI Testing**:
1. Open http://localhost:3000
2. Verify Cinemagic movies appear in grid ✅
3. Verify "Cinemagic Theater" appears in dropdown ✅
4. Select "Cinemagic Theater" from dropdown
5. Verify only Cinemagic movies shown ✅ (SC-003: <1 second filter)
6. Verify theater labels show "Cinemagic Theater" ✅
7. Click "Tickets" button on a Cinemagic movie
8. Verify correct movie page opens at tickets.thecinemagictheater.com ✅ (SC-004)
9. Verify "Limited Information Available" badge appears for unmatched movies ✅

---

## Step 9: Verify Success Criteria

**SC-001**: Scraper extracts ≥10 unique movies ✅  
**SC-002**: 90% TMDB match rate ✅  
**SC-003**: Filter to Cinemagic in <1 second ✅  
**SC-004**: Tickets button opens correct page 100% ✅  
**SC-005**: Scrape completes in <3 minutes ✅  
**SC-006**: UI updates without code changes ✅  

---

## Step 10: Documentation

**Create**: `docs/cinemagic-theater-specs.md` (already created)

**Verify**:
- [ ] Theater overview section complete
- [ ] Website structure documented
- [ ] Film format detection explained
- [ ] Scraper implementation details provided
- [ ] Data storage format examples included

---

## Common Issues & Troubleshooting

### Issue: Selectors not finding elements

**Symptom**: Scraper returns 0 movies or showtimes  
**Solution**: 
1. Inspect live Cinemagic website HTML with browser DevTools
2. Update selectors in `scraper.go` to match actual structure
3. Use Colly's debug mode: `colly.Debugger(&debug.LogDebugger{})`

### Issue: Rate limiting errors (429 Too Many Requests)

**Symptom**: HTTP 429 errors in logs  
**Solution**:
1. Verify LimitRule is set to 3 seconds
2. Increase delay to 5 seconds if needed
3. Check that Parallelism is set to 1

### Issue: TMDB match rate below 90%

**Symptom**: Many movies show "Limited Information Available" badge  
**Solution**:
1. Review title normalization logic (research.md Task 4)
2. Check for special characters or encoding issues in scraped titles
3. Add more normalization rules for common patterns

### Issue: Scraper takes longer than 3 minutes

**Symptom**: Scrape exceeds SC-005 time limit  
**Solution**:
1. Verify rate limiting is 3 seconds (not higher)
2. Check for inefficient loops or redundant page visits
3. Reduce months scraped from 3 to 2 if acceptable

---

## Performance Benchmarks

Expected performance on standard hardware:

| Metric | Target | Actual | Status |
|--------|--------|--------|--------|
| Scrape time (45 movie pages) | <3 min | _TBD_ | ⏳ |
| API response (GET /api/showtimes) | <500ms | _TBD_ | ⏳ |
| Filter response (select theater) | <1s | _TBD_ | ⏳ |
| TMDB match rate | ≥90% | _TBD_ | ⏳ |
| Movies extracted | ≥10 | _TBD_ | ⏳ |

_TBD values filled in during implementation and testing._

---

## Next Steps After Implementation

1. **Code Review**: Submit PR with scraper implementation
2. **Integration Testing**: Verify no regressions to Clinton Street Theater scraper
3. **User Acceptance Testing**: Have stakeholders review Cinemagic listings
4. **Deployment**: Merge to main, deploy to production
5. **Monitoring**: Set up alerts for scraper failures or low TMDB match rates
6. **Documentation**: Update main README.md to mention Cinemagic support

---

## Resources

- **Spec**: [spec.md](spec.md)
- **Research**: [research.md](research.md)
- **Data Model**: [data-model.md](data-model.md)
- **Cinemagic Specs**: `docs/cinemagic-theater-specs.md`
- **Reference Implementation**: `backend/internal/scrapers/clinton_street_theater/scraper.go`
- **Colly Documentation**: https://github.com/gocolly/colly
- **TMDB API**: https://developers.themoviedb.org/

---

**Ready to begin implementation!** Start with Step 1 and work through systematically.
