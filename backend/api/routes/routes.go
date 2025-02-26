package routes

import (
	"bookhub/backend/api/handlers"
	"bookhub/backend/internal/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes sets up all API routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	// Auth routes
	authGroup := api.Group("/")
	authGroup.Use(auth.AuthMiddleware())
	{
		authGroup.GET("/users", handlers.GetUsers)
		authGroup.GET("/users/:id", handlers.GetUserByID)
		authGroup.PUT("/users/:id", handlers.UpdateUser)
		authGroup.DELETE("/users/:id", handlers.DeleteUser)
		authGroup.GET("/user/profile", handlers.GetUserProfile)
		authGroup.GET("/user/books", handlers.GetUserBooks)

		authGroup.POST("/books", handlers.CreateBook)
		authGroup.PUT("/books/:id", handlers.UpdateBook)
		authGroup.DELETE("/books/:id", handlers.DeleteBook)

		authGroup.POST("/orders", handlers.CreateOrder)

		authGroup.GET("/cart", handlers.GetUserCart)                     // Get all cart items for a user
		authGroup.POST("/cart", handlers.AddBookToCart)                  // Add an item to the cart
		authGroup.PUT("/cart/:book_id", handlers.UpdateCartItemQuantity) // Update quantity of a cart item
		authGroup.DELETE("/cart/:book_id", handlers.RemoveFromCart)      // Remove item
		authGroup.DELETE("/cart/remove/:book_id", handlers.RemoveBookFromCart)
		authGroup.POST("/cart/checkout", handlers.Checkout) // Checkout the cart

		authGroup.GET("/user/events", handlers.GetUserEvents) // Event route

		// Review routes
		reviewHandler := handlers.NewReviewHandler(db)
		authGroup.GET("/books/:id/reviews", reviewHandler.GetBookReviews)
		authGroup.GET("/user/reviews", reviewHandler.GetUserReviews)
		authGroup.POST("/user/reviews", reviewHandler.CreateReview)
		authGroup.PUT("/user/reviews/:id", reviewHandler.UpdateReview)    // Edit Review
		authGroup.DELETE("/user/reviews/:id", reviewHandler.DeleteReview) // Delete Review

		// Balance routes
		authGroup.GET("/users/:id/balance", handlers.GetUserBalance)
		authGroup.POST("/users/:id/balance/topup", handlers.TopUpBalance)

	}

	// Auth routes
	api.POST("/auth/register", auth.Register)
	api.POST("/auth/login", auth.Login)

	// Book routes
	api.GET("/books", handlers.GetBooks)
	api.GET("/books/trend", handlers.GetTrendingBooks)
	api.GET("/books/:id", handlers.GetBookByID)

	// Order routes
	api.GET("/orders", handlers.GetOrders)
	api.GET("/orders/:id", handlers.GetOrderByID)
}
