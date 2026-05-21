package main

import (
	"fmt"
	"log"
	"net/http"

	"bookmanagement/internal/handlers"
	"bookmanagement/internal/middleware"
)

func main() {
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
