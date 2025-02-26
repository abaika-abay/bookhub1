package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bookhub/backend/database/models" // Path to your models
)

type reviewHandler struct {
	db *gorm.DB
}

func NewReviewHandler(db *gorm.DB) *reviewHandler {
	return &reviewHandler{db: db}
}

// GetBookReviews retrieves all reviews for a specific book.
func (h *reviewHandler) GetUserReviews(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.Query("user_id")

	var reviews []struct {
		ID        uint   `json:"id"`
		BookID    uint   `json:"book_id"`
		BookTitle string `json:"book_title"`
		Content   string `json:"content"`
		Rating    int    `json:"rating"`
	}

	err := db.Table("reviews").
		Select("reviews.id, reviews.book_id, books.title as book_title, reviews.text as content, reviews.rating").
		Joins("JOIN books ON books.id = reviews.book_id").
		Where("reviews.user_id = ?", userID).
		Find(&reviews).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *reviewHandler) GetBookReviews(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	bookID := c.Param("id") // Get book ID from URL

	var reviews []models.Review
	if err := db.Where("book_id = ?", bookID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *reviewHandler) CreateReview(c *gin.Context) {
	var review models.Review

	// Bind JSON request
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert userID to int
	userIDInt, ok := userID.(uint)
	if !ok {
		userIDFloat, ok := userID.(float64) // Check if it's float64 (from JSON)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}
		userIDInt = uint(userIDFloat) // Convert to uint
	}

	review.UserID = userIDInt

	// Validate book exists
	var book models.Book
	if err := h.db.First(&book, review.BookID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	// Save review
	if err := h.db.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"review":  review,
	})
}

func (h *reviewHandler) UpdateReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	reviewID := c.Param("id")

	var review models.Review
	if err := db.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	var updateData struct {
		Text   string `json:"text"`
		Rating int    `json:"rating"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	review.Text = updateData.Text
	review.Rating = updateData.Rating
	if err := db.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully", "review": review})
}

func (h *reviewHandler) DeleteReview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	reviewID := c.Param("id")

	if err := db.Delete(&models.Review{}, reviewID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
