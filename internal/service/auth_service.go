package service

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
	"Bookstore/internal/wrong"
	"errors"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthServ interface {
	RegisterUser(user *models.User) error
	Login(username, password string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUserByName(username string) (*models.User, error)
	GetByUserID(id int64) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
}

type AuthService struct {
	UserRepo *repository.UserRepository
	Log      *zap.Logger
}

var ErrUserNotFound = errors.New("user not found")

func NewUserService(userRepo *repository.UserRepository, logger *zap.Logger) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
		Log:      logger,
	}
}

func (s *AuthService) RegisterUser(user *models.User) error {
	if err := s.validateUserFields(user); err != nil {
		s.Log.Warn("Ошибка валидации при создании пользователя", zap.Error(err))
		return err
	}

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
	if username == "" {
		s.Log.Warn("Попытка получить пользователя с пустым именем")
		return nil, wrong.ErrEmptyUsername
	}

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

func (s *AuthService) GetByUserID(id int) (*models.User, error) {
	if id <= 0 {
		s.Log.Warn("Trying to get with empty user ID")
		return nil, wrong.ErrUserIDZero
	}

	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			zap.L().Warn("User not found", zap.Int("id", id))
			return nil, ErrUserNotFound
		}
		zap.L().Error("Error fetching user by ID", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *AuthService) UpdateUser(user *models.User) error {
	if err := s.validateUpdateFields(user); err != nil {
		s.Log.Warn("Ошибка валидации при обновлении пользователя", zap.Error(err))
		return err
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

	if user.ID <= 0 {
		zap.L().Warn("Invalid user ID for deletion")
		return errors.New("invalid user ID")
	}

	err := s.UserRepo.DeleteUser(user.ID)
	if err != nil {
		zap.L().Error("Error deleting user", zap.String("user ID", strconv.Itoa(user.ID)), zap.Error(err))
		return err
	}

	zap.L().Info("User deleted successfully", zap.String("username", user.Username))
	return nil
}

// Вспомогательная функция для валидации полей пользователя при создании
func (s *AuthService) validateUserFields(user *models.User) error {
	if user.Username == "" {
		return wrong.ErrEmptyUsername
	}
	if user.Password == "" {
		return wrong.ErrEmptyPassword
	}
	if user.Role == "" {
		return wrong.ErrEmptyRole
	}
	return nil
}

// Валидация полей при обновлении пользователя
func (s *AuthService) validateUpdateFields(user *models.User) error {
	if user.ID == 0 {
		return wrong.ErrUserIDZero
	}
	if user.Username == "" {
		return wrong.ErrEmptyUsername
	}
	if user.Role == "" {
		return wrong.ErrEmptyRole
	}
	return nil
}
