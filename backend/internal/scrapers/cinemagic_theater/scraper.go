package cinemagic_theater

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"theater-showtimes/internal/models"
)

const (
	// Base URLs for Cinemagic Theater
	baseURL    = "https://tickets.thecinemagictheater.com"
	nowShowing = "https://tickets.thecinemagictheater.com/now-showing"
)

// Scraper implements the scraper for Cinemagic Theater
type Scraper struct {
	theater models.Theater
}

// NewScraper creates a new Cinemagic Theater scraper
func NewScraper() *Scraper {
	return &Scraper{
		theater: models.Theater{
			ID:      "cinemagic-theater",
			Name:    "Cinemagic Theater",
			Address: "2021 SE Hawthorne Blvd, Portland, OR 97214",
			City:    "Portland",
			Zip:     "97214",
			Website: "https://www.thecinemagictheater.com",
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
	movieLinks := make(map[string]bool) // Track unique movie URLs

	// Create Colly collector with allowed domains
	c := colly.NewCollector(
		colly.AllowedDomains("tickets.thecinemagictheater.com", "www.thecinemagictheater.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; TheaterShowtimesBot/1.0)"),
	)

	// Rate limiting: 3 seconds between requests (T016 - already complete)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*.thecinemagictheater.com",
		Delay:       3 * time.Second,
		RandomDelay: 500 * time.Millisecond,
		Parallelism: 1,
	})

	// T009: Extract movie links from now-showing page
	// Try multiple possible selectors for movie cards/links
	c.OnHTML("a[href*='/movie/']", func(e *colly.HTMLElement) {
		movieURL := e.Attr("href")
		if movieURL == "" {
			return
		}

		// Convert relative URLs to absolute
		if !strings.HasPrefix(movieURL, "http") {
			movieURL = e.Request.AbsoluteURL(movieURL)
		}

		// Track unique movies
		if !movieLinks[movieURL] {
			movieLinks[movieURL] = true
			fmt.Printf("Found movie: %s\n", movieURL)
			
			// Visit movie page to extract showtimes (T011)
			e.Request.Visit(movieURL)
		}
	})

	// T011 & T012: Extract showtimes from individual movie pages
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Only process movie detail pages
		if !strings.Contains(e.Request.URL.String(), "/movie/") {
			return
		}

		// Extract movie title
		rawTitle := s.extractMovieTitle(e)
		if rawTitle == "" {
			fmt.Printf("Warning: Could not extract title from %s\n", e.Request.URL)
			return
		}

		// Extract film format (using helper from T006)
		format := extractFilmFormat(e)
		
		// Extract showtimes from the page
		movieShowtimes := s.extractShowtimes(e, rawTitle, format, e.Request.URL.String())
		
		if len(movieShowtimes) > 0 {
			fmt.Printf("  Found %d showtimes for \"%s\" (format: %s)\n", 
				len(movieShowtimes), rawTitle, format)
			showtimes = append(showtimes, movieShowtimes...)
		}
	})

	// T015: Error handling and logging
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error scraping %s: %v\n", r.Request.URL, err)
		// Continue scraping other pages even if one fails
	})

	// Request logging (for debugging)
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s\n", r.URL)
	})

	// Start scraping from now-showing page
	fmt.Printf("Starting Cinemagic Theater scrape for 3-month window...\n")
	
	// T010: Calendar navigation - visit now-showing and navigate through months
	// Visit now-showing page
	err := c.Visit(nowShowing)
	if err != nil {
		return nil, fmt.Errorf("failed to visit now-showing page: %w", err)
	}

	// TODO: Implement calendar month navigation if needed
	// For now, we're scraping from the now-showing page which may include future dates
	// If the website uses calendar navigation, we would implement month-by-month iteration here
	
	// Wait for all async requests to complete
	c.Wait()

	fmt.Printf("Scraping complete. Found %d showtimes from %d movies\n", 
		len(showtimes), len(movieLinks))
	
	return showtimes, nil
}

// extractMovieTitle extracts the movie title from a movie page
func (s *Scraper) extractMovieTitle(e *colly.HTMLElement) string {
	// Try multiple possible selectors for movie title
	selectors := []string{
		"h1.movie-title",
		"h1.title",
		"h1",
		".movie-title",
		"[data-title]",
	}

	for _, selector := range selectors {
		title := strings.TrimSpace(e.ChildText(selector))
		if title != "" {
			return title
		}
	}

	return ""
}

// extractShowtimes extracts all showtimes from a movie page (T012)
func (s *Scraper) extractShowtimes(e *colly.HTMLElement, movieTitle, format, link string) []models.Showtime {
	showtimes := []models.Showtime{}
	
	// Try to find showtime elements
	// This is a flexible approach that tries multiple patterns
	
	// Pattern 1: Showtimes grouped by date
	e.ForEach(".showtime-date, .date-group, [data-date]", func(_ int, dateElem *colly.HTMLElement) {
		dateStr := s.extractDate(dateElem)
		if dateStr == "" {
			return
		}

		// Find time elements within this date group
		dateElem.ForEach(".time, .showtime, [data-time]", func(_ int, timeElem *colly.HTMLElement) {
			timeStr := s.extractTime(timeElem)
			if timeStr != "" {
				showtime := s.createShowtime(movieTitle, dateStr, timeStr, format, link)
				if showtime != nil {
					showtimes = append(showtimes, *showtime)
				}
			}
		})
	})

	// Pattern 2: Flat list of showtime elements with date and time
	if len(showtimes) == 0 {
		e.ForEach(".showtime-item, .screening, [data-showtime]", func(_ int, elem *colly.HTMLElement) {
			dateStr := s.extractDate(elem)
			timeStr := s.extractTime(elem)
			
			if dateStr != "" && timeStr != "" {
				showtime := s.createShowtime(movieTitle, dateStr, timeStr, format, link)
				if showtime != nil {
					showtimes = append(showtimes, *showtime)
				}
			}
		})
	}

	return showtimes
}

// extractDate extracts and normalizes a date from an element
func (s *Scraper) extractDate(e *colly.HTMLElement) string {
	// Try datetime attribute first
	if dt := e.Attr("datetime"); dt != "" {
		// Check if it's already in YYYY-MM-DD format
		if len(dt) >= 10 && dt[4] == '-' && dt[7] == '-' {
			return dt[:10]
		}
	}

	// Try data-date attribute
	if dateData := e.Attr("data-date"); dateData != "" {
		return s.normalizeDate(dateData)
	}

	// Try text content
	dateText := strings.TrimSpace(e.Text)
	return s.normalizeDate(dateText)
}

// extractTime extracts and normalizes a time from an element
func (s *Scraper) extractTime(e *colly.HTMLElement) string {
	// Try datetime attribute first
	if dt := e.Attr("datetime"); dt != "" && strings.Contains(dt, "T") {
		parts := strings.Split(dt, "T")
		if len(parts) == 2 {
			timePart := strings.Split(parts[1], "-")[0] // Remove timezone
			if len(timePart) >= 5 {
				return timePart[:5] // HH:MM
			}
		}
	}

	// Try data-time attribute
	if timeData := e.Attr("data-time"); timeData != "" {
		return s.normalizeTime(timeData)
	}

	// Try text content
	timeText := strings.TrimSpace(e.Text)
	return s.normalizeTime(timeText)
}

// normalizeDate converts various date formats to YYYY-MM-DD
func (s *Scraper) normalizeDate(dateStr string) string {
	// Already in correct format
	if len(dateStr) == 10 && dateStr[4] == '-' && dateStr[7] == '-' {
		return dateStr
	}

	// Try to parse common date formats and convert to YYYY-MM-DD
	layouts := []string{
		"2006-01-02",
		"01/02/2006",
		"1/2/2006",
		"January 2, 2006",
		"Jan 2, 2006",
		"Monday, January 2, 2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	return ""
}

// normalizeTime converts various time formats to HH:MM (24-hour)
func (s *Scraper) normalizeTime(timeStr string) string {
	timeStr = strings.TrimSpace(timeStr)
	
	// Already in HH:MM format
	if len(timeStr) == 5 && timeStr[2] == ':' {
		return timeStr
	}

	// Try to parse common time formats
	layouts := []string{
		"3:04 PM",
		"3:04PM",
		"15:04",
		"3:04 pm",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t.Format("15:04")
		}
	}

	return ""
}

// createShowtime creates a Showtime model from extracted data
func (s *Scraper) createShowtime(movieTitle, date, timeStr, format, link string) *models.Showtime {
	// Validate date and time
	if date == "" || timeStr == "" {
		return nil
	}

	// Validate date is within 3-month window
	// Parse date in local timezone to match comparison
	now := time.Now()
	showtimeDate, err := time.ParseInLocation("2006-01-02", date, now.Location())
	if err != nil {
		return nil
	}

	// Compare dates only (ignore time of day)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	maxDate := today.AddDate(0, 3, 0)
	
	if showtimeDate.Before(today) || showtimeDate.After(maxDate) {
		return nil
	}

	// Generate unique ID
	id := fmt.Sprintf("%s-%s-%s", 
		s.theater.ID, 
		strings.ReplaceAll(date, "-", ""), 
		strings.ReplaceAll(timeStr, ":", ""))

	return &models.Showtime{
		ID:         id,
		TheaterID:  s.theater.ID,
		MovieTitle: movieTitle,
		Date:       date,
		Time:       timeStr,
		Format:     format,
		Link:       link,
		TMDBID:     0, // Will be filled by TMDB enrichment (T013)
	}
}

// extractShowtime parses showtime data from an HTML element
// This will be implemented in Phase 3 (T012)
func (s *Scraper) extractShowtime(movieTitle, date, time, format, link string) *models.Showtime {
	return s.createShowtime(movieTitle, date, time, format, link)
}
