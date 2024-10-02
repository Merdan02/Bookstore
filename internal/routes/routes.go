package routes

import (
	"Bookstore/internal/handler"
	"Bookstore/internal/middleware"
	"github.com/gin-gonic/gin"
	"log" // Для логирования
)

// SetupRoutes — это функция, которая регистрирует маршруты в Gin
func SetupRoutes(authHandler *handler.AuthHandler, bookHandler *handler.BookHandler) *gin.Engine {
	log.Println("Initializing routes...") // Логирование начала настройки маршрутов

	router := gin.Default()

	// Маршруты аутентификации
	router.POST("auth/register", authHandler.Register)
	router.POST("auth/login", authHandler.Login)
	router.Use(middleware.AuthRequired())

	// Маршруты для книг
	books := router.Group("/books")
	books.GET("/", bookHandler.GetAllBook)
	books.GET("/:id", bookHandler.GetBookByID)

	// Админские маршруты для управления книгами
	admin := router.Group("/admin")
	admin.Use(middleware.AdminOnly())

	log.Println("Setting up /admin routes...") // Логирование настройки админских маршрутов
	admin.POST("/books", bookHandler.CreateBookHandler)
	admin.PUT("/books/:id", bookHandler.UpdateBookHandler)
	admin.DELETE("/books/:id", bookHandler.DeleteBookHandler)

	log.Println("Routes setup completed.") // Логирование завершения настройки маршрутов
	return router
}
