# Cinemagic Theater - Scraper Specifications

## Theater Overview

- **Website**: https://www.thecinemagictheater.com
- **Ticketing/Showtimes URL**: https://tickets.thecinemagictheater.com/now-showing
- **Theater ID**: `cinemagic-theater`
- **Film Format**: Varies by screening - extract from movie page badges (digital, 35mm, 70mm, etc.)

## Website Structure

### Now Showing View
The Cinemagic Theater ticketing site displays current and upcoming movies:
- Now Showing URL: `https://tickets.thecinemagictheater.com/now-showing`
- Movies are displayed in a grid/list format
- Each movie is a link to its dedicated movie page
- Calendar navigation component allows browsing future dates

### Movie-Specific Pages
Individual movie pages follow this pattern:
- Format: `https://tickets.thecinemagictheater.com/movie/{movie-slug}`
- Example: `https://tickets.thecinemagictheater.com/movie/sentimental-value`
- Contains all showtimes for that specific movie across multiple dates
- May include film format badge (e.g., "Digital", "35mm", "70mm")

### Calendar Navigation
The calendar component allows date selection:
- Navigate through days/weeks to access future showtimes
- Scraper should navigate through current month and next 2 months (3 months total)
- Each date may show different movies

## Film Format Detection

### Format Badges
Cinemagic movie pages may display format badges indicating projection type:
- **Digital**: Most common format, standard digital projection
- **35mm**: Traditional film projection
- **70mm**: Large format film projection
- Other specialty formats as labeled

### Scraper Handling
When extracting film format:
1. Look for format badges or indicators on the movie page
2. Common selectors might include: `.format-badge`, `.film-format`, or similar classes
3. If format badge is found, extract and normalize the text (e.g., "Digital" → "digital")
4. If NO format badge is found, default to `"digital"` (most common)
5. Store normalized format in the showtime record

**Default Behavior**: If format cannot be determined from the page, default to `"digital"`

## Movie Title Parsing

### Standard Movie Titles
Most Cinemagic movies have straightforward titles that can be directly searched in TMDB.

### Special Event Formats
Cinemagic may host special screenings, festivals, or director's cuts. Handle these appropriately:

**Potential Patterns**:
- Standard: `"Movie Title"`
- Director's Cut: `"Movie Title (Director's Cut)"`
- Anniversary: `"Movie Title - 25th Anniversary"`
- Special Edition: `"Movie Title: Special Edition"`

**Scraper Handling**:
1. Extract the core movie title
2. For parenthetical or suffix additions, consider extracting the base title for TMDB matching
3. Store full title as-is for display purposes
4. Use cleaned/normalized title for TMDB API search

**Example Processing**:
- Event title: `"Blade Runner (The Final Cut)"`
- Extracted for TMDB search: `"Blade Runner"`
- TMDB search: Find the movie with additional metadata
- Stored title: Keep full title for context

## Scraper Implementation Details

### Date Range
The scraper collects showtimes for **3 months** starting from the current date.

### Data Extraction
For each movie showing, the scraper extracts:
- **Movie title** (clean title without special characters)
- **Date** (in YYYY-MM-DD format)
- **Time** (in HH:MM 24-hour format, converted to Pacific Time)
- **Theater ID** (`cinemagic-theater`)
- **Film format** (extracted from badge or default to "digital")
- **Event page link** - The direct URL to the movie page on Cinemagic

#### Extracting Movie Links
From the "Now Showing" page:
1. Identify each movie element/card
2. Extract the `href` attribute from the movie title link
3. Common selectors might include:
   - `.movie-card a`
   - `.movie-title a`
   - `a[href*="/movie/"]`
4. Links follow pattern: `https://tickets.thecinemagictheater.com/movie/{movie-slug}`

#### Extracting Showtimes
From individual movie pages:
1. Navigate to movie-specific URL
2. Find showtime listings (may be organized by date)
3. Extract date and time for each showing
4. Parse time formats (may be "7:00 PM" → convert to "19:00")
5. Associate all showtimes with the movie from that page

### TMDB Integration
After scraping Cinemagic Theater events:
1. Movie titles are searched in TMDB API
2. Matches are enriched with:
   - Poster images (w500 for cards, w780 for detail views)
   - Overview/synopsis
   - TMDB rating and vote count
   - Runtime
   - Official rating (PG, PG-13, R, etc.)
   - Genres
3. TMDB ID is stored to link showtimes to movie metadata
4. Cache TMDB data for 7 days to reduce API calls

### Title Normalization
The scraper should normalize titles before TMDB search:
- Remove special edition suffixes (e.g., "(Director's Cut)")
- Trim whitespace
- Remove year indicators if present (e.g., "(2023)")
- Handle special characters appropriately
- Case-insensitive matching may improve results

## Calendar Scraping Process

### Navigation Strategy
To scrape 3 months of showtimes:
1. Start at "Now Showing" page
2. Identify calendar navigation controls
3. For each date in the range:
   - Select/navigate to that date
   - Extract all movies showing on that date
   - For each movie, extract showtimes
4. Handle pagination if needed
5. Navigate to next month when current month is complete

### HTML Structure (To Be Determined)
The exact HTML structure will need to be examined when implementing:
- Movie grid/list structure
- Calendar date selectors
- Showtime display format
- Format badge location

**Note**: Initial implementation should inspect live site to determine actual selectors

## Data Storage

### Showtimes Format
Each showtime entry in `backend/data/showtimes.json`:
```json
{
  "id": "unique-id",
  "tmdb_id": 123456,
  "title": "Movie Title",
  "date": "2026-02-11",
  "time": "19:00",
  "theater_id": "cinemagic-theater",
  "format": "digital",
  "link": "https://tickets.thecinemagictheater.com/movie/movie-title-slug"
}
```

### Theater Entry
Theater information in `backend/data/theaters.json`:
```json
{
  "id": "cinemagic-theater",
  "name": "Cinemagic Theater",
  "address": "2021 SE Hawthorne Blvd, Portland, OR 97214",
  "website": "https://www.thecinemagictheater.com"
}
```

## Important Notes

1. **Film formats vary** - Unlike some theaters with consistent formats, Cinemagic may show digital, 35mm, 70mm, or other formats. Always check for format badges.

2. **Default to digital** - If format cannot be determined, default to `"digital"` as it's the most common format.

3. **Ticket links** - The movie page URL serves as the ticket link, allowing users to purchase tickets directly from Cinemagic's site.

4. **Time zones** - All times should be stored and displayed in **Pacific Time (PT)**.

5. **TMDB matching** - Most mainstream movies should match TMDB. Handle gracefully when no match is found (use placeholder poster, store title only).

6. **Calendar navigation** - May require JavaScript interaction or dynamic page loading. Handle appropriately in scraper.

7. **Rate limiting** - Respect 1 request per 2-3 seconds to avoid overloading Cinemagic's servers.

8. **Website changes** - Theater websites can change structure. Implement robust error handling and logging.

## Scraper Architecture

### Module Location
Following the project's modular scraper pattern:
```
backend/internal/scrapers/cinemagic_theater/
└── scraper.go
```

### Interface Requirements
The scraper MUST implement the standard scraper interface:
- `Scrape()` method returning showtimes and errors
- Theater ID: `"cinemagic-theater"`
- Configurable via `backend/configs/config.yaml`

### Error Handling
- Log all scraping errors with context
- Continue scraping on individual movie failures
- Return partial results if some movies succeed
- Provide clear error messages for debugging

## Testing Strategy

### Manual Testing
1. Run scraper against live Cinemagic website
2. Verify extracted showtimes match website display
3. Check TMDB matching accuracy
4. Verify format detection works correctly
5. Test calendar navigation through multiple months

### Automated Testing
1. Unit tests for title normalization
2. Unit tests for time parsing (various formats)
3. Integration tests with mock HTML responses
4. Validation of data structure (JSON schema)

## Future Enhancements

Potential improvements for the scraper:
- Detect sold-out shows (if Cinemagic displays this information)
- Extract ticket pricing if available
- Support for special event information (Q&A sessions, guest appearances)
- Handle series/marathon events
- Improved format detection with multiple sources
- Fallback strategies if TMDB match fails
- Support for canceled or rescheduled shows
