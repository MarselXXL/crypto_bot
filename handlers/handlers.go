package handlers

import (
	"crypto_bot/crypto"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleMessage обрабатывает входящие сообщения от пользователя
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil { // Если получено сообщение
		text := update.Message.Text
		log.Printf("[%s] %s", update.Message.From.UserName, text)

		if text == "/price" {
			price, err := crypto.GetBitcoinPrice()
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить курс биткоина.")
				bot.Send(msg)
				log.Printf("Для [%s] отправляено: %s", update.Message.From.UserName, msg.Text)
				return
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Курс биткоина: $%.2f", price))
			bot.Send(msg)
			log.Printf("Для [%s] отправляено: %s", update.Message.From.UserName, msg.Text)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Напиши /price, чтобы узнать курс биткоина.")
			bot.Send(msg)
			log.Printf("Для [%s] отправляено: %s", update.Message.From.UserName, msg.Text)
		}
	}
}
