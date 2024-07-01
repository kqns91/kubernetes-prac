package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	engine.GET("/metrics", func(c *gin.Context) {
		c.String(200, "Hello, Prometheus!")
	})

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
