package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bookmanagement/internal/models"
)

// HandleGetBooks corresponds to GET /books
func HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	books := models.GetAllBooks()
	
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

	book, err := models.GetBookByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	book := models.CreateBook(reqBody.Title, reqBody.Author)
	
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

	book, err := models.UpdateBook(id, reqBody.Title, reqBody.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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

	if err := models.DeleteBook(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Return HTTP 204 No Content
}
