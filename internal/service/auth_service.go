package services

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(username, role string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}
