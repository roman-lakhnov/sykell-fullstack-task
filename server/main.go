package main

import (
	"gin/analyzer"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	analyzer.InitDB()
	analyzer.StartAnalyzerWorker()

	router := gin.Default()

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},  // Your React app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/links", analyzer.AddLinks)
	router.GET("/links", analyzer.GetLinks)
	router.PUT("/links", analyzer.UpdateLink)
	router.Run(":8080")
}
