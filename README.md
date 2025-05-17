# Online Store API

A robust and scalable RESTful API for an online store built with Go, featuring user authentication, product management, categories, reviews, and wishlist functionality.

## Overview

This application provides a complete backend solution for an online store with the following features:
- User authentication and authorization (JWT-based)
- Product management with categories
- Product reviews and ratings
- User wishlists
- RESTful API endpoints
- PostgreSQL database with migrations
- Docker containerization

## Technologies Used

- **Backend Framework**: Go (Golang)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Tokens)
- **API Router**: Chi Router
- **Database Migrations**: golang-migrate
- **Environment Variables**: godotenv
- **Testing**: testify
- **Session Management**: scs

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL 14 (if running locally without Docker)

## Setup Instructions

1. Clone the repository:
```bash
git clone https://github.com/MinhNHHH/online-store.git
cd online-store
```

2. Create a `.env` file in the root directory with the following variables:
```env
SO_DB_CONNECTION_URI=your_db_uri
SO_POSTGRES_USER=your_db_user
SO_POSTGRES_PASSWORD=your_db_password
SO_POSTGRES_DB=your_db_name
SO_JWT_SECRET=your_jwt_secret
SO_DOMAIN=localhost
```


3. Run database migrations:
   1. Create migration.
	```bash
		go run cmd/store.go create [name_migration]
	``` 
   2. Migrate
   ```bash
		// TODO: step = 0 means run all migrations
		// step = n means run only the next n migrations
		// step = -n means run only the previous n migrations
		go run cmd/store.go migrate
	```


4. Run application:
```bash
	docker-compose run up -d
```

5. Running Tests

To run the unit tests:
```bash
go test ./... -v
```

## Database Schema

The database consists of the following tables:

### Users
- id (Primary Key)
- name
- email (Unique)
- password
- is_admin
- created_at
- updated_at

### Products
- id (Primary Key)
- name
- description
- price
- stock_quantity
- status (in_stock, out_of_stock, draft)
- created_at
- updated_at

### Categories
- id (Primary Key)
- name
- description
- created_at

### Product Categories (Junction Table)
- product_id (Foreign Key)
- category_id (Foreign Key)

### Reviews
- id (Primary Key)
- product_id (Foreign Key)
- user_id (Foreign Key)
- rating (1-5)
- comment
- created_at
- updated_at

### Wishlist
- user_id (Foreign Key)
- product_id (Foreign Key)
- added_at

## API Documentation

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword"
}
```

#### Login
```http
POST /api/v1/auth
Content-Type: application/json

{
    "email": "john@example.com",
    "password": "securepassword"
}
```

### Products

#### Get All Products
```http
GET /api/v1/products
Authorization: Bearer <jwt_token>
```

#### Create Product
```http
POST /api/v1/products
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "name": "Product Name",
    "description": "Product Description",
    "price": 99.99,
    "stock_quantity": 100,
    "categories": [1, 2]
}
```

### Categories

#### Get All Categories
```http
GET /api/v1/categories
Authorization: Bearer <jwt_token>
```

#### Create Category
```http
POST /api/v1/categories
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "name": "Category Name",
    "description": "Category Description"
}
```

#### Update Category
```http
PUT /api/v1/categories/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "name": "Updated Category Name",
    "description": "Updated Category Description"
}
```

#### Delete Category
```http
DELETE /api/v1/categories/{id}
Authorization: Bearer <jwt_token>
```

### Reviews

#### Get Product Reviews
```http
GET /api/v1/reviews/{product_id}/
Authorization: Bearer <jwt_token>
```

#### Create Review
```http
POST /api/v1/reviews/{product_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
    "rating": 5,
    "comment": "Great product! Highly recommended."
}
```

#### Delete Review
```http
DELETE /api/v1/reviews/{id}
Authorization: Bearer <jwt_token>
```

### Wishlist

#### Get User's Wishlist
```http
GET /api/v1/users/wishlist
Authorization: Bearer <jwt_token>
```

#### Add Product to Wishlist
```http
POST /api/v1/users/wishlist
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
	"user_id": 123,
    "product_id": 123
}
```

#### Remove Product from Wishlist
```http
DELETE /api/v1/users/wishlist
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
	"product_id": 123,
	"user_id":123
}
```
