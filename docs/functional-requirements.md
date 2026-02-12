# Functional Requirements: Portland Movie Theater Showtimes

## Project Overview

**Project Title:** Portland Movie Theater Showtimes

A web application that scrapes movie showtimes from local Portland theater websites and displays them in a retro theater-themed interface with neon color aesthetics. Integrates with TMDB (The Movie Database) via MCP server for rich movie metadata including posters, ratings, and descriptions.

## System Architecture

### Backend (Go + Colly)
Web scrapers that collect showtime data from multiple local theater websites and expose the data via REST API.

### TMDB Integration (MCP Server)
Model Context Protocol server connecting to TMDB API for movie metadata, posters, ratings, and additional information.

### Frontend (React)
Single-page application that displays aggregated showtimes with movie details in a retro theater aesthetic.

---

## 1. Core Components

### 1.1 Web Scrapers (Go)
**Purpose:** Collect movie showtime data from various theater websites

**Requirements:**
- Individual scraper modules for each theater website
- Built using Colly framework
- Configurable scraper settings (URLs, selectors, rate limiting)
- Error handling and retry logic
- On-demand execution via API endpoint or CLI command
- Data validation and normalization
- **Scrape current month plus next 2 months of showtimes** (3 months total) to provide comprehensive upcoming schedule

**Data to Extract:**
- Theater name and location
- Movie title
- Movie overview/description
- Show dates and times
- Movie format (2D, 3D, IMAX, Dolby, etc.)
- Ticket pricing (if available)
- **Event page link** - direct link to the theater's event page for the movie screening
- Booking URLs (if available)

### 1.2 TMDB MCP Server Integration
**Purpose:** Enrich scraped showtime data with comprehensive movie information

**MCP Server:** `mcp-server-tmdb` (https://github.com/Laksh-star/mcp-server-tmdb)

**Capabilities:**
- Search movies by title
- Retrieve detailed movie information (plot, cast, crew, runtime, ratings)
- Fetch high-quality movie posters and backdrops
- Get movie genres and tags
- Access user ratings and review scores
- Discover trending and popular movies

**Integration Points:**
- Match scraped movie titles with TMDB database
- **Fetch high-quality movie poster images for display on website** (primary requirement)
- Retrieve runtime, rating (G, PG, PG-13, R), and genre
- Get TMDB ratings and vote counts
- Pull movie descriptions/overviews
- Provide fallback placeholder image for movies without TMDB posters

**Data Flow:**
1. Scrapers extract movie titles from theater websites
2. Backend queries TMDB MCP server with movie titles
3. MCP server returns matching movie data from TMDB
4. Backend merges TMDB data with showtime data
5. API serves enriched data to frontend

### 1.3 Data Storage
**Requirements:**
- Store scraped showtime data with timestamps
- Cache TMDB movie data to reduce API calls
- Support data deduplication
- Maintain historical data for trend analysis (optional)
- Fast read access for API queries

**Storage Format:**
- JSON files for data persistence
- Separate JSON cache for TMDB movie metadata
- Structure organized by theater, movie, and date

### 1.4 REST API (Go)
**Purpose:** Serve enriched showtime data to the React frontend

**Endpoints:**
- `GET /api/theaters` - List all theaters
- `GET /api/showtimes` - Get all showtimes (with date/theater filters)
- `GET /api/showtimes/:theater` - Get showtimes for specific theater
- `GET /api/movies` - Get unique list of movies currently showing (with TMDB data)
- `GET /api/movies/:id` - Get detailed movie information from TMDB
- `POST /api/scrape` - Trigger scraper on-demand (optional theater filter)
- `GET /api/health` - Health check endpoint
- `GET /api/last-updated` - Get last scrape timestamp

**Query Parameters:**
- `date` - Filter by specific date (default: today)
- `movie` - Filter by movie title
- `format` - Filter by format (2D, 3D, IMAX, etc.)
- `theater` - Filter by theater ID
- `genre` - Filter by genre

### 1.5 React Frontend
**Purpose:** Display showtimes with rich movie data in user-friendly, visually appealing interface

**Core Features:**
- Display all theaters and their current showtimes
- **Display TMDB movie posters for every movie showing** (essential visual element)
- Display movie ratings, runtime, and descriptions
- Filter by date (today, tomorrow, this week)
- Filter by theater, genre, or rating
- Search by movie title
- View showtimes organized by movie or by theater
- Responsive design (mobile and desktop)
- Last updated timestamp display
- Movie detail modal/page with full TMDB information

**UI Components:**
- Theater list/grid
- **Movie card featuring TMDB poster image** as primary visual element, plus title, rating, runtime, theater label, and showtimes
- **Expandable movie synopsis** - truncated overview with "More" button to reveal full text
- **Theater label display** - shows which theater(s) are showing the movie
- **Ticket link button** - clickable "Tickets" button linking to the theater's event page for purchasing tickets
- Movie detail modal (expanded view with full-size poster, description, cast, etc.)
- Date selector with **next week navigation** - allows browsing 7 days at a time with next/previous controls
- Search/filter bar
- Theater selector/tabs
- Genre filter chips
- Loading states
- Error states (when data unavailable)

---

## 2. Retro Theater Color Scheme

### 2.1 Color Palette
**Primary Colors:**
- Neon Blue: `#00D9FF`, `#0099FF`
- Neon Pink: `#FF1493`, `#FF69B4`
- Neon Orange: `#FF6B35`, `#FFA500`
- Neon Red: `#FF073A`, `#FF4444`

**Background Colors:**
- Dark Navy/Black: `#0A0E27`, `#1A1A2E`
- Deep Purple: `#16213E`, `#2E1A47`

**Accent Colors:**
- Neon Yellow: `#FFD700`, `#FFF44F`
- Neon Green: `#39FF14`, `#00FF41`
- White/Light: `#FFFFFF`, `#F0F0F0`

### 2.2 Design Elements
- Neon glow effects (CSS box-shadow with multiple layers)
- Vintage marquee-style headers with animated lights
- Art deco-inspired borders and dividers
- Retro font styles (e.g., "Bebas Neue", "Monoton", "Righteous", "Press Start 2P")
- Gradient backgrounds with dark base colors
- Animated neon flickering effects (subtle)
- Old film grain texture overlay (optional)
- Movie posters with neon-glow frames
- Retro ticket stub design elements
- Vintage cinema curtain motifs

---

## 3. Technical Requirements

### 3.1 Backend (Go)
**Dependencies:**
- `github.com/gocolly/colly/v2` - Web scraping
- `github.com/gin-gonic/gin` or `net/http` - HTTP server
- `encoding/json` - JSON handling
- MCP client library for Go (or HTTP client for MCP server)
- Standard library packages for utilities

**Configuration:**
- Environment-based configuration
- Configurable scraper targets (theater URLs)
- TMDB MCP server connection settings
- Rate limiting settings
- CORS configuration for frontend
- Cache TTL settings

**Performance:**
- Concurrent scraping where appropriate
- Caching of TMDB responses (24-48 hours)
- Efficient data serialization
- Connection pooling for MCP server

### 3.2 TMDB MCP Server Setup
**Requirements:**
- Node.js environment for MCP server
- TMDB API key (free tier sufficient for most use cases)
- MCP server configuration
- Environment variables for API credentials

**Configuration:**
```bash
TMDB_API_KEY=your_api_key_here
MCP_SERVER_PORT=3001
```

### 3.3 Frontend (React)
**Dependencies:**
- React 18+
- React Router (for navigation)
- Axios or Fetch API (for HTTP requests)
- CSS Modules or Styled Components
- Date library (date-fns or dayjs)
- Animation library (Framer Motion or React Spring) for neon effects

**Browser Support:**
- Modern browsers (Chrome, Firefox, Safari, Edge)
- Mobile responsive (iOS Safari, Chrome Mobile)

---

## 4. Data Model

### 4.1 Theater
```json
{
  "id": "string",
  "name": "string",
  "address": "string",
  "city": "string",
  "zip": "string",
  "website": "string",
  "phone": "string (optional)"
}
```

### 4.2 Movie (Enriched with TMDB)
```json
{
  "tmdb_id": "number",
  "title": "string",
  "original_title": "string",
  "overview": "string",
  "runtime": "number (minutes)",
  "rating": "string (G, PG, PG-13, R, NR)",
  "genres": ["string"],
  "release_date": "string (YYYY-MM-DD)",
  "poster_path": "string (TMDB URL - REQUIRED for display)",
  "poster_url_small": "string (w342 - for thumbnails)",
  "poster_url_medium": "string (w500 - for movie cards)",
  "poster_url_large": "string (w780 - for detail views)",
  "backdrop_path": "string (TMDB URL - optional)",
  "tmdb_rating": "number (0-10)",
  "vote_count": "number",
  "popularity": "number",
  "cast": ["string"] (top billed),
  "director": "string"
}
```

**Note:** Poster images are a primary requirement. Every movie displayed on the website should show its TMDB poster image. If no poster is available from TMDB, a placeholder image must be used.

### 4.3 Showtime
```json
{
  "id": "string",
  "theater_id": "string",
  "movie_title": "string",
  "tmdb_id": "number",
  "date": "string (YYYY-MM-DD)",
  "time": "string (HH:MM)",
  "format": "string (35mm, digital, 70mm, IMAX - defaults to 'digital' if not specified)",
  "price": "number (optional)",
  "booking_url": "string (optional)",
  "screen": "string (optional)"
}
```

**Film Format Standards:**
- Valid formats: `35mm`, `digital`, `70mm`, `IMAX`
- Default format: `digital` (when not specified by theater)
- Legacy formats like "35mm/Digital" should be normalized to the primary format

### 4.4 Scrape Metadata
```json
{
  "last_updated": "timestamp",
  "theater_id": "string",
  "status": "success | error",
  "error_message": "string (optional)",
  "movies_scraped": "number",
  "showtimes_scraped": "number"
}
```

---

## 5. User Stories

### As a User:
1. I want to see all movies playing at local theaters today **with movie posters** and descriptions from TMDB
2. I want to **see which theater is showing each movie** directly on the movie card
3. I want to **read full movie overviews** by clicking "More" when the synopsis is truncated
4. I want to **browse future weeks** by clicking "Next" to see showtimes 7+ days ahead
5. I want to filter showtimes by specific theater
6. I want to view showtimes for upcoming dates
7. I want to search for a specific movie
8. I want to see movie ratings and reviews from TMDB
9. I want to see what **film formats** are available (35mm, digital, 70mm, IMAX)
10. I want to view movie details including **full-size poster**, cast, director, and plot
11. I want to filter movies by genre
12. I want to know when the data was last updated
10. I want to access the site on my phone
11. I want an aesthetically pleasing retro theater experience
12. I want to click through to purchase tickets

### As an Administrator:
1. I want to add new theater scrapers easily
2. I want to see scraper health status
3. I want to trigger scraper execution on-demand
4. I want to handle scraper failures gracefully
5. I want to monitor TMDB API usage
6. I want to manage TMDB data cache

---

## 6. Non-Functional Requirements

### 6.1 Performance
- API response time < 500ms (with cached TMDB data)
- Frontend initial load < 3s
- **TMDB movie poster images optimized for web delivery** (use TMDB's image CDN with appropriate sizes: w500 for cards, w780 for detail views)
- Image loading optimized (lazy loading, responsive images)
- Scraping should not overload theater websites (respect robots.txt, rate limiting)
- TMDB cache hit ratio > 90%

### 6.2 Reliability
- Graceful degradation when scrapers fail
- Fallback to basic data if TMDB unavailable
- Data freshness indicators
- Error logging and monitoring
- Retry logic for failed TMDB requests

### 6.3 Scalability
- Easy to add new theater scrapers
- Modular scraper architecture
- TMDB caching to minimize API calls
- Support for 5-10 theaters initially, scalable to 20+
- Handle 100+ unique movies per week

### 6.4 Maintainability
- Clean, documented code
- Separation of concerns
- Configuration-driven scraper definitions
- Unit tests for critical functions
- Clear separation between scraping, enrichment, and API layers

### 6.5 Security
- No storage of personal user data
- CORS properly configured
- Rate limiting on API endpoints
- Secure storage of TMDB API credentials
- Respect theater website terms of service
- HTTPS for production deployment

---

## 7. TMDB Integration Workflow

### 7.1 Movie Matching Process
1. **Scraper extracts movie title** from theater website
2. **Normalize title** (remove special characters, trim whitespace)
3. **Search TMDB** via MCP server using normalized title
4. **Match results** (exact match preferred, fuzzy match with confidence score)
5. **Fetch full movie details** if match found
6. **Cache TMDB data** locally with timestamp
7. **Merge data** with showtime information

### 7.2 Cache Strategy
- **Initial cache duration:** 7 days for movie metadata
- **Poster images:** Cache TMDB poster URLs (TMDB CDN handles image delivery), let browser cache handle actual images
- **Poster sizes:** Store multiple poster URLs (w342 for thumbnails, w500 for cards, w780 for detail modals)
- **Daily refresh:** Update TMDB ratings and vote counts
- **Manual invalidation:** Admin can force refresh for specific movies

### 7.3 Error Handling
- If TMDB match fails, display movie with title only and **generic movie poster placeholder image**
- If poster URL is invalid or image fails to load, show **fallback placeholder** (film reel or theater icon)
- If MCP server unavailable, use cached data
- Log all TMDB failures for manual review
- Provide admin interface to manually link movies to TMDB IDs

---

## 8. Future Enhancements

### Phase 1 (Post-MVP):
- Movie trailers from TMDB/YouTube
- User favorites/bookmarking (localStorage)
- Email/push notifications for new showtimes
- Calendar export functionality (.ics files)
- Share showtimes via social media

### Phase 2 (Advanced):
- User accounts and preferences
- Personalized recommendations based on viewing history
- Advanced filtering (by actor, director, TMDB rating)
- Movie reviews and ratings from multiple sources
- Interactive theater seat maps (if APIs available)
- Mobile app version (React Native)

### Phase 3 (Community):
- User-submitted theater data
- Community ratings and reviews
- Discussion forums for movies
- Local movie club features

---

## 9. Project Structure

```
theater-showtimes/
├── backend/
│   ├── cmd/
│   │   ├── scraper/          # Scraper CLI
│   │   └── api/              # API server
│   ├── internal/
│   │   ├── scrapers/         # Theater-specific scrapers
│   │   │   ├── theater1.go
│   │   │   ├── theater2.go
│   │   │   └── base.go       # Common scraper utilities
│   │   ├── models/           # Data models
│   │   ├── storage/          # Data persistence
│   │   ├── tmdb/             # TMDB MCP client
│   │   ├── enricher/         # Data enrichment logic
│   │   └── api/              # API handlers
│   ├── configs/              # Configuration files
│   ├── pkg/                  # Shared utilities
│   └── go.mod
├── mcp-server/
│   ├── package.json          # TMDB MCP server setup
│   └── .env                  # TMDB API credentials
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/       # React components
│   │   │   ├── MovieCard/
│   │   │   ├── TheaterList/
│   │   │   ├── DateSelector/
│   │   │   ├── FilterBar/
│   │   │   └── MovieModal/
│   │   ├── pages/            # Page components
│   │   │   ├── Home.jsx
│   │   │   ├── Theater.jsx
│   │   │   └── Movie.jsx
│   │   ├── services/         # API services
│   │   │   └── api.js
│   │   ├── styles/           # CSS/styling
│   │   │   ├── colors.css    # Retro color scheme
│   │   │   ├── neon.css      # Neon effects
│   │   │   └── global.css
│   │   ├── utils/            # Utility functions
│   │   └── App.jsx
│   └── package.json
├── docs/
│   ├── functional-requirements.md
│   ├── api-documentation.md
│   └── scraper-guide.md
└── README.md
```

---

## 10. Development Phases

### Phase 1: Foundation (Week 1-2)
- Set up project structure
- Configure TMDB MCP server
- Implement first theater scraper
- Test TMDB integration and movie matching
- Create basic data storage with caching

### Phase 2: Backend Core (Week 2-3)
- Implement remaining theater scrapers (2-3 minimum)
- Build REST API with all endpoints
- Implement data enrichment pipeline
- Add comprehensive error handling
- Implement on-demand scraping endpoint

### Phase 3: Frontend Core (Week 3-4)
- Build React app structure
- Implement core UI components
- Integrate with backend API
- Basic filtering and search functionality
- Display movie posters and TMDB data

### Phase 4: Styling & Polish (Week 4-5)
- Apply retro color scheme throughout
- Add neon glow effects and animations
- Implement smooth transitions
- Mobile responsiveness
- Loading and error states
- Performance optimization

### Phase 5: Testing & Deployment (Week 5-6)
- Integration testing
- User acceptance testing
- Bug fixes and refinements
- Deployment configuration
- Monitoring setup
- Documentation completion

---

## 11. Success Criteria

1. ✅ Successfully scrape showtimes from at least 3 local theaters
2. ✅ Integrate TMDB data for 95%+ of movies
3. ✅ Display accurate, up-to-date showtime information
4. ✅ Show high-quality movie posters and descriptions
5. ✅ Responsive design works on mobile and desktop
6. ✅ Retro aesthetic is visually appealing and consistent
7. ✅ Data updates automatically at configured intervals
8. ✅ User can filter by theater, date, movie, and genre
9. ✅ TMDB data loads quickly via caching
10. ✅ Application is stable and handles errors gracefully
11. ✅ API response times under 500ms
12. ✅ Zero user-facing errors during normal operation

---

## 12. API Usage & Rate Limits

### TMDB API (via MCP Server)
- **Free Tier:** 40 requests per 10 seconds
- **Daily Limit:** Effectively unlimited for our use case
- **Strategy:** Aggressive caching (7-day TTL) to stay well under limits
- **Estimated usage:** ~50-100 movies/week = 50-100 requests/week after initial cache

### Theater Websites
- **Respect robots.txt** on all theater websites
- **Rate limiting:** 1 request per 2-3 seconds per domain
- **User agent:** Custom user agent identifying the scraper
- **Execution:** On-demand via API or CLI (manual trigger)

---

## 13. Configuration Files

### Backend Config (config.yaml)
```yaml
scraper:
  timeout: 30s
  rate_limit: 2s
  
tmdb:
  mcp_server_url: http://localhost:3001
  cache_ttl: 168h # 7 days
  
api:
  port: 8080
  cors_origins: ["http://localhost:3000"]
  
storage:
  type: json
  path: ./data

theaters:
  - id: theater1
    name: "Local Cinema 1"
    url: "https://example.com/showtimes"
    scraper: theater1
```

### Frontend Config (.env)
```
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENABLE_ANIMATIONS=true
REACT_APP_TMDB_IMAGE_BASE=https://image.tmdb.org/t/p/w500
```

---

## 14. Known Challenges & Solutions

### Challenge 1: Movie Title Matching
- **Problem:** Theater websites may use different titles than TMDB
- **Solution:** Fuzzy matching algorithm, manual override system, common title variations database

### Challenge 2: Website Structure Changes
- **Problem:** Theater websites redesign breaks scrapers
- **Solution:** Versioned selectors, health monitoring, alerts, graceful fallbacks

### Challenge 3: Performance with Many Theaters
- **Problem:** Scraping/enrichment takes too long
- **Solution:** Concurrent scraping, incremental updates, background jobs

### Challenge 4: TMDB API Rate Limits
- **Problem:** Too many requests during initial setup
- **Solution:** Request queuing, aggressive caching, batch processing

### Challenge 5: Missing TMDB Data
- **Problem:** Indie or very new films may not be in TMDB
- **Solution:** Graceful degradation, manual entry option, multiple data sources
