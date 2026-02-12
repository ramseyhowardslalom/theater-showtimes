package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"theater-showtimes/internal/models"
	"theater-showtimes/internal/scrapers"
	"theater-showtimes/internal/scrapers/clinton_street_theater"
	"theater-showtimes/internal/scrapers/example_theater"
	"theater-showtimes/internal/scrapers/local_cinema"
	"theater-showtimes/internal/storage"
	"theater-showtimes/internal/tmdb"
)

func main() {
	// Initialize scraper registry
	registry := scrapers.NewRegistry()
	registry.Register(clinton_street_theater.NewScraper())
	registry.Register(example_theater.NewScraper())
	registry.Register(local_cinema.NewScraper())

	// Initialize storage
	store, err := storage.NewStorage("./data")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize TMDB client
	tmdbClient := tmdb.NewClient(168 * time.Hour)

	// Parse command line arguments
	args := os.Args[1:]
	
	var scrapersToRun []scrapers.Scraper
	
	if len(args) == 0 {
		// Run all scrapers
		fmt.Println("No scrapers specified, running all...")
		for _, scraper := range registry.GetAll() {
			scrapersToRun = append(scrapersToRun, scraper)
		}
	} else {
		// Run specified scrapers
		for _, arg := range args {
			scraper, exists := registry.Get(arg)
			if !exists {
				fmt.Printf("Warning: Scraper '%s' not found, skipping...\n", arg)
				continue
			}
			scrapersToRun = append(scrapersToRun, scraper)
		}
	}

	if len(scrapersToRun) == 0 {
		fmt.Println("No valid scrapers to run")
		fmt.Printf("Available scrapers: %s\n", strings.Join(registry.GetIDs(), ", "))
		os.Exit(1)
	}

	// Run scrapers
	fmt.Printf("Running %d scraper(s)...\n", len(scrapersToRun))
	
	allShowtimes := []models.Showtime{}
	allMovies := make(map[string]*models.Movie)
	
	for _, scraper := range scrapersToRun {
		fmt.Printf("\n=== Scraping %s ===\n", scraper.GetTheaterInfo().Name)
		
		showtimes, err := scraper.Scrape()
		if err != nil {
			log.Printf("Error scraping %s: %v", scraper.GetID(), err)
			continue
		}
		
		fmt.Printf("Found %d showtimes\n", len(showtimes))
		
		// Enrich with TMDB data
		if len(showtimes) > 0 {
			fmt.Println("\nEnriching with TMDB data...")
			enriched, movieData := tmdbClient.EnrichShowtimes(showtimes)
			showtimes = enriched
			
			// Merge movie data
			for title, movie := range movieData {
				if movie != nil {
					allMovies[title] = movie
				}
			}
			
			// Display scraped showtimes with TMDB data
			fmt.Println("\nShowtimes found:")
			for _, st := range showtimes {
				movie := movieData[st.MovieTitle]
				if movie != nil {
					fmt.Printf("  • %s (%s) - %s @ %s\n", st.MovieTitle, movie.ReleaseDate[:4], st.Date, st.Time)
					fmt.Printf("    TMDB ID: %d, Rating: %.1f/10, Poster: %s\n", movie.TMDBID, movie.TMDBRating, movie.PosterPath)
				} else {
					fmt.Printf("  • %s - %s @ %s (No TMDB data)\n", st.MovieTitle, st.Date, st.Time)
				}
			}
		}
		
		allShowtimes = append(allShowtimes, showtimes...)
	}

	// Save results
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total showtimes: %d\n", len(allShowtimes))
	fmt.Printf("Unique movies: %d\n", len(allMovies))
	
	// Save to JSON files
	if len(allShowtimes) > 0 {
		showtimesJSON, _ := json.MarshalIndent(allShowtimes, "", "  ")
		if err := os.WriteFile("./data/showtimes.json", showtimesJSON, 0644); err != nil {
			log.Printf("Failed to save showtimes: %v", err)
		} else {
			fmt.Println("\n✓ Saved showtimes to ./data/showtimes.json")
		}
	}
	
	if len(allMovies) > 0 {
		moviesJSON, _ := json.MarshalIndent(allMovies, "", "  ")
		if err := os.WriteFile("./data/movies.json", moviesJSON, 0644); err != nil {
			log.Printf("Failed to save movies: %v", err)
		} else {
			fmt.Println("✓ Saved movie data to ./data/movies.json")
		}
	}
	
	_ = store // Storage interface for future use
}
