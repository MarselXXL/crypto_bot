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
		case text == "/price" || userStates[chatID][0] == "price":
			HandlePrice(bot, chatID, update)

			// Обработка команды /history
		case text == "/history" || userStates[chatID][0] == "history":
			HandleHistory(bot, chatID, dbConn, text)

			// Обработка команды /wallet
		case text == "/wallet":
			HandleWallet(bot, chatID, dbConn, update)

			// Обработка команды /add_test_balance
		case text == "/add_test_balance" || userStates[chatID][0] == "add_test_balance":
			HandleAdd_test_balance(bot, chatID, dbConn, update)
			// Обработка команды /balance
		case text == "/balance":
			HandleBalance(bot, chatID, dbConn, update)
			// Обработка команды /buy
		case text == "/buy" || userStates[chatID][0] == "buy":
			HandleBuy(bot, chatID, dbConn, update)
			// Обработка команды /sell
		case text == "/sell" || userStates[chatID][0] == "sell":
			HandleSell(bot, chatID, dbConn, update)
			// Обработка команды /send
		case text == "/send" || userStates[chatID][0] == "send":
			HandleSend(bot, chatID, dbConn, update)

			// Ответ на другие команды
		default:
			HandleHello(bot, chatID)
		}
	}
}
