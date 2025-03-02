# BookHub - Online Bookstore

## 📖 Overview
BookHub is an **online bookstore** where users can browse books, purchase them, leave reviews, and manage their shopping carts. The platform also features an **admin panel** for managing users, tracking events, and handling book inventory.

## 🚀 Features
### **User Features**
- 📚 **View Books** - Browse a collection of books.
- 🛒 **Shopping Cart** - Add, remove, and manage books in the cart.
- 💰 **Balance System** - Users can top-up their balance and make purchases.
- ✍️ **Leave Reviews** - Users can leave reviews and ratings (1-5) for books.
- 🔐 **Authentication** - Secure login and registration.

### **Admin Features**
- 🛠️ **Manage Books** - CRUD operations (Create, Read, Update, Delete).
- 👥 **Manage Users** - View user actions and delete/edit their books/reviews.
- 📊 **Event Logs** - Monitor user activities and track events.

## 🏗️ Tech Stack
### **Frontend**
- **HTML, CSS, JavaScript** (Vanilla JS)
- **Fetch API** for making HTTP requests

### **Backend**
- **Golang (Gin Framework)** - REST API
- **PostgreSQL** - Database
- **GORM** - ORM for database interactions
- **JWT (JSON Web Token)** - Secure authentication

## 📂 Folder Structure
```
BookHub/
│── backend/       # Golang Backend API
│   ├── cmd/bookhub       # The Main
│   ├── api/
│   │   ├── handlers/     # Business logic (Books, Users,Cart, etc.)
│   │   ├── routes/       # API Endpoints
│   ├── internal/auth/    # JWT Authentication Middleware
│   ├── database/        
│   │   ├── models/       # Database models (User, Book, Cart, etc.)
│   ├── config/           # Database migrations
│   ├── utils/            # Utils
│── frontend/      # Frontend (HTML, CSS, JS)
│   ├── index.html        # Homepage
│   ├── cart.html         # Shopping Cart
│   ├── profile.html      # User Profile & Reviews
│   ├── odds.html         # Review Dashboard
│   ├── book.html         # Book Details
│   ├── registration.html # Book Details
│── README.md       # Project Documentation
```

## ⚙️ Setup & Installation
### **1️⃣ Clone the Repository**
```bash
git clone https://github.com/trippingcoin/BookHub.git
cd BookHub
```

### **2️⃣ Backend Setup**
#### Install Dependencies
```bash
cd backend
go mod tidy
```
#### Run the Backend Server
```bash
go run cmd/bookhub/main.go
```

### **3️⃣ Frontend Setup**
Start a simple local server (e.g., using Python):
```bash
cd frontend
python3 -m http.server 8000
```
Then, open `http://localhost:8000/index.html` in your browser.

## 📡 API Endpoints
### **User Authentication**
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and get a JWT token

### **Books**
- `GET /api/books` - Get all books
- `GET /api/books/:id` - Get book details
- `POST /api/books` - Add a new book (Admin only)
- `PUT /api/books/:id` - Update book details (Admin only)
- `DELETE /api/books/:id` - Delete a book (Admin only)

### **Cart & Orders**
- `GET /api/cart` - View user cart
- `POST /api/cart` - Add book to cart
- `PUT /api/cart/:id` - Update book quantity in cart
- `DELETE /api/cart/:id` - Remove book from cart
- `POST /api/cart/checkout` - Checkout and place an order

### **Reviews**
- `GET /api/user/reviews` - View user reviews
- `POST /api/user/reviews` - Add a review
- `PUT /api/user/reviews/:id` - Edit a review
- `DELETE /api/user/reviews/:id` - Delete a review

## ✅ Future Improvements
- 📦 **Implement Book Categories**
- 📊 **Add Sales & Order History Analytics**
- 📜 **Implement Wishlist & Favorites**

## 🎉 Contributing
Contributions are welcome! Feel free to submit issues or pull requests.

## 📜 License
MIT License - You are free to use and modify this project.

#   b o o k h u b 1  
 