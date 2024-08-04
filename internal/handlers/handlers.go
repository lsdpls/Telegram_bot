package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleStart обрабатывает команду /start
func HandleStart(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Сводка по доступным командам
	summary := `
Привет! Я Жопич. Вот список доступных команд:

/who - Выбрать случайного человека из чата ("/who кто виноват?", а также другие склонения "кто" или "чей")
/or - Выбрать один из вариантов, разделенных "или"/"," ("/or macOS, Windows или Linux")
/meme - Создать мем из текста сообщения (смех не гарантируется)
/wiki - Найти статью в Википедии (поиск умный, но не всесильный)
/yesnot - Отправить случайную гифку "да", "нет" или "возможно"
/fuck - Отправить случайное оскорбление
/joke - Рассказать шутку
/t9 - Проверить орфографию текста

Команды вызывать реплаем на сообщение или писать "/команда текст_команды"!
`

	msg := tgbotapi.NewMessage(message.Chat.ID, summary)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func HandleTest(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var text string
	if message.ReplyToMessage != nil {
		text = message.ReplyToMessage.Text
	} else {
		text = message.Text
	}

	var messageID int
	if message.ReplyToMessage != nil {
		messageID = message.ReplyToMessage.MessageID
	} else {
		messageID = message.MessageID
	}

	// Отправляем сообщение
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = messageID
	bot.Send(msg)
}
