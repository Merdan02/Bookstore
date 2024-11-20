package app

import (
	"Bookstore/internal/handler"
	"Bookstore/internal/repository"
	"Bookstore/internal/routes"
	"Bookstore/internal/service"
	config "Bookstore/pkg/database"
	"database/sql"
	"go.uber.org/zap"
	"log" // Для логирования
)

// InitApp инициализирует все зависимости (репозитории, сервисы, обработчики)
func InitApp(db *sql.DB, logger *zap.Logger) (*handler.AuthHandler, *handler.BookHandler) {
	log.Println("Initializing repositories, services, and handlers...")

	// Initialize dependencies for users
	userRepo := repository.NewUserRepository(db, logger)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewAuthHandler(userService)
	log.Println("User dependencies initialized successfully.")

	// Initialize dependencies for books
	bookRepo := repository.NewBookRepository(db, logger)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)
	log.Println("Book dependencies initialized successfully.")

	return userHandler, bookHandler
}

func Run() {
	log.Println("Connecting to the database...")

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Database connection established successfully.")

	// Initialize the logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger) // Flushes buffer, if any

	// Initialize dependencies
	authHandler, bookHandler := InitApp(db, logger)

	log.Println("Setting up routes...")
	r := routes.SetupRoutes(authHandler, bookHandler)

	log.Println("Starting the server on port :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
