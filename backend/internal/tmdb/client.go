package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"theater-showtimes/internal/models"
)

// Client handles TMDB API communication
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	cache      *Cache
}

// NewClient creates a new TMDB client
func NewClient(cacheTTL time.Duration) *Client {
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		apiKey = "4574f5c4ac74dc981d3e6a6491ebf58d" // Default from config
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.themoviedb.org/3",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: NewCache(cacheTTL),
	}
}

// SearchMovie searches for a movie by title
func (c *Client) SearchMovie(title string) (*models.Movie, error) {
	// Check cache first
	if cached := c.cache.Get(title); cached != nil {
		return cached, nil
	}

	// Normalize title for search
	normalizedTitle := c.normalizeTitle(title)

	// TODO: Implement actual MCP server communication
	// This is a placeholder that should be replaced with actual MCP server calls
	movie, err := c.searchTMDB(normalizedTitle)
	if err != nil {
		return nil, err
	}

	// Cache the result
	c.cache.Set(title, movie)

	return movie, nil
}

// normalizeTitle cleans up the title for better matching
func (c *Client) normalizeTitle(title string) string {
	// Remove special characters, extra whitespace, etc.
	cleaned := strings.TrimSpace(title)
	cleaned = strings.ToLower(cleaned)
	return cleaned
}

// searchTMDB performs the actual TMDB search
func (c *Client) searchTMDB(title string) (*models.Movie, error) {
	// Build request URL
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("query", title)
	params.Add("language", "en-US")
	params.Add("page", "1")
	
	requestURL := fmt.Sprintf("%s/search/movie?%s", c.baseURL, params.Encode())

	// Make request
	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search TMDB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TMDB search failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Results []struct {
			ID           int     `json:"id"`
			Title        string  `json:"title"`
			OriginalTitle string  `json:"original_title"`
			Overview     string  `json:"overview"`
			ReleaseDate  string  `json:"release_date"`
			PosterPath   string  `json:"poster_path"`
			BackdropPath string  `json:"backdrop_path"`
			VoteAverage  float64 `json:"vote_average"`
			VoteCount    int     `json:"vote_count"`
			Popularity   float64 `json:"popularity"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no results found for: %s", title)
	}

	// Convert first result to our Movie model
	tmdbMovie := result.Results[0]
	movie := &models.Movie{
		TMDBID:        tmdbMovie.ID,
		Title:         tmdbMovie.Title,
		OriginalTitle: tmdbMovie.OriginalTitle,
		Overview:      tmdbMovie.Overview,
		ReleaseDate:   tmdbMovie.ReleaseDate,
		PosterPath:    c.buildPosterURL(tmdbMovie.PosterPath),
		BackdropPath:  c.buildBackdropURL(tmdbMovie.BackdropPath),
		TMDBRating:    tmdbMovie.VoteAverage,
		VoteCount:     tmdbMovie.VoteCount,
		Popularity:    tmdbMovie.Popularity,
	}

	// Fetch additional details
	if err := c.enrichMovieDetails(movie); err != nil {
		fmt.Printf("Warning: failed to enrich movie details: %v\n", err)
	}

	return movie, nil
}

// buildPosterURL constructs the full TMDB poster image URL
func (c *Client) buildPosterURL(posterPath string) string {
	if posterPath == "" {
		return ""
	}
	// TMDB image base URL with w500 size for movie cards
	return fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", posterPath)
}

// buildBackdropURL constructs the full TMDB backdrop image URL
func (c *Client) buildBackdropURL(backdropPath string) string {
	if backdropPath == "" {
		return ""
	}
	return fmt.Sprintf("https://image.tmdb.org/t/p/w1280%s", backdropPath)
}

// enrichMovieDetails fetches additional details like runtime, rating, genres, cast
func (c *Client) enrichMovieDetails(movie *models.Movie) error {
	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("append_to_response", "credits,release_dates")
	
	requestURL := fmt.Sprintf("%s/movie/%d?%s", c.baseURL, movie.TMDBID, params.Encode())

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var details struct {
		Runtime int      `json:"runtime"`
		Genres  []struct {
			Name string `json:"name"`
		} `json:"genres"`
		Credits struct {
			Cast []struct {
				Name      string `json:"name"`
				Character string `json:"character"`
				Order     int    `json:"order"`
			} `json:"cast"`
			Crew []struct {
				Name string `json:"name"`
				Job  string `json:"job"`
			} `json:"crew"`
		} `json:"credits"`
		ReleaseDates struct {
			Results []struct {
				ISO31661     string `json:"iso_3166_1"`
				ReleaseDates []struct {
					Certification string `json:"certification"`
				} `json:"release_dates"`
			} `json:"results"`
		} `json:"release_dates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return err
	}

	// Set runtime
	movie.Runtime = details.Runtime

	// Set genres
	movie.Genres = make([]string, 0, len(details.Genres))
	for _, genre := range details.Genres {
		movie.Genres = append(movie.Genres, genre.Name)
	}

	// Set top cast (first 5)
	movie.Cast = make([]string, 0)
	for i, castMember := range details.Credits.Cast {
		if i >= 5 {
			break
		}
		movie.Cast = append(movie.Cast, castMember.Name)
	}

	// Find director
	for _, crewMember := range details.Credits.Crew {
		if crewMember.Job == "Director" {
			movie.Director = crewMember.Name
			break
		}
	}

	// Find US rating
	for _, result := range details.ReleaseDates.Results {
		if result.ISO31661 == "US" && len(result.ReleaseDates) > 0 {
			movie.Rating = result.ReleaseDates[0].Certification
			break
		}
	}
	if movie.Rating == "" {
		movie.Rating = "NR"
	}

	return nil
}

// GetMovieDetails fetches detailed information for a movie by TMDB ID
func (c *Client) GetMovieDetails(tmdbID int) (*models.Movie, error) {
	cacheKey := fmt.Sprintf("id-%d", tmdbID)
	
	// Check cache
	if cached := c.cache.Get(cacheKey); cached != nil {
		return cached, nil
	}

	params := url.Values{}
	params.Add("api_key", c.apiKey)
	params.Add("append_to_response", "credits,release_dates")
	
	requestURL := fmt.Sprintf("%s/movie/%d?%s", c.baseURL, tmdbID, params.Encode())

	resp, err := c.httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TMDB get movie failed with status %d: %s", resp.StatusCode, string(body))
	}

	var details struct {
		ID           int     `json:"id"`
		Title        string  `json:"title"`
		OriginalTitle string  `json:"original_title"`
		Overview     string  `json:"overview"`
		Runtime      int     `json:"runtime"`
		ReleaseDate  string  `json:"release_date"`
		PosterPath   string  `json:"poster_path"`
		BackdropPath string  `json:"backdrop_path"`
		VoteAverage  float64 `json:"vote_average"`
		VoteCount    int     `json:"vote_count"`
		Popularity   float64 `json:"popularity"`
		Genres       []struct {
			Name string `json:"name"`
		} `json:"genres"`
		Credits struct {
			Cast []struct {
				Name      string `json:"name"`
				Character string `json:"character"`
				Order     int    `json:"order"`
			} `json:"cast"`
			Crew []struct {
				Name string `json:"name"`
				Job  string `json:"job"`
			} `json:"crew"`
		} `json:"credits"`
		ReleaseDates struct {
			Results []struct {
				ISO31661     string `json:"iso_3166_1"`
				ReleaseDates []struct {
					Certification string `json:"certification"`
				} `json:"release_dates"`
			} `json:"results"`
		} `json:"release_dates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to decode movie details: %w", err)
	}

	movie := &models.Movie{
		TMDBID:        details.ID,
		Title:         details.Title,
		OriginalTitle: details.OriginalTitle,
		Overview:      details.Overview,
		Runtime:       details.Runtime,
		ReleaseDate:   details.ReleaseDate,
		PosterPath:    c.buildPosterURL(details.PosterPath),
		BackdropPath:  c.buildBackdropURL(details.BackdropPath),
		TMDBRating:    details.VoteAverage,
		VoteCount:     details.VoteCount,
		Popularity:    details.Popularity,
	}

	// Set genres
	movie.Genres = make([]string, 0, len(details.Genres))
	for _, genre := range details.Genres {
		movie.Genres = append(movie.Genres, genre.Name)
	}

	// Set top cast (first 5)
	movie.Cast = make([]string, 0)
	for i, castMember := range details.Credits.Cast {
		if i >= 5 {
			break
		}
		movie.Cast = append(movie.Cast, castMember.Name)
	}

	// Find director
	for _, crewMember := range details.Credits.Crew {
		if crewMember.Job == "Director" {
			movie.Director = crewMember.Name
			break
		}
	}

	// Find US rating
	for _, result := range details.ReleaseDates.Results {
		if result.ISO31661 == "US" && len(result.ReleaseDates) > 0 {
			movie.Rating = result.ReleaseDates[0].Certification
			break
		}
	}
	if movie.Rating == "" {
		movie.Rating = "NR"
	}

	// Cache the result
	c.cache.Set(cacheKey, movie)

	return movie, nil
}

// EnrichShowtimes takes a slice of showtimes and enriches each with TMDB data
func (c *Client) EnrichShowtimes(showtimes []models.Showtime) ([]models.Showtime, map[string]*models.Movie) {
	enriched := make([]models.Showtime, 0, len(showtimes))
	movieCache := make(map[string]*models.Movie)

	for _, showtime := range showtimes {
		// Check if we already looked up this movie
		if movie, exists := movieCache[showtime.MovieTitle]; exists {
			if movie != nil {
				showtime.TMDBID = movie.TMDBID
			}
			enriched = append(enriched, showtime)
			continue
		}

		// Search for the movie
		movie, err := c.SearchMovie(showtime.MovieTitle)
		if err != nil {
			// Create placeholder movie with limited_info flag (T014)
			fmt.Printf("Failed to find TMDB data for '%s': %v\n", showtime.MovieTitle, err)
			placeholderMovie := c.createPlaceholderMovie(showtime.MovieTitle)
			movieCache[showtime.MovieTitle] = placeholderMovie
			showtime.TMDBID = 0 // Placeholder movies have TMDB ID 0
			enriched = append(enriched, showtime)
			continue
		}

		// Store movie data and update showtime
		movieCache[showtime.MovieTitle] = movie
		showtime.TMDBID = movie.TMDBID
		enriched = append(enriched, showtime)
	}

	return enriched, movieCache
}

// createPlaceholderMovie creates a minimal movie record for TMDB match failures
// This implements the limited_info flag requirement from data-model.md
func (c *Client) createPlaceholderMovie(title string) *models.Movie {
	return &models.Movie{
		TMDBID:       0,
		Title:        title,
		Overview:     "",
		Runtime:      0,
		Rating:       "NR",
		Genres:       []string{},
		PosterPath:   "/assets/placeholder-poster.png",
		TMDBRating:   0,
		VoteCount:    0,
		LimitedInfo:  true, // Flag for frontend badge display
	}
}

