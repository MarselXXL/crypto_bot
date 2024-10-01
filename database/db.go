package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type CryptoPrice struct {
	ID        int
	Currency  string
	Price     float64
	CreatedAt time.Time
}

// Connect подключается к базе данных PostgreSQL
func Connect(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
