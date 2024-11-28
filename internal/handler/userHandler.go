package handler

import (
	"Bookstore/internal/middleware"
	"Bookstore/internal/models"
	"Bookstore/internal/service"
	"Bookstore/internal/wrong"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log" // Для логирования
	"net/http"
	"strconv"
)

// Хелпер для ответов с ошибкой
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

// Хелпер для успешных ответов
func respondWithSuccess(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}

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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON in RefreshToken:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(input.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return wrong.JwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token or expired token"})
		return
	}

	username := claims["username"].(string)
	role := claims["role"].(string)
	newAccessToken, err := middleware.GenerateAccessToken(username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": newAccessToken})
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
		log.Println("Error during login for user:", input.Username, err) // Логируем ошибку при входе
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateAccessToken(user.Username, user.Role)
	if err != nil {
		log.Println("Error generating JWT token for user:", user.Username, err) // Логируем ошибку генерации токена
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Println("User logged in successfully:", user.Username) // Логируем успешный вход
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) GetAllUser(c *gin.Context) {
	users, err := h.AuthService.GetAllUsers()
	if err != nil {
		log.Println("Error fetching all users:", err) // Логируем ошибку при получении всех пользователей
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Successfully fetched all users") // Логируем успешное получение всех пользователей
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AuthHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	// Запрос на получение пользователя по имени
	user, err := h.AuthService.GetUserByName(username)
	if err != nil {
		if errors.Is(err, wrong.ErrUserNotFound) {
			log.Printf("User not found with username: %s", username)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by username: %s, error: %v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	log.Printf("Successfully fetched user by username: %s", user.Username)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Invalid user ID:", c.Param("id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := h.AuthService.GetByUserID(id)
	if err != nil {
		if errors.Is(err, wrong.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Printf("Error fetching user by id: %d, error: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	log.Printf("Successfully fetched user by id: %d", id)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	var user models.User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid user id: %v", err)
		c.JSON(400, gin.H{"error": "Invalid user id"})
		return
	}
	user.ID = id

	// Привязка JSON данных к модели пользователя
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON in UpdateUser:", err) // Логируем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Вызов метода обновления пользователя в AuthService
	if err := h.AuthService.UpdateUser(&user); err != nil {
		if err.Error() == "user not found" {
			log.Printf("User not found for update, ID: %d", user.ID)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		log.Println("Error updating user:", user.ID, err) // Логируем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	log.Println("User updated successfully:", user.ID) // Логируем успешное обновление
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println("Invalid user ID:", c.Param("id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.AuthService.GetByUserID(id)
	if err != nil {
		log.Println("Error fetching user for deletion:", id, err) // Логируем ошибку при поиске пользователя для удаления
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := h.AuthService.DeleteUser(user); err != nil {
		log.Println("Error deleting user:", user.Username, err) // Логируем ошибку при удалении пользователя
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("User deleted successfully:", user.Username) // Логируем успешное удаление пользователя
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
