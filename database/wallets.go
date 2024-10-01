package database

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func CreateUser(conn *pgx.Conn, update tgbotapi.Update) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO wallets (tg_name) VALUES ($1)", update.Message.From.UserName)
	return err
}
