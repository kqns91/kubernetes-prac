package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin でサーバー立てる
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
