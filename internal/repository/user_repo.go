package repository

import (
	"database/sql"
	"log"

	"Bookstore/internal/models"
)

// RegisterUser Register a new user
func RegisterUser(db *sql.DB) (*models.User, error) {
	query := `INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, user.Username, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		log.Println("Error registering user:", err)
		return nil, err
	}
	return nil, err
}

// GetUserByUsername Get a user by username
func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password, role FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, err
	}
	return &user, nil
}

// DeleteBookByID Delete a book by ID
func DeleteUserByUsername(db *sql.DB, username string) error {
	query := `DELETE FROM books WHERE username = $1`
	_, err := db.Exec(query, username)
	if err != nil {
		log.Println("Error deleting book:", err)
		return err
	}
	return nil
}
