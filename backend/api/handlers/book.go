package handlers

import (
	"bookhub/backend/database/models"
	"bookhub/backend/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// GetUserBooks fetches books owned by the user
func GetUserBooks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDStr := c.Query("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	// Fetch books
	var books []models.Book
	if err := db.Where("user_id = ?", userID).Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}

	// Ensure response is valid JSON
	if len(books) == 0 {
		c.JSON(http.StatusOK, []models.Book{}) // Return empty array instead of null
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetUserBooks fetches books owned by the user
func GetBooksByUserID(db *gorm.DB, userID int) ([]models.Book, error) {
	var books []models.Book
	if err := db.Where("user_id = ?", userID).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// GetBooks retrieves all books from the database
func GetBooks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var books []models.Book

	if err := db.Preload("Reviews").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	for i := range books {
		books[i].Price = books[i].CalculateFinalPrice() // Apply discount
	}

	c.JSON(http.StatusOK, books)
}

// AddBook adds a new book to the database
func AddBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetBookByID retrieves a book by ID from the database
func GetBookByID(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var book models.Book
	if err := db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook handles the creation of a new book
func CreateBook(c *gin.Context) {
	// Extract user ID from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert userID to uint
	userIDInt, ok := userID.(uint)
	if !ok {
		userIDFloat, ok := userID.(float64) // Check if it's float64 (JWT stores numbers as float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}
		userIDInt = uint(userIDFloat) // Convert to uint
	}

	// Get the request body
	var bookRequest struct {
		Title       string  `json:"title"`
		Author      string  `json:"author"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

	// Bind the request body to the bookRequest struct
	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get a reference to the database (passed via middleware)
	db := c.MustGet("db").(*gorm.DB)

	// Call the CreateBook function with the proper arguments
	book, err := CreateBookInDB(db, userIDInt, bookRequest.Title, bookRequest.Author, bookRequest.Description, bookRequest.Price, bookRequest.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.LogEvent(db, userIDInt, "BOOK_CREATED", "A new book was added: "+book.Title)

	// Return the newly created book
	c.JSON(http.StatusOK, book)
}

// CreateBookInDB creates a book record in the database and links it to the user
func CreateBookInDB(db *gorm.DB, userID uint, title, author, description string, price float64, stock int) (*models.Book, error) {
	book := models.Book{
		UserID:      userID, // 🔹 Store the user ID
		Title:       title,
		Author:      author,
		Description: description,
		Price:       decimal.NewFromFloat(price), //Convert float64 to decimal
		Stock:       stock,
	}

	// Attempt to create the book record in the database
	if err := db.Create(&book).Error; err != nil {
		return nil, err
	}

	// Return the newly created book object
	return &book, nil
}

// UpdateBook updates a book in the database by ID
func UpdateBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var book models.Book
	if err := db.First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book"})
		}
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book from the database by ID
func DeleteBook(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := db.Delete(&models.Book{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// GetOrders retrieves all orders from the database
func GetOrders(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var orders []models.Order

	if err := db.Preload("Books").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder retrieves an order by ID from the database
func GetTrendingBooks(c *gin.Context) {
	var books []models.Book
	db := c.MustGet("db").(*gorm.DB)

	// Ensure 'sales' column exists in the books table
	result := db.Order("sales DESC").Limit(10).Find(&books)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Return books in JSON format
	c.JSON(http.StatusOK, books)
}
