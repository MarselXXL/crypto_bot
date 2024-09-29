package main

import (
	"crypto_bot/handlers"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := "7976608763:AAEIvnuRONaEfO6UOR8QSjKYAXV1_LL8eKY"
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not provided")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handlers.HandleMessage(bot, update) // Обрабатываем сообщение
	}
}
