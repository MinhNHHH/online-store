package dbrepo

import (
	"context"
	"fmt"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
)

func (p *DBRepo) AllCategories(name string, page, pageSize int) ([]*schema.Category, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Base query for counting total records
	countQuery := `select count(*) from categories where 1=1`

	// Base query for fetching records
	query := `select id, name, description from categories where 1=1`

	args := []interface{}{}
	argCount := 1

	if name != "" {
		query += fmt.Sprintf(" and name ILIKE $%d", argCount)
		countQuery += fmt.Sprintf(" and name ILIKE $%d", argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}

	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" order by created_at desc LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pageSize, offset)

	var total int
	err := p.SqlConn.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := p.SqlConn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	categories := []*schema.Category{}
	for rows.Next() {
		var category schema.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, &category)
	}
	return categories, total, nil
}

func (p *DBRepo) InsertCategory(category *schema.Category) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into categories (name, description) values ($1, $2) returning id`
	rows, err := p.SqlConn.QueryContext(ctx, query, category.Name, category.Description)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var newID int
	if !rows.Next() {
		return 0, fmt.Errorf("no rows returned from insert")
	}
	err = rows.Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (p *DBRepo) UpdateCategory(category *schema.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `update categories set name = $1, description = $2 where id = $3`
	_, err := p.SqlConn.ExecContext(ctx, query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *DBRepo) DeleteCategory(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from categories where id = $1`
	_, err := p.SqlConn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
