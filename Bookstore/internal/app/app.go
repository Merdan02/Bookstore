package app

import (
	"Bookstore/internal/handler"
	"Bookstore/internal/repository"
	"Bookstore/internal/routes"
	"Bookstore/internal/service"
	config "Bookstore/pkg/database"
	"database/sql"
	"github.com/gin-gonic/gin"
)

// InitApp инициализирует все зависимости (репозитории, сервисы, обработчики)
func InitApp(db *sql.DB) (*handler.UserHandler, *handler.BookHandler) {
	// Инициализация зависимостей для пользователей
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Инициализация зависимостей для книг
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	return userHandler, bookHandler
}

func Run() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Инициализация зависимостей
	userHandler, bookHandler := InitApp(db)

	// Создание нового Gin роутера
	r := gin.Default()

	// Настройка маршрутов
	routes.SetupRoutes(r, db, userHandler, bookHandler)

	// Запуск сервера
	if err := r.Run(); err != nil {
		panic(err)
	}
}
