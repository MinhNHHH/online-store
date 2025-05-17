package dbrepo

import (
	"context"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
)

func (p *DBRepo) AddToWishlist(userID, productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into wishlist (user_id, product_id) values ($1, $2)`
	_, err := p.SqlConn.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (p *DBRepo) RemoveFromWishlist(userID, productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from wishlist where user_id = $1 and product_id = $2`
	_, err := p.SqlConn.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (p *DBRepo) GetWishlist(userID int) ([]*schema.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select p.id, p.name, p.price, p.stock_quantity, p.status, c.name from products p
		inner join product_categories pc on p.id = pc.product_id
		inner join categories c on pc.category_id = c.id
		inner join wishlist w on p.id = w.product_id
		where w.user_id = $1`

	rows, err := p.SqlConn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*schema.Product
	for rows.Next() {
		var product schema.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.StockQuantity, &product.Status, &product.CategoryName)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
