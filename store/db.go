package store

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func LoadDBWithBunClient() (*bun.DB, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	sqldb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}
	db := bun.NewDB(sqldb, pgdialect.New())
	return db, nil
}
