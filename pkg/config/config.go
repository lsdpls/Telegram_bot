package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func LoadConfig() Config {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Возвращаем конфигурацию
	return Config{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}
