package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("secret_key")

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		c.Abort()
		return
	}

	tokenStr := strings.Split(authHeader, " ")[1]
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("username", claims.Subject)
	c.Next()
}
