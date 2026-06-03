package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"bookmanagement/internal/models"
)

// HandleGetBooks corresponds to GET /books
func HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetAllBooks(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(books) is equivalent to writing JSON.stringify(books)
	json.NewEncoder(w).Encode(books)
}

// HandleGetBookByID corresponds to GET /books/{id}
func HandleGetBookByID(w http.ResponseWriter, r *http.Request) {
	// r.PathValue gets the {id} from the URL defined in our router
	idStr := r.PathValue("id")

	// Convert the string ID to an integer (like parseInt in JS)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	book, err := models.GetBookByID(r.Context(), id)
	if errors.Is(err, models.ErrBookNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to fetch book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// HandleCreateBook corresponds to POST /books
func HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	// A temporary struct to hold the incoming JSON (like destructuring a JS object)
	var reqBody struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}

	// json.NewDecoder(r.Body).Decode() is equivalent to JSON.parse(request.body)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	book, err := models.CreateBook(r.Context(), reqBody.Title, reqBody.Author)
	if err != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Return HTTP 201 Created
	json.NewEncoder(w).Encode(book)
}

// HandleUpdateBook corresponds to PUT /books/{id}
func HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	book, err := models.UpdateBook(r.Context(), id, reqBody.Title, reqBody.Author)
	if errors.Is(err, models.ErrBookNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// HandleDeleteBook corresponds to DELETE /books/{id}
func HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = models.DeleteBook(r.Context(), id)
	if errors.Is(err, models.ErrBookNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Return HTTP 204 No Content
}
