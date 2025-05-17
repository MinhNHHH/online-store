package dbrepo

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

type DBRepo struct {
	SqlConn *sql.DB
}

func (p *DBRepo) SQLConnection() *sql.DB {
	return p.SqlConn
}
