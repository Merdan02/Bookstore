package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Используем безопасный генератор токенов
var jwtKey = []byte("your_secret_key")

func GenerateAccessToken(username, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		zap.L().Error("Error generating access token", zap.Error(err))
		return "", err
	}
	return tokenString, nil
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			zap.L().Warn("Unauthorized access attempt", zap.Any("role", role))
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
			zap.L().Warn("Missing Authorization header")
			c.JSON(http.StatusForbidden, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		claims := &jwt.StandardClaims{}
		tkn, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !tkn.Valid {
			zap.L().Error("Invalid token", zap.Error(err))
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Вытаскиваем данные пользователя
		c.Set("username", claims.Subject)
		c.Next()
	}
}
