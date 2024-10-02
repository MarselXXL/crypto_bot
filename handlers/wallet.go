package handlers

import (
	"crypto_bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// HandleWallet создает запись в БД wallets и отвечает
func HandleWallet(bot *tgbotapi.BotAPI, chatID int64, dbConn *pgx.Conn, update tgbotapi.Update) {
	err := database.CreateUser(dbConn, update)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		//log.Printf("%T %v", pgErr, pgErr)
		msg := tgbotapi.NewMessage(chatID, "Действия в кошельке:\n/add_test_balance - добавить тествый баланс в кошелек,\n/balance - Ваш баланс,\n/buy - купить биткоин по текущей цене,\n/sell - продать по текущей,\n/send - оправить крипту другому пользователю")
		bot.Send(msg)
		return
	} else if err != nil {
		//log.Printf("%T %v", pgErr, pgErr)
		msg := tgbotapi.NewMessage(chatID, "Ошибка при запросе к таблице wallets")
		bot.Send(msg)
		return
	}
	msg := tgbotapi.NewMessage(chatID, "Кошелек создан\nДействия в кошельке:\n/add_test_balance - добавить тествый баланс в кошелек,\n/balance - Ваш баланс,\n/buy - купить биткоин по текущей цене,\n/sell - продать по текущей,\n/send - оправить крипту другому пользователю")
	bot.Send(msg)
}
