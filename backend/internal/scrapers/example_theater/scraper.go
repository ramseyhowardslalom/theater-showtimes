package example_theater

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"theater-showtimes/internal/models"
)

// Scraper implements the scraper for Example Theater
type Scraper struct {
	theater models.Theater
}

// NewScraper creates a new Example Theater scraper
func NewScraper() *Scraper {
	return &Scraper{
		theater: models.Theater{
			ID:      "example-theater",
			Name:    "Example Theater",
			Address: "123 Main St",
			City:    "Your City",
			Zip:     "12345",
			Website: "https://example-theater.com",
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
		colly.AllowedDomains("example-theater.com"),
	)

	// Rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 2 * time.Second,
	})

	// Set up callbacks
	c.OnHTML(".showtime", func(e *colly.HTMLElement) {
		// Example parsing logic - customize based on actual website structure
		showtime := models.Showtime{
			ID:         fmt.Sprintf("%s-%s", s.theater.ID, e.Attr("data-id")),
			TheaterID:  s.theater.ID,
			MovieTitle: e.ChildText(".movie-title"),
			Date:       e.ChildText(".date"),
			Time:       e.ChildText(".time"),
			Format:     e.ChildText(".format"),
			Link: e.ChildAttr(".booking-link", "href"),
		}
		showtimes = append(showtimes, showtime)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error scraping %s: %v\n", r.Request.URL, err)
	})

	// Start scraping
	err := c.Visit(s.theater.Website + "/showtimes")
	if err != nil {
		return nil, fmt.Errorf("failed to visit website: %w", err)
	}

	return showtimes, nil
}
