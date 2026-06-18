package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"bookmanagement/internal/database"
	"bookmanagement/internal/handlers"
	"bookmanagement/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env into the environment, like require('dotenv').config() in Node.
	// In production you set real env vars instead, so a missing file is fine.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Open the shared connection pool and verify the DB is reachable.
	if err := database.Connect(); err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer database.Close()
	fmt.Println("Connected to Postgres!")

	// Create tables (if missing) and seed starter data.
	if err := database.InitSchema(context.Background()); err != nil {
		log.Fatal("schema init failed: ", err)
	}

	router := gin.Default()

	// Apply CORS middleware globally.
	router.Use(middleware.EnableCORS())

	// === Public Routes ===

	// Liveness/readiness check used by hosting platforms.
	router.GET("/health", handlers.HandleHealth)

	// Authentication
	router.POST("/register", handlers.HandleRegister)
	router.POST("/login", handlers.HandleLogin)

	// Anyone can see the list of books
	router.GET("/books", handlers.HandleGetBooks)

	// === Protected Routes ===
	// We use a Gin route group with our RequireAuth middleware.
	protected := router.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		protected.GET("/books/:id", handlers.HandleGetBookByID)
		protected.POST("/books", handlers.HandleCreateBook)
		protected.PUT("/books/:id", handlers.HandleUpdateBook)
		protected.DELETE("/books/:id", handlers.HandleDeleteBook)
	}

	// Read PORT from the environment — most hosting platforms (Render, Railway,
	// Fly.io, Cloud Run, etc.) inject this and require the app to bind to it.
	// Fall back to 8080 for local development.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Printf("Server is running on port %s...\n", port)
	fmt.Printf("Try visiting: http://localhost:%s/books\n", port)

	// Start the server
	log.Fatal(router.Run(addr))
}
