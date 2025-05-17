package schema

import "time"

type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"is_admin,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	Status        string  `json:"status,omitempty"`
	CategoryID    int     `json:"category_id,omitempty"`
	CategoryName  string  `json:"category_name,omitempty"`
}

type Category struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type ProductCategory struct {
	ProductID  int `json:"product_id"`
	CategoryID int `json:"category_id"`
}

type Review struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	UserID      int       `json:"user_id"`
	UserName    string    `json:"user_name,omitempty"`
	ProductName string    `json:"product_name,omitempty"`
	Rating      int       `json:"rating" validate:"required,min=1,max=5"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
