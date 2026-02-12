package cinemagic_theater

import (
	"testing"
	"time"
)

func TestNewScraper(t *testing.T) {
	scraper := NewScraper()

	if scraper == nil {
		t.Fatal("NewScraper() returned nil")
	}

	theater := scraper.GetTheaterInfo()
	if theater.ID != "cinemagic-theater" {
		t.Errorf("Theater ID = %q, want %q", theater.ID, "cinemagic-theater")
	}

	if theater.Name != "Cinemagic Theater" {
		t.Errorf("Theater Name = %q, want %q", theater.Name, "Cinemagic Theater")
	}

	if theater.Address != "2021 SE Hawthorne Blvd, Portland, OR 97214" {
		t.Errorf("Theater Address = %q, want expected address", theater.Address)
	}

	if scraper.GetID() != "cinemagic-theater" {
		t.Errorf("GetID() = %q, want %q", scraper.GetID(), "cinemagic-theater")
	}
}

func TestNormalizeDate(t *testing.T) {
	scraper := NewScraper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already normalized YYYY-MM-DD",
			input:    "2026-02-12",
			expected: "2026-02-12",
		},
		{
			name:     "MM/DD/YYYY format",
			input:    "02/12/2026",
			expected: "2026-02-12",
		},
		{
			name:     "single digit month and day",
			input:    "2/12/2026",
			expected: "2026-02-12",
		},
		{
			name:     "full month name",
			input:    "February 12, 2026",
			expected: "2026-02-12",
		},
		{
			name:     "abbreviated month",
			input:    "Feb 12, 2026",
			expected: "2026-02-12",
		},
		{
			name:     "invalid date returns empty",
			input:    "not a date",
			expected: "",
		},
		{
			name:     "empty string returns empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scraper.normalizeDate(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeDate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeTime(t *testing.T) {
	scraper := NewScraper()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already normalized HH:MM",
			input:    "19:30",
			expected: "19:30",
		},
		{
			name:     "12-hour with PM",
			input:    "7:30 PM",
			expected: "19:30",
		},
		{
			name:     "12-hour with AM",
			input:    "10:15 AM",
			expected: "10:15",
		},
		{
			name:     "12-hour uppercase PM no space",
			input:    "7:30PM",
			expected: "19:30",
		},
		{
			name:     "12-hour lowercase pm",
			input:    "7:30 pm",
			expected: "19:30",
		},
		{
			name:     "noon 12-hour format",
			input:    "12:00 PM",
			expected: "12:00",
		},
		{
			name:     "midnight 12-hour format",
			input:    "12:00 AM",
			expected: "00:00",
		},
		{
			name:     "invalid time returns empty",
			input:    "not a time",
			expected: "",
		},
		{
			name:     "empty string returns empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scraper.normalizeTime(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTime(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateShowtime(t *testing.T) {
	scraper := NewScraper()

	tests := []struct {
		name       string
		movieTitle string
		date       string
		time       string
		format     string
		link       string
		wantNil    bool
		checkID    bool
	}{
		{
			name:       "valid showtime",
			movieTitle: "The Matrix",
			date:       "2026-02-15",
			time:       "19:30",
			format:     "digital",
			link:       "https://tickets.thecinemagictheater.com/movie/the-matrix",
			wantNil:    false,
			checkID:    true,
		},
		{
			name:       "valid showtime with 35mm",
			movieTitle: "Blade Runner",
			date:       "2026-03-01",
			time:       "20:00",
			format:     "35mm",
			link:       "https://tickets.thecinemagictheater.com/movie/blade-runner",
			wantNil:    false,
			checkID:    true,
		},
		{
			name:       "empty date returns nil",
			movieTitle: "Movie",
			date:       "",
			time:       "19:30",
			format:     "digital",
			link:       "https://example.com",
			wantNil:    true,
		},
		{
			name:       "empty time returns nil",
			movieTitle: "Movie",
			date:       "2026-02-15",
			time:       "",
			format:     "digital",
			link:       "https://example.com",
			wantNil:    true,
		},
		{
			name:       "invalid date format returns nil",
			movieTitle: "Movie",
			date:       "not-a-date",
			time:       "19:30",
			format:     "digital",
			link:       "https://example.com",
			wantNil:    true,
		},
		{
			name:       "date in past returns nil",
			movieTitle: "Movie",
			date:       "2020-01-01",
			time:       "19:30",
			format:     "digital",
			link:       "https://example.com",
			wantNil:    true,
		},
		{
			name:       "date beyond 3-month window returns nil",
			movieTitle: "Movie",
			date:       "2026-12-31",
			time:       "19:30",
			format:     "digital",
			link:       "https://example.com",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scraper.createShowtime(tt.movieTitle, tt.date, tt.time, tt.format, tt.link)

			if tt.wantNil {
				if result != nil {
					t.Errorf("createShowtime() = %v, want nil", result)
				}
				return
			}

			if result == nil {
				t.Fatal("createShowtime() returned nil, want valid showtime")
			}

			if result.TheaterID != "cinemagic-theater" {
				t.Errorf("TheaterID = %q, want %q", result.TheaterID, "cinemagic-theater")
			}

			if result.MovieTitle != tt.movieTitle {
				t.Errorf("MovieTitle = %q, want %q", result.MovieTitle, tt.movieTitle)
			}

			if result.Date != tt.date {
				t.Errorf("Date = %q, want %q", result.Date, tt.date)
			}

			if result.Time != tt.time {
				t.Errorf("Time = %q, want %q", result.Time, tt.time)
			}

			if result.Format != tt.format {
				t.Errorf("Format = %q, want %q", result.Format, tt.format)
			}

			if result.Link != tt.link {
				t.Errorf("Link = %q, want %q", result.Link, tt.link)
			}

			if result.TMDBID != 0 {
				t.Errorf("TMDBID = %d, want 0 (enrichment happens later)", result.TMDBID)
			}

			if tt.checkID && result.ID == "" {
				t.Error("ID is empty, want non-empty ID")
			}
		})
	}
}

func TestShowtimeIDGeneration(t *testing.T) {
	scraper := NewScraper()

	// Test that IDs are unique for different times
	showtime1 := scraper.createShowtime("Movie", "2026-02-15", "19:00", "digital", "link")
	showtime2 := scraper.createShowtime("Movie", "2026-02-15", "21:00", "digital", "link")

	if showtime1.ID == showtime2.ID {
		t.Errorf("IDs should be unique for different times, both got %q", showtime1.ID)
	}

	// Test that IDs are unique for different dates
	showtime3 := scraper.createShowtime("Movie", "2026-02-16", "19:00", "digital", "link")

	if showtime1.ID == showtime3.ID {
		t.Errorf("IDs should be unique for different dates, both got %q", showtime1.ID)
	}
}

func TestDateWindowValidation(t *testing.T) {
	scraper := NewScraper()
	now := time.Now()

	tests := []struct {
		name    string
		date    string
		wantNil bool
	}{
		{
			name:    "today is valid",
			date:    now.Format("2006-01-02"),
			wantNil: false,
		},
		{
			name:    "tomorrow is valid",
			date:    now.AddDate(0, 0, 1).Format("2006-01-02"),
			wantNil: false,
		},
		{
			name:    "1 month ahead is valid",
			date:    now.AddDate(0, 1, 0).Format("2006-01-02"),
			wantNil: false,
		},
		{
			name:    "2 months ahead is valid",
			date:    now.AddDate(0, 2, 0).Format("2006-01-02"),
			wantNil: false,
		},
		{
			name:    "exactly 3 months ahead is valid",
			date:    now.AddDate(0, 3, 0).Format("2006-01-02"),
			wantNil: false,
		},
		{
			name:    "just over 3 months is invalid",
			date:    now.AddDate(0, 3, 1).Format("2006-01-02"),
			wantNil: true,
		},
		{
			name:    "yesterday is invalid",
			date:    now.AddDate(0, 0, -1).Format("2006-01-02"),
			wantNil: true,
		},
		{
			name:    "1 month ago is invalid",
			date:    now.AddDate(0, -1, 0).Format("2006-01-02"),
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scraper.createShowtime("Test Movie", tt.date, "19:00", "digital", "link")

			if tt.wantNil && result != nil {
				t.Errorf("createShowtime() with date %q should return nil (outside 3-month window)", tt.date)
			}

			if !tt.wantNil && result == nil {
				t.Errorf("createShowtime() with date %q should return valid showtime (within 3-month window)", tt.date)
			}
		})
	}
}

func TestExtractMovieTitle(t *testing.T) {
	// This would require mocking colly.HTMLElement
	// For now, we test the logic indirectly through integration tests
	// or by creating a testable wrapper function
	t.Skip("Requires colly.HTMLElement mocking - covered by integration tests")
}

func TestExtractShowtimes(t *testing.T) {
	// This would require mocking colly.HTMLElement
	t.Skip("Requires colly.HTMLElement mocking - covered by integration tests")
}

func TestScrape(t *testing.T) {
	// This would require mocking HTTP responses
	// Integration test that hits live website or uses recorded fixtures
	t.Skip("Integration test - requires HTTP mocking or live website")
}
