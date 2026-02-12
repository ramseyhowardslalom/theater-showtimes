# Data Model: Cinemagic Theater Support

**Date**: 2026-02-11  
**Feature**: Cinemagic Theater Scraper  
**Branch**: 001-cinemagic-scraper

## Purpose

This document defines the data entities, fields, relationships, and validation rules for the Cinemagic Theater integration. All entities align with existing data structures used by Clinton Street Theater to ensure consistency.

---

## Entity: Theater

**Description**: Represents a physical movie theater location. Cinemagic Theater is a new instance of this existing entity type.

### Fields

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|------------------|
| `id` | string | Yes | Unique identifier for the theater | Must be "cinemagic-theater" (kebab-case standard) |
| `name` | string | Yes | Display name of the theater | Must be "Cinemagic Theater" |
| `address` | string | Yes | Full street address | Must be "2021 SE Hawthorne Blvd, Portland, OR 97214" |
| `city` | string | No | City name | Optional: "Portland" |
| `zip` | string | No | ZIP code | Optional: "97214" |
| `website` | string | Yes | Main website URL | Must be valid URL: "https://www.thecinemagictheater.com" |
| `phone` | string | No | Contact phone number | Optional, not required for Cinemagic |

### Relationships

- **One-to-Many**: Theater → Showtimes (one theater has many showtimes)
- **No Dependencies**: Theater entity is independent and can be created without showtimes

### State Transitions

Theaters are static entities (no state changes once created).

### Storage Format

```json
{
  "id": "cinemagic-theater",
  "name": "Cinemagic Theater",
  "address": "2021 SE Hawthorne Blvd, Portland, OR 97214",
  "website": "https://www.thecinemagictheater.com"
}
```

**File**: `backend/data/theaters.json` (array of Theater objects)

---

## Entity: Showtime

**Description**: Represents a single movie screening at a specific theater, date, and time. Cinemagic showtimes extend this existing entity.

### Fields

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|------------------|
| `id` | string | Yes | Unique identifier for the showtime | Auto-generated UUID or {theater_id}_{tmdb_id}_{date}_{time} |
| `theater_id` | string | Yes | Foreign key to Theater |  Must be "cinemagic-theater" for Cinemagic showtimes |
| `tmdb_id` | number | Yes | Foreign key to Movie (TMDB ID) | Must be valid TMDB ID (positive integer) or 0 for unmatched movies |
| `title` | string | Yes | Movie title (original from theater) | Must match scraped title from Cinemagic |
| `date` | string | Yes | Screening date | Must be YYYY-MM-DD format, within 3-month window from current date |
| `time` | string | Yes | Screening time | Must be HH:MM format (24-hour), Pacific Time |
| `format` | string | Yes | Film projection format | Must be one of: "digital", "35mm", "70mm", "IMAX" (default: "digital") |
| `link` | string | Yes | Link to theater's event page  | Must be valid URL: "https://tickets.thecinemagictheater.com/movie/{slug}" |
| `price` | number | No | Ticket price (if available) | Optional, not implemented for Cinemagic MVP |
| `booking_url` | string | No | Direct booking link | Optional, can use `link` field for now |
| `screen` | string | No | Screen number/name | Optional, not available on Cinemagic site |

### Relationships

- **Many-to-One**: Showtime → Theater (many showtimes belong to one theater via `theater_id`)
- **Many-to-One**: Showtime → Movie (many showtimes reference one movie via `tmdb_id`)

### State Transitions

Showtimes are immutable once created. Updates occur via full replacement (all Cinemagic showtimes deleted and re-populated each scraper run).

### Validation Rules

- **Date Range**: `date` must be between current date and current date + 3 months
- **Time Format**: `time` must be valid 24-hour time (00:00-23:59)
- **Format Default**: If `format` not extracted from page, default to "digital"
- **Required Fields**: All fields except `price`, `booking_url`, and `screen` are mandatory

### Storage Format

```json
{
  "id": "cinemagic-theater_12345_2026-02-15_19:00",
  "theater_id": "cinemagic-theater",
  "tmdb_id": 12345,
  "title": "Example Movie",
  "date": "2026-02-15",
  "time": "19:00",
  "format": "35mm",
  "link": "https://tickets.thecinemagictheater.com/movie/example-movie"
}
```

**File**: `backend/data/showtimes.json` (array of Showtime objects)

**Refresh Strategy**: On each Cinemagic scraper run:
1. Load all showtimes from `showtimes.json`
2. Filter out all entries where `theater_id == "cinemagic-theater"`
3. Append newly scraped Cinemagic showtimes
4. Save updated array back to `showtimes.json`

---

## Entity: Movie

**Description**: Represents movie metadata enriched from TMDB API. Movies are shared across theaters (not theater-specific).

### Fields

| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|------------------|
| `tmdb_id` | number | Yes | TMDB movie ID (primary key) | Must be valid TMDB ID or 0 for unmatched |
| `title` | string | Yes | Official movie title from TMDB | From TMDB API or original scraped title |
| `original_title` | string | No | Original language title | From TMDB API (optional) |
| `overview` | string | Yes | Movie synopsis/description | From TMDB API or empty string for unmatched |
| `runtime` | number | Yes | Runtime in minutes | From TMDB API or 0 for unmatched |
| `rating` | string | Yes | MPAA rating (G, PG, PG-13, R, NR) | From TMDB API or "NR" for unmatched |
| `genres` | array | Yes | List of genre strings | From TMDB API or empty array |
| `release_date` | string | No | Original release date | YYYY-MM-DD format from TMDB |
| `poster_path` | string | Yes | Full URL to poster image | TMDB CDN URL or "/assets/placeholder-poster.png" |
| `poster_url_small` | string | No | Thumbnail size (w342) | TMDB CDN URL |
| `poster_url_medium` | string | No | Card size (w500) | TMDB CDN URL |
| `poster_url_large` | string | No | Detail view size (w780) | TMDB CDN URL |
| `backdrop_path` | string | No | Background image URL | TMDB CDN URL (optional) |
| `tmdb_rating` | number | Yes | TMDB user rating (0-10) | From TMDB API or 0 for unmatched |
| `vote_count` | number | No | Number of votes | From TMDB API |
| `popularity` | number | No | TMDB popularity score | From TMDB API |
| `cast` | array | No | Top-billed cast members | From TMDB API (optional) |
| `director` | string | No | Director name | From TMDB API (optional) |
| `limited_info` | boolean | Yes | Flag for TMDB match failure | `true` if tmdb_id == 0, else `false` |

### Relationships

- **One-to-Many**: Movie → Showtimes (one movie can have many showtimes across different theaters/dates)
- **No Dependencies**: Movie can exist without showtimes (cached TMDB data)

### State Transitions

Movies are cached entities:
1. **Initial State**: Not in cache
2. **Fetched**: Scraped title triggers TMDB search → movie data stored
3. **Cached**: Movie remains in cache for 7 days
4. **Stale**: After 7 days, can be re-fetched to update ratings/metadata
5. **Unmatched**: If TMDB search fails, create minimal movie record with `limited_info: true`

### Validation Rules for Unmatched Movies

When TMDB match fails (`tmdb_id == 0`):
- `title`: Use original scraped title
- `poster_path`: Use placeholder: "/assets/placeholder-poster.png"
- `overview`: Empty string ""
- `runtime`: 0
- `rating`: "NR" (Not Rated)
- `genres`: Empty array []
- `tmdb_rating`: 0
- `limited_info`: `true` ⚠️ **Flag for frontend badge display**

### Storage Format (TMDB Matched)

```json
{
  "tmdb_id": 12345,
  "title": "Blade Runner",
  "overview": "A blade runner must pursue...",
  "runtime": 117,
  "rating": "R",
  "genres": ["Science Fiction", "Thriller"],
  "release_date": "1982-06-25",
  "poster_path": "https://image.tmdb.org/t/p/w500/abc123.jpg",
  "tmdb_rating": 8.1,
  "vote_count": 12543,
  "limited_info": false
}
```

### Storage Format (Unmatched)

```json
{
  "tmdb_id": 0,
  "title": "Local Independent Film",
  "overview": "",
  "runtime": 0,
  "rating": "NR",
  "genres": [],
  "poster_path": "/assets/placeholder-poster.png",
  "tmdb_rating": 0,
  "limited_info": true
}
```

**File**: `backend/data/movies.json` (array of Movie objects, indexed by `tmdb_id`)

---

## Entity Relationships Diagram

```
Theater (cinemagic-theater)
    ↓ 1:N
Showtime (theater_id, tmdb_id, date, time)
    ↓ N:1
Movie (tmdb_id, title, poster, rating, etc.)
    ↓ Optional
TMDB API Cache (7-day TTL)
```

**Key Points**:
- Theaters are independent (can exist without showtimes)
- Showtimes reference both Theater (via `theater_id`) and Movie (via `tmdb_id`)
- Movies can be shared across theaters (same TMDB ID for same movie at different venues)
- TMDB cache is separate from movie data (cache layer in `backend/internal/tmdb/cache.go`)

---

## Data Integrity Constraints

### Theater-Showtime Relationship
- **Constraint**: Every Showtime MUST reference a valid Theater
- **Enforcement**: Scraper creates Theater entity before creating Showtimes
- **Violation Handling**: If theater_id not found, reject showtime (log error)

### Showtime-Movie Relationship
- **Constraint**: Every Showtime MUST reference a Movie (even if TMDB match failed)
- **Enforcement**: Create minimal Movie record (tmdb_id=0) for unmatched titles
- **Violation Handling**: Never create orphan Showtimes

### Date Range Constraint
- **Constraint**: Showtimes MUST be within 3-month window (current month + 2 future months)
- **Enforcement**: Scraper filters out dates beyond 3-month range
- **Violation Handling**: Log warning, skip out-of-range showtimes

### Duplicate Prevention
- **Constraint**: No duplicate Showtimes (same theater, movie, date, time)
- **Enforcement**: Use composite ID or uniqueness check during scraping
- **Violation Handling**: If duplicate detected, keep later-scraped version

---

## Storage Layer Updates

### Required Changes to `backend/internal/storage/storage.go`

**New Method**: `ReplaceTheaterShowtimes(theaterID string, showtimes []Showtime) error`

```go
// ReplaceTheaterShowtimes removes all showtimes for a specific theater
// and replaces them with the provided showtimes.
func (s *JSONStorage) ReplaceTheaterShowtimes(theaterID string, newShowtimes []Showtime) error {
    // 1. Load existing showtimes
    allShowtimes, err := s.LoadShowtimes()
    if err != nil {
        return fmt.Errorf("failed to load showtimes: %w", err)
    }
    
    // 2. Filter out showtimes for the specified theater
    filtered := []Showtime{}
    for _, st := range allShowtimes {
        if st.TheaterID != theaterID {
            filtered = append(filtered, st)
        }
    }
    
    // 3. Append new showtimes
    filtered = append(filtered, newShowtimes...)
    
    // 4. Save back to file (atomic write)
    return s.SaveShowtimes(filtered)
}
```

**Rationale**: Implements full refresh strategy (clear all Cinemagic data, add fresh data) while preserving other theaters' showtimes.

---

## Summary

All entities use existing data structures from the codebase:
- **Theater**: Standard theater entity, Cinemagic is a new instance
- **Showtime**: Standard showtime entity with required fields (theater_id, tmdb_id, date, time, format, link)
- **Movie**: Standard TMDB-enriched movie entity with new `limited_info` flag for unmatched movies

**Key Design Decisions**:
1. **Limited Info Flag**: New `limited_info: true` field on Movie entity enables "Limited Information Available" badge in UI
2. **Full Refresh**: Showtimes replaced entirely per scraper run (simplest, most reliable)
3. **Placeholder Handling**: Unmatched movies use tmdb_id=0 and placeholder assets
4. **Format Default**: Film format defaults to "digital" if not detected on page

**Next**: Generate contracts/ for API and quickstart.md for developer guide.
