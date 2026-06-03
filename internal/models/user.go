package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"bookmanagement/internal/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
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

// tokenTTL is how long a session token stays valid after issue.
const tokenTTL = "7 days"

// RegisterUser creates a new user with a bcrypt-hashed password.
func RegisterUser(ctx context.Context, username, password string) (User, error) {
	// bcrypt salts and hashes the password. The result is a single string
	// like "$2a$10$..." that already contains the salt — store it as-is.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	var u User
	err = database.Pool.
		QueryRow(ctx,
			"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username",
			username, string(hash),
		).
		Scan(&u.ID, &u.Username)
	if err != nil {
		// Postgres error 23505 is a unique-constraint violation — username taken.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrUserExists
		}
		return User{}, err
	}
	return u, nil
}

// LoginUser verifies the bcrypt hash and returns a fresh token valid for tokenTTL.
func LoginUser(ctx context.Context, username, password string) (string, error) {
	var storedHash string
	err := database.Pool.
		QueryRow(ctx, "SELECT password FROM users WHERE username = $1", username).
		Scan(&storedHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidLogin
	}
	if err != nil {
		return "", err
	}

	// bcrypt does the comparison in constant time and handles the embedded salt.
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return "", ErrInvalidLogin
	}

	// Generate a simple 16-byte random token.
	bytes := make([]byte, 16)
	rand.Read(bytes)
	token := hex.EncodeToString(bytes)

	// Persist the session token with an expiry timestamp.
	if _, err := database.Pool.Exec(ctx,
		"INSERT INTO tokens (token, username, expires_at) VALUES ($1, $2, NOW() + INTERVAL '"+tokenTTL+"')",
		token, username,
	); err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken checks that a token exists AND hasn't expired, returning the owning username.
func ValidateToken(ctx context.Context, token string) (string, error) {
	var username string
	err := database.Pool.
		QueryRow(ctx,
			"SELECT username FROM tokens WHERE token = $1 AND expires_at > NOW()",
			token,
		).
		Scan(&username)
	// pgx.ErrNoRows covers both "no such token" and "token expired" — both look the same to the caller.
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidToken
	}
	if err != nil {
		return "", err
	}
	return username, nil
}
