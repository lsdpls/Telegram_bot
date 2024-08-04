package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Joke представляет структуру шутки
type Joke struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

// HandleJoke обрабатывает команду /joke
func HandleJoke(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Запрос к API для получения случайной шутки
	resp, err := http.Get("https://official-joke-api.appspot.com/jokes/random")
	if err != nil {
		sendMessage(bot, message, "Не удалось получить шутку")
		return
	}
	defer resp.Body.Close()

	// Обрабатываем ответ от API
	if resp.StatusCode != http.StatusOK {
		sendMessage(bot, message, "Не удалось получить шутку")
		return
	}

	var joke Joke
	if err := json.NewDecoder(resp.Body).Decode(&joke); err != nil {
		sendMessage(bot, message, "Ошибка при обработке ответа")
		return
	}

	// Формируем сообщение с шуткой
	jokeMessage := fmt.Sprintf("%s\n\n<tg-spoiler>%s</tg-spoiler>", joke.Setup, joke.Punchline) // Панчлайн под спойлером

	sendMessageParse(bot, message, jokeMessage)
}
