package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/yourusername/themanager/database"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	PinHash   string    `json:"-"` // Don't expose the hash in JSON
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int) (*User, error) {
	var user User
	err := database.DB.QueryRow(
		"SELECT id, name, pin_hash, role, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.PinHash, &user.Role, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// ValidateUserPin checks if the provided PIN is valid for a user
func ValidateUserPin(pin string) (*User, error) {
	// Hash the provided PIN
	pinHash := hashPin(pin)

	// Query the database for a user with this PIN hash
	var user User
	err := database.DB.QueryRow(
		"SELECT id, name, pin_hash, role, created_at FROM users WHERE pin_hash = ?",
		pinHash,
	).Scan(&user.ID, &user.Name, &user.PinHash, &user.Role, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Invalid PIN, no user found
		}
		return nil, err
	}

	return &user, nil
}

// Helper function to hash a PIN
func hashPin(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(hash[:])
}
