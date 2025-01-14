package handlers

import (
	"crypto_bot/database/crypto_prices"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

var userHistoryCurrency = make(map[int64]string)

// HandleHistory запрашивает количество дней
func HandleHistory(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, text string) {
	switch userStates[chatID][1] {
	case "":
		msg := tgbotapi.NewMessage(chatID, "Выберите валюту:\n/bitcoin")
		bot.Send(msg)
		userStates[chatID] = [2]string{"history", "recieved_currency"}
	case "recieved_currency":
		userHistoryCurrency[chatID] = text[1:]
		msg := tgbotapi.NewMessage(chatID, "За сколько минут вы хотите получить данные?")
		bot.Send(msg)
		//log.Printf("In Chat [%v] sent: %s", chatID, msg.Text)

		// Сохраняем состояние пользователя
		userStates[chatID] = [2]string{"history", "awaiting_minutes"}
	case "awaiting_minutes":
		if text == "/exit" {
			delete(userStates, chatID)
			HandleHello(bot, chatID)
			return
		}
		// Попробуем преобразовать текст в число
		days, err := strconv.Atoi(text)
		if err != nil || days <= 0 {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректное количество минут или /exit чтобы выйти")
			bot.Send(msg)
			return
		}

		// Запрашиваем данные из базы данных за указанный период
		prices, err := crypto_prices.GetCryptoPrices(dbConn, userHistoryCurrency[chatID], days)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Не удалось получить исторические данные.")
			bot.Send(msg)
			delete(userStates, chatID)
			delete(userHistoryCurrency, chatID)
			return
		}

		// Формируем ответ с историческими данными
		if len(prices) == 0 {
			msg := tgbotapi.NewMessage(chatID, "Нет данных за указанный период.")
			bot.Send(msg)
		} else {
			response := fmt.Sprintf("Исторические данные по курсу %v за последние %d мин.:\n", userHistoryCurrency[chatID], days)
			for _, price := range prices {
				response += fmt.Sprintf("%s: $%.2f\n", price.CreatedAt.Format("2006-01-02 15:04:05"), price.Price)
			}
			msg := tgbotapi.NewMessage(chatID, response)
			bot.Send(msg)
		}

		// Сбрасываем состояние пользователя после обработки
		delete(userStates, chatID)
		delete(userHistoryCurrency, chatID)
	}

}
