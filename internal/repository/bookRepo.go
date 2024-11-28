package repository

import (
	"Bookstore/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"strconv"
)

type BookRepository interface {
	CreateBook(book *models.Book) error
	GetAllBooks() ([]*models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	Update(book *models.Book) error
	DeleteBook(id int) error
}

type bookRepository struct {
	db  *sql.DB
	Log *zap.Logger
}

func NewBookRepository(db *sql.DB, logger *zap.Logger) BookRepository {
	return &bookRepository{
		db:  db,
		Log: logger,
	}
}

const (
	queryCreateBook  = "INSERT INTO books (id, title, author, price, quantity) VALUES ($1, $2, $3, $4, $5)"
	queryGetAllBooks = "SELECT id, title, author, price, quantity from books"
	queryGetBookByID = "SELECT id, title, author, price, quantity from books where id = $1"
	queryUpdateBook  = "UPDATE books SET title = $1, author = $2, price = $3, quantity = $4 WHERE id = $5"
	queryDeleteBook  = "DELETE FROM books WHERE id = $1"
)

func (r *bookRepository) CreateBook(book *models.Book) error {

	_, err := r.db.Exec(queryCreateBook, book.ID, book.Title, book.Author, book.Price, book.Quantity)
	if err != nil {
		r.Log.Error("Error when creating book", zap.String("bookTitle", book.Title), zap.Error(err))
		return err
	}
	return nil
}

func (r *bookRepository) GetAllBooks() ([]*models.Book, error) {
	rows, err := r.db.Query(queryGetAllBooks)
	if err != nil {
		r.Log.Error("Error when querying books", zap.String("query", queryGetAllBooks), zap.Error(err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}(rows)

	var books []*models.Book

	for rows.Next() {
		book := &models.Book{}
		// Используем указатели на поля для сканирования данных
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Quantity)
		if err != nil {
			r.Log.Error("Error when scanning books", zap.String("query", queryGetAllBooks), zap.Error(err))
			return nil, err
		}
		books = append(books, book)
	}

	// Проверяем, если ошибка при итерации по строкам
	if err = rows.Err(); err != nil {
		log.Printf("Error when iterating over rows: %v", err)
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) GetBookByID(id int) (*models.Book, error) {
	book := &models.Book{}

	err := r.db.QueryRow(queryGetBookByID, id).Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Log.Warn("book not found", zap.String("title", book.Title))
			return nil, err
		}
		r.Log.Error("Error when getting book", zap.String("title", book.Title), zap.Error(err))
		return nil, err
	}
	//r.Log.Info("successfully got book", zap.String("title", book.Title))
	return book, nil
}

func (r *bookRepository) Update(book *models.Book) error {
	_, err := r.db.Exec(queryUpdateBook, book.Title, book.Author, book.Price, book.Quantity, book.ID)
	if err != nil {
		log.Printf("Error when updating book: %v", err)
		return fmt.Errorf("failed to update book: %v", err)
	}
	return nil
}

// DeleteBook Delete book
func (r *bookRepository) DeleteBook(id int) error {
	_, err := r.db.Exec(queryDeleteBook, id)
	if err != nil {
		r.Log.Error("Error when deleting book", zap.String("id", strconv.Itoa(id)))
		return fmt.Errorf("unsuccess to delete book: %v", err)
	}
	return nil
}
