package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", handler.Health)
		api.GET("/last-updated", handler.GetLastUpdated)
		
		api.GET("/theaters", handler.GetTheaters)
		api.GET("/showtimes", handler.GetShowtimes)
		api.GET("/showtimes/:theater", handler.GetTheaterShowtimes)
		
		api.GET("/movies", handler.GetMovies)
		api.GET("/movies/:id", handler.GetMovieDetails)
		
		api.POST("/scrape", handler.TriggerScrape)
	}

	return router
}
