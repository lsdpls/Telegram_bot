package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InsultResponse представляет структуру ответа от API генерации оскорблений
type InsultResponse struct {
	Insult string `json:"insult"`
}

// HandleFuck обрабатывает команду /fuck
func HandleFuck(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Проверяем, есть ли сообщение, на которое был дан ответ
	if message.ReplyToMessage == nil {
		sendMessage(bot, message, "Используйте команду в ответ на сообщение")
		return
	}

	// Отправляем GET-запрос к API генерации оскорблений
	resp, err := http.Get("https://evilinsult.com/generate_insult.php?lang=ru&type=json")
	if err != nil {
		sendMessage(bot, message, "Не удалось отправить запрос")
		return
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	if resp.StatusCode != http.StatusOK {
		sendMessage(bot, message, "Не удалось придумать оскорбление")
		return
	}

	var insultResponse InsultResponse
	if err := json.NewDecoder(resp.Body).Decode(&insultResponse); err != nil {
		sendMessage(bot, message, "Ошибка при обработке ответа")
		return
	}

	// Формируем сообщение с оскорблением и тегом пользователя
	username := message.ReplyToMessage.From.UserName
	if username == "" {
		username = fmt.Sprintf("%d", message.ReplyToMessage.From.ID)
	}
	msgText := fmt.Sprintf("@%s %s", username, insultResponse.Insult)

	// Отправляем сообщение
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	msg.ReplyToMessageID = message.ReplyToMessage.MessageID
	bot.Send(msg)
}
