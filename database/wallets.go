package database

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// Создает запись о новом юзере с балансом 0
func CreateUser(conn *pgx.Conn, update tgbotapi.Update) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO wallets (tg_name) VALUES ($1)", update.Message.From.UserName)
	return err
}

// Обновляет баланс на заданную сумму
func UpdateBalance(conn *pgx.Conn, update tgbotapi.Update, sign bool, amount float64) error {
	var query string
	// Открываем транзакцию
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при открытии транзакции: %v", err)
	}
	defer func() {
		// Откатываем транзакцию в случае ошибки
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()
	// Определяем будем прибавлять или отнимать
	if sign {
		query = "UPDATE wallets SET balance_usd = balance_usd + $1 WHERE tg_name = $2"
	} else {
		query = "UPDATE wallets SET balance_usd = balance_usd - $1 WHERE tg_name = $2"
	}
	// Выполняем запрос
	_, err = tx.Exec(context.Background(), query, amount, update.Message.From.UserName)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса: %v", err)
	}
	// Коммит
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %v", err)
	}
	return nil

}

// Запрашивает баланс юзера
func Balance(conn *pgx.Conn, update tgbotapi.Update) (map[string]float64, error) {
	rows, err := conn.Query(context.Background(), "SELECT balance_usd,balance_btc FROM wallets WHERE tg_name = $1", update.Message.From.UserName)
	if err != nil {
		return nil, fmt.Errorf("1 %v", err)
	}
	defer rows.Close()
	balance := make(map[string]float64)

	rows.Next()

	var a, b float64
	err = rows.Scan(&a, &b)
	balance["USD"] = a
	balance["BTC"] = b
	if err != nil {
		return nil, fmt.Errorf("2 %v", err)
	}

	return balance, nil
}
