# Portland Movie Theater Showtimes - Backend

Go backend for scraping and serving Portland theater showtime data.

## Structure

```
backend/
├── cmd/
│   ├── api/              # API server
│   └── scraper/          # CLI scraper tool
├── internal/
│   ├── api/              # API handlers and routing
│   ├── models/           # Data models
│   ├── scrapers/         # Theater scrapers
│   │   ├── example_theater/
│   │   └── local_cinema/
│   ├── storage/          # JSON storage
│   └── tmdb/             # TMDB client
├── configs/              # Configuration files
└── go.mod
```

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Create data directory:
```bash
mkdir -p data
```

3. Configure TMDB MCP server in `configs/config.yaml`

## Running

### API Server
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### CLI Scraper

Run all scrapers:
```bash
go run cmd/scraper/main.go
```

Run specific scrapers:
```bash
go run cmd/scraper/main.go example-theater local-cinema
```

## Adding New Scrapers

1. Create a new directory under `internal/scrapers/your_theater/`
2. Implement the `Scraper` interface in `scraper.go`
3. Register your scraper in `cmd/api/main.go` and `cmd/scraper/main.go`
4. Add theater configuration to `configs/config.yaml`

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/theaters` - List all theaters
- `GET /api/showtimes` - Get all showtimes (filterable)
- `GET /api/showtimes/:theater` - Get theater-specific showtimes
- `GET /api/movies` - List all movies
- `GET /api/movies/:id` - Get movie details
- `POST /api/scrape` - Trigger scraper on-demand
- `GET /api/last-updated` - Get last scrape timestamp
