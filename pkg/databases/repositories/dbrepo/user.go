package dbrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"golang.org/x/crypto/bcrypt"
)

func (p *DBRepo) AllUsers() ([]*schema.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, is_admin, created_at, updated_at
	from users order by last_name`

	rows, err := p.SqlConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*schema.User

	for rows.Next() {
		var user schema.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Password,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (p *DBRepo) GetUser(id int) (*schema.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			u.id, u.email, u.name, u.password, u.is_admin, u.created_at, u.updated_at
		from 
			users u
		where
		    u.id = $1`

	var user schema.User
	row := p.SqlConn.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *DBRepo) GetUserByEmail(email string) (*schema.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
			u.id, u.email, u.name, u.password, u.is_admin, u.created_at, u.updated_at
		from 
			users u
		where
		    u.email = $1`

	var user schema.User
	row := p.SqlConn.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *DBRepo) UpdateUser(u schema.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		name = $2,
		is_admin = $4,
		updated_at = $5
		where id = $6
	`

	_, err := p.SqlConn.ExecContext(ctx, stmt,
		u.Email,
		u.Name,
		u.IsAdmin,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *DBRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := p.SqlConn.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *DBRepo) InsertUser(user schema.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, name, password, is_admin, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id`

	rows, err := p.SqlConn.QueryContext(ctx, stmt,
		user.Email,
		user.Name,
		hashedPassword,
		user.IsAdmin,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("no rows returned from insert")
	}

	err = rows.Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (p *DBRepo) ResetPassword(id int, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `update users set password = $1 where id = $2`
	_, err = p.SqlConn.ExecContext(ctx, stmt, hashedPassword, id)
	if err != nil {
		return err
	}

	return nil
}
