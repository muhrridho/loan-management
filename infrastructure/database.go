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
	CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_id INTEGER,
    loan_id INTEGER NOT NULL,
    due_date DATE,
    payment_no INTEGER,
    amount REAL,
    interest REAL,
    total_amount REAL,
    status INTEGER,
    paid_at TIMESTAMP,
    created_at TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id)
		FOREIGN KEY (transaction_id) REFERENCES transactions(id)
	);
	CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    total_amount REAL,
    penalty REAL,
    status INTEGER,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
