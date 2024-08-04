package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"telegram_bot/internal/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleWebhook(bot *tgbotapi.BotAPI, w http.ResponseWriter, r *http.Request) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Error decoding update: %v", err)
		http.Error(w, "Error decoding update", http.StatusBadRequest)
		return
	}

	if update.Message != nil && update.Message.IsCommand() {
		handleCommand(bot, update.Message)
	}

	// Отправляем положительный ответ Telegram, чтобы он не повторял запрос
	w.WriteHeader(http.StatusOK)
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch strings.ToLower(message.Command()) {
	case "who":
		handlers.HandleWho(bot, message)
	case "or":
		handlers.HandleOr(bot, message)
	case "meme":
		handlers.HandleMeme(bot, message)
	case "wiki":
		handlers.HandleWikit9(bot, message)
	case "yesnot":
		handlers.HandleYesNoT(bot, message)
	case "fuck":
		handlers.HandleFuck(bot, message)
	case "joke":
		handlers.HandleJoke(bot, message)
	case "t9":
		handlers.HandleT9(bot, message)
	case "when":
		handlers.HandleWhen(bot, message)
	case "test":
		handlers.HandleTest(bot, message)
	case "start":
		handlers.HandleStart(bot, message)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command")
		msg.ReplyToMessageID = message.MessageID
		bot.Send(msg)
	}
}
