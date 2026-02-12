package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"theater-showtimes/internal/models"
)

// Storage handles JSON file-based data persistence
type Storage struct {
	dataPath string
	mu       sync.RWMutex
}

// NewStorage creates a new storage instance
func NewStorage(dataPath string) (*Storage, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &Storage{
		dataPath: dataPath,
	}, nil
}

// SaveTheaters saves theaters to JSON
func (s *Storage) SaveTheaters(theaters []models.Theater) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dataPath, "theaters.json")
	return s.writeJSON(path, theaters)
}

// LoadTheaters loads theaters from JSON
func (s *Storage) LoadTheaters() ([]models.Theater, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.dataPath, "theaters.json")
	var theaters []models.Theater
	err := s.readJSON(path, &theaters)
	return theaters, err
}

// SaveShowtimes saves showtimes to JSON
func (s *Storage) SaveShowtimes(showtimes []models.Showtime) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dataPath, "showtimes.json")
	return s.writeJSON(path, showtimes)
}

// LoadShowtimes loads showtimes from JSON
func (s *Storage) LoadShowtimes() ([]models.Showtime, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.dataPath, "showtimes.json")
	var showtimes []models.Showtime
	err := s.readJSON(path, &showtimes)
	return showtimes, err
}

// SaveMovies saves movies to JSON
func (s *Storage) SaveMovies(movies []models.Movie) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.dataPath, "movies.json")
	return s.writeJSON(path, movies)
}

// LoadMovies loads movies from JSON
func (s *Storage) LoadMovies() ([]models.Movie, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.dataPath, "movies.json")
	
	// Try loading as map first (saved by scraper)
	var moviesMap map[string]*models.Movie
	if err := s.readJSON(path, &moviesMap); err == nil {
		// Convert map to slice
		movies := make([]models.Movie, 0, len(moviesMap))
		for _, movie := range moviesMap {
			if movie != nil {
				movies = append(movies, *movie)
			}
		}
		return movies, nil
	}
	
	// Fallback to array format
	var movies []models.Movie
	err := s.readJSON(path, &movies)
	return movies, err
}

// SaveMetadata saves scrape metadata
func (s *Storage) SaveMetadata(metadata models.ScrapeMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Load existing metadata
	path := filepath.Join(s.dataPath, "metadata.json")
	var allMetadata []models.ScrapeMetadata
	_ = s.readJSON(path, &allMetadata)

	// Add new metadata
	allMetadata = append(allMetadata, metadata)

	// Keep only last 100 entries
	if len(allMetadata) > 100 {
		allMetadata = allMetadata[len(allMetadata)-100:]
	}

	return s.writeJSON(path, allMetadata)
}

// GetLastUpdate returns the most recent scrape timestamp
func (s *Storage) GetLastUpdate() (time.Time, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.dataPath, "metadata.json")
	var allMetadata []models.ScrapeMetadata
	
	if err := s.readJSON(path, &allMetadata); err != nil {
		return time.Time{}, err
	}

	if len(allMetadata) == 0 {
		return time.Time{}, nil
	}

	return allMetadata[len(allMetadata)-1].LastUpdated, nil
}

// Helper function to write JSON to file
func (s *Storage) writeJSON(path string, data interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// Helper function to read JSON from file
func (s *Storage) readJSON(path string, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Return empty data if file doesn't exist
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}
