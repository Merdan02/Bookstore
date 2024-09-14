package service

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
	"fmt"
	"log"
)

type USerService interface {
	CreateUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService создает новый сервис для пользователей
func NewUserService(repo repository.UserRepository) USerService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(user *models.User) error {
	// Логируем полученные данные
	log.Printf("Полученные данные для создания пользователя: %+v", user)

	if user.ID == 0 || user.Username == "" || user.Password == "" || user.Role == "" {
		log.Printf("Не все обязательные поля заполнены: %+v", user)
		return fmt.Errorf("все поля обязательны для заполнения")
	}

	log.Printf("Попытка создать пользователя: %+v", user)
	return s.repo.CreateUser(user)
}

// GetUserByID возвращает пользователя по его ID
func (s *userService) GetUserByID(id int64) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

// GetUserByName возвращает пользователя по имени
func (s *userService) GetUserByName(name string) (*models.User, error) {
	return s.repo.GetUserByName(name)
}

// UpdateUser обновляет пользователя
func (s *userService) UpdateUser(user *models.User) error {
	if user.ID == 0 {
		return fmt.Errorf("ID пользователя обязателен для обновления")
	}
	return s.repo.UpdateUser(user)
}

// DeleteUser удаляет пользователя по ID
func (s *userService) DeleteUser(id int64) error {
	return s.repo.DeleteUser(id)
}
