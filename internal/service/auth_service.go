package service

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
	"errors"
	"log" // Добавляем для логирования

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// RegisterUser Регистрация пользователя
func (s *AuthService) RegisterUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err) // Логируем ошибку
		return err
	}
	user.Password = string(hashedPassword)
	if err := s.UserRepo.CreateUser(user); err != nil {
		log.Println("Error creating user:", err) // Логируем ошибку
		return err
	}
	return nil
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		log.Println("Error fetching user by username:", err) // Логируем ошибку
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Invalid credentials:", err) // Логируем ошибку при неверных данных
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
