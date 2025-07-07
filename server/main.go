package main

import (
	"gin/analyzer"

	"github.com/gin-gonic/gin"
)

func main() {
	analyzer.InitDB()

	router := gin.Default()
	router.POST("/links", analyzer.AddLinks)
	router.GET("/links", analyzer.GetLinks)
	router.PUT("/links", analyzer.UpdateLink)
	router.Run(":8080")
}
