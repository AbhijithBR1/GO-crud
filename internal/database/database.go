package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is the shared connection pool used across the app.
// Think of it like a single exported db client you'd import everywhere in Node.
var Pool *pgxpool.Pool

// Connect reads DATABASE_URL, opens a pool, and verifies it with a ping.
func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("creating connection pool: %w", err)
	}

	// A pool is lazy — it doesn't actually connect until first use.
	// Ping forces a real connection so we fail fast on startup if the DB is unreachable.
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return fmt.Errorf("pinging database: %w", err)
	}

	Pool = pool
	return nil
}

// Close shuts down the pool. Call it with defer in main.
func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

// schemaSQL creates the tables if they don't already exist and seeds a couple
// of starter books. It's safe to run on every startup.
const schemaSQL = `
CREATE TABLE IF NOT EXISTS books (
	id     SERIAL PRIMARY KEY,
	title  TEXT NOT NULL,
	author TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
	id       SERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
	token      TEXT PRIMARY KEY,
	username   TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
	expires_at TIMESTAMPTZ NOT NULL
);

-- Idempotent migration: add expires_at to a pre-existing tokens table.
-- The DEFAULT only applies to rows that already exist; new inserts set it explicitly.
ALTER TABLE tokens
	ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '7 days');

INSERT INTO books (title, author)
SELECT seed.title, seed.author
FROM (VALUES
	('The Go Programming Language', 'Alan A. A. Donovan'),
	('JavaScript: The Good Parts', 'Douglas Crockford')
) AS seed(title, author)
WHERE NOT EXISTS (SELECT 1 FROM books);
`

// InitSchema runs the schema SQL above. Call it once on startup after Connect.
func InitSchema(ctx context.Context) error {
	_, err := Pool.Exec(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("initializing schema: %w", err)
	}
	return nil
}
