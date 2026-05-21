package models

import (
	"errors"
	"sync"
)

// Book represents our data structure.
// The `json:"id"` part tells the Go JSON encoder how to map fields to JSON,
// which is similar to how object keys work in JavaScript.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// In-memory data store.
// A mutex (sync.Mutex) is used to prevent issues when multiple requests
// try to read/write the 'books' array at the same time.
var (
	books = []Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan"},
		{ID: 2, Title: "JavaScript: The Good Parts", Author: "Douglas Crockford"},
	}
	nextID = 3
	mu     sync.Mutex
)

// ErrBookNotFound is a custom error, similar to throwing a new Error() in JS.
var ErrBookNotFound = errors.New("book not found")

// GetAllBooks returns a copy of our books slice.
func GetAllBooks() []Book {
	mu.Lock()
	defer mu.Unlock() // 'defer' ensures this runs right before the function returns
	return books
}

// GetBookByID searches for a book.
func GetBookByID(id int) (Book, error) {
	mu.Lock()
	defer mu.Unlock()
	
	for _, b := range books {
		if b.ID == id {
			return b, nil
		}
	}
	return Book{}, ErrBookNotFound
}

// CreateBook adds a new book to our slice.
func CreateBook(title, author string) Book {
	mu.Lock()
	defer mu.Unlock()
	
	book := Book{
		ID:     nextID,
		Title:  title,
		Author: author,
	}
	nextID++
	books = append(books, book) // append is like Array.push() in JS
	return book
}

// UpdateBook finds and modifies a book.
func UpdateBook(id int, title, author string) (Book, error) {
	mu.Lock()
	defer mu.Unlock()
	
	for i, b := range books {
		if b.ID == id {
			books[i].Title = title
			books[i].Author = author
			return books[i], nil
		}
	}
	return Book{}, ErrBookNotFound
}

// DeleteBook removes a book from the slice.
func DeleteBook(id int) error {
	mu.Lock()
	defer mu.Unlock()
	
	for i, b := range books {
		if b.ID == id {
			// This is Go's way to 'splice' an array
			// We take everything before the item, and append everything after it
			books = append(books[:i], books[i+1:]...)
			return nil
		}
	}
	return ErrBookNotFound
}
