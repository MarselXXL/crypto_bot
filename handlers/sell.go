package handlers

import (
	"crypto_bot/cryptoapi"
	"crypto_bot/database"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

var userSellCurrency = make(map[int64]string)

func HandleSell(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "Выберите, какую валюту продать:\n/bitcoin")
		bot.Send(msg)
		userStates[chatID] = [2]string{"sell", "recieved_currency"}
	case "recieved_currency":
		userSellCurrency[chatID] = update.Message.Text[1:]
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Какую сумму в %v вы хотите продать?", userSellCurrency[chatID]))
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"sell", "awaiting_amount"}
	case "awaiting_amount":
		// Попробуем преобразовать текст в число
		amountSell, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil || amountSell <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректную сумму.")
			bot.Send(msg)
			return
		}
		//Запрашиваем актуальный курс биткоина
		price, err := cryptoapi.GetCryptoPrice(userSellCurrency[chatID])
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при получении цены %v: %v", userSellCurrency[chatID], err))
			bot.Send(msg)
			return
		}
		//Считаем сколько биткоинов нужно зачислить
		amountBuy := amountSell * price
		//Обновляем баланс
		err = database.UpdateBalanceBuy(dbConn, update, userSellCurrency[chatID], "usd", amountSell, amountBuy)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при обновлении баланса: %v", err))
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Продано %.6f %v по курсу %v", amountSell, userSellCurrency[chatID], price))
		bot.Send(msg)
		//Чистим состояние пользователя
		delete(userStates, chatID)
		delete(userBuyCurrency, chatID)
	}

}
