package service

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthServ interface {
	RegisterUser(user *models.User) error
	Login(username, password string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUserByName(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
}

type AuthService struct {
	UserRepo *repository.UserRepository
}

var ErrUserNotFound = errors.New("user not found")

func NewUserService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) RegisterUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Error hashing password", zap.Error(err))
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.UserRepo.CreateUser(user); err != nil {
		zap.L().Error("Error creating user", zap.String("username", user.Username), zap.Error(err))
		return err
	}

	zap.L().Info("User successfully registered", zap.String("username", user.Username))
	return nil
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		zap.L().Error("Error fetching user by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		zap.L().Warn("Invalid credentials", zap.String("username", username))
		return nil, errors.New("invalid credentials")
	}

	zap.L().Info("User logged in successfully", zap.String("username", username))
	return user, nil
}

func (s *AuthService) GetAllUsers() ([]*models.User, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		zap.L().Error("Error fetching all users", zap.Error(err))
		return nil, err
	}
	zap.L().Info("Fetched all users successfully")
	return users, nil
}

func (s *AuthService) GetUserByName(username string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			zap.L().Warn("User not found", zap.String("username", username))
			return nil, ErrUserNotFound
		}
		zap.L().Error("Error fetching user by username", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	zap.L().Info("Fetched user successfully", zap.String("username", username))
	return user, nil
}

func (s *AuthService) UpdateUser(user *models.User) error {
	if user.ID == 0 {
		zap.L().Warn("Invalid user ID for update", zap.Int("userID", user.ID))
		return errors.New("invalid user ID")
	}

	err := s.UserRepo.UpdateUser(user)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			zap.L().Warn("User not found for update", zap.String("username", user.Username))
			return ErrUserNotFound
		}
		zap.L().Error("Error updating user", zap.String("username", user.Username), zap.Error(err))
		return err
	}

	zap.L().Info("User updated successfully", zap.String("username", user.Username))
	return nil
}

func (s *AuthService) DeleteUser(user *models.User) error {
	if user.Username == "" {
		zap.L().Warn("Invalid username for deletion")
		return errors.New("invalid username")
	}

	err := s.UserRepo.DeleteUser(user)
	if err != nil {
		zap.L().Error("Error deleting user", zap.String("username", user.Username), zap.Error(err))
		return err
	}

	zap.L().Info("User deleted successfully", zap.String("username", user.Username))
	return nil
}
