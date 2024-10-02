package handler

import (
	"Bookstore/internal/middleware"
	"Bookstore/internal/models"
	"Bookstore/internal/service"
	"github.com/gin-gonic/gin"
	"log" // Для логирования
	"net/http"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

// NewAuthHandler создаёт новый AuthHandler с указанной AuthService
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON in Register:", err) // Логируем ошибку привязки JSON
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//user.Role = "user" // по умолчанию роль обычного пользователя
	err := h.AuthService.RegisterUser(&user)
	if err != nil {
		log.Println("Error registering user:", err) // Логируем ошибку регистрации пользователя
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("User registered successfully:", user.Username) // Логируем успешную регистрацию
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON in Login:", err) // Логируем ошибку привязки JSON
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.AuthService.Login(input.Username, input.Password)
	if err != nil {
		log.Println("Error during login:", err) // Логируем ошибку при входе
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(user.Username, user.Role)
	if err != nil {
		log.Println("Error generating JWT token:", err) // Логируем ошибку генерации токена
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Println("User logged in successfully:", user.Username) // Логируем успешный вход
	c.JSON(http.StatusOK, gin.H{"token": token})
}
