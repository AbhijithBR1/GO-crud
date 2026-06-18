package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"bookmanagement/internal/models"

	"github.com/gin-gonic/gin"
)

// HandleGetBooks corresponds to GET /books
func HandleGetBooks(c *gin.Context) {
	books, err := models.GetAllBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	// json.NewEncoder(w).Encode(books) is equivalent to writing JSON.stringify(books)
	c.JSON(http.StatusOK, books)
}

// HandleGetBookByID corresponds to GET /books/:id
func HandleGetBookByID(c *gin.Context) {
	// c.Param gets the :id from the URL defined in our router
	idStr := c.Param("id")

	// Convert the string ID to an integer (like parseInt in JS)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	book, err := models.GetBookByID(c.Request.Context(), id)
	if errors.Is(err, models.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// HandleCreateBook corresponds to POST /books
func HandleCreateBook(c *gin.Context) {
	// A temporary struct to hold the incoming JSON (like destructuring a JS object)
	var reqBody struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	// c.ShouldBindJSON is equivalent to JSON.parse(request.body)
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	book, err := models.CreateBook(c.Request.Context(), reqBody.Title, reqBody.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book) // Return HTTP 201 Created
}

// HandleUpdateBook corresponds to PUT /books/:id
func HandleUpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var reqBody struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	book, err := models.UpdateBook(c.Request.Context(), id, reqBody.Title, reqBody.Author)
	if errors.Is(err, models.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// HandleDeleteBook corresponds to DELETE /books/:id
func HandleDeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = models.DeleteBook(c.Request.Context(), id)
	if errors.Is(err, models.ErrBookNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.Status(http.StatusNoContent) // Return HTTP 204 No Content
}
