package repository

import (
	"Bookstore/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый репозиторий для пользователей
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *models.User) error {
	// Логируем, что делаем запрос в базу данных
	log.Printf("Попытка добавить пользователя в БД: %+v", user)

	query := "INSERT INTO users (id, username, password, role) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, user.ID, user.Username, user.Password, user.Role)
	if err != nil {
		// Логируем ошибку запроса
		log.Printf("Ошибка при добавлении пользователя в БД: %v", err)
		return err
	}

	// Лог успеха
	log.Println("Пользователь успешно добавлен в БД")
	return nil
}

// GetUserByID возвращает пользователя по его ID
func (r *userRepository) GetUserByID(id int64) (*models.User, error) {
	query := "SELECT id, username, password, role FROM users WHERE id = $1"
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %v", err)
	}
	return user, nil
}

// GetUserByName возвращает пользователя по имени
func (r *userRepository) GetUserByName(name string) (*models.User, error) {
	query := "SELECT id, name, password, role FROM users WHERE name = $1"
	user := &models.User{}
	err := r.db.QueryRow(query, name).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить пользователя: %v", err)
	}
	return user, nil
}

// UpdateUser обновляет информацию о пользователе
func (r *userRepository) UpdateUser(user *models.User) error {
	query := "UPDATE users SET name = $1, password = $2, role = $3 WHERE id = $4"
	_, err := r.db.Exec(query, user.Username, user.Password, user.Role, user.ID)
	if err != nil {
		return fmt.Errorf("не удалось обновить пользователя: %v", err)
	}
	return nil
}

// DeleteUser удаляет пользователя по ID
func (r *userRepository) DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить пользователя: %v", err)
	}
	return nil
}
