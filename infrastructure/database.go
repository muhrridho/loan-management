package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Initialize() (*sql.DB, error) {
	var err error

	DB, err = sql.Open("sqlite", "loans")
	if err != nil {
		return nil, fmt.Errorf("Failed accessing db: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return DB, nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		name TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS loans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    interest REAL,
    interest_type INTEGER,
    tenure INTEGER,
    tenure_type INTEGER,
    amount REAL,
    outstanding REAL,
    status INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    billing_start_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("Migration is failed: %w", err)
	}

	log.Println("Migration is success")
	return nil
}

func Destroy() error {
	CloseDB()

	err := os.Remove("loans")
	if err != nil {
		return fmt.Errorf("Reinitialized DB: %w", err)
	}

	return nil
}

func Seed() error {
	// Insert a test user into the users table
	query := `
	INSERT INTO users (email, name) VALUES (?, ?)
	`
	_, err := DB.Exec(query, "test@test", "test")
	if err != nil {
		return fmt.Errorf("Failed to seed database: %w", err)
	}

	log.Println("Database seeded successfully")
	return nil
}
