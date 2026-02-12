package main

import (
	"fmt"
	"log"
	"time"

	"theater-showtimes/internal/api"
	"theater-showtimes/internal/scrapers"
	"theater-showtimes/internal/scrapers/clinton_street_theater"
	"theater-showtimes/internal/scrapers/example_theater"
	"theater-showtimes/internal/scrapers/local_cinema"
	"theater-showtimes/internal/storage"
	"theater-showtimes/internal/tmdb"
)

func main() {
	// Initialize storage
	store, err := storage.NewStorage("./data")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize TMDB client (7-day cache)
	tmdbClient := tmdb.NewClient(168 * time.Hour)

	// Initialize scraper registry and register all scrapers
	registry := scrapers.NewRegistry()
	registry.Register(clinton_street_theater.NewScraper())
	registry.Register(example_theater.NewScraper())
	registry.Register(local_cinema.NewScraper())

	// Initialize theaters in storage
	theaters := []interface{}{}
	for _, scraper := range registry.GetAll() {
		theaters = append(theaters, scraper.GetTheaterInfo())
	}

	// Initialize API handler
	handler := api.NewHandler(store, registry, tmdbClient)

	// Setup and start server
	router := api.SetupRouter(handler)

	port := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", port)
	fmt.Printf("Available scrapers: %v\n", registry.GetIDs())
	
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
