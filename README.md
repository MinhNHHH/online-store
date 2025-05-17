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



### RECOMMENDATIONS TO OPTIMIZE PERFORMANCE:
    Use Redis to Cache HTTP Requests
    To reduce redundant processing and improve response time, it's recommended to use Redis for caching HTTP request results. This is especially useful for endpoints with expensive database operations or third-party API calls.


```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

var (
	ctx   = context.Background()
	rdb   = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

type Response struct {
	Query string `json:"query"`
	Data  string `json:"data"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	cacheKey := fmt.Sprintf("data:%s", query)

	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	time.Sleep(2 * time.Second)
	result := Response{
		Query: query,
		Data:  fmt.Sprintf("Result for %s", query),
	}

	// Encode result as JSON and cache it
	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Store in Redis for 60 seconds
	err = rdb.Set(ctx, cacheKey, jsonData, 60*time.Second).Err()
	if err != nil {
		log.Printf("Failed to write to Redis: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
```
### BE PART2

Product Search Optimization (Focus on data search performance):
● Task: Implement a search API that allows users to search for products by description. You should generate 10 million product records and demonstrate how your system can search for products by description efficiently.
● Objective: This question tests your ability to optimize search operations and handle large datasets.
● Expectations: You need to write an API that can generate 10 million product records, and describe in the README file how someone reviewing your solution can use the API to generate these records. Additionally, explain your approach for efficiently searching the products by description in this large dataset.

I think Elasticsearch it okay for this problem. If I switch to ElasticSearch, I need to refactor a lot my codes. such as: 
1. Create new data model
2. Implement new repository layer for ES
3. Modify Business logic like: update product handler to use the new ES reporsitoy, implement some function that just support for search and insert, ....
4. Update testing.
5. And product have relationship with catergory. Therefore, some action insert or update need to satisfy ACID. So Psql good to me.

I'm using PostgreSQL. I think use PostgreSQL’s built-in tsvector and tsquery functionality to perform full-text search efficiently. 

```bash
	description (TEXT): Stores the raw product description

	description_tsv (TSVECTOR): Stores the tokenized, indexed form for search
```

```sql
-- Automatically update tsvector using trigger
CREATE FUNCTION update_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  NEW.description_tsv := to_tsvector('english', NEW.description);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tsv
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_tsvector_trigger();

-- Index for search performance
CREATE INDEX idx_description_tsv ON products USING GIN(description_tsv);
```

I would say an easy way about GIN
- The description_tsv column (type tsvector) contains tokens (words) extracted from the description.
- The GIN index maps each unique word → the rows that contain it.
Example:
```
| id   | description                                |
| ---- | ------------------------------------------ |
| 1    | A gaming laptop                            |
| 2    | A lightweight ultrabook                    |
| 5    | A lightweight device for travel            |
| 10   | A lightweight gaming laptop for developers |
| 104  | Budget gaming laptop                       |
| 2024 | Best gaming laptop of 2024                 |
```
GIN	Index:
```
| Word          | Row IDs             |
| ------------- | ------------------- |
| `gaming`      | \[1, 10, 104, 2024] |
| `laptop`      | \[1, 10, 104, 2024] |
| `lightweight` | \[2, 5, 10]         |
| `developers`  | \[10]               |
| `budget`      | \[104]              |
```
We can see it is quite similar to the mechanism of ES.

When we search 
```
SELECT * FROM products
WHERE description_tsv @@ to_tsquery('english', 'gaming & laptop');
```
That mean: 
- Look up "gaming" → [1, 10, 104, 2024]
- Look up "laptop" → [1, 10, 104, 2024]

And about `update_tsvector_trigger` that help to automatically keep the description_tsv column up to date with the tokenized version of the description field whenever a product is inserted or updated.

For example:

```sql
INSERT INTO products (name, description)
VALUES ('Gaming Laptop', 'A fast, lightweight gaming laptop');
```

```
Automatically computed:
description_tsv = to_tsvector('english', 'A fast, lightweight gaming laptop')

Resulting tsvector (something like):
'fast':2 'gaming':4 'laptop':5 'lightweight':3
```