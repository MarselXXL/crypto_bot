package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

// Словарь состояний диалогов. Ключь chatID, Значение [Топик, Состояние]
var userStates = make(map[int64][2]string)

// HandleMessage обрабатывает входящие сообщения от пользователя
func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, dbConn *pgx.Conn) {
	if update.Message != nil {
		chatID := update.Message.Chat.ID
		text := update.Message.Text
		//log.Printf("User [%s] Chat [%v]: %s", update.Message.From.UserName, chatID, text)

		//Обработка команд
		switch {
		// Обработка команды /price
		case text == "/price":
			HandlePrice(bot, chatID)

			// Обработка команды /history
		case text == "/history" || userStates[chatID][0] == "history":
			HandleHistory(bot, chatID, dbConn, text)

			// Обработка команды /wallet
		case text == "/wallet":
			HandleWallet(bot, chatID, dbConn, update)

			// Ответ на другие команды
		default:
			msg := tgbotapi.NewMessage(
				chatID,
				"Привет! Напиши:\n/price, чтобы узнать текущий курс биткоина\n/history, чтобы запросить исторические данные\n/wallet, чтобы открыть кошелек",
			)
			bot.Send(msg)
		}
	}
}
