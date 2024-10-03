package handlers

import (
	"crypto_bot/crypto"
	"crypto_bot/database"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func HandleBuy(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "На какую сумму в USD вы хотите купить биткоин?")
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"buy", "awaiting_amount"}
	case "awaiting_amount":
		// Попробуем преобразовать текст в число
		amountSell, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil || amountSell <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректную сумму в USD.")
			bot.Send(msg)
			return
		}
		//Запрашиваем актуальный курс биткоина
		priceBTC, err := crypto.GetBitcoinPrice()
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при получении цены биткоина: %v", err))
			bot.Send(msg)
			return
		}
		//Считаем сколько биткоинов нужно зачислить
		amountBuy := amountSell / priceBTC
		//Обновляем баланс
		err = database.UpdateBalanceBuy(dbConn, update, "balance_usd", "balance_btc", amountSell, amountBuy)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при обновлении баланса: %v", err))
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Куплено %.6f BTC по курсу %v", amountBuy, priceBTC))
		bot.Send(msg)
		//Чистим состояние пользователя
		delete(userStates, chatID)
	}

}
