package handler

import (
	"Bookstore/internal/models"
	"Bookstore/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type BookHandler struct {
	service service.BOokService
}

func NewBookHandler(s service.BOokService) *BookHandler {
	return &BookHandler{service: s}

}

func (h *BookHandler) CreateBookHandler(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Printf("Error connection to Bindjson: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong data"})
		return
	}

	if err := h.service.CreateBook(&book); err != nil {
		log.Printf("Error create Book in service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})

}

func (h *BookHandler) GetAllBook(c *gin.Context) {
	books, err := h.service.GetAllBook()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (h *BookHandler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong ID"})
		return
	}
	book, err := h.service.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if book == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}
func (h *BookHandler) UpdateBookHandler(c *gin.Context) {
	var book models.Book

	// Получаем ID из параметра пути и конвертируем его в int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid book ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	book.ID = id // Устанавливаем ID для обновления

	if err := c.ShouldBindJSON(&book); err != nil {
		log.Printf("Error in ShouldBindJSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong data"})
		return
	}

	log.Printf("Received data for update: %+v", book)
	if err := h.service.UpdateBook(&book); err != nil {
		log.Printf("Error in UpdateBook service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *BookHandler) DeleteBookHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong ID"})
		return
	}
	if err := h.service.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": true})
}
