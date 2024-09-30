package handlers

import (
	"crypto_bot/crypto"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandlePrice отправляет сообщение с текущей ценой биткоина
func HandlePrice(bot *tgbotapi.BotAPI, chatID int64) {
	price, err := crypto.GetBitcoinPrice()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Не удалось получить курс биткоина.")
		bot.Send(msg)
		//log.Printf("In Chat [%v] sent: %s", chatID, msg.Text)
		return
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Курс биткоина: $%.2f", price))
	bot.Send(msg)
	//log.Printf("In Chat [%v] sent: %s", chatID, msg.Text)
}
