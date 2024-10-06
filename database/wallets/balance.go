package wallets

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// Запрашивает баланс юзера
func Balance(conn *pgx.Conn, update tgbotapi.Update) (map[string]float64, error) {
	rows, err := conn.Query(context.Background(), "SELECT usd,bitcoin FROM wallets WHERE tg_name = $1", update.Message.From.UserName)
	if err != nil {
		return nil, fmt.Errorf("1 %v", err)
	}
	defer rows.Close()
	balance := make(map[string]float64)

	rows.Next()

	var a, b float64
	err = rows.Scan(&a, &b)
	balance["usd"] = a
	balance["bitcoin"] = b
	if err != nil {
		return nil, fmt.Errorf("2 %v", err)
	}

	return balance, nil
}
