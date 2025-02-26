package utils

import (
	"bookhub/backend/database/models"
	"time"

	"gorm.io/gorm"
)

// LogEvent logs an event into the database
func LogEvent(db *gorm.DB, userID uint, eventType, message string) error {
	event := models.Event{
		UserID:    userID,
		EventType: eventType,
		Message:   message,
		Timestamp: time.Now(),
	}

	return db.Create(&event).Error
}
