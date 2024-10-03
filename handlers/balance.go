package handlers

import (
	"crypto_bot/database"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func HandleBalance(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	balance, err := database.Balance(dbConn, update)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ошибка при получении баланса: %v", err))
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Баланс:\nUSD: %.2f\nBitcoin: %.6f", balance["usd"], balance["bitcoin"]))
	bot.Send(msg)
}
