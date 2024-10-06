package wallets

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

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
