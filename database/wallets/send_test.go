package wallets

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

// Подгатавливаеи получателя
func setupReciever(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO wallets (tg_name, usd, bitcoin) VALUES 
		('test_user2', 0, 0)
		ON CONFLICT (tg_name) DO UPDATE SET usd = EXCLUDED.usd, bitcoin = EXCLUDED.bitcoin;
	`)
	return err
}

// Удаление тестовых данных из таблицы wallets
func cleanupReciever(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM wallets WHERE tg_name = 'test_user2';")
	return err
}

// Юнит-тест для функции Balance
func TestSend(t *testing.T) {
	// Подключаемся к тестовой базе данных
	conn, err := connectToTestDB()
	assert.NoError(t, err, "Не удалось подключиться к тестовой базе данных")
	defer conn.Close(context.Background())

	// Подготавливаем тестовые данные
	err = setupTestData(conn)
	assert.NoError(t, err, "Не удалось подготовить тестовые данные")
	defer cleanupTestData(conn) // Удаляем данные после теста

	// Подготавливаем получаетя
	err = setupReciever(conn)
	assert.NoError(t, err, "Не удалось подготовить получаетя")
	defer cleanupReciever(conn) // Удаляем данные после теста

	// Мокаем данные для update
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				UserName: "testuser",
			},
		},
	}

	// Вызываем функцию Send
	err = Send(conn, update, "usd", 500, "test_user2")

	// Проверки
	assert.NoError(t, err, "Функция вернула ошибку")
}
