package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url+"&application_name=stock-api&prefer_simple_protocol=true")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, db.PingContext(context.Background())
}
