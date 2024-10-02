package repository

import (
	"Bookstore/internal/models"
	"database/sql"
	"errors"
	"log" // Для логирования
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser Создание нового пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
	_, err := r.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Role)
	if err != nil {
		log.Println("Error creating user:", err) // Логируем ошибку создания пользователя
	}
	return err
}

// GetUserByUsername Поиск пользователя по имени
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("User not found:", username) // Логируем отсутствие пользователя
			return nil, errors.New("user not found")
		}
		log.Println("Error fetching user by username:", err) // Логируем ошибку получения данных пользователя
		return nil, err
	}
	// Логируем успешный запрос
	log.Println("User found:", user.Username)
	return user, nil
}
