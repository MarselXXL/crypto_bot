package handlers

import (
	"crypto_bot/crypto"
	"crypto_bot/database"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// Словарь для хранения состояния диалога с пользователем
var userStates = make(map[int64]string)

// HandleMessage обрабатывает входящие сообщения от пользователя
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, dbConn *pgx.Conn) {
	if update.Message != nil {
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		// Обработка команды /price
		if text == "/price" {
			price, err := crypto.GetBitcoinPrice()
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Не удалось получить курс биткоина.")
				bot.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Курс биткоина: $%.2f", price))
			bot.Send(msg)

			// Обработка команды /history
		} else if text == "/history" {
			msg := tgbotapi.NewMessage(chatID, "За сколько дней вы хотите получить данные?")
			bot.Send(msg)

			// Сохраняем состояние пользователя
			userStates[chatID] = "awaiting_days"

			// Обработка числа после /history
		} else if state, exists := userStates[chatID]; exists && state == "awaiting_days" {
			// Попробуем преобразовать текст в число
			days, err := strconv.Atoi(text)
			if err != nil || days <= 0 {
				msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректное количество дней.")
				bot.Send(msg)
				return
			}

			// Запрашиваем данные из базы данных за указанный период
			prices, err := database.GetCryptoPrices(dbConn, "bitcoin", days)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Не удалось получить исторические данные.")
				bot.Send(msg)
				return
			}

			// Формируем ответ с историческими данными
			if len(prices) == 0 {
				msg := tgbotapi.NewMessage(chatID, "Нет данных за указанный период.")
				bot.Send(msg)
			} else {
				response := fmt.Sprintf("Исторические данные по курсу биткоина за последние %d дней:\n", days)
				for _, price := range prices {
					response += fmt.Sprintf("%s: $%.2f\n", price.CreatedAt.Format("2006-01-02 15:04:05"), price.Price)
				}
				msg := tgbotapi.NewMessage(chatID, response)
				bot.Send(msg)
			}

			// Сбрасываем состояние пользователя после обработки
			delete(userStates, chatID)

			// Ответ на другие команды
		} else {
			msg := tgbotapi.NewMessage(chatID, "Привет! Напиши /price, чтобы узнать текущий курс биткоина, или /history, чтобы запросить исторические данные.")
			bot.Send(msg)
		}
	}
}
