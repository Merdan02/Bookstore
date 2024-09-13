package repository

import (
	"database/sql"
	"log"

	"Bookstore/internal/models"
)

func CreateBook(db *sql.DB, book *models.Book) error {
	Query := `INSERT INTO books (title, author, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.QueryRow(Query, book.Title, book.Author, book.Price, book.Quantity).Scan(&book.ID)
	if err != nil {
		log.Println("Error creating book:", err)
		return err
	}
	return nil
}

func GetBooks(db *sql.DB) ([]models.Book, error) {
	rows, err := db.Query(`SELECT id, title, author, price, quantity FROM books`)
	if err != nil {
		log.Println("Error fetching books:", err)
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Quantity); err != nil {
			log.Println("Error scanning book:", err)
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// GetBookByID retrieves a book by its ID.
func GetBookByID(db *sql.DB, id int) (*models.Book, error) {
	var book models.Book
	query := `SELECT id, title, author, price, quantity FROM books WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Quantity)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// UpdateBook updates the details of a book.
func UpdateBook(db *sql.DB, book *models.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, price = $3, quantity = $4
		WHERE id = $5
	`
	_, err := db.Exec(query, book.Title, book.Author, book.Price, book.Quantity, book.ID)
	return err
}

// DeleteBookByID deletes a book by its ID.
func DeleteBookByID(db *sql.DB, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
