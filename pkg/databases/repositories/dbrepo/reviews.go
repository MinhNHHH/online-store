package dbrepo

import (
	"context"
	"fmt"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
)

func (p *DBRepo) ReviewsByProductID(productID int) ([]*schema.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select r.id, p.name, u.name, r.rating, r.comment 
		from reviews r
		inner join products p on r.product_id = p.id
		inner join users u on r.user_id = u.id
		where r.product_id = $1`

	rows, err := p.SqlConn.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*schema.Review
	for rows.Next() {
		var review schema.Review
		err := rows.Scan(&review.ID, &review.ProductName, &review.UserName, &review.Rating, &review.Comment)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}
	return reviews, nil
}

func (p *DBRepo) InsertReview(review *schema.Review) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into reviews (product_id, user_id, rating, comment) values ($1, $2, $3, $4) returning id`
	rows, err := p.SqlConn.QueryContext(ctx, query, review.ProductID, review.UserID, review.Rating, review.Comment)
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

func (p *DBRepo) DeleteReview(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from reviews where id = $1`
	_, err := p.SqlConn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
