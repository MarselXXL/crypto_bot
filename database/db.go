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

// SaveCryptoPrice сохраняет курс криптовалюты в базе данных
func SaveCryptoPrice(conn *pgx.Conn, currency string, price float64) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO crypto_prices (currency, price) VALUES ($1, $2)", currency, price)
	return err
}

// GetCryptoPrices получает курсы криптовалюты за последние n дней
func GetCryptoPrices(conn *pgx.Conn, currency string, days int) ([]CryptoPrice, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, currency, price, created_at FROM crypto_prices WHERE currency = $1 AND created_at >= NOW() - INTERVAL '1 day' * $2", currency, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []CryptoPrice
	for rows.Next() {
		var price CryptoPrice
		err := rows.Scan(&price.ID, &price.Currency, &price.Price, &price.CreatedAt)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	return prices, nil
}
