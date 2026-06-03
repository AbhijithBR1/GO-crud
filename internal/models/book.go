package models

import (
	"context"
	"errors"

	"bookmanagement/internal/database"

	"github.com/jackc/pgx/v5"
)

// Book represents our data structure.
// The `json:"id"` part tells the Go JSON encoder how to map fields to JSON,
// which is similar to how object keys work in JavaScript.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// ErrBookNotFound is a custom error, similar to throwing a new Error() in JS.
var ErrBookNotFound = errors.New("book not found")

// GetAllBooks returns every book, ordered by id.
func GetAllBooks(ctx context.Context) ([]Book, error) {
	rows, err := database.Pool.Query(ctx, "SELECT id, title, author FROM books ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // always close rows to release the connection back to the pool

	books := []Book{}
	for rows.Next() {
		var b Book
		// Scan copies the columns of the current row into our struct fields,
		// in the same order they appear in the SELECT.
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

// GetBookByID searches for a single book.
func GetBookByID(ctx context.Context, id int) (Book, error) {
	var b Book
	// $1 is a placeholder — pgx sends the value separately, which prevents SQL injection.
	err := database.Pool.
		QueryRow(ctx, "SELECT id, title, author FROM books WHERE id = $1", id).
		Scan(&b.ID, &b.Title, &b.Author)
	if errors.Is(err, pgx.ErrNoRows) {
		return Book{}, ErrBookNotFound
	}
	if err != nil {
		return Book{}, err
	}
	return b, nil
}

// CreateBook inserts a new book and returns it with its DB-generated id.
func CreateBook(ctx context.Context, title, author string) (Book, error) {
	var b Book
	// RETURNING gives us back the inserted row (including the SERIAL id) in one round trip.
	err := database.Pool.
		QueryRow(ctx,
			"INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id, title, author",
			title, author,
		).
		Scan(&b.ID, &b.Title, &b.Author)
	if err != nil {
		return Book{}, err
	}
	return b, nil
}

// UpdateBook modifies an existing book and returns the updated row.
func UpdateBook(ctx context.Context, id int, title, author string) (Book, error) {
	var b Book
	err := database.Pool.
		QueryRow(ctx,
			"UPDATE books SET title = $1, author = $2 WHERE id = $3 RETURNING id, title, author",
			title, author, id,
		).
		Scan(&b.ID, &b.Title, &b.Author)
	if errors.Is(err, pgx.ErrNoRows) {
		return Book{}, ErrBookNotFound
	}
	if err != nil {
		return Book{}, err
	}
	return b, nil
}

// DeleteBook removes a book by id.
func DeleteBook(ctx context.Context, id int) error {
	// Exec is used for statements that don't return rows.
	tag, err := database.Pool.Exec(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrBookNotFound
	}
	return nil
}
