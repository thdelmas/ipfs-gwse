package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thdelmas/ipfs-gwse/api/handlers"
)

func main() {
	router := gin.Default()
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	router.GET("/:cid", handlers.HandleMetadata)
	router.Run(":8000")
}
