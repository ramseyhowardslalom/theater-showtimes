package scrapers

import (
	"theater-showtimes/internal/models"
)

// Scraper interface that all theater scrapers must implement
type Scraper interface {
	// GetTheaterInfo returns basic theater information
	GetTheaterInfo() models.Theater
	
	// Scrape performs the scraping and returns showtimes
	Scrape() ([]models.Showtime, error)
	
	// GetID returns the unique identifier for this scraper
	GetID() string
}

// Registry holds all available scrapers
type Registry struct {
	scrapers map[string]Scraper
}

// NewRegistry creates a new scraper registry
func NewRegistry() *Registry {
	return &Registry{
		scrapers: make(map[string]Scraper),
	}
}

// Register adds a scraper to the registry
func (r *Registry) Register(scraper Scraper) {
	r.scrapers[scraper.GetID()] = scraper
}

// Get retrieves a scraper by ID
func (r *Registry) Get(id string) (Scraper, bool) {
	scraper, exists := r.scrapers[id]
	return scraper, exists
}

// GetAll returns all registered scrapers
func (r *Registry) GetAll() map[string]Scraper {
	return r.scrapers
}

// GetIDs returns all scraper IDs
func (r *Registry) GetIDs() []string {
	ids := make([]string, 0, len(r.scrapers))
	for id := range r.scrapers {
		ids = append(ids, id)
	}
	return ids
}
