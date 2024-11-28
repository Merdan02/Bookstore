package repository

import (
	"Bookstore/internal/models"
	"Bookstore/internal/wrong"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

// UserRepo Интерфейс для операций с пользователями
type UserRepo interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(ID int) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

// UserRepository Структура репозитория пользователей
type UserRepository struct {
	DB  *sql.DB
	Log *zap.Logger
}

// NewUserRepository Конструктор UserRepository
func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{DB: db, Log: logger}
}

// CreateUser Создание нового пользователя
func (r *UserRepository) CreateUser(user *models.User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		r.Log.Error("Ошибка хеширования пароля", zap.Error(err))
		return errors.New("внутренняя ошибка: не удалось создать пользователя")
	}

	_, err = r.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", user.Username, hashedPassword, user.Role)
	if err != nil {
		r.Log.Error("Ошибка базы данных при создании пользователя", zap.String("username", user.Username), zap.Error(err))
		return errors.New("не удалось создать пользователя")
	}
	return nil
}

// GetUserByUsername Получение пользователя по имени
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {

	user := &models.User{}
	err := r.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Log.Warn("Пользователь не найден", zap.String("username", username))
			return nil, wrong.ErrUserNotFound
		}
		r.Log.Error("Ошибка базы данных при получении пользователя", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(ID int) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRow("SELECT id, username, password, role FROM users WHERE id = $1", ID).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Log.Warn("User not found", zap.String("user ID", strconv.Itoa(ID)))
			return nil, wrong.ErrUserNotFound
		}
		r.Log.Error("Error database while getting user", zap.String("username", strconv.Itoa(ID)), zap.Error(err))
		return nil, err
	}
	return user, nil
}

// GetAllUsers Получение всех пользователей
func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := "SELECT id, username, password, role FROM users"
	rows, err := r.DB.Query(query)
	if err != nil {
		r.Log.Error("Ошибка базы данных при получении всех пользователей", zap.Error(err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Role); err != nil {
			r.Log.Error("Ошибка при сканировании строки", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		r.Log.Error("Ошибка при итерации по строкам", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.Log.Error("Ошибка при закрытии строк результата", zap.Error(err))
		}
	}()

	return users, nil
}

// UpdateUser Обновление данных пользователя
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := "UPDATE users SET username = $1, role = $2 WHERE id = $3"
	args := []interface{}{user.Username, user.Role, user.ID}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			r.Log.Error("Ошибка хеширования пароля", zap.Error(err))
			return errors.New("внутренняя ошибка: не удалось обновить пользователя")
		}
		query = "UPDATE users SET username = $1, role = $2, password = $3 WHERE id = $4"
		args = append(args[:2], hashedPassword, user.ID)
	}

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		r.Log.Error("Ошибка базы данных при обновлении пользователя", zap.String("username", user.Username), zap.Error(err))
		return err
	}
	return nil
}

// DeleteUser Удаление пользователя
func (r *UserRepository) DeleteUser(id int) error {

	_, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		r.Log.Error("Ошибка базы данных при удалении пользователя", zap.String("id", strconv.Itoa(id)), zap.Error(err))
		return errors.New("не удалось удалить пользователя")
	}

	return nil
}
