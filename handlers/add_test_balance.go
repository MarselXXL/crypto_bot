package handlers

import (
	"crypto_bot/database/wallets"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func HandleAdd_test_balance(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "Сколько USD вы хотите добавить на счёт?")
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"add_test_balance", "awaiting_amount"}
	case "awaiting_amount":
		if update.Message.Text == "/exit" {
			delete(userStates, chatID)
			return
		}
		// Попробуем преобразовать текст в число
		amount, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil || amount < 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректную сумму в USD или /exit")
			bot.Send(msg)
			return
		}
		//Обновляем баланс
		err = wallets.UpdateBalance(dbConn, update, "usd", true, float64(amount))
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при обновлении баланса: %v", err))
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(chatID, "Баланс обновлён")
		bot.Send(msg)
		//Чистим состояние пользователя
		delete(userStates, chatID)
	}

}
