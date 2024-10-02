package database

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func CreateUser(conn *pgx.Conn, update tgbotapi.Update) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO wallets (tg_name) VALUES ($1)", update.Message.From.UserName)
	return err
}

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
