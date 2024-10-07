package main

import (
	"context"
	"crypto_bot/database"         // Ваш пакет для подключения к базе данных
	"crypto_bot/database/wallets" // Ваш пакет с функцией Balance
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

// Интеграционный тест для проверки полного процесса взаимодействия с ботом и базой данных
func TestBotBalanceIntegration(t *testing.T) {
	// Подключение к тестовой базе данных
	conn, err := database.Connect("postgres://postgres:111111@localhost:5432/test_db") // Измените строку подключения
	assert.NoError(t, err, "Не удалось подключиться к тестовой базе данных")
	defer conn.Close(context.Background())

	// Подготовка тестовых данных
	_, err = conn.Exec(context.Background(), `
        INSERT INTO wallets (tg_name, usd, bitcoin) VALUES 
        ('testuser', 1000.0, 0.5)
        ON CONFLICT (tg_name) DO UPDATE SET usd = EXCLUDED.usd, bitcoin = EXCLUDED.bitcoin;
    `)
	assert.NoError(t, err, "Не удалось подготовить тестовые данные")

	// Создание обновления с тестовыми данными
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{
				UserName: "testuser",
			},
		},
	}

	// Вызов функции Balance для получения баланса пользователя
	balance, err := wallets.Balance(conn, update)
	assert.NoError(t, err, "Ошибка при вызове функции Balance")
	assert.NotNil(t, balance, "Баланс не должен быть nil")
	assert.Equal(t, 1000.0, balance["usd"], "Неверное значение usd")
	assert.Equal(t, 0.5, balance["bitcoin"], "Неверное значение bitcoin")

	// Удаление тестовых данных
	_, err = conn.Exec(context.Background(), "DELETE FROM wallets WHERE tg_name = 'testuser';")
	assert.NoError(t, err, "Не удалось очистить тестовые данные")
}
