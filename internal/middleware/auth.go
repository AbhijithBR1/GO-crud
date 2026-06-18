package middleware

import (
	"net/http"
	"strings"

	"bookmanagement/internal/models"

	"github.com/gin-gonic/gin"
)

// RequireAuth is our authentication middleware for Gin.
// It extracts the Bearer token from the Authorization header, validates it,
// and aborts the request if the token is missing or invalid.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header from the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		// 2. Check if it follows the "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format. Use 'Bearer <token>'"})
			return
		}

		// 3. Validate the token against the database
		token := parts[1]
		username, err := models.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return // Return early to stop execution (like NOT calling next() in Express)
		}

		// Store the authenticated username in the Gin context for downstream handlers.
		c.Set("username", username)

		// 4. Token is valid! Call the next handler in the chain.
		// This is equivalent to calling `next()` in Express.js.
		c.Next()
	}
}
