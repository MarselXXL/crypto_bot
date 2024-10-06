package wallets

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// Отправляет указанное кол-во валюты юзеру
func Send(conn *pgx.Conn, update tgbotapi.Update, sendCurrency string, sendAmount float64, recieverName string) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при открытии транзакции: %w", err)
	}
	// Откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()
	// Запрос на уменьшение баланса отправителя
	query1 := fmt.Sprintf(`
        UPDATE wallets 
        SET %v = %v - $1 
        WHERE tg_name = $2`, sendCurrency, sendCurrency)
	// Выполняем запрос 1
	_, err = tx.Exec(context.Background(), query1, sendAmount, update.Message.From.UserName)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении балансоа отправителя: %w", err)
	}
	// Запрос на увеличение баланса получателя
	query2 := fmt.Sprintf(`
        UPDATE wallets 
        SET %v = %v + $1 
        WHERE tg_name = $2`, sendCurrency, sendCurrency)
	// Выполняем запрос 2
	_, err = tx.Exec(context.Background(), query2, sendAmount, recieverName)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса полуателя: %w", err)
	}
	// Фиксируем транзакцию
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %v", err)
	}
	// Все ок
	return nil
}
