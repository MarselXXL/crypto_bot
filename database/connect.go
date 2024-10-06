package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Connect подключается к базе данных PostgreSQL
func Connect(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
