package handlers

import (
	"crypto_bot/crypto"
	"crypto_bot/database"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// HandleMessage обрабатывает входящие сообщения от пользователя
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, dbConn *pgx.Conn) {
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
		} else if text == "/history" {
			// Запрашиваем данные из БД за последние 7 дней
			prices, err := database.GetCryptoPrices(dbConn, "bitcoin", 7)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить исторические данные.")
				bot.Send(msg)
				return
			}

			// Формируем сообщение с историческими данными
			response := "Исторические данные по курсу биткоина:\n"
			for _, price := range prices {
				response += fmt.Sprintf("%s: $%.2f\n", price.CreatedAt.Format("2006-01-02 15:04:05"), price.Price)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			bot.Send(msg)
			log.Printf("Для [%s] отправляено: %s", update.Message.From.UserName, msg.Text)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Напиши /price, чтобы узнать курс биткоина, или /history, чтобы получить исторические данные.")
			bot.Send(msg)
			log.Printf("Для [%s] отправляено: %s", update.Message.From.UserName, msg.Text)
		}
	}
}
