package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string          `gorm:"type:varchar(100);not null" json:"name"`
	Email     string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string          `gorm:"type:varchar(255);not null" json:"password"`
	Role      string          `gorm:"type:varchar(50);default:'user'" json:"role"` // Default role is "user"
	Balance   decimal.Decimal `gorm:"type:numeric(10,2);not null;default:0.00" json:"balance"`
	CreatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
type Balance struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint            `gorm:"not null;uniqueIndex" json:"user_id"`                  // ID пользователя
	Balance   decimal.Decimal `gorm:"type:numeric(19,2);not null;default:0" json:"balance"` // Баланс
	CreatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Book represents a book in the system
type Book struct {
	ID              uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string          `gorm:"type:varchar(255);not null" json:"title"`
	Author          string          `gorm:"type:varchar(255);not null" json:"author"`
	Description     string          `gorm:"type:text" json:"description"`
	Price           decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`             // Changed from float64 to decimal.Decimal
	DiscountType    string          `gorm:"type:varchar(10);default:'none'" json:"discount_type"` // 'percentage' or 'fixed'
	DiscountValue   float64         `gorm:"type:decimal(10,2);default:0" json:"discount_value"`
	FlashSaleActive bool            `gorm:"default:false" json:"flash_sale_active"` // Flash sale active?
	FlashSaleExpiry time.Time       `json:"flash_sale_expiry,omitempty"`            // Flash sale end time
	Stock           int             `gorm:"not null" json:"stock"`
	UserID          uint            `gorm:"not null;index" json:"user_id"`                           // Foreign key linking to User
	User            User            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"` // Ensure cascading delete
	CreatedAt       time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Reviews         []Review        `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE;" json:"reviews"`
}

// Order represents an order in the system
type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	UserID     uint        `gorm:"not null" json:"user_id"`
	Total      float64     `gorm:"not null" json:"total"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

type OrderItem struct {
	ID        uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint            `gorm:"not null" json:"order_id"`
	BookID    uint            `gorm:"not null" json:"book_id"`
	Quantity  int             `gorm:"not null" json:"quantity"`
	UnitPrice decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"unit_price"` // Changed from float64 to decimal.Decimal
}

type Review struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookID    uint      `gorm:"not null" json:"book_id"`
	Book      Book      `gorm:"foreignKey:BookID" json:"book"`
	UserID    uint      `gorm:"not null;index" json:"user_id"` // Foreign key linking to User
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Rating    int       `gorm:"not null" json:"rating"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CartItem struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	UserID     uint    `gorm:"not null" json:"user_id"`
	BookID     uint    `gorm:"not null" json:"book_id"`
	Book       Book    `gorm:"foreignKey:BookID" json:"book"`
	Quantity   int     `gorm:"not null;default:1" json:"quantity"`
	TotalPrice float64 `gorm:"not null" json:"total_price"`
}

type Event struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	EventType string    `json:"event_type"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func (b *Book) CalculateFinalPrice() decimal.Decimal {
	discount := decimal.NewFromFloat(0.1)                     // Example: 10% discount
	return b.Price.Mul(decimal.NewFromFloat(1).Sub(discount)) // Use decimal operations
}

func (u *User) SetUserRole(role string) {
	u.Role = role
}
