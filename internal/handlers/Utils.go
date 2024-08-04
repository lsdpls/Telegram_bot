package handlers

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// extractTextFromMessage извлекает текст из сообщения и удаляет упоминания бота и команды
func extractTextFromMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, command string) string {
	var text string
	if message.ReplyToMessage != nil {
		text = message.ReplyToMessage.Text
	} else {
		text = message.Text
	}

	text = strings.Replace(text, command+"@"+bot.Self.UserName, "", -1)
	text = strings.Replace(text, command, "", -1)
	return strings.TrimSpace(text) // Удаление начальных и конечных пробелов
}

// sendMessage отправляет сообщение пользователю
func sendMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, response string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

// sendMessageParse отправляет сообщение пользователю с парсингом HTML
func sendMessageParse(bot *tgbotapi.BotAPI, message *tgbotapi.Message, response string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func sendMessagePhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message, photoURL string) {
	photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(photoURL))
	photo.ReplyToMessageID = message.MessageID
	bot.Send(photo)
}

func sendMessageParsePhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message, photoURL string, caption string) {
	photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileURL(photoURL))
	photo.Caption = caption
	photo.ParseMode = "HTML"
	photo.ReplyToMessageID = message.MessageID
	bot.Send(photo)
}
