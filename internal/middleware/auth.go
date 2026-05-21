package middleware

import (
	"net/http"
	"strings"

	"bookmanagement/internal/models"
)

// RequireAuth is our authentication middleware.
// In Go, middleware is a function that takes an http.HandlerFunc and returns an http.HandlerFunc.
// This is exactly like Express.js where middleware takes (req, res, next).
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	// We return a new anonymous function that matches the http.HandlerFunc signature
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the Authorization header from the request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Check if it follows the "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format. Use 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		// 3. Validate the token against our map in models
		token := parts[1]
		_, err := models.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return // Return early to stop execution (like NOT calling next() in Express)
		}

		// 4. Token is valid! Call the next handler in the chain.
		// This is equivalent to calling `next()` in Express.js.
		next(w, r)
	}
}
