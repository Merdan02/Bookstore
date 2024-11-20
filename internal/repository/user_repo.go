package repository

import (
	"Bookstore/internal/models"
	"Bookstore/internal/wrong"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserRepo Интерфейс для операций с пользователями
type UserRepo interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(username string) error
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
	if err := r.validateUserFields(user); err != nil {
		r.Log.Warn("Ошибка валидации при создании пользователя", zap.Error(err))
		return err
	}

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

	//r.Log.Info("Пользователь успешно создан", zap.String("username", user.Username))
	return nil
}

// GetUserByUsername Получение пользователя по имени
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	if username == "" {
		r.Log.Warn("Попытка получить пользователя с пустым именем")
		return nil, wrong.ErrEmptyUsername
	}

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

	r.Log.Info("Пользователь найден", zap.String("username", user.Username))
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

	r.Log.Info("Пользователи успешно получены", zap.Int("count", len(users)))
	return users, nil
}

// UpdateUser Обновление данных пользователя
func (r *UserRepository) UpdateUser(user *models.User) error {
	if err := r.validateUpdateFields(user); err != nil {
		r.Log.Warn("Ошибка валидации при обновлении пользователя", zap.Error(err))
		return err
	}

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

	r.Log.Info("Пользователь успешно обновлен", zap.String("username", user.Username))
	return nil
}

// DeleteUser Удаление пользователя
func (r *UserRepository) DeleteUser(user *models.User) error {
	if user.Username == "" {
		r.Log.Warn("Попытка удалить пользователя с пустым именем")
		return wrong.ErrEmptyUsername
	}

	_, err := r.DB.Exec("DELETE FROM users WHERE username = $1", user.Username)
	if err != nil {
		r.Log.Error("Ошибка базы данных при удалении пользователя", zap.String("username", user.Username), zap.Error(err))
		return errors.New("не удалось удалить пользователя")
	}

	r.Log.Info("Пользователь успешно удален", zap.String("username", user.Username))
	return nil
}

// Вспомогательная функция для валидации полей пользователя при создании
func (r *UserRepository) validateUserFields(user *models.User) error {
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
func (r *UserRepository) validateUpdateFields(user *models.User) error {
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
