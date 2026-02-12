package models

import "time"

// Theater represents a movie theater
type Theater struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	Zip     string `json:"zip"`
	Website string `json:"website"`
	Phone   string `json:"phone,omitempty"`
}

// Movie represents a movie with TMDB enrichment
type Movie struct {
	TMDBID        int      `json:"tmdb_id"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"original_title"`
	Overview      string   `json:"overview"`
	Runtime       int      `json:"runtime"`
	Rating        string   `json:"rating"`
	Genres        []string `json:"genres"`
	ReleaseDate   string   `json:"release_date"`
	PosterPath    string   `json:"poster_path"`
	BackdropPath  string   `json:"backdrop_path"`
	TMDBRating    float64  `json:"tmdb_rating"`
	VoteCount     int      `json:"vote_count"`
	Popularity    float64  `json:"popularity"`
	Cast          []string `json:"cast"`
	Director      string   `json:"director"`
}

// Showtime represents a movie showtime
type Showtime struct {
	ID         string  `json:"id"`
	TheaterID  string  `json:"theater_id"`
	MovieTitle string  `json:"movie_title"`
	TMDBID     int     `json:"tmdb_id"`
	Date       string  `json:"date"`
	Time       string  `json:"time"`
	Format     string  `json:"format"`
	Price      float64 `json:"price,omitempty"`
	Link       string  `json:"link,omitempty"`
	Screen     string  `json:"screen,omitempty"`
}

// ScrapeMetadata tracks scraping status
type ScrapeMetadata struct {
	LastUpdated      time.Time `json:"last_updated"`
	TheaterID        string    `json:"theater_id"`
	Status           string    `json:"status"`
	ErrorMessage     string    `json:"error_message,omitempty"`
	MoviesScraped    int       `json:"movies_scraped"`
	ShowtimesScraped int       `json:"showtimes_scraped"`
}

// ScraperResult contains the data returned by a scraper
type ScraperResult struct {
	Theater   Theater
	Showtimes []Showtime
	Metadata  ScrapeMetadata
}
