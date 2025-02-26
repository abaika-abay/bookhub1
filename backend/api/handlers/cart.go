package handlers

import (
	"bookhub/backend/database/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCartItems – Fetch all items in the user's cart
func GetCartItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDFloat, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDFloat.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var cartItems []models.CartItem
	db.Preload("Book").Where("user_id = ?", uint(userID)).Find(&cartItems)

	c.JSON(http.StatusOK, cartItems)
}

// GetUserCart - Fetch all items in a user's cart
func GetUserCart(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var cartItems []models.CartItem
	if err := db.Preload("Book").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	c.JSON(http.StatusOK, cartItems)
}

// AddBookToCart – Adds a book to the cart
func AddBookToCart(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDFloat, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDFloat.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cartItem.UserID = uint(userID)

	// Check if item already exists in cart
	var existingItem models.CartItem
	if err := db.Where("user_id = ? AND book_id = ?", cartItem.UserID, cartItem.BookID).First(&existingItem).Error; err == nil {
		// Update quantity if book already in cart
		existingItem.Quantity += cartItem.Quantity
		db.Save(&existingItem)
	} else {
		// Add new item
		db.Create(&cartItem)
	}

	var book models.Book

	LogEvent(db, uint(userID), "Cart Update", "Added book "+book.Title+" to cart")

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

func UpdateCartItemQuantity(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Extract book ID from URL
	bookIDStr := c.Param("book_id")
	fmt.Println("Received book ID:", bookIDStr) // Debugging log

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 { // Ensure book ID is valid
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book_id", "received": bookIDStr})
		return
	}

	// Extract user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// Find cart item
	var cartItem models.CartItem
	if err := db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	// Update quantity
	cartItem.Quantity = input.Quantity
	if err := db.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart updated", "cart_item": cartItem})
}

// RemoveBookFromCart - Deletes a book from cart
func RemoveBookFromCart(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if err := db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove book from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book removed from cart"})
}

// UpdateCartItem updates the quantity of a cart item
func UpdateCartItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	cartItemID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Quantity int `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var cartItem models.CartItem
	if err := db.First(&cartItem, cartItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	cartItem.Quantity = input.Quantity
	db.Save(&cartItem)

	c.JSON(http.StatusOK, cartItem)
}

// RemoveFromCart – Deletes an item from the cart
func RemoveFromCart(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userIDFloat, _ := c.Get("user_id")
	userID := uint(userIDFloat.(float64))

	bookID, _ := strconv.Atoi(c.Param("book_id"))

	db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.CartItem{})
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}
