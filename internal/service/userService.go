package service

//
//import (
//	"Bookstore/internal/models"
//	"Bookstore/internal/repository"
//	"database/sql"
//	"errors"
//	"fmt"
//	"log"
//	"time"
//
//	"github.com/dgrijalva/jwt-go"
//	"golang.org/x/crypto/bcrypt"
//)
//
//type UserService interface {
//	RegisterUser(user *models.User) error
//	Login(username, password string) (string, error)
//}
//
//type userService struct {
//	repo repository.UserRepository
//}
//
//// NewUserService создает новый сервис для пользователей
//func NewUserService(repo repository.UserRepository) UserService {
//	return &userService{
//		repo: repo,
//	}
//}
//
//var jwtKey = []byte("your_secret_key") // секретный ключ для JWT
//
//type Claims struct {
//	Username string `json:"username"`
//	RoleID   int    `json:"role_id"`
//	jwt.StandardClaims
//}
//
//func (s *userService) RegisterUser(user *models.User) error {
//	// Проверяем, что все обязательные поля заполнены
//	if user.ID == 0 || user.Username == "" || user.Password == "" || user.Role == "" {
//		log.Printf("Не все обязательные поля заполнены: %+v", user)
//		return fmt.Errorf("все поля обязательны для заполнения")
//	}
//
//	_, err := s.repo.GetUserByName(user.Username)
//	if err == nil {
//		return errors.New("user already exists")
//	} //else if err != nil || err != sql.ErrNoRows {
//	//	// Если ошибка не связана с отсутствием записи, возвращаем ошибку
//	//	return fmt.Errorf("database error: %v", err)
//	//}
//
//	// Хешируем пароль
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
//	if err != nil {
//		return fmt.Errorf("error hashing password: %v", err)
//	}
//
//	// Обновляем пользователя с захэшированным паролем
//	user.Password = string(hashedPassword)
//
//	return s.repo.CreateUser(user)
//}
//
//func (s *userService) Login(username, password string) (string, error) {
//	user, err := s.repo.GetUserByName(username)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return "", errors.New("user not found")
//		}
//		return "", fmt.Errorf("database error: %v", err)
//	}
//
//	// Сравниваем пароли
//	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//	if err != nil {
//		return "", errors.New("incorrect password")
//	}
//
//	// Генерируем JWT токен
//	expirationTime := time.Now().Add(24 * time.Hour)
//	claims := &Claims{
//		Username: user.Username,
//		// RoleID:   user.Role,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: expirationTime.Unix(),
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(jwtKey)
//	if err != nil {
//		return "", fmt.Errorf("error signing token: %v", err)
//	}
//
//	return tokenString, nil
//}
//
////// UpdateUser обновляет пользователя
////func (s *userService) UpdateUser(user *models.User) error {
////	if user.ID == 0 {
////		return fmt.Errorf("ID пользователя обязателен для обновления")
////	}
////	return s.repo.UpdateUser(user)
////}
////
////// DeleteUser удаляет пользователя по ID
////func (s *userService) DeleteUser(id int64) error {
////	return s.repo.DeleteUser(id)
////}
//
////func (s *userService) GetAllUsers(models.User) ([]*models.User, error) {
////	return s.repo.GetAllUsers()
////}
//
////// GetUserByID возвращает пользователя по его ID
////func (s *userService) GetUserByID(id int64) (*models.User, error) {
////	return s.repo.GetUserByID(id)
////}
//
////func (s *userService) CreateUser(user *models.User) error {
////	// Логируем полученные данные
////	log.Printf("Полученные данные для создания пользователя: %+v", user)
////
////	if user.ID == 0 || user.Username == "" || user.Password == "" || string(user.RoleID) == "" {
////		log.Printf("Не все обязательные поля заполнены: %+v", user)
////		return fmt.Errorf("все поля обязательны для заполнения")
////	}
////
////	log.Printf("Попытка создать пользователя: %+v", user)
////	return s.repo.CreateUser(user)
////}
