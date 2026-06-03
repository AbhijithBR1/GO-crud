package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"bookmanagement/internal/database"
	"bookmanagement/internal/handlers"
	"bookmanagement/internal/middleware"

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

	mux := http.NewServeMux()

	// === Public Routes ===
	
	// Authentication
	mux.HandleFunc("POST /register", handlers.HandleRegister)
	mux.HandleFunc("POST /login", handlers.HandleLogin)
	
	// Anyone can see the list of books
	mux.HandleFunc("GET /books", handlers.HandleGetBooks)


	// === Protected Routes ===
	// We wrap these handler functions with our RequireAuth middleware.
	
	mux.HandleFunc("GET /books/{id}", middleware.RequireAuth(handlers.HandleGetBookByID))
	mux.HandleFunc("POST /books", middleware.RequireAuth(handlers.HandleCreateBook))
	mux.HandleFunc("PUT /books/{id}", middleware.RequireAuth(handlers.HandleUpdateBook))
	mux.HandleFunc("DELETE /books/{id}", middleware.RequireAuth(handlers.HandleDeleteBook))

	fmt.Println("Server is running on port 8080...")
	fmt.Println("Try visiting: http://localhost:8080/books")
	
	// Start the server
	log.Fatal(http.ListenAndServe(":8080", mux))
}
