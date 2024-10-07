package wallets

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

// Функция для подключения к тестовой базе данных
func connectToTestDB() (*pgx.Conn, error) {
	connString := "postgres://postgres:111111@localhost:5432/test_db" // Измените на ваши данные
	return pgx.Connect(context.Background(), connString)
}

// Подготовка тестовых данных в таблице wallets
func setupTestData(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		INSERT INTO wallets (tg_name, usd, bitcoin) VALUES 
		('testuser', 1000.0, 0.5)
		ON CONFLICT (tg_name) DO UPDATE SET usd = EXCLUDED.usd, bitcoin = EXCLUDED.bitcoin;
	`)
	return err
}

// Удаление тестовых данных из таблицы wallets
func cleanupTestData(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM wallets WHERE tg_name = 'testuser';")
	return err
}

// Юнит-тест для функции Balance
func TestBalance(t *testing.T) {
	// Подключаемся к тестовой базе данных
	conn, err := connectToTestDB()
	assert.NoError(t, err, "Не удалось подключиться к тестовой базе данных")
	defer conn.Close(context.Background())

	// Подготавливаем тестовые данные
	err = setupTestData(conn)
	assert.NoError(t, err, "Не удалось подготовить тестовые данные")
	defer cleanupTestData(conn) // Удаляем данные после теста

	// Мокаем данные для update
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				UserName: "testuser",
			},
		},
	}

	// Вызываем функцию Balance
	balance, err := Balance(conn, update)

	// Проверки
	assert.NoError(t, err, "Функция вернула ошибку")
	assert.NotNil(t, balance, "Баланс не должен быть nil")
	assert.Equal(t, 1000.0, balance["usd"], "Неверное значение usd")
	assert.Equal(t, 0.5, balance["bitcoin"], "Неверное значение bitcoin")
}
