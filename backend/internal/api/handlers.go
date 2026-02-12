package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"theater-showtimes/internal/models"
	"theater-showtimes/internal/scrapers"
	"theater-showtimes/internal/storage"
	"theater-showtimes/internal/tmdb"
)

// Handler contains all API handlers
type Handler struct {
	storage  *storage.Storage
	registry *scrapers.Registry
	tmdb     *tmdb.Client
}

// NewHandler creates a new API handler
func NewHandler(storage *storage.Storage, registry *scrapers.Registry, tmdb *tmdb.Client) *Handler {
	return &Handler{
		storage:  storage,
		registry: registry,
		tmdb:     tmdb,
	}
}

// GetTheaters returns all theaters
func (h *Handler) GetTheaters(c *gin.Context) {
	theaters, err := h.storage.LoadTheaters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, theaters)
}

// GetShowtimes returns all showtimes with optional filters
func (h *Handler) GetShowtimes(c *gin.Context) {
	showtimes, err := h.storage.LoadShowtimes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply filters
	date := c.Query("date")
	theater := c.Query("theater")
	movie := c.Query("movie")

	filtered := h.filterShowtimes(showtimes, date, theater, movie)

	c.JSON(http.StatusOK, filtered)
}

// GetTheaterShowtimes returns showtimes for a specific theater
func (h *Handler) GetTheaterShowtimes(c *gin.Context) {
	theaterID := c.Param("theater")

	showtimes, err := h.storage.LoadShowtimes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter by theater
	var filtered []models.Showtime
	for _, st := range showtimes {
		if st.TheaterID == theaterID {
			filtered = append(filtered, st)
		}
	}

	c.JSON(http.StatusOK, filtered)
}

// GetMovies returns all unique movies
func (h *Handler) GetMovies(c *gin.Context) {
	movies, err := h.storage.LoadMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetMovieDetails returns details for a specific movie
func (h *Handler) GetMovieDetails(c *gin.Context) {
	_ = c.Param("id") // movieID - TODO: Convert to int and fetch from TMDB

	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not yet implemented"})
}

// TriggerScrape triggers the scraper on-demand
func (h *Handler) TriggerScrape(c *gin.Context) {
	var request struct {
		TheaterIDs []string `json:"theater_ids,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If no specific theaters requested, scrape all
	scrapersToRun := h.registry.GetAll()
	if len(request.TheaterIDs) > 0 {
		scrapersToRun = make(map[string]scrapers.Scraper)
		for _, id := range request.TheaterIDs {
			if scraper, exists := h.registry.Get(id); exists {
				scrapersToRun[id] = scraper
			}
		}
	}

	// Run scrapers
	results := h.runScrapers(scrapersToRun)

	c.JSON(http.StatusOK, gin.H{
		"message": "Scraping completed",
		"results": results,
	})
}

// Health returns health status
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now(),
	})
}

// GetLastUpdated returns the last scrape timestamp
func (h *Handler) GetLastUpdated(c *gin.Context) {
	lastUpdate, err := h.storage.GetLastUpdate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"last_updated": lastUpdate,
	})
}

// Helper functions

func (h *Handler) filterShowtimes(showtimes []models.Showtime, date, theater, movie string) []models.Showtime {
	var filtered []models.Showtime

	for _, st := range showtimes {
		if date != "" && st.Date != date {
			continue
		}
		if theater != "" && st.TheaterID != theater {
			continue
		}
		if movie != "" && st.MovieTitle != movie {
			continue
		}
		filtered = append(filtered, st)
	}

	return filtered
}

func (h *Handler) runScrapers(scrapers map[string]scrapers.Scraper) []models.ScrapeMetadata {
	results := []models.ScrapeMetadata{}

	for _, scraper := range scrapers {
		metadata := models.ScrapeMetadata{
			LastUpdated: time.Now(),
			TheaterID:   scraper.GetID(),
			Status:      "success",
		}

		showtimes, err := scraper.Scrape()
		if err != nil {
			metadata.Status = "error"
			metadata.ErrorMessage = err.Error()
		} else {
			metadata.ShowtimesScraped = len(showtimes)
			// TODO: Save showtimes, enrich with TMDB data, etc.
		}

		results = append(results, metadata)
		h.storage.SaveMetadata(metadata)
	}

	return results
}
