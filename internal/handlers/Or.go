package handlers

import (
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleOr обрабатывает команду /or
func HandleOr(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Извлекаем текст из сообщения и удаляем упоминания бота и команды
	var text string
	text = extractTextFromMessage(bot, message, "/or")
	text = strings.Replace(text, "?", "", -1)

	// Заменяем запятые на " или "
	text = strings.ReplaceAll(text, ",", " или ")

	// Разделяем на варианты
	options := strings.Split(text, " или ")
	if len(options) < 2 {
		sendMessage(bot, message, "Не удалось найти варианты для выбора")
		return
	}

	// Выбираем случайный вариант
	randomOption := options[rand.Intn(len(options))]
	response := randomOption

	sendMessage(bot, message, response)
}
