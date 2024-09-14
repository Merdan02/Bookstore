package repository

import (
	"Bookstore/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type BookRepository interface {
	CreateBook(book *models.Book) error
	GetAllBooks() ([]*models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{
		db: db,
	}
}
func (r *bookRepository) CreateBook(book *models.Book) error {
	query := "INSERT INTO books (id, title, author, price, quantity) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, book.ID, book.Title, book.Author, book.Price, book.Quantity)
	if err != nil {
		log.Printf("Error when creating book: %v", err)
		return err
	}
	log.Printf("Created book with id: %v", book.ID)
	return nil
}

func (r *bookRepository) GetAllBooks() ([]*models.Book, error) {
	return nil, nil
}

func (r *bookRepository) GetBookByID(id int) (*models.Book, error) {
	query := "SELECT id, title, author, price, quantity from books where id = $1"
	book := &models.Book{}
	err := r.db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error when getting book: %v", err)
		return nil, fmt.Errorf("Unsuccess to get book: %v", err)
	}
	return book, nil
}

// Update book
func (r *bookRepository) UpdateBook(book *models.Book) error {
	query := "UPDATE books SET title = $1, author = $2, quantity = $3 WHERE id = $4"
	_, err := r.db.Exec(query, book.Title, book.Author, book.Quantity, book.ID)
	if err != nil {
		log.Printf("Error when updating book: %v", err)
		return fmt.Errorf("unsuccess to update book: %v", err)
	}
	return nil
}

// Delete book
func (r *bookRepository) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error when deleting book: %v", err)
		return fmt.Errorf("unsuccess to delete book: %v", err)
	}
	return nil
}
