package wrong

import (
	"errors"
	"os"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrEmptyUsername  = errors.New("username cannot be empty")
	ErrEmptyPassword  = errors.New("password cannot be empty")
	ErrEmptyRole      = errors.New("role cannot be empty")
	ErrUserIDZero     = errors.New("user ID cannot be zero")
	ErrBookNotFound   = errors.New("book not found")
	ErrEmptyBook      = errors.New("book cannot be empty")
	JwtKey            = os.Getenv("JWT_SECRET")
	ErrInvalidRequest = "Invalid request"
	ErrInvalidToken   = "Invalid or expired token"
	ErrInternalServer = "Internal server error"
	SuccessMessage    = "User registered successfully"
	ErrEmptyTitle     = errors.New("title cannot be empty")
	ErrEmptyAuthor    = errors.New("author cannot be empty")
	ErrEmptyPrice     = errors.New("price cannot be empty")
	ErrBookIDZero     = errors.New("book ID cannot be zero")
	ErrEmptyQuantity  = errors.New("quantity cannot be empty")
)
