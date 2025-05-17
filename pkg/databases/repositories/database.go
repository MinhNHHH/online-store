package databases

import (
	"database/sql"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
)

type DatabaseRepo interface {
	SQLConnection() *sql.DB
	AllCategories(name string, page, pageSize int) ([]*schema.Category, int, error)
	InsertCategory(category *schema.Category) (int, error)
	UpdateCategory(category *schema.Category) error
	DeleteCategory(id int) error
	AllProducts(name, description, categoryName, status string, page, pageSize int) ([]*schema.Product, int, error)
	InsertProduct(product *schema.Product) (int, error)
	UpdateProduct(product *schema.Product) error
	DeleteProduct(id int) error
	ReviewsByProductID(productID int) ([]*schema.Review, error)
	InsertReview(review *schema.Review) (int, error)
	DeleteReview(id int) error
	AddToWishlist(userID, productID int) error
	RemoveFromWishlist(userID, productID int) error
	GetWishlist(userID int) ([]*schema.Product, error)
	AllUsers() ([]*schema.User, error)
	GetUser(id int) (*schema.User, error)
	GetUserByEmail(email string) (*schema.User, error)
	UpdateUser(user schema.User) error
	DeleteUser(id int) error
	InsertUser(user schema.User) (int, error)
	ResetPassword(id int, password string) error
}
