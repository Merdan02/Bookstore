package routes

import (
	"Bookstore/internal/handler"
	"database/sql"
	"github.com/gin-gonic/gin"
)

// SetupRoutes — это функция, которая регистрирует маршруты в Gin
func SetupRoutes(r *gin.Engine, db *sql.DB, userHandler *handler.UserHandler, bookHandler *handler.BookHandler) {
	// Маршруты для пользователей
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUserHandler)
		userRoutes.GET("/", userHandler.GetAllUserHandler)
		userRoutes.GET("/:id", userHandler.GetUserHandlerByID)
		userRoutes.PUT("/:id", userHandler.UpdateUserHandler)
		userRoutes.DELETE("/:id", userHandler.DeleteUserHandler)
	}

	// Маршруты для книг
	bookRoutes := r.Group("/books")
	{
		bookRoutes.POST("/", bookHandler.CreateBookHandler)
		bookRoutes.GET("/", bookHandler.GetAllBook)
		bookRoutes.GET("/:id", bookHandler.GetBookByID)
		bookRoutes.PUT("/:id", bookHandler.UpdateBookHandler)
		bookRoutes.DELETE("/:id", bookHandler.DeleteBookHandler)
	}
}
