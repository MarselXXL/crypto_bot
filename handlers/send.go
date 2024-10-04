package handlers

import (
	"crypto_bot/database"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

var userSendCurrency = make(map[int64]string)
var userSendAmount = make(map[int64]float64)

func HandleSend(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "Выберите, какую валюту отправить:\n/usd\n/bitcoin")
		bot.Send(msg)
		userStates[chatID] = [2]string{"send", "awaiting_currency"}
	case "awaiting_currency":
		userSendCurrency[chatID] = update.Message.Text[1:]
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
			"Какое количество %v вы хотите отправить?",
			userSendCurrency[chatID],
		))
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"send", "awaiting_amount"}
	case "awaiting_amount":
		if update.Message.Text == "/exit" {
			delete(userStates, chatID)
			delete(userSendCurrency, chatID)
			HandleHello(bot, chatID)
			return
		}
		// Попробуем преобразовать текст в число
		sendAmount, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil || sendAmount <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректную сумму или /exit чтобы выйти")
			bot.Send(msg)
			return
		}
		userSendAmount[chatID] = sendAmount
		msg := tgbotapi.NewMessage(chatID, "Введите через @ имя пользователя, которому хотите отправить валюту")
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"send", "awaiting_reciever_name"}
	case "awaiting_reciever_name":
		recieverName := update.Message.Text[1:]
		err := database.Send(dbConn, update, userSendCurrency[chatID], userSendAmount[chatID], recieverName)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при отправке:\n%v", err))
			bot.Send(msg)
			delete(userStates, chatID)
			delete(userSendCurrency, chatID)
			delete(userSendAmount, chatID)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v отправлено %v %v", recieverName, userSendAmount[chatID], userSendCurrency[chatID]))
		bot.Send(msg)
		delete(userStates, chatID)
		delete(userSendCurrency, chatID)
		delete(userSendAmount, chatID)
		return
	}
}
