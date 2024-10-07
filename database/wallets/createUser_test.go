package wallets

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

// Юнит-тест для функции Balance
func TestCreateUSer(t *testing.T) {
	// Подключаемся к тестовой базе данных
	conn, err := connectToTestDB()
	assert.NoError(t, err, "Не удалось подключиться к тестовой базе данных")
	defer conn.Close(context.Background())

	// Мокаем данные для update
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				UserName: "testuser",
			},
		},
	}

	// Вызываем функцию Balance
	err = CreateUser(conn, update)

	// Проверки
	assert.NoError(t, err, "Функция вернула ошибку")
}
