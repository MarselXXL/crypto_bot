package wallets

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// Обновляет баланс на заданную сумму
func UpdateBalance(conn *pgx.Conn, update tgbotapi.Update, ticker string, sign bool, amount float64) error {
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
		query = fmt.Sprintf("UPDATE wallets SET %v = %v + $1 WHERE tg_name = $2", ticker, ticker)
	} else {
		query = fmt.Sprintf("UPDATE wallets SET %v = %v - $1 WHERE tg_name = $2", ticker, ticker)
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
