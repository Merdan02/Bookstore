package controllers

import (
	"Bookstore/internal/models"
	"Bookstore/internal/repository"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

var db *sql.DB

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := repository.RegisterUser(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	users, err := repository.RegisterUser(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
