package main

import (
	"flag"
	"log"
	"net/http"
	"retroblog/config"
	"retroblog/handlers"
	"retroblog/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse flags
	flag.StringVar(&config.SrcDir, "src", "./apple_notes_export", "Notes source directory")
	flag.StringVar(&config.OutDir, "out", "./public", "Output directory")
	flag.StringVar(&config.Addr, "addr", ":8080", "HTTP listen address")
	flag.StringVar(&config.BaseURL, "url", "https://firsh.me", "Base URL for sitemap")
	flag.Parse()

	// Initialize configuration
	config.Initialize()

	// Initial build
	if err := services.RebuildAll(); err != nil {
		log.Fatalf("Initial build failed: %v", err)
	}
	log.Println("Initial build done.")

	// Start file watcher
	go services.Watch()

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Setup routes
	setupRoutes(router)

	log.Printf("Serving %s at http://%s", config.OutDir, config.Addr)
	log.Printf("Site base URL: %s", config.BaseURL)
	log.Fatal(router.Run(config.Addr))
}

func setupRoutes(router *gin.Engine) {
	// API routes
	api := router.Group("/api")
	{
		api.POST("/unlock", handlers.HandleUnlock)
	}
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static") {
			c.Redirect(http.StatusMovedPermanently, "https://firsh.me/")
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})

	// Static file serving with authentication middleware
	router.Use(handlers.AttachmentAuthMiddleware())
	router.Static("/", config.OutDir)
}
