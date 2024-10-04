package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

//Приветствует пользователя
func HandleHello(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(
		chatID,
		"Привет! Напиши:\n/price, чтобы узнать текущий курс биткоина\n/history, чтобы запросить исторические данные\n/wallet, чтобы открыть кошелек",
	)
	bot.Send(msg)
}
