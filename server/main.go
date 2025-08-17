package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lnwdevelopers007/job-applier-3000/server/internal/controller"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	router := gin.Default()
	router.GET("/", getAlbums)
	router.GET("/jobs", controller.GetJobs)
	router.POST("/jobs", controller.PostJobs)

	// router.POST("/movies/aggregations", mongo.AggregateMovies)
	_ = godotenv.Load()
	url := os.Getenv("URL")
	if url == "" {
		url = "localhost:8080"
	}

	router.Run(url)
}
