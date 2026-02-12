package clinton_street_theater

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"theater-showtimes/internal/models"
)

// Scraper implements the scraper for Clinton Street Theater
type Scraper struct {
	theater models.Theater
}

// NewScraper creates a new Clinton Street Theater scraper
func NewScraper() *Scraper {
	return &Scraper{
		theater: models.Theater{
			ID:      "clinton-street-theater",
			Name:    "Clinton Street Theater",
			Address: "2522 SE Clinton Street",
			City:    "Portland",
			Zip:     "97202",
			Website: "https://cstpdx.com",
			Phone:   "(971) 808-3331",
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
		colly.AllowedDomains("cstpdx.com", "www.cstpdx.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; TheaterShowtimesBot/1.0)"),
	)

	// Rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*cstpdx.com*",
		RandomDelay: 2 * time.Second,
	})

	// Extract event data from calendar month view
	c.OnHTML("article.tribe-events-calendar-month__calendar-event", func(eventElem *colly.HTMLElement) {
		showtime := s.extractCalendarShowtime(eventElem)
		if showtime != nil {
			showtimes = append(showtimes, *showtime)
		}
	})

	// Fallback: Extract from list view if calendar doesn't work
	c.OnHTML(".tribe-events-calendar-list__event", func(e *colly.HTMLElement) {
		showtime := s.extractShowtime(e)
		if showtime != nil {
			showtimes = append(showtimes, *showtime)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error scraping %s: %v\n", r.Request.URL, err)
	})

	// Scrape current month and next 2 months for complete schedule
	now := time.Now()
	monthsToScrape := 3
	
	for i := 0; i < monthsToScrape; i++ {
		targetDate := now.AddDate(0, i, 0)
		monthURL := fmt.Sprintf("%s/schedule/month/%d-%02d/", 
			s.theater.Website, 
			targetDate.Year(), 
			targetDate.Month())
		
		fmt.Printf("Scraping month: %s\n", targetDate.Format("January 2006"))
		err := c.Visit(monthURL)
		if err != nil {
			fmt.Printf("Warning: failed to scrape %s: %v\n", monthURL, err)
		}
	}

	return showtimes, nil
}

// extractShowtime parses an event element and returns a Showtime if it's a movie screening
func (s *Scraper) extractShowtime(e *colly.HTMLElement) *models.Showtime {
	// Extract movie title (clean up year and special tags)
	rawTitle := e.ChildText(".tribe-events-calendar-list__event-title-link")
	if rawTitle == "" {
		rawTitle = e.ChildText("h3")
	}
	if rawTitle == "" {
		return nil
	}

	movieTitle := s.cleanMovieTitle(rawTitle)
	
	// Skip non-movie events
	if s.isNonMovieEvent(movieTitle, rawTitle) {
		return nil
	}

	// Extract date and time
	dateTimeStr := e.ChildText(".tribe-events-calendar-list__event-datetime")
	if dateTimeStr == "" {
		dateTimeStr = e.ChildText("time")
	}
	
	date, showTime := s.parseDateTime(dateTimeStr)
	if date == "" || showTime == "" {
		return nil
	}

	// Extract event page link
	eventLink := e.ChildAttr(".tribe-events-calendar-list__event-title-link", "href")
	if eventLink == "" {
		eventLink = e.ChildAttr("a", "href")
	}

	// Extract price
	priceStr := e.ChildText(".tribe-events-c-small-cta__price")
	if priceStr == "" {
		priceStr = e.ChildText("a[href*='square.site']")
	}
	price := s.parsePrice(priceStr)

	// Generate unique ID
	id := fmt.Sprintf("%s-%s-%s", s.theater.ID, s.sanitizeForID(movieTitle), strings.ReplaceAll(date+showTime, ":", ""))

	return &models.Showtime{
		ID:         id,
		TheaterID:  s.theater.ID,
		MovieTitle: movieTitle,
		Date:       date,
		Time:       showTime,
		Format:     "digital", // CST shows digital format
		Price:      price,
		Link:       eventLink,
	}
}

// extractCalendarShowtime parses an event from the calendar view
func (s *Scraper) extractCalendarShowtime(e *colly.HTMLElement) *models.Showtime {
	// Extract the event link and title
	rawTitle := e.ChildText(".tribe-events-calendar-month__calendar-event-title a")
	if rawTitle == "" {
		rawTitle = e.ChildText("a")
	}
	if rawTitle == "" {
		return nil
	}

	movieTitle := s.cleanMovieTitle(rawTitle)
	
	// Skip non-movie events
	if s.isNonMovieEvent(movieTitle, rawTitle) {
		return nil
	}

	// Extract time from the time element's datetime attribute
	timeStr := e.ChildAttr(".tribe-events-calendar-month__calendar-event-datetime time", "datetime")
	if timeStr == "" {
		timeStr = "19:00" // Default to 7 PM
	}

	// Extract date from parent day cell
	// The parent .tribe-events-calendar-month__day has the date in a time element
	dateStr := ""
	e.DOM.ParentsFiltered(".tribe-events-calendar-month__day").Each(func(_ int, sel *goquery.Selection) {
		sel.Find("time[datetime]").Each(func(_ int, timeSel *goquery.Selection) {
			if dt, exists := timeSel.Attr("datetime"); exists && len(dt) == 10 { // YYYY-MM-DD format
				dateStr = dt
			}
		})
	})

	if dateStr == "" {
		return nil
	}

	// Extract event page link
	eventLink := e.ChildAttr(".tribe-events-calendar-month__calendar-event-title a", "href")
	if eventLink == "" {
		eventLink = e.ChildAttr("a", "href")
	}

	// Generate unique ID
	id := fmt.Sprintf("%s-%s-%s", s.theater.ID, s.sanitizeForID(movieTitle), strings.ReplaceAll(dateStr+timeStr, ":", ""))

	return &models.Showtime{
		ID:         id,
		TheaterID:  s.theater.ID,
		MovieTitle: movieTitle,
		Date:       dateStr,
		Time:       timeStr,
		Format:     "digital",
		Price:      10.0, // Default price for CST
		Link:       eventLink,
	}
}

// cleanMovieTitle removes year annotations and special tags from the title
func (s *Scraper) cleanMovieTitle(title string) string {
	cleaned := title
	
	// Handle festival presentations: "Festival presents: Movie Title" -> "Movie Title"
	if strings.Contains(cleaned, "presents:") {
		parts := strings.Split(cleaned, "presents:")
		if len(parts) > 1 {
			cleaned = strings.TrimSpace(parts[1])
		}
	}
	
	// Handle special guest screenings: "Movie Title with Guest" -> "Movie Title"
	if strings.Contains(strings.ToLower(cleaned), " with ") {
		// Find the position of " with " (case-insensitive)
		lowerCleaned := strings.ToLower(cleaned)
		withIndex := strings.Index(lowerCleaned, " with ")
		if withIndex > 0 {
			cleaned = strings.TrimSpace(cleaned[:withIndex])
		}
	}
	
	// Remove year in parentheses: "Movie Title (1999)" -> "Movie Title"
	yearPattern := regexp.MustCompile(`\s*\(\d{4}\)\s*`)
	cleaned = yearPattern.ReplaceAllString(cleaned, "")
	
	// Remove special tags like "(Church of Film)"
	tagPattern := regexp.MustCompile(`\s*\([^)]*\)\s*$`)
	cleaned = tagPattern.ReplaceAllString(cleaned, "")
	
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}

// isNonMovieEvent checks if this is a live event rather than a movie screening
func (s *Scraper) isNonMovieEvent(title, rawTitle string) bool {
	nonMovieKeywords := []string{
		"comedy night",
		"live concert",
		"drag show",
		"stand up",
		"standup",
		"performance",
	}
	
	lowerTitle := strings.ToLower(title)
	lowerRaw := strings.ToLower(rawTitle)
	
	for _, keyword := range nonMovieKeywords {
		if strings.Contains(lowerTitle, keyword) || strings.Contains(lowerRaw, keyword) {
			return true
		}
	}
	
	return false
}

// parseDateTime extracts date and time from strings like "Wednesday, February 11 @ 7:00 PM"
func (s *Scraper) parseDateTime(dateTimeStr string) (string, string) {
	if dateTimeStr == "" {
		return "", ""
	}

	// Pattern: "Day, Month DD @ HH:MM AM/PM"
	parts := strings.Split(dateTimeStr, "@")
	if len(parts) != 2 {
		return "", ""
	}

	datePart := strings.TrimSpace(parts[0])
	timePart := strings.TrimSpace(parts[1])

	// Parse the date (e.g., "Wednesday, February 11")
	date := s.parseDate(datePart)
	
	// Parse the time (e.g., "7:00 PM")
	showTime := s.parseTime(timePart)

	return date, showTime
}

// parseDate converts "Wednesday, February 11" to "2026-02-11" format
func (s *Scraper) parseDate(datePart string) string {
	// Remove day of week
	parts := strings.Split(datePart, ",")
	if len(parts) < 2 {
		return ""
	}

	monthDay := strings.TrimSpace(parts[1])
	
	// Parse month and day
	fields := strings.Fields(monthDay)
	if len(fields) < 2 {
		return ""
	}

	monthName := fields[0]
	dayStr := fields[1]

	// Map month name to number
	monthMap := map[string]string{
		"January": "01", "February": "02", "March": "03", "April": "04",
		"May": "05", "June": "06", "July": "07", "August": "08",
		"September": "09", "October": "10", "November": "11", "December": "12",
	}

	monthNum, exists := monthMap[monthName]
	if !exists {
		return ""
	}

	// Pad day with leading zero if needed
	day := dayStr
	if len(day) == 1 {
		day = "0" + day
	}

	// Assume current year (could be improved with logic to handle year boundaries)
	year := time.Now().Year()

	return fmt.Sprintf("%d-%s-%s", year, monthNum, day)
}

// parseTime converts "7:00 PM" to "19:00" format
func (s *Scraper) parseTime(timePart string) string {
	timePart = strings.TrimSpace(timePart)
	
	// Check for AM/PM
	isPM := strings.HasSuffix(strings.ToUpper(timePart), "PM")
	isAM := strings.HasSuffix(strings.ToUpper(timePart), "AM")
	
	if !isAM && !isPM {
		return timePart // Return as-is if no AM/PM
	}

	// Remove AM/PM
	timePart = strings.TrimSuffix(strings.TrimSuffix(timePart, " PM"), " AM")
	timePart = strings.TrimSuffix(strings.TrimSuffix(timePart, " pm"), " am")
	timePart = strings.TrimSpace(timePart)

	// Split hour and minute
	parts := strings.Split(timePart, ":")
	if len(parts) != 2 {
		return ""
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return ""
	}

	minute := parts[1]

	// Convert to 24-hour format
	if isPM && hour != 12 {
		hour += 12
	} else if isAM && hour == 12 {
		hour = 0
	}

	return fmt.Sprintf("%02d:%s", hour, minute)
}

// parsePrice extracts price from strings like "$10" or "[$10]"
func (s *Scraper) parsePrice(priceStr string) float64 {
	if priceStr == "" {
		return 0
	}

	// Remove everything except digits and decimal point
	re := regexp.MustCompile(`[^\d.]`)
	cleaned := re.ReplaceAllString(priceStr, "")

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}

	return price
}

// sanitizeForID creates a URL-safe version of the title for IDs
func (s *Scraper) sanitizeForID(title string) string {
	// Convert to lowercase
	safe := strings.ToLower(title)
	
	// Replace spaces and special chars with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	safe = re.ReplaceAllString(safe, "-")
	
	// Remove leading/trailing hyphens
	safe = strings.Trim(safe, "-")
	
	return safe
}
