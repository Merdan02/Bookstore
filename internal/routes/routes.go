package routes

import (
	"Bookstore/internal/handler"
	"Bookstore/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap" // Используем zap для логирования
)

// SetupRoutes — это функция, которая регистрирует маршруты в Gin
func SetupRoutes(authHandler *handler.AuthHandler, bookHandler *handler.BookHandler) *gin.Engine {
	// Инициализация логгера zap
	logger, _ := zap.NewProduction()
	defer logger.Sync() // Закрываем логгер после завершения работы

	router := gin.Default()

	// Маршруты аутентификации
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Middleware аутентификации для защищенных маршрутов
	router.Use(middleware.AuthRequired())

	// Маршруты для работы с пользователями
	usersGroup := router.Group("/users")
	{
		usersGroup.GET("/", authHandler.GetAllUser)
		router.GET("/users/id/:id", authHandler.GetUserByID)
		router.GET("/users/username/:username", authHandler.GetUserByUsername)
		usersGroup.PUT("/:id", authHandler.UpdateUser)
		usersGroup.DELETE("/:id", authHandler.DeleteUser)
	}

	// Маршруты для работы с книгами
	booksGroup := router.Group("/books")
	{
		booksGroup.GET("/", bookHandler.GetAllBook)
		booksGroup.GET("/:id", bookHandler.GetBookByID)
	}

	// Админские маршруты для управления книгами, только для администраторов
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthRequired()) // Добавляем отдельное middleware для проверки роли администратора
	{
		logger.Info("Setting up /admin routes...")
		adminGroup.POST("/books", bookHandler.CreateBookHandler)
		adminGroup.PUT("/books/:id", bookHandler.UpdateBookHandler)
		adminGroup.DELETE("/books/:id", bookHandler.DeleteBookHandler)
	}
	return router
}
