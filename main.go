package main

import (
	"context"
	"crypto_bot/crypto"
	"crypto_bot/database"
	"crypto_bot/handlers"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := "7976608763:AAEIvnuRONaEfO6UOR8QSjKYAXV1_LL8eKY"
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not provided")
	}

	connString := "postgres://postgres:111111@localhost:5432/crypto_db" // Замените на свои данные
	dbConn, err := database.Connect(connString)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	log.Println("Подключено к crypto_db")
	defer dbConn.Close(context.Background())

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Запускаем горутину для периодического сохранения курса биткоина
	go func() {
		for {
			price, err := crypto.GetBitcoinPrice()
			if err != nil {
				log.Println("Error fetching Bitcoin price:", err)
				time.Sleep(1 * time.Minute) // Ждем минуту перед следующей попыткой
				continue
			}

			// Сохраняем курс в базу данных
			err = database.SaveCryptoPrice(dbConn, "bitcoin", price)
			if err != nil {
				log.Println("Error saving Bitcoin price:", err)
			}

			time.Sleep(1 * time.Minute) // Ждем минуту перед следующим сохранением
		}
	}()

	for update := range updates {
		handlers.HandleMessage(bot, update) // Обрабатываем сообщение
	}
}
