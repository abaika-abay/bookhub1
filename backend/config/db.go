package config

import (
	"bookhub/backend/database/models"

	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Order{},
		&models.Review{},
		&models.Event{},
		&models.Balance{},
		&models.CartItem{},
		&models.OrderItem{},
	)
}
