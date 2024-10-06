package crypto_prices

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// GetCryptoPrices получает курсы криптовалюты за последние n минут
func GetCryptoPrices(conn *pgx.Conn, currency string, minutes int) ([]CryptoPrice, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, currency, price, created_at FROM crypto_prices WHERE currency = $1 AND created_at >= NOW() - INTERVAL '1 minute' * $2", currency, minutes)
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
