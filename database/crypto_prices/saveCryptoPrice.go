package crypto_prices

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// SaveCryptoPrice сохраняет курс криптовалюты в базе данных
func SaveCryptoPrice(conn *pgx.Conn, currency string, price float64) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO crypto_prices (currency, price) VALUES ($1, $2)", currency, price)
	return err
}
