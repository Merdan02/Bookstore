package app

import (
	"Bookstore/internal/handler"
	"Bookstore/internal/repository"
	"Bookstore/internal/routes"
	"Bookstore/internal/service"
	config "Bookstore/pkg/database"
	"database/sql"
	"log" // Для логирования
)

// InitApp инициализирует все зависимости (репозитории, сервисы, обработчики)
func InitApp(db *sql.DB) (*handler.AuthHandler, *handler.BookHandler) {
	// Логирование процесса инициализации зависимостей
	log.Println("Initializing repositories, services, and handlers...")

	// Инициализация зависимостей для пользователей
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewAuthHandler(userService)
	log.Println("User dependencies initialized successfully.")

	// Инициализация зависимостей для книг
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)
	log.Println("Book dependencies initialized successfully.")

	return userHandler, bookHandler
}

func Run() {
	// Логирование подключения к базе данных
	log.Println("Connecting to the database...")

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err) // Логируем фатальную ошибку при подключении к БД
		//panic(err)
	}
	log.Println("Database connection established successfully.")

	// Инициализация зависимостей
	authHandler, bookHandler := InitApp(db)

	// Настройка маршрутов
	log.Println("Setting up routes...")
	r := routes.SetupRoutes(authHandler, bookHandler)

	log.Println("Starting the server on port :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err) // Логируем фатальную ошибку при запуске сервера
		//panic(err)
	}
}
