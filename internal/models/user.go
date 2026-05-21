package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
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
	// Maps are Go's built-in key-value stores (like new Map() or {} in JS)
	users   = make(map[string]User) // Keyed by username
	nextUID = 1
	tokens  = make(map[string]string) // Key: token string, Value: username
	userMu  sync.Mutex
)

var (
	ErrUserExists   = errors.New("username already exists")
	ErrInvalidLogin = errors.New("invalid username or password")
	ErrInvalidToken = errors.New("invalid or expired token")
)

// RegisterUser creates a new user.
func RegisterUser(username, password string) (User, error) {
	userMu.Lock()
	defer userMu.Unlock()

	// Check if the key already exists in the map
	if _, exists := users[username]; exists {
		return User{}, ErrUserExists
	}

	u := User{
		ID:       nextUID,
		Username: username,
		// Note: In a real app, you MUST hash this password using bcrypt!
		Password: password, 
	}
	nextUID++
	users[username] = u
	return u, nil
}

// LoginUser verifies credentials and returns a token string.
func LoginUser(username, password string) (string, error) {
	userMu.Lock()
	defer userMu.Unlock()

	user, exists := users[username]
	if !exists || user.Password != password {
		return "", ErrInvalidLogin
	}

	// Generate a simple 16-byte random token
	bytes := make([]byte, 16)
	rand.Read(bytes) // Fills the byte array with random data
	token := hex.EncodeToString(bytes)

	// Save token session in our map
	tokens[token] = username
	return token, nil
}

// ValidateToken checks if a token exists in our map.
func ValidateToken(token string) (string, error) {
	userMu.Lock()
	defer userMu.Unlock()

	username, exists := tokens[token]
	if !exists {
		return "", ErrInvalidToken
	}
	return username, nil
}
