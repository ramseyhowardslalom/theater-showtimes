# Clinton Street Theater - Scraper Specifications

## Theater Overview

- **Website**: https://cstpdx.com
- **Theater ID**: `clinton-street-theater`
- **Film Format**: All Clinton Street Theater screenings are **digital** format

## Website Structure

### Calendar View
The Clinton Street Theater website uses a calendar-based event system:
- Calendar URL: `https://cstpdx.com/calendar/`
- Events are displayed in a monthly calendar grid
- Each event is a separate link to an event detail page
- Calendar can be navigated month-by-month

### Event URLs
Individual events follow this pattern:
- Format: `https://cstpdx.com/event/{event-slug}/`
- Example: `https://cstpdx.com/event/the-rocky-horror-picture-show-with-sinophelia-14/`

## Movie Title Parsing

### Standard Movie Events
Most events have straightforward movie titles that can be directly searched in TMDB.

### Special Guest Events
Clinton Street Theater sometimes hosts screenings with special guests introducing films. These events follow a specific naming pattern:

**Format**: `{movie name} with {person}`

**Examples**:
- "Rocky Horror Picture Show with Sinophelia"
- "The Room with Tommy Wiseau"

**Scraper Handling**:
When encountering event titles containing "with", the scraper should:
1. Extract the movie title portion **before** "with"
2. Trim whitespace from the extracted title
3. Search TMDB using only the movie title portion
4. If a match is found, include the movie's info in the database
5. The full event title (including "with {person}") can be stored as metadata, but TMDB matching should use only the movie title

**Example Processing**:
- Event title: `"Rocky Horror Picture Show with Sinophelia"`
- Extracted for TMDB search: `"Rocky Horror Picture Show"`
- TMDB search succeeds: Include movie data
- Stored event: Links the showtime to the matched TMDB movie

### Festival Presentation Events
Clinton Street Theater sometimes hosts films as part of festivals or special presentations. These events follow a specific naming pattern:

**Format**: `{Festival/Organization} presents: {movie name}`

**Examples**:
- "The Portland EcoFilm Festival presents: Runa Simi"
- "The Portland EcoFilm Festival presents: Cradled by the Earth"
- "Church of Film presents: The Singing Ringing Tree"

**Scraper Handling**:
When encountering event titles containing "presents:", the scraper should:
1. Extract the movie title portion **after** "presents:"
2. Trim whitespace from the extracted title
3. Search TMDB using only the movie title portion
4. If a match is found, include the movie's info in the database
5. The full event title (including festival/organization name) can be stored as metadata, but TMDB matching should use only the movie title

**Example Processing**:
- Event title: `"The Portland EcoFilm Festival presents: Runa Simi"`
- Extracted for TMDB search: `"Runa Simi"`
- TMDB search: Attempt to find match
- Stored event: Links the showtime to the matched TMDB movie (if found)

**Priority of Title Extraction**:
If an event title contains both patterns, apply in this order:
1. Check for "presents:" first - extract text after the colon
2. Then check for "with" - extract text before "with"
3. This handles complex cases like "Festival presents: Movie with Guest"

## Scraper Implementation Details

### Date Range
The scraper collects showtimes for **3 months** starting from the current date.

### Data Extraction
For each event, the scraper extracts:
- **Event title** (movie name, potentially with guest info)
- **Date** (in YYYY-MM-DD format)
- **Time** (in HH:MM 24-hour format, converted to Pacific Time)
- **Theater ID** (`clinton-street-theater`)
- **Film format** (`digital`)
- **Event page link** - The direct URL to the event page on the Clinton Street Theater website

#### Extracting Event Links
Each event in the calendar view has a link to its detail page. The scraper should:

1. **Calendar View**: Extract the `href` attribute from the event title link
   - Selector: `.tribe-events-calendar-month__calendar-event-title a`
   - Fallback selector: `a` (first anchor tag in event element)

2. **List View**: Extract the `href` attribute from the event title link
   - Selector: `.tribe-events-calendar-list__event-title-link`
   - Fallback selector: `a` (first anchor tag in event element)

3. **Link Format**: Links follow the pattern `https://cstpdx.com/event/{event-slug}/`
   - Example: `https://cstpdx.com/event/ghost-dog-the-way-of-the-samurai/`
   - Example: `https://cstpdx.com/event/the-rocky-horror-picture-show-with-sinophelia-14/`

4. **Storage**: Store the full URL in the `link` field of the showtime record

**Purpose**: The event page link is displayed as a "Tickets" button on the frontend, allowing users to visit the theater's event page for more information and ticket purchasing.

### TMDB Integration
After scraping Clinton Street Theater events:
1. Movie titles are searched in TMDB API
2. Matches are enriched with:
   - Poster images
   - Overview/synopsis
   - TMDB rating
   - Runtime
   - Official rating (PG, PG-13, R, etc.)
3. TMDB ID is stored to link showtimes to movie metadata

### Title Normalization
The scraper should normalize titles before TMDB search:
- Remove "with {person}" suffix for special events
- Trim whitespace
- Handle special characters appropriately
- Case-insensitive matching may improve results

## Calendar Scraping Process

### HTML Structure
The calendar uses the following structure:
- Calendar grid with date cells
- Event links within date cells
- Each event has a title and time
- Events may have multiple showtimes per day

### Navigation
To scrape 3 months:
1. Start at current month's calendar
2. Parse all events in current month
3. Navigate to next month (usually via "Next Month" button/link)
4. Repeat for 3 total months

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
  "theater_id": "clinton-street-theater",
  "format": "digital",
  "link": "https://cstpdx.com/event/movie-title-slug/"
}
```

### Theater Entry
Theater information in `backend/data/theaters.json`:
```json
{
  "id": "clinton-street-theater",
  "name": "Clinton Street Theater",
  "address": "2522 SE Clinton St, Portland, OR 97202",
  "website": "https://cstpdx.com"
}
```

## Important Notes

1. **All films are digital** - Unlike some theaters that show 35mm or 70mm, Clinton Street Theater screenings should default to `"digital"` format

2. **Guest appearances** - The "with {person}" pattern is important for special events. Don't skip these events; parse them correctly to match the movie

3. **Event uniqueness** - The same movie may appear multiple times in the calendar (different dates/times or with different guests). Each is a separate showtime entry

4. **Time zones** - All times should be stored and displayed in **Pacific Time (PT)**

5. **TMDB matching** - Not all events may have TMDB matches (especially local/independent films). Handle gracefully when no match is found

6. **Calendar updates** - The theater's calendar is updated regularly. Re-scraping will capture new events and showtimes

## Future Enhancements

Potential improvements for the scraper:
- Handle films that don't match TMDB (local films, special events)
- Extract guest information as separate metadata
- Support for series/marathon events
- Ticket pricing information if available
- Detect sold-out shows
- Handle rescheduled or cancelled events
