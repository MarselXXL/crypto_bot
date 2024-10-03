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

// Обновляет балансы на заданную сумму при покупке/продаже
func UpdateBalanceBuy(conn *pgx.Conn, update tgbotapi.Update, tickerSell string, tickerBuy string, amountSell float64, amountBuy float64) error {
	// Открываем транзакцию
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при открытии транзакции: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background()) // Откат транзакции в случае ошибки
		}
	}()

	// SQL-запрос для обновления сразу двух валют
	query := fmt.Sprintf(`
        UPDATE wallets 
        SET %v = %v - $1, 
            %v = %v + $2
        WHERE tg_name = $3`, tickerSell, tickerSell, tickerBuy, tickerBuy)

	// Выполняем запрос
	_, err = tx.Exec(context.Background(), query, amountSell, amountBuy, update.Message.From.UserName)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении балансов: %v", err)
	}

	// Фиксируем транзакцию
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %v", err)
	}

	return nil
}

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
