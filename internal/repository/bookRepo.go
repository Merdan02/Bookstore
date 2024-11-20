package repository

import (
	"Bookstore/internal/models"
	"Bookstore/internal/wrong"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
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
func (r *bookRepository) CreateBook(book *models.Book) error {
	if err := r.validateBookFields(book); err != nil {
		r.Log.Warn("Error validating book", zap.Error(err))
		return err
	}

	query := "INSERT INTO books (id, title, author, price, quantity) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, book.ID, book.Title, book.Author, book.Price, book.Quantity)
	if err != nil {
		r.Log.Error("Error when creating book", zap.String("bookTitle", book.Title), zap.Error(err))
		return err
	}
	//r.Log.Info("successfully created book", zap.String("title", book.Title), zap.String("author", book.Author))
	return nil
}

func (r *bookRepository) GetAllBooks() ([]*models.Book, error) {
	query := "SELECT id, title, author, price, quantity FROM books"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error when getting all books: %v", err)
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
			log.Printf("Error when scanning book row: %v", err)
			return nil, err
		}
		books = append(books, book)
	}

	// Проверяем, если ошибка при итерации по строкам
	if err = rows.Err(); err != nil {
		log.Printf("Error when iterating over rows: %v", err)
		return nil, err
	}

	//log.Printf("Retrieved all books: %v", books)
	return books, nil
}

func (r *bookRepository) GetBookByID(id int) (*models.Book, error) {
	if id <= 0 {
		r.Log.Warn("your book id is empty")
		return nil, wrong.ErrBookIDZero
	}

	book := &models.Book{}

	query := "SELECT id, title, author, price, quantity from books where id = $1"
	err := r.db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Quantity)
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
	query := "UPDATE books SET title = $1, author = $2, price = $3, quantity = $4 WHERE id = $5"

	_, err := r.db.Exec(query, book.Title, book.Author, book.Price, book.Quantity, book.ID)
	if err != nil {
		log.Printf("Error when updating book: %v", err)
		return fmt.Errorf("failed to update book: %v", err)
	}
	return nil
}

// DeleteBook Delete book
func (r *bookRepository) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error when deleting book: %v", err)
		return fmt.Errorf("unsuccess to delete book: %v", err)
	}
	return nil
}

func (r *bookRepository) validateBookFields(book *models.Book) error {
	if book.ID <= 0 {
		return wrong.ErrBookIDZero
	}
	if book.Title == "" {
		return wrong.ErrEmptyTitle
	}
	if book.Author == "" {
		return wrong.ErrEmptyAuthor
	}
	if book.Price <= 0 {
		return wrong.ErrBookIDZero
	}
	if book.Quantity <= 0 {
		return wrong.ErrBookIDZero
	}
	return nil
}
