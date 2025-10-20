package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"retroblog/config"
	"retroblog/handlers"
	"retroblog/services"
	"strings"
	"syscall"

	"github.com/fvbock/endless"

	"github.com/gin-gonic/gin"
)

func preSigUsr1() {
	log.Println("pre SIGUSR1")
}

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
	if os.Getenv("BLOG_WEB_BRANCH") == "release" {
		server := endless.NewServer(config.Addr, router)
		server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR1] = append(
			server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGUSR1],
			preSigUsr1)
		server.BeforeBegin = func(add string) {
			log.Printf("Actual pid is %d", syscall.Getpid())
			// save it somehow
			//
		}
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	} else {
		log.Fatal(router.Run(config.Addr))
	}
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
	router.Use(handlers.AttachmentAuthMiddleware(), handlers.AutoVideoMiddleware(config.OutDir), handlers.SharedArrayBufferMiddleware())
	router.Static("/", config.OutDir)

}
