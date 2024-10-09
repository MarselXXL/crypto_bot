-- Создание таблицы wallets
CREATE TABLE wallets (
    user_id SERIAL PRIMARY KEY,          -- Автоинкрементное поле user_id
    tg_name TEXT NOT NULL UNIQUE,        -- Поле tg_name, уникальное и обязательное
    usd NUMERIC DEFAULT 0 CHECK (usd >= 0),   -- Поле usd, с проверкой на неотрицательное значение
    bitcoin NUMERIC DEFAULT 0 CHECK (bitcoin >= 0)  -- Поле bitcoin, с проверкой на неотрицательное значение
);

-- Создание таблицы crypto_prices
CREATE TABLE crypto_prices (
    id SERIAL PRIMARY KEY,                -- Автоинкрементное поле id
    currency VARCHAR(50) NOT NULL,       -- Поле currency, обязательное
    price NUMERIC NOT NULL,               -- Поле price, обязательное
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Поле created_at, по умолчанию текущее время
);
