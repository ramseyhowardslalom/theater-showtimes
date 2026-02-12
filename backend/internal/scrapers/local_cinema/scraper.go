package local_cinema

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"theater-showtimes/internal/models"
)

// Scraper implements the scraper for Local Cinema
type Scraper struct {
	theater models.Theater
}

// NewScraper creates a new Local Cinema scraper
func NewScraper() *Scraper {
	return &Scraper{
		theater: models.Theater{
			ID:      "local-cinema",
			Name:    "Local Cinema",
			Address: "456 Cinema Blvd",
			City:    "Your City",
			Zip:     "12345",
			Website: "https://local-cinema.com",
		},
	}
}

// GetTheaterInfo returns theater information
func (s *Scraper) GetTheaterInfo() models.Theater {
	return s.theater
}

// GetID returns the scraper ID
func (s *Scraper) GetID() string {
	return s.theater.ID
}

// Scrape performs the actual scraping
func (s *Scraper) Scrape() ([]models.Showtime, error) {
	showtimes := []models.Showtime{}

	c := colly.NewCollector(
		colly.AllowedDomains("local-cinema.com"),
	)

	// Rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 2 * time.Second,
	})

	// Set up callbacks - customize based on actual website structure
	c.OnHTML(".movie-listing", func(e *colly.HTMLElement) {
		// Example parsing logic
		showtime := models.Showtime{
			ID:         fmt.Sprintf("%s-%d", s.theater.ID, time.Now().Unix()),
			TheaterID:  s.theater.ID,
			MovieTitle: e.ChildText("h3.title"),
			Date:       e.ChildText(".show-date"),
			Time:       e.ChildText(".show-time"),
			Format:     e.ChildText(".format"),
		}
		showtimes = append(showtimes, showtime)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error scraping %s: %v\n", r.Request.URL, err)
	})

	err := c.Visit(s.theater.Website + "/now-showing")
	if err != nil {
		return nil, fmt.Errorf("failed to visit website: %w", err)
	}

	return showtimes, nil
}
