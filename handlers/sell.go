package handlers

import (
	"crypto_bot/crypto"
	"crypto_bot/database"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func HandleSell(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "Какую сумму BTC вы хотите продать?")
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"sell", "awaiting_amount"}
	case "awaiting_amount":
		// Попробуем преобразовать текст в число
		amountSell, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil || amountSell <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректную сумму в BTC.")
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
		//Считаем сколько USD нужно зачислить
		amountBuy := amountSell * priceBTC
		//Обновляем баланс
		err = database.UpdateBalanceBuy(dbConn, update, "balance_btc", "balance_usd", amountSell, amountBuy)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при обновлении баланса: %v", err))
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Зачислено %.6f USD по курсу %.2f за BTC", amountBuy, priceBTC))
		bot.Send(msg)
		//Чистим состояние пользователя
		delete(userStates, chatID)
	}

}
