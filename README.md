# Portland Movie Theater Showtimes

A web application that scrapes movie showtimes from local Portland theater websites and displays them in a retro theater-themed interface with neon color aesthetics. Integrates with TMDB (The Movie Database) for rich movie metadata.

## Project Structure

```
theater-showtimes/
├── backend/              # Go backend with scrapers and API
│   ├── cmd/
│   │   ├── api/         # API server
│   │   └── scraper/     # CLI scraper tool
│   ├── internal/
│   │   ├── api/         # API handlers
│   │   ├── models/      # Data models
│   │   ├── scrapers/    # Theater scrapers (modular)
│   │   │   ├── example_theater/
│   │   │   └── local_cinema/
│   │   ├── storage/     # JSON storage
│   │   └── tmdb/        # TMDB client
│   └── configs/         # Configuration
├── frontend/            # React frontend
│   ├── src/
│   │   ├── components/  # UI components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API services
│   │   └── styles/      # Global styles
│   └── public/
└── docs/                # Documentation
    ├── functional-requirements.md
    └── coding-guidelines.md
```

## Features

### Backend
- **Modular Scraper Architecture**: Each theater has its own scraper module
- **On-Demand Scraping**: Trigger scraping via API or CLI
- **JSON Storage**: Simple file-based data persistence
- **TMDB Integration**: Enrich movie data with posters, ratings, and metadata
- **REST API**: Serve data to frontend with filtering capabilities

### Frontend
- **Retro Neon Aesthetic**: Beautiful theater-themed UI
- **Movie Browsing**: View all movies and showtimes
- **Multiple View Modes**: Browse by movie or by theater
- **Filtering**: Filter by date, theater, and genre
- **TMDB Integration**: Display movie posters, ratings, and details
- **Responsive Design**: Works on mobile and desktop

## Quick Start

### Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create data directory:
```bash
mkdir -p data
```

4. Run the API server:
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start development server:
```bash
npm run dev
```

The app will be available at `http://localhost:3000`

## Adding New Scrapers

The scraper system is designed to be modular. Each theater has its own folder under `backend/internal/scrapers/`.

### Steps to Add a New Scraper:

1. **Create a new directory**:
```bash
mkdir backend/internal/scrapers/your_theater_name
```

2. **Create scraper.go** implementing the `Scraper` interface:
```go
package your_theater_name

import (
    "theater-showtimes/internal/models"
    "github.com/gocolly/colly/v2"
)

type Scraper struct {
    theater models.Theater
}

func NewScraper() *Scraper {
    return &Scraper{
        theater: models.Theater{
            ID:      "your-theater-id",
            Name:    "Your Theater Name",
            Address: "123 Theater St",
            City:    "Your City",
            Zip:     "12345",
            Website: "https://yourtheater.com",
        },
    }
}

func (s *Scraper) GetTheaterInfo() models.Theater {
    return s.theater
}

func (s *Scraper) GetID() string {
    return s.theater.ID
}

func (s *Scraper) Scrape() ([]models.Showtime, error) {
    // Implement scraping logic using Colly
    // See example_theater or local_cinema for reference
}
```

3. **Register your scraper** in both:
   - `backend/cmd/api/main.go`
   - `backend/cmd/scraper/main.go`

```go
import "theater-showtimes/internal/scrapers/your_theater_name"

// In main()
registry.Register(your_theater_name.NewScraper())
```

4. **Add configuration** to `backend/configs/config.yaml`:
```yaml
theaters:
  - id: your-theater-id
    name: "Your Theater Name"
    url: "https://yourtheater.com"
    scraper: your_theater_name
```

## Running Scrapers

### Via CLI

Run all scrapers:
```bash
cd backend
go run cmd/scraper/main.go
```

Run specific scrapers:
```bash
go run cmd/scraper/main.go example-theater local-cinema
```

### Via API

Trigger all scrapers:
```bash
curl -X POST http://localhost:8080/api/scrape
```

Trigger specific scrapers:
```bash
curl -X POST http://localhost:8080/api/scrape \
  -H "Content-Type: application/json" \
  -d '{"theater_ids": ["example-theater", "local-cinema"]}'
```

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/theaters` - List all theaters
- `GET /api/showtimes` - Get all showtimes (with filters)
- `GET /api/showtimes/:theater` - Get theater-specific showtimes
- `GET /api/movies` - List all movies
- `GET /api/movies/:id` - Get movie details
- `POST /api/scrape` - Trigger scraper on-demand
- `GET /api/last-updated` - Get last scrape timestamp

### Query Parameters

- `date` - Filter by date (YYYY-MM-DD)
- `theater` - Filter by theater ID
- `movie` - Filter by movie title
- `format` - Filter by format (2D, 3D, IMAX, etc.)

## Development

### Code Quality Tools

Both frontend and backend are set up with linting and formatting tools:

**Frontend:**
- ESLint for code quality
- Prettier for formatting
- Run: `npm run lint` and `npm run format`

**Backend:**
- Follow Go conventions
- Use `go fmt` for formatting
- Use `go vet` for static analysis

See [docs/coding-guidelines.md](docs/coding-guidelines.md) for detailed guidelines.

## Configuration

### Backend (`backend/configs/config.yaml`)
- Scraper settings (timeout, rate limits)
- TMDB MCP server configuration
- API server settings
- Storage configuration

### Frontend (`frontend/.env`)
- API URL configuration
- Feature flags

## TMDB Integration

The project includes a TMDB MCP (Model Context Protocol) server that provides AI-assisted access to The Movie Database API. This enables GitHub Copilot and other AI assistants to:

- Search for movies by title or keywords
- Get movie recommendations
- Retrieve trending movies
- Access detailed movie information

The MCP server is located in `mcp-servers/mcp-server-tmdb/` and is configured in `.vscode/mcp.json`.

**Quick Setup:**
1. Get a free TMDB API key from [TMDB](https://www.themoviedb.org/)
2. The API key is already configured in `.vscode/mcp.json`
3. The server runs automatically when using GitHub Copilot in VS Code

For detailed information, see [docs/tmdb-mcp-integration.md](docs/tmdb-mcp-integration.md).

## Data Storage

All data is stored in JSON format in the `backend/data/` directory:
- `theaters.json` - Theater information
- `showtimes.json` - All showtimes
- `movies.json` - Movies with TMDB data
- `metadata.json` - Scrape history and status

## Documentation

- [Functional Requirements](docs/functional-requirements.md) - Detailed project requirements
- [Coding Guidelines](docs/coding-guidelines.md) - Code standards and best practices
- [TMDB MCP Integration](docs/tmdb-mcp-integration.md) - TMDB MCP server setup and usage
- [Backend README](backend/README.md) - Backend-specific documentation
- [Frontend README](frontend/README.md) - Frontend-specific documentation

## Tech Stack

**Backend:**
- Go 1.21+
- Colly (web scraping)
- Gin (HTTP framework)
- JSON file storage

**Frontend:**
- React 18 with TypeScript
- React Router
- Vite
- Axios
- Framer Motion
- date-fns

## License

MIT
