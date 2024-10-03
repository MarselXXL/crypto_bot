package handlers

import (
	"crypto_bot/cryptoapi"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandlePrice отправляет сообщение с текущей ценой биткоина
func HandlePrice(bot *tgbotapi.BotAPI, chatID int64, update tgbotapi.Update) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "Выберите валюту:\n/bitcoin\n/ethereum")
		bot.Send(msg)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"price", "awaiting_currency"}
	case "awaiting_currency":
		price, err := cryptoapi.GetCryptoPrice(update.Message.Text[1:])
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить курс.")
			bot.Send(msg)
			//log.Printf("In Chat [%v] sent: %s", chatID, msg.Text)
			return
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Курс %v: $%.2f", update.Message.Text[1:], price))
		bot.Send(msg)
		//log.Printf("In Chat [%v] sent: %s", chatID, msg.Text)
		delete(userStates, chatID)
	}
}
