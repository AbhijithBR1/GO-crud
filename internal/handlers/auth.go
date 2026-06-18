package handlers

import (
 	"errors"
	"net/http"

	"bookmanagement/internal/models"

	"github.com/gin-gonic/gin"
)

// HandleRegister corresponds to POST /register
func HandleRegister(c *gin.Context) {
	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := models.RegisterUser(c.Request.Context(), reqBody.Username, reqBody.Password)
	if errors.Is(err, models.ErrUserExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// HandleLogin corresponds to POST /login
func HandleLogin(c *gin.Context) {
	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	token, err := models.LoginUser(c.Request.Context(), reqBody.Username, reqBody.Password)
	if errors.Is(err, models.ErrInvalidLogin) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()}) // 401 Unauthorized
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log in"})
		return
	}

	// We return the token as JSON: { "token": "abc123def456" }
	c.JSON(http.StatusOK, gin.H{"token": token})
}
