package service

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
)

type BOokService interface {
	CreateBook(book *models.Book) error
	GetBookByID(id int) (*models.Book, error)
	GetBookByName(name string) (*models.Book, error)
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}

type bookService struct {
	repo repository.BookRepository
}

func (s *bookService) GetBookByID(id int) (*models.Book, error) {
	
	panic("implement me")
}

func (s *bookService) GetBookByName(name string) (*models.Book, error) {
	//TODO implement me
	panic("implement me")
}

func (s *bookService) UpdateBook(book *models.Book) error {
	//TODO implement me
	panic("implement me")
}

func (s *bookService) DeleteBook(id int) error {
	//TODO implement me
	panic("implement me")
}

func NewBookService(repo repository.BookRepository) BOokService {
	return &bookService{
		repo: repo,
	}
}

func (s *bookService) CreateBook(book *models.Book) error {
	return s.repo.CreateBook(book)

}
