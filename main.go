package main

import (
	"net/http"
	"web-service-gin-docker/goCache"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	// We set the API endpoints
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	// We start the server
	router.Run()
}

// getAlbums returns a JSON representation of all albums in the album manager.
// The JSON result is written to the http response.
func getAlbums(c *gin.Context) {
	c.JSON(http.StatusOK, goCache.GetAlbumManager().GetAlbums())
}

// The postAlbums function is called to process an HTTP POST request
// to add a new album to the application.
// The new album is added to the in-memory data structure.
// The new album is then returned to the caller as JSON.
func postAlbums(c *gin.Context) {
	// Declare a variable to hold the new album data.
	var newAlbum goCache.Album

	// Call BindJSON to decode the JSON request data.
	// If there are any errors, return an HTTP 400 Bad Request error.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the album manager.
	goCache.GetAlbumManager().AddAlbum(&newAlbum)

	// Return the new album as JSON in the response.
	c.JSON(http.StatusCreated, newAlbum)
}

// Gets an album by ID.
// The ID is taken from the URL, and the album is retrieved from the cache.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	album, err := goCache.GetAlbumManager().GetAlbum(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.JSON(http.StatusOK, album)
}
