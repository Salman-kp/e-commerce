package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"e-commerce/config"
	"e-commerce/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to DB
	config.ConnectDatabase()
	config.MigrateAll()

	// Create router
	router := gin.Default()
   
	// Make DB accessible in handlers via context if desired
	// router.Use(func(c *gin.Context) {
	// 	c.Set("db", config.DB)
	// 	c.Next()
	// })

	// Load templates & static
	//router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/view/login")
	})
   
	routes.AuthRoutes(router)
    routes.UserRoutes(router)
	routes.AdminRoutes(router)
    routes.AdminViewRoutes(router)
    routes.ProductRoutes(router)


	
	// Server port from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running at http://localhost%s", addr)
	router.Run(addr)
}
