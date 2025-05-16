package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB() error {
	// Get the user's home directory for storing the database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create the application data directory if it doesn't exist
	appDataDir := filepath.Join(homeDir, ".themanager")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return fmt.Errorf("failed to create app data directory: %w", err)
	}

	// Database file path
	dbPath := filepath.Join(appDataDir, "store.db")
	log.Printf("Database path: %s", dbPath)

	// Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connection established")

	// Create tables if they don't exist
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Create default admin user if no users exist
	if err := createDefaultUser(); err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	return nil
}

// createTables creates all necessary tables if they don't exist
func createTables() error {
	// Users table
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		pin_hash TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Products table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		barcode TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Stock table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS stock (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		stock_group TEXT NOT NULL,
		quantity REAL NOT NULL,
		cost_price REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id)
	)`)
	if err != nil {
		return err
	}

	// Clients table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS clients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		contact TEXT,
		address TEXT,
		credit_limit REAL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Transactions table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_type TEXT NOT NULL,
		client_id INTEGER,
		total_amount REAL NOT NULL,
		is_paid BOOLEAN NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (client_id) REFERENCES clients(id)
	)`)
	if err != nil {
		return err
	}

	// Transaction Items table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS transaction_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		stock_id INTEGER NOT NULL,
		quantity REAL NOT NULL,
		price REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (transaction_id) REFERENCES transactions(id),
		FOREIGN KEY (product_id) REFERENCES products(id),
		FOREIGN KEY (stock_id) REFERENCES stock(id)
	)`)
	if err != nil {
		return err
	}

	// Payments table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS payments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_id INTEGER NOT NULL,
		amount REAL NOT NULL,
		payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (transaction_id) REFERENCES transactions(id)
	)`)
	if err != nil {
		return err
	}

	log.Println("All tables created successfully")
	return nil
}

// createDefaultUser creates a default admin user if no users exist
func createDefaultUser() error {
	// Check if any users exist
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	// If no users exist, create a default admin user with PIN 1234
	if count == 0 {
		// In a real app, you would use a proper password hashing function
		// For simplicity, we're using a simple hash here
		defaultPinHash := "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4" // SHA-256 hash of "1234"
		_, err := DB.Exec(
			"INSERT INTO users (name, pin_hash, role) VALUES (?, ?, ?)",
			"Admin", defaultPinHash, "admin",
		)
		if err != nil {
			return err
		}
		log.Println("Default admin user created with PIN: 1234")
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
