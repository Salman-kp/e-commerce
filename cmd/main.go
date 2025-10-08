package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"e-commerce/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	config.ConnectDatabase()
	//Auto migrate all models
	config.MigrateAll()

	// Create router
	router := gin.Default()

	// Load HTML templates and static files
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Root route (renders homepage)
	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "E-Commerce Home",
			"msg":   "Welcome to E-Commerce 🚀",
		})
	})

	// Server port from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running at http://localhost%s", addr)
	router.Run(addr)
}
