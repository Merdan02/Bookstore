package service

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
	"Bookstore/internal/wrong"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

type BOokService interface {
	CreateBook(book *models.Book) error
	GetBookByID(id int) (*models.Book, error)
	GetAllBook() ([]*models.Book, error)
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}

type bookService struct {
	repo repository.BookRepository
	Log  *zap.Logger
}

func NewBookService(repo repository.BookRepository, logger *zap.Logger) BOokService {
	return &bookService{
		repo: repo,
		Log:  logger,
	}
}

func (s *bookService) CreateBook(book *models.Book) error {
	if err := s.validateBookFields(book); err != nil {
		s.Log.Warn("Error validating book", zap.Error(err))
		return err
	}
	err := s.repo.CreateBook(book)
	if err != nil {
		s.Log.Error("Failed to create book", zap.Error(err))
		return fmt.Errorf("could not save book: %w", err)
	}
	return nil
}

func (s *bookService) GetBookByID(id int) (*models.Book, error) {
	if id <= 0 {
		s.Log.Warn("your book id is empty")
		return nil, wrong.ErrBookIDZero
	}

	return s.repo.GetBookByID(id)
}

func (s *bookService) GetAllBook() ([]*models.Book, error) {
	return s.repo.GetAllBooks()
}

func (s *bookService) UpdateBook(book *models.Book) error {
	if book.ID == 0 {
		s.Log.Error("Error when creating book", zap.String("id", strconv.Itoa(book.ID)), zap.Error(wrong.ErrBookIDZero))
		return wrong.ErrBookIDZero
	}
	return s.repo.Update(book)
}

func (s *bookService) DeleteBook(id int) error {
	if id <= 0 {
		s.Log.Warn("your book id is empty")
		return wrong.ErrBookIDZero
	}
	return s.repo.DeleteBook(id)
}

func (s *bookService) validateBookFields(book *models.Book) error {
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
