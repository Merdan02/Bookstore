package handler

//
//import (
//	"Bookstore/internal/repository"
//	"net/http"
//	"strconv"
//
//	"Bookstore/internal/models"
//	"github.com/gin-gonic/gin"
//)
//
//type BookHandler struct {
//	bookService service.BookService
//}
//
//func NewBookHandler(bookService service.BookService) *BookHandler {
//	return &BookHandler{bookService: bookService}
//}
//
//func (h *BookHandler) CreateBook(c *gin.Context) {
//	var book *models.Book
//	if err := c.BindJSON(&book); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	if err := repository.CreateBook(db, &book); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create book"})
//		return
//	}
//
//	c.JSON(http.StatusOK, book)
//}
//
//func (h *BookHandler) GetBooks(c *gin.Context) {
//	books, err := repository.GetBooks(db)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch books"})
//		return
//	}
//	c.JSON(http.StatusOK, books)
//}
//
//func (h *BookHandler) UpdateBook(c *gin.Context) {
//	id := c.Param("id")
//	bookID, err := strconv.Atoi(id)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
//		return
//	}
//
//	var book models.Book
//	if err := c.BindJSON(&book); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	book.ID = bookID
//	if err := repository.UpdateBook(db, &book); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update book"})
//		return
//	}
//
//	c.JSON(http.StatusOK, book)
//}
//func (h *BookHandler) DeleteBook(c *gin.Context) {
//	id := c.Param("id")
//	bookID, err := strconv.Atoi(id)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
//		return
//	}
//
//	if err := repository.DeleteBookByID(db, bookID); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete book"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
//}
