package handlers

import (
	"net/http"
	"strconv"
	"time"

	"bookhub/backend/database/models"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal" // Используем decimal для работы с деньгами
	"gorm.io/gorm"
)

// GetUserBalance возвращает баланс пользователя
func GetUserBalance(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": user.Balance})
}

// TopUpBalance пополняет баланс пользователя
func TopUpBalance(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Преобразуем amount в decimal.Decimal перед сложением
	amountDecimal := decimal.NewFromFloat(request.Amount)
	user.Balance = user.Balance.Add(amountDecimal)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance updated successfully", "new_balance": user.Balance})
}

// Checkout обрабатывает заказ
func Checkout(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var cartItems []models.CartItem
	if err := db.Where("user_id = ?", userID).Preload("Book").Find(&cartItems).Error; err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Calculate total price using decimal
	var totalPrice decimal.Decimal
	for _, item := range cartItems {
		price := item.Book.Price.Mul(decimal.NewFromInt(int64(item.Quantity))) // 🛠 Corrected calculation using decimal
		totalPrice = totalPrice.Add(price)
	}

	// Get user balance
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user balance"})
		return
	}

	// Check if the user has enough balance
	if user.Balance.LessThan(totalPrice) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Deduct balance
	user.Balance = user.Balance.Sub(totalPrice)
	db.Save(&user)

	// Create order
	fTotalPrice, _ := totalPrice.Float64() // Get only the float64 value
	order := models.Order{
		UserID: userID.(uint),
		Total:  fTotalPrice, // 🛠 Now using decimal.Decimal
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items
	for _, item := range cartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			BookID:    item.BookID,
			Quantity:  item.Quantity,
			UnitPrice: item.Book.Price, // 🛠 Now using decimal.Decimal
		}
		db.Create(&orderItem)

		// Reduce stock from books
		item.Book.Stock -= item.Quantity
		db.Save(&item.Book)
	}

	// Clear the cart
	db.Where("user_id = ?", userID).Delete(&models.CartItem{})

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful", "order_id": order.ID})
}

// StartFlashSale starts a flash sale for a book
func StartFlashSale(bookID uint, db *gorm.DB) {
	db.Model(&models.Book{}).Where("id = ?", bookID).
		Updates(map[string]interface{}{
			"flash_sale_active": true,
			"flash_sale_expiry": time.Now().Add(24 * time.Hour), // 24-hour flash sale
		})
}
