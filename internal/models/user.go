package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"bookmanagement/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// User represents our user data structure.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	// `json:"-"` tells the JSON encoder to never include this field in responses.
	// We don't want to accidentally leak passwords!
	Password string `json:"-"`
}

var (
	ErrUserExists   = errors.New("username already exists")
	ErrInvalidLogin = errors.New("invalid username or password")
	ErrInvalidToken = errors.New("invalid or expired token")
)

// RegisterUser creates a new user.
func RegisterUser(ctx context.Context, username, password string) (User, error) {
	var u User
	// Note: In a real app, you MUST hash this password using bcrypt before storing it!
	err := database.Pool.
		QueryRow(ctx,
			"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username",
			username, password,
		).
		Scan(&u.ID, &u.Username)
	if err != nil {
		// Postgres error 23505 is a unique-constraint violation — the username is taken.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrUserExists
		}
		return User{}, err
	}
	return u, nil
}

// LoginUser verifies credentials and returns a token string.
func LoginUser(ctx context.Context, username, password string) (string, error) {
	var storedPassword string
	err := database.Pool.
		QueryRow(ctx, "SELECT password FROM users WHERE username = $1", username).
		Scan(&storedPassword)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidLogin
	}
	if err != nil {
		return "", err
	}
	if storedPassword != password {
		return "", ErrInvalidLogin
	}

	// Generate a simple 16-byte random token.
	bytes := make([]byte, 16)
	rand.Read(bytes)
	token := hex.EncodeToString(bytes)

	// Persist the session token so it survives restarts.
	if _, err := database.Pool.Exec(ctx,
		"INSERT INTO tokens (token, username) VALUES ($1, $2)", token, username,
	); err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken checks that a token exists and returns the owning username.
func ValidateToken(ctx context.Context, token string) (string, error) {
	var username string
	err := database.Pool.
		QueryRow(ctx, "SELECT username FROM tokens WHERE token = $1", token).
		Scan(&username)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidToken
	}
	if err != nil {
		return "", err
	}
	return username, nil
}
