package crypto_prices

import "time"

type CryptoPrice struct {
	ID        int
	Currency  string
	Price     float64
	CreatedAt time.Time
}
