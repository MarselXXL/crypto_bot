package handlers

import (
	"crypto_bot/database"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// HandleHistory запрашивает количество дней
func HandleHistory(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, text string) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "За сколько дней вы хотите получить данные?")
		bot.Send(msg)
		//log.Printf("In Chat [%v] sent: %s", chatID, msg.Text)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"history", "awaiting_days"}
	case "awaiting_days":
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
	}

}
