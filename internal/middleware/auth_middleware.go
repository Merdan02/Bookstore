package middleware

import (
	"github.com/gin-gonic/gin"
	"log" // Для логирования
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

func GenerateJWT(username, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("Error generating JWT:", err) // Логируем ошибку генерации токена
		return "", err
	}

	return tokenString, nil
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, exists := c.Get("role")
		log.Println("User role:", role) // Логируем роль пользователя

		if !exists || role != "admin" {
			log.Println("Unauthorized access attempt by user with role:", role) // Логируем ошибку доступа
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization header") // Логируем отсутствие заголовка авторизации
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this resource"})
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		claims := &jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !tkn.Valid {
			log.Println("Invalid token or parsing error:", err) // Логируем ошибку токена
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this resource"})
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set("username", (*claims)["username"])
		c.Set("role", (*claims)["role"])
		c.Next()
	}
}
