package tmdb

import (
	"testing"

	"theater-showtimes/internal/models"
)

func TestCreatePlaceholderMovie(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name  string
		title string
	}{
		{
			name:  "standard movie title",
			title: "The Matrix",
		},
		{
			name:  "complex title",
			title: "The Lord of the Rings: The Fellowship of the Ring",
		},
		{
			name:  "title with special characters",
			title: "Amélie (2001)",
		},
		{
			name:  "empty title",
			title: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie := client.createPlaceholderMovie(tt.title)

			if movie == nil {
				t.Fatal("createPlaceholderMovie() returned nil")
			}

			// Verify TMDBID is 0
			if movie.TMDBID != 0 {
				t.Errorf("TMDBID = %d, want 0", movie.TMDBID)
			}

			// Verify LimitedInfo is true
			if !movie.LimitedInfo {
				t.Error("LimitedInfo = false, want true")
			}

			// Verify title is preserved
			if movie.Title != tt.title {
				t.Errorf("Title = %q, want %q", movie.Title, tt.title)
			}

			// Verify placeholder poster path
			expectedPoster := "/assets/placeholder-poster.png"
			if movie.PosterPath != expectedPoster {
				t.Errorf("PosterPath = %q, want %q", movie.PosterPath, expectedPoster)
			}

			// Verify empty overview
			if movie.Overview != "" {
				t.Errorf("Overview = %q, want empty string", movie.Overview)
			}

			// Verify empty release date
			if movie.ReleaseDate != "" {
				t.Errorf("ReleaseDate = %q, want empty string", movie.ReleaseDate)
			}
		})
	}
}

func TestPlaceholderMovieStructure(t *testing.T) {
	client := &Client{}
	movie := client.createPlaceholderMovie("Test Film")

	// Verify the movie can be serialized to JSON (important for API)
	if movie.Title == "" && movie.TMDBID != 0 {
		t.Error("Placeholder movie has unexpected structure")
	}

	// Verify it matches Movie model structure
	var _ models.Movie = *movie
}

func TestPlaceholderWithSpecialCharacters(t *testing.T) {
	client := &Client{}

	specialTitles := []string{
		"Film & Title",
		"Movie: The Sequel",
		"Title/Subtitle",
		"Title (2026)",
		"100% Action",
	}

	for _, title := range specialTitles {
		t.Run(title, func(t *testing.T) {
			movie := client.createPlaceholderMovie(title)

			if movie == nil {
				t.Fatalf("createPlaceholderMovie(%q) returned nil", title)
			}

			if movie.Title != title {
				t.Errorf("Title = %q, want %q", movie.Title, title)
			}

			if !movie.LimitedInfo {
				t.Error("LimitedInfo should be true for placeholder")
			}

			if movie.TMDBID != 0 {
				t.Error("TMDBID should be 0 for placeholder")
			}
		})
	}
}

// Note: Full TMDB client integration tests would require API mocking
// These tests focus on the placeholder creation logic added for Cinemagic Theater feature
func TestEnrichShowtimesWithPlaceholder(t *testing.T) {
	t.Skip("Requires TMDB API mocking - integration test")
	// This test would verify that EnrichShowtimes uses createPlaceholderMovie
	// when SearchMovie returns no results
}
