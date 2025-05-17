package dbrepo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
)

func (p *DBRepo) AllProducts(name, categoryName, status string, page, pageSize int) ([]*schema.Product, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Base query for counting total records
	countQuery := `select count(*) from products as p 
		inner join product_categories as pc on p.id = pc.product_id
		inner join categories as c on pc.category_id = c.id 
		where 1=1`

	// Base query for fetching records
	query := `select p.id, p.name, p.description, p.price, p.stock_quantity, p.status, c.name
		from products as p 
		inner join product_categories as pc on p.id = pc.product_id
		inner join categories as c on pc.category_id = c.id 
		where 1=1`

	args := []interface{}{}
	argCount := 1

	if name != "" {
		query += fmt.Sprintf(" AND p.name ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" AND p.name ILIKE $%d", argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}

	if categoryName != "" {
		query += fmt.Sprintf(" AND c.name ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" AND c.name ILIKE $%d", argCount)
		args = append(args, "%"+categoryName+"%")
		argCount++
	}

	if status != "" {
		query += fmt.Sprintf(" AND p.status ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" AND p.status ILIKE $%d", argCount)
		args = append(args, "%"+status+"%")
		argCount++
	}

	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" order by p.created_at desc LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pageSize, offset)

	// Get total count
	var total int
	err := p.SqlConn.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	rows, err := p.SqlConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []*schema.Product{}
	for rows.Next() {
		var product schema.Product
		var priceStr string
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&priceStr,
			&product.StockQuantity,
			&product.Status,
			&product.CategoryName,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, 0, err
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Println("Error parsing price", err)
			return nil, 0, err
		}
		product.Price = price

		products = append(products, &product)
	}
	return products, total, nil
}

func (p *DBRepo) InsertProduct(product *schema.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	tx, err := p.SqlConn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt := `insert into products (name, description, price, stock_quantity, status, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = p.SqlConn.QueryRowContext(ctx, stmt,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Status,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	stmt = `insert into product_categories (product_id, category_id)
		values ($1, $2)`

	_, err = tx.ExecContext(ctx, stmt, newID, product.CategoryID)

	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newID, nil
}

func (p *DBRepo) UpdateProduct(product *schema.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update products set
		name = $1,
		description = $2,
		price = $3,
		stock_quantity = $4,
		status = $5,
		updated_at = $6
		where id = $7
	`

	_, err := p.SqlConn.ExecContext(ctx, stmt,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Status,
		time.Now(),
		product.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *DBRepo) DeleteProduct(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from products where id = $1`

	_, err := p.SqlConn.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	return nil
}
