package main

import (
	"net/http"
	"web-service-gin-docker/goCache"

	"github.com/gin-gonic/gin"
)

func main() {
	
	router := gin.Default()

	// Añadimos los endpoints de nuestra API
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	// Iniciamos el servidor
	router.Run()
}

// getAlbums devuelve la lista de álbums en JSON
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, goCache.GetAlbumManager().GetAlbums())
}

// postAlbums añade un objeto álbum al listado de memoria
func postAlbums(c *gin.Context) {
	var newAlbum goCache.Album

	// BindJSON completa el objeto newAlbum con los datos del JSON
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Añadimos el álbum al listado en memoria
	goCache.GetAlbumManager().AddAlbum(&newAlbum)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID devuelve un álbum dado su ID
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	album, err := goCache.GetAlbumManager().GetAlbum(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}
