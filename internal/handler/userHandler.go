package handler

import (
	"Bookstore/internal/models"
	"Bookstore/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service service.USerService
}

// NewUserHandler создает новый обработчик для пользователей
func NewUserHandler(s service.USerService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var user models.User
	// Логируем, если ошибка при получении данных
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Ошибка привязки данных (BindJSON): %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Логируем, если ошибка на уровне сервиса
	if err := h.service.CreateUser(&user); err != nil {
		log.Printf("Ошибка при создании пользователя в сервисе: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Лог успеха
	log.Printf("Пользователь успешно создан: %+v", user)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно создан"})
}

// GetUserHandler возвращает пользователя по ID
func (h *UserHandler) GetUserHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserHandler обновляет пользователя
func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	if err := h.service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUserHandler удаляет пользователя
func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь удален"})
}
