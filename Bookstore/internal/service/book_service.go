package service

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
	"errors"
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
}

func NewBookService(repo repository.BookRepository) BOokService {
	return &bookService{
		repo: repo,
	}
}

func (s *bookService) CreateBook(book *models.Book) error {
	return s.repo.CreateBook(book)
}

func (s *bookService) GetBookByID(id int) (*models.Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *bookService) GetAllBook() ([]*models.Book, error) {
	return s.repo.GetAllBooks()
}

func (s *bookService) UpdateBook(book *models.Book) error {
	if book.ID == 0 {
		return errors.New("invalid id")
	}
	return s.repo.UpdateBook(book)
}

func (s *bookService) DeleteBook(id int) error {
	return s.repo.DeleteBook(id)
}
