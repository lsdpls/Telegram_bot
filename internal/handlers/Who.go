package handlers

import (
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var pronouns = []string{"кто", "кого", "кому", "кем", "чей", "чья", "чье", "чьё", "чьи", "чью"}

// HandleWho обрабатывает команду /who
func HandleWho(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Получаем список участников чата
	chatMembers, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: message.Chat.ID},
	})
	if err != nil || len(chatMembers) == 0 {
		sendMessage(bot, message, "Не удалось получить список участников чата")
		return
	}

	// Выбираем случайного участника
	randomMember := chatMembers[rand.Intn(len(chatMembers))]
	username := "@" + randomMember.User.UserName

	// Извлекаем текст из сообщения и удаляем упоминания бота и команды
	var text string
	text = extractTextFromMessage(bot, message, "/who")
	text = strings.ToLower(text)

	// Убираем местояимения и "?" из текста
	for _, pronoun := range pronouns {
		text = strings.Replace(text, pronoun, "", -1)
	}
	text = strings.Replace(text, "?", "", -1)

	response := username + " " + text

	sendMessage(bot, message, response)
}
