package handlers

import (
	"crypto_bot/cryptoapi"
	"crypto_bot/database"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

var userBuyCurrency = make(map[int64]string)

func HandleBuy(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "Выберите, какую валюту купить:\n/bitcoin")
		bot.Send(msg)
		userStates[chatID] = [2]string{"buy", "recieved_currency"}
	case "recieved_currency":
		userBuyCurrency[chatID] = update.Message.Text[1:]
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("На какую сумму в USD вы хотите купить %v?", userBuyCurrency[chatID]))
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"buy", "awaiting_amount"}
	case "awaiting_amount":
		if update.Message.Text == "/exit" {
			delete(userStates, chatID)
			delete(userBuyCurrency, chatID)
			HandleHello(bot, chatID)
			return
		}
		// Попробуем преобразовать текст в число
		amountSell, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil || amountSell <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректную сумму в USD или /exit чтобы выйти")
			bot.Send(msg)
			return
		}
		//Запрашиваем актуальный курс биткоина
		price, err := cryptoapi.GetCryptoPrice(userBuyCurrency[chatID])
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при получении цены %v: %v", userBuyCurrency[chatID], err))
			bot.Send(msg)
			delete(userStates, chatID)
			delete(userBuyCurrency, chatID)
			HandleHello(bot, chatID)
			return
		}
		//Считаем сколько биткоинов нужно зачислить
		amountBuy := amountSell / price
		//Обновляем баланс
		err = database.UpdateBalanceBuy(dbConn, update, "usd", userBuyCurrency[chatID], amountSell, amountBuy)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка. Проверьте, достаточно ли USD на балансе с помощью /balance")
			bot.Send(msg)
			delete(userStates, chatID)
			delete(userBuyCurrency, chatID)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Куплено %.6f %v по курсу %v", amountBuy, userBuyCurrency[chatID], price))
		bot.Send(msg)
		//Чистим состояние пользователя
		delete(userStates, chatID)
		delete(userBuyCurrency, chatID)
	}

}
