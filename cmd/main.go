package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"telegram_bot/internal/bot"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Подключаем бота
	botapi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Режим отладки
	botapi.Debug = false

	// Устанавливаем обработчик для обновлений
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bot.HandleWebhook(botapi, w, r)
	})

	// Запускаем HTTP сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
