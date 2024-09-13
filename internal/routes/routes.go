package routes

import (
	"Bookstore/internal/handler"
	"Bookstore/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes
	router.POST("/register", handler.Register)
	router.POST("/login", bookHandler.Login)

	// Protected routes
	auth := router.Group("/admin")
	auth.Use(middlewares.AuthMiddleware)
	{
		auth.POST("/books", controllers.CreateBook)
		auth.GET("/books", controllers.GetBooks)
	}

	router.GET("/books", controllers.GetBooks)

	return router
}
