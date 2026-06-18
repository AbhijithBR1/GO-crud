package handlers

import (
	"net/http"

	"bookmanagement/internal/database"

	"github.com/gin-gonic/gin"
)

// HandleHealth corresponds to GET /health.
// Hosting platforms use this to confirm the app is alive AND its DB is reachable.
// Returns 200 + {"status":"ok"} on success, 503 + {"status":"db unreachable"} otherwise.
func HandleHealth(c *gin.Context) {
	if err := database.Pool.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db unreachable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
