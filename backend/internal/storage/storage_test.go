package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"theater-showtimes/internal/models"
)

func TestNewStorage(t *testing.T) {
	tempDir := t.TempDir()
	dataPath := filepath.Join(tempDir, "data")

	storage, err := NewStorage(dataPath)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	if storage == nil {
		t.Fatal("NewStorage() returned nil")
	}

	// Verify directory was created
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		t.Errorf("Data directory was not created: %s", dataPath)
	}
}

func TestReplaceTheaterShowtimes(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// Create initial showtimes from different theaters
	initialShowtimes := []models.Showtime{
		{
			ID:         "show1",
			TheaterID:  "theater-a",
			MovieTitle: "Movie A",
			Date:       "2026-02-15",
			Time:       "19:00",
			Format:     "digital",
		},
		{
			ID:         "show2",
			TheaterID:  "theater-a",
			MovieTitle: "Movie B",
			Date:       "2026-02-15",
			Time:       "21:00",
			Format:     "35mm",
		},
		{
			ID:         "show3",
			TheaterID:  "theater-b",
			MovieTitle: "Movie C",
			Date:       "2026-02-16",
			Time:       "20:00",
			Format:     "digital",
		},
	}

	// Save initial showtimes
	err = storage.SaveShowtimes(initialShowtimes)
	if err != nil {
		t.Fatalf("SaveShowtimes() error = %v", err)
	}

	// Create new showtimes for theater-a (replacement)
	newShowtimes := []models.Showtime{
		{
			ID:         "show4",
			TheaterID:  "theater-a",
			MovieTitle: "Movie D",
			Date:       "2026-02-17",
			Time:       "19:30",
			Format:     "digital",
		},
		{
			ID:         "show5",
			TheaterID:  "theater-a",
			MovieTitle: "Movie E",
			Date:       "2026-02-17",
			Time:       "22:00",
			Format:     "70mm",
		},
	}

	// Replace theater-a showtimes
	err = storage.ReplaceTheaterShowtimes("theater-a", newShowtimes)
	if err != nil {
		t.Fatalf("ReplaceTheaterShowtimes() error = %v", err)
	}

	// Load all showtimes and verify
	allShowtimes, err := storage.LoadShowtimes()
	if err != nil {
		t.Fatalf("LoadShowtimes() error = %v", err)
	}

	// Should have 3 showtimes total: 2 new theater-a + 1 theater-b
	if len(allShowtimes) != 3 {
		t.Errorf("Expected 3 showtimes, got %d", len(allShowtimes))
	}

	// Verify old theater-a showtimes are gone
	for _, st := range allShowtimes {
		if st.ID == "show1" || st.ID == "show2" {
			t.Errorf("Old theater-a showtime still present: %s", st.ID)
		}
	}

	// Verify new theater-a showtimes are present
	foundShow4 := false
	foundShow5 := false
	foundShow3 := false

	for _, st := range allShowtimes {
		switch st.ID {
		case "show4":
			foundShow4 = true
			if st.MovieTitle != "Movie D" {
				t.Errorf("show4 MovieTitle = %q, want %q", st.MovieTitle, "Movie D")
			}
		case "show5":
			foundShow5 = true
			if st.Format != "70mm" {
				t.Errorf("show5 Format = %q, want %q", st.Format, "70mm")
			}
		case "show3":
			foundShow3 = true
			if st.TheaterID != "theater-b" {
				t.Errorf("show3 TheaterID = %q, want %q", st.TheaterID, "theater-b")
			}
		}
	}

	if !foundShow4 {
		t.Error("New showtime show4 not found")
	}
	if !foundShow5 {
		t.Error("New showtime show5 not found")
	}
	if !foundShow3 {
		t.Error("Theater-b showtime show3 was removed (should be preserved)")
	}
}

func TestReplaceTheaterShowtimes_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// Replace showtimes when file doesn't exist yet
	newShowtimes := []models.Showtime{
		{
			ID:         "show1",
			TheaterID:  "theater-a",
			MovieTitle: "Movie A",
			Date:       "2026-02-15",
			Time:       "19:00",
			Format:     "digital",
		},
	}

	err = storage.ReplaceTheaterShowtimes("theater-a", newShowtimes)
	if err != nil {
		t.Fatalf("ReplaceTheaterShowtimes() error = %v", err)
	}

	// Load and verify
	allShowtimes, err := storage.LoadShowtimes()
	if err != nil {
		t.Fatalf("LoadShowtimes() error = %v", err)
	}

	if len(allShowtimes) != 1 {
		t.Errorf("Expected 1 showtime, got %d", len(allShowtimes))
	}

	if allShowtimes[0].ID != "show1" {
		t.Errorf("Showtime ID = %q, want %q", allShowtimes[0].ID, "show1")
	}
}

func TestReplaceTheaterShowtimes_EmptyReplacement(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// Create initial showtimes
	initialShowtimes := []models.Showtime{
		{
			ID:         "show1",
			TheaterID:  "theater-a",
			MovieTitle: "Movie A",
			Date:       "2026-02-15",
			Time:       "19:00",
			Format:     "digital",
		},
		{
			ID:         "show2",
			TheaterID:  "theater-b",
			MovieTitle: "Movie B",
			Date:       "2026-02-15",
			Time:       "21:00",
			Format:     "35mm",
		},
	}

	err = storage.SaveShowtimes(initialShowtimes)
	if err != nil {
		t.Fatalf("SaveShowtimes() error = %v", err)
	}

	// Replace with empty slice (remove all theater-a showtimes)
	err = storage.ReplaceTheaterShowtimes("theater-a", []models.Showtime{})
	if err != nil {
		t.Fatalf("ReplaceTheaterShowtimes() error = %v", err)
	}

	// Load and verify
	allShowtimes, err := storage.LoadShowtimes()
	if err != nil {
		t.Fatalf("LoadShowtimes() error = %v", err)
	}

	// Should only have theater-b showtime
	if len(allShowtimes) != 1 {
		t.Errorf("Expected 1 showtime, got %d", len(allShowtimes))
	}

	if len(allShowtimes) > 0 && allShowtimes[0].TheaterID != "theater-b" {
		t.Errorf("Remaining showtime TheaterID = %q, want %q", allShowtimes[0].TheaterID, "theater-b")
	}
}

func TestSaveAndLoadShowtimes(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// Create test showtimes
	testShowtimes := []models.Showtime{
		{
			ID:         "show1",
			TheaterID:  "theater-a",
			MovieTitle: "The Matrix",
			Date:       "2026-02-15",
			Time:       "19:00",
			Format:     "digital",
			TMDBID:     603,
			Link:       "https://example.com/tickets/1",
		},
		{
			ID:         "show2",
			TheaterID:  "theater-a",
			MovieTitle: "Blade Runner",
			Date:       "2026-02-16",
			Time:       "21:00",
			Format:     "35mm",
			TMDBID:     78,
			Link:       "https://example.com/tickets/2",
		},
	}

	// Save showtimes
	err = storage.SaveShowtimes(testShowtimes)
	if err != nil {
		t.Fatalf("SaveShowtimes() error = %v", err)
	}

	// Load showtimes
	loadedShowtimes, err := storage.LoadShowtimes()
	if err != nil {
		t.Fatalf("LoadShowtimes() error = %v", err)
	}

	// Verify count
	if len(loadedShowtimes) != len(testShowtimes) {
		t.Errorf("Loaded %d showtimes, want %d", len(loadedShowtimes), len(testShowtimes))
	}

	// Verify data integrity
	for i, st := range loadedShowtimes {
		expected := testShowtimes[i]
		if st.ID != expected.ID {
			t.Errorf("Showtime[%d].ID = %q, want %q", i, st.ID, expected.ID)
		}
		if st.MovieTitle != expected.MovieTitle {
			t.Errorf("Showtime[%d].MovieTitle = %q, want %q", i, st.MovieTitle, expected.MovieTitle)
		}
		if st.TMDBID != expected.TMDBID {
			t.Errorf("Showtime[%d].TMDBID = %d, want %d", i, st.TMDBID, expected.TMDBID)
		}
		if st.Format != expected.Format {
			t.Errorf("Showtime[%d].Format = %q, want %q", i, st.Format, expected.Format)
		}
	}
}

func TestSaveAndLoadMovies(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// Create test movies
	testMovies := []models.Movie{
		{
			TMDBID:      603,
			Title:       "The Matrix",
			ReleaseDate: "1999-03-31",
			PosterPath:  "/poster1.jpg",
			Overview:    "A computer hacker learns...",
			LimitedInfo: false,
		},
		{
			TMDBID:      0,
			Title:       "Unknown Film",
			ReleaseDate: "",
			PosterPath:  "/placeholder.png",
			Overview:    "",
			LimitedInfo: true,
		},
	}

	// Save movies
	err = storage.SaveMovies(testMovies)
	if err != nil {
		t.Fatalf("SaveMovies() error = %v", err)
	}

	// Load movies
	loadedMovies, err := storage.LoadMovies()
	if err != nil {
		t.Fatalf("LoadMovies() error = %v", err)
	}

	// Verify count
	if len(loadedMovies) != len(testMovies) {
		t.Errorf("Loaded %d movies, want %d", len(loadedMovies), len(testMovies))
	}

	// Verify LimitedInfo flag is preserved
	var foundLimitedInfo bool
	for _, movie := range loadedMovies {
		if movie.LimitedInfo {
			foundLimitedInfo = true
			if movie.TMDBID != 0 {
				t.Errorf("LimitedInfo movie should have TMDBID=0, got %d", movie.TMDBID)
			}
		}
	}

	if !foundLimitedInfo {
		t.Error("LimitedInfo flag was not preserved")
	}
}

func TestLoadShowtimes_NonexistentFile(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// Load from nonexistent file should return empty slice
	showtimes, err := storage.LoadShowtimes()
	if err != nil {
		t.Fatalf("LoadShowtimes() error = %v", err)
	}

	if len(showtimes) != 0 {
		t.Errorf("Expected empty slice, got %d showtimes", len(showtimes))
	}
}

func TestConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewStorage(tempDir)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}

	// Initial save
	initialShowtimes := []models.Showtime{
		{
			ID:         "show1",
			TheaterID:  "theater-a",
			MovieTitle: "Movie A",
			Date:       "2026-02-15",
			Time:       "19:00",
			Format:     "digital",
		},
	}
	err = storage.SaveShowtimes(initialShowtimes)
	if err != nil {
		t.Fatalf("SaveShowtimes() error = %v", err)
	}

	// Test concurrent reads and writes
	done := make(chan bool)

	// Multiple concurrent readers
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				_, err := storage.LoadShowtimes()
				if err != nil {
					t.Errorf("Concurrent LoadShowtimes() error = %v", err)
				}
				time.Sleep(time.Millisecond)
			}
			done <- true
		}()
	}

	// Multiple concurrent writers
	for i := 0; i < 3; i++ {
		theaterID := i
		go func() {
			for j := 0; j < 5; j++ {
				showtimes := []models.Showtime{
					{
						ID:         fmt.Sprintf("show-%d-%d", theaterID, j),
						TheaterID:  fmt.Sprintf("theater-%d", theaterID),
						MovieTitle: "Test Movie",
						Date:       "2026-02-15",
						Time:       "19:00",
						Format:     "digital",
					},
				}
				err := storage.ReplaceTheaterShowtimes(fmt.Sprintf("theater-%d", theaterID), showtimes)
				if err != nil {
					t.Errorf("Concurrent ReplaceTheaterShowtimes() error = %v", err)
				}
				time.Sleep(time.Millisecond * 2)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 8; i++ {
		<-done
	}

	// Final verification - should have data from 3 theaters
	finalShowtimes, err := storage.LoadShowtimes()
	if err != nil {
		t.Fatalf("Final LoadShowtimes() error = %v", err)
	}

	if len(finalShowtimes) < 3 {
		t.Errorf("Expected at least 3 showtimes after concurrent writes, got %d", len(finalShowtimes))
	}
}
