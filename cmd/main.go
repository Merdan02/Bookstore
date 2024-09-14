package main

import (
	"Bookstore/internal/handler"
	"Bookstore/internal/repository"
	"Bookstore/internal/service"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=merdan password=0800 dbname=bookstore sslmode=disable")
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUserHandler)
		userRoutes.GET("/:id", userHandler.GetUserHandler)
		userRoutes.PUT("/", userHandler.UpdateUserHandler)
		userRoutes.DELETE("/:id", userHandler.DeleteUserHandler)
	}

	err = r.Run()
	if err != nil {
		return
	}
}
