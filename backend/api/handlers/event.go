package handlers

import (
	"bookhub/backend/database/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LogEvent(db *gorm.DB, userID uint, eventType, message string) {
	event := models.Event{
		UserID:    userID,
		EventType: eventType,
		Message:   message,
		Timestamp: time.Now(),
	}

	if err := db.Create(&event).Error; err != nil {
		log.Printf("Failed to log event: %v", err)
	}
}

func GetUserEvents(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID, exists := c.Get("user_id") // Get user ID from JWT

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var events []models.Event
	if err := db.Where("user_id = ?", userID).Order("timestamp DESC").Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user events"})
		return
	}

	c.JSON(http.StatusOK, events)
}
