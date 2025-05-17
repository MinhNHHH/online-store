package dbrepo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
)

type TestDBRepo struct{}

func (p *TestDBRepo) SQLConnection() *sql.DB {
	return nil
}

func (p *TestDBRepo) AllCategories(name string, page, pageSize int) ([]*schema.Category, int, error) {
	var categories []*schema.Category
	mocks := []*schema.Category{
		{
			ID:          1,
			Name:        "test",
			Description: "test",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "test2",
			Description: "test2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	categories = append(categories, mocks...)

	return categories, 0, nil
}

func (p *TestDBRepo) InsertCategory(category *schema.Category) (int, error) {
	category.ID = 1
	return category.ID, nil
}

func (p *TestDBRepo) UpdateCategory(category *schema.Category) error {
	return nil
}

func (p *TestDBRepo) DeleteCategory(id int) error {
	return nil
}

func (p *TestDBRepo) AllProducts(name, categoryName, status string, page, pageSize int) ([]*schema.Product, int, error) {
	return nil, 0, nil
}

func (p *TestDBRepo) InsertProduct(product *schema.Product) (int, error) {
	return 0, nil
}

func (p *TestDBRepo) UpdateProduct(product *schema.Product) error {
	return nil
}

func (p *TestDBRepo) DeleteProduct(id int) error {
	return nil
}

func (p *TestDBRepo) DeleteReview(id int) error {
	return nil
}

func (p *TestDBRepo) AddToWishlist(userID, productID int) error {
	return nil
}

func (p *TestDBRepo) RemoveFromWishlist(userID, productID int) error {
	return nil
}

func (p *TestDBRepo) GetWishlist(userID int) ([]*schema.Product, error) {
	return nil, nil
}

func (p *TestDBRepo) AllUsers() ([]*schema.User, error) {
	return nil, nil
}

func (p *TestDBRepo) GetUser(id int) (*schema.User, error) {
	return nil, nil
}

func (p *TestDBRepo) GetUserByEmail(email string) (*schema.User, error) {
	if email == "admin@example.com" {
		user := schema.User{
			ID:        1,
			Name:      "Admin",
			Email:     "admin@example.com",
			Password:  "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return &user, nil
	}
	return nil, errors.New("not found")
}

func (p *TestDBRepo) UpdateUser(user schema.User) error {
	return nil
}

func (p *TestDBRepo) DeleteUser(id int) error {
	return nil
}

func (p *TestDBRepo) InsertUser(user schema.User) (int, error) {
	return 0, nil
}

func (p *TestDBRepo) ResetPassword(id int, password string) error {
	return nil
}

func (p *TestDBRepo) ReviewsByProductID(productID int) ([]*schema.Review, error) {
	return nil, nil
}

func (p *TestDBRepo) InsertReview(review *schema.Review) (int, error) {
	return 0, nil
}
